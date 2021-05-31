package notion

import (
	"fmt"
	"log"
)

const (
	TypeBlockParagraph        = "paragraph"
	TypeBlockHeading1         = "heading_1"
	TypeBlockHeading2         = "heading_2"
	TypeBlockHeading3         = "heading_3"
	TypeBlockBulletedListItem = "bulleted_list_item"
	TypeBlockNumberedListItem = "numbered_list_item"
	TypeBlockTodo             = "to_do"
	TypeBlockToggle           = "toggle"
	TypeBlockChildPage        = "child_page"
	TypeBlockUnsupported      = "unsupported"
)

type Block interface {
	Object() string
	ID() string
	CreatedTime() string
	LastEditedTime() string
	HasChildren() bool
	Type() string

	Interface() interface{}

	Json() JSON
}

func AssignBlock(json JSON) (Block, error) {
	var block Block

	t := json.Get("type")

	switch json.GetString("type") {
	case TypeBlockParagraph:
		block = &BlockParagraph{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}
	case TypeBlockHeading1:
		block = &BlockHeading1{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}
	case TypeBlockHeading2:
		block = &BlockHeading2{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}
	case TypeBlockHeading3:
		block = &BlockHeading3{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}
	case TypeBlockBulletedListItem:
		block = &BlockBulletedListItem{&ChildrenBlock{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}}
	case TypeBlockNumberedListItem:
		block = &BlockNumberedListItem{&ChildrenBlock{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}}
	case TypeBlockTodo:
		block = &BlockTodo{&ChildrenBlock{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}}
	case TypeBlockToggle:
		block = &BlockToggle{&ChildrenBlock{&RichTextBlock{&CustomBlock{id: json.GetString("id"), JSON: json}}}}
	case TypeBlockChildPage:
		block = &BlockChildPage{&CustomBlock{id: json.GetString("id"), JSON: json}}
	case TypeBlockUnsupported:
		block = &BlockUnsupported{&CustomBlock{id: json.GetString("id"), JSON: json}}
	default:
		return nil, fmt.Errorf("invalid type: '%s'", t)
	}
	return block, nil
}

type CustomBlock struct {
	id string

	JSON JSON
}

func (block *CustomBlock) Object() string {
	return "block"
}

func (block *CustomBlock) ID() string {
	return block.id
}

func (block *CustomBlock) CreatedTime() string {
	return block.JSON.GetString("created_time")
}

func (block *CustomBlock) LastEditedTime() string {
	return block.JSON.GetString("last_edited_time")
}

func (block *CustomBlock) HasChildren() bool {
	return block.JSON.GetBool("has_children")
}

func (block *CustomBlock) Type() string {
	return block.JSON.GetString("type")
}

func (block *CustomBlock) Json() JSON {
	return block.JSON
}

func (block *CustomBlock) Interface() interface{} {
	return block
}

func NewBlock(j JSON) Block {
	block := &CustomBlock{
		id:   j.GetString("id"),
		JSON: j,
	}

	return block
}

type RichTextBlock struct {
	*CustomBlock
}

func (block *RichTextBlock) Interface() interface{} {
	return block
}

func (block *RichTextBlock) AddText(text []RichText) error {
	t := block.Type()
	if len(t) == 0 {
		return fmt.Errorf("not found type")
	}

	v, ok := block.JSON.GetJSON(t)
	if !ok {
		return fmt.Errorf("not found '%s' field", t)
	}

	_, ok = v.GetJSONList("text")
	if !ok {
		return fmt.Errorf("not found text field")
	}

	for _, t := range text {
		v.Append("text", t.JSON)
	}

	return nil
}

func (block *RichTextBlock) ListText() ([]RichText, error) {
	t := block.Type()
	if len(t) == 0 {
		return nil, fmt.Errorf("not found type")
	}

	v, ok := block.JSON.GetJSON(t)
	if !ok {
		return nil, fmt.Errorf("not found '%s' field", t)
	}

	vv, ok := v.GetJSONList("text")
	if !ok {
		return nil, fmt.Errorf("not found text field")
	}

	list := make([]RichText, 0)

	for _, j := range vv {
		rt := RichText{}
		if j.Unmarshal(&rt.JSON) == nil {
			list = append(list, rt)
		}
	}

	return list, nil
}

func newRichTextBlock(Type string, Text []RichText) *RichTextBlock {
	block := &RichTextBlock{
		CustomBlock: &CustomBlock{
			id: "",
			JSON: JSON{
				"object": "block",
				"type":   Type,
				Type: JSON{
					"text": make([]JSON, 0),
				},
			},
		},
	}

	if err := block.AddText(Text); err != nil {
		log.Print(err)
	}
	return block
}

type ChildrenBlock struct {
	*RichTextBlock
}

func (block *ChildrenBlock) Interface() interface{} {
	return block
}

func newChildrenBlock(Type string, Text []RichText, Children ...Block) *ChildrenBlock {
	base := newRichTextBlock(Type, Text)

	block := &ChildrenBlock{
		RichTextBlock: base,
	}

	block.JSON.Set("children", []JSON{})
	block.AddChildren(Children)

	return block
}

func (block *ChildrenBlock) AddChildren(children []Block) error {
	t := block.Type()
	if len(t) == 0 {
		return fmt.Errorf("not found type")
	}

	v, ok := block.JSON.GetJSON(t)
	if !ok {
		return fmt.Errorf("not found '%s' field", t)
	}

	_, ok = v.GetJSONList("children")
	if !ok {
		return fmt.Errorf("not found 'children' field")
	}

	for _, t := range children {
		v.Append("text", t.Json())
	}

	return nil
}

func (block *ChildrenBlock) Children() ([]Block, error) {
	t := block.Type()
	if len(t) == 0 {
		return nil, fmt.Errorf("not found type")
	}

	v, ok := block.JSON.GetJSON(t)
	if !ok {
		return nil, fmt.Errorf("not found '%s' field", t)
	}

	vv, ok := v.GetJSONList("children")
	if !ok {
		return nil, fmt.Errorf("not found 'children' field")
	}

	blocks := make([]Block, 0)

	for _, j := range vv {
		block, err := AssignBlock(j)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

type BlockParagraph struct {
	*RichTextBlock
}

func (block *BlockParagraph) Interface() interface{} {
	return block
}

func NewBlockParagraph(Text []RichText) Block {
	block := &BlockParagraph{
		RichTextBlock: newRichTextBlock("paragraph", Text),
	}
	return block
}

type BlockHeading1 struct {
	*RichTextBlock
}

func (block *BlockHeading1) Interface() interface{} {
	return block
}

func NewBlockHeading1(Text []RichText) Block {
	block := &BlockHeading1{
		RichTextBlock: newRichTextBlock("heading_1", Text),
	}

	return block
}

type BlockHeading2 struct {
	*RichTextBlock
}

func (block *BlockHeading2) Interface() interface{} {
	return block
}

func NewBlockHeading2(Text []RichText) Block {
	block := &BlockHeading2{
		RichTextBlock: newRichTextBlock("heading_2", Text),
	}

	return block
}

type BlockHeading3 struct {
	*RichTextBlock
}

func (block *BlockHeading3) Interface() interface{} {
	return block
}

func NewBlockHeading3(Text []RichText) Block {
	block := &BlockHeading3{
		RichTextBlock: newRichTextBlock("heading_3", Text),
	}

	return block
}

type BlockBulletedListItem struct {
	*ChildrenBlock
}

func (block *BlockBulletedListItem) Interface() interface{} {
	return block
}

func NewBlockBulletedListItem(Text []RichText, Children ...Block) Block {
	block := &BlockBulletedListItem{
		ChildrenBlock: newChildrenBlock("bulleted_list_item", Text, Children...),
	}

	return block
}

type BlockNumberedListItem struct {
	*ChildrenBlock
}

func (block *BlockNumberedListItem) Interface() interface{} {
	return block
}

func NewBlockNumberedListItem(Text []RichText, Children ...Block) Block {
	block := &BlockBulletedListItem{
		ChildrenBlock: newChildrenBlock("bulleted_list_item", Text, Children...),
	}

	return block
}

type BlockTodo struct {
	*ChildrenBlock
}

func (block *BlockTodo) Interface() interface{} {
	return block
}

func NewBlockTodo(Checked bool, Text []RichText, Children ...Block) Block {
	block := &BlockBulletedListItem{
		ChildrenBlock: newChildrenBlock("to_do", Text, Children...),
	}
	if j, ok := block.JSON.GetJSON("to_do"); ok {
		j["checked"] = Checked
	}

	return block
}

func (block *BlockTodo) IsChecked() bool {
	if j, ok := block.JSON.GetJSON("to_do"); ok {
		return j.GetBool("checked")
	}

	return false
}

func (block *BlockTodo) Checked() {
	if j, ok := block.JSON.GetJSON("to_do"); ok {
		j.Set("checked", true)
	}
}

func (block *BlockTodo) Unchecked() {
	if j, ok := block.JSON.GetJSON("to_do"); ok {
		j.Set("checked", false)
	}
}

type BlockToggle struct {
	*ChildrenBlock
}

func (block *BlockToggle) Interface() interface{} {
	return block
}

func NewBlockToggle(Text []RichText, Children ...Block) Block {
	block := &BlockBulletedListItem{
		ChildrenBlock: newChildrenBlock("toggle", Text, Children...),
	}

	return block
}

type BlockChildPage struct {
	*CustomBlock
}

func (block *BlockChildPage) Interface() interface{} {
	return block
}

func NewBlockChildPage(PageID string) Block {
	block := &BlockChildPage{
		CustomBlock: &CustomBlock{
			id: "",
			JSON: JSON{
				"object": "block",
				"type":   "child_page",
				"child_page": JSON{
					"title": PageID,
				},
			},
		},
	}

	return block
}

type BlockUnsupported struct {
	*CustomBlock
}

func (block *BlockUnsupported) Interface() interface{} {
	return block
}

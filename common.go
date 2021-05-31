package notion

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PaginationRequest struct {
	StartCursor string `json:"start_cursor"`
	PageSize    int    `json:"page_size"`
}

func (request *PaginationRequest) QueryString() string {
	param := make(map[string]interface{})

	paramString := make([]string, 0)

	if len(request.StartCursor) > 0 {
		param["start_cursor"] = request.StartCursor
	}
	if request.PageSize > 0 {
		param["page_size"] = request.PageSize
	}

	for k, v := range param {
		kv := fmt.Sprintf("%s=%v", k, v)

		paramString = append(paramString, kv)
	}

	return strings.Join(paramString, "&")
}

func (request *PaginationRequest) Json() JSON {
	j := JSON{}
	if len(request.StartCursor) > 0 {
		j["start_cursor"] = request.StartCursor
	}
	if request.PageSize > 0 {
		j["page_size"] = request.PageSize
	}

	return j
}

type PaginationResponse struct {
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
	Object     string `json:"object"` //always 'list'

	Results []interface{} `json:"results"`
}

func (response *PaginationResponse) Unmarshal(v interface{}) error {
	b, err := json.Marshal(response.Results)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (response *PaginationResponse) Users() ([]User, error) {
	users := make([]User, 0)

	results := make([]JSON, 0)

	if err := response.Unmarshal(&results); err != nil {
		return nil, err
	}
	for _, result := range results {
		user := User{
			ID:   result.GetString("id"),
			JSON: result,
		}
		if user.JSON.GetString("object") == "user" {
			users = append(users, user)
		}
	}

	return users, nil
}

func (response *PaginationResponse) Blocks() ([]Block, error) {
	blocks := make([]Block, 0)

	results := make([]JSON, 0)

	if err := response.Unmarshal(&results); err != nil {
		return nil, err
	}
	for _, result := range results {
		block, err := AssignBlock(result)
		if err != nil {
			continue
		}
		if block.Json().GetString("object") == "block" {
			blocks = append(blocks, block)
		}
	}

	return blocks, nil
}

func (response *PaginationResponse) Pages() ([]Page, error) {
	pages := make([]Page, 0)

	results := make([]JSON, 0)

	if err := response.Unmarshal(&results); err != nil {
		return nil, err
	}
	for _, result := range results {
		page := Page{}
		if err := result.Unmarshal(&page.JSON); err != nil {
			continue
		}
		if page.JSON.GetString("object") == "page" {
			pages = append(pages, page)
		}
	}

	return pages, nil
}

func (response *PaginationResponse) Databases() ([]Database, error) {
	databases := make([]Database, 0)

	results := make([]JSON, 0)
	if err := response.Unmarshal(&results); err != nil {
		return nil, err
	}

	for _, result := range results {
		database := Database{}
		if err := result.Unmarshal(&database.JSON); err != nil {
			continue
		}

		if database.JSON.GetString("object") == "database" {
			databases = append(databases, database)
		}
	}

	return databases, nil
}

type Date struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type RichText struct {
	JSON JSON
}

func NewRichText(PlainText string) *RichText {
	return &RichText{
		JSON: JSON{
			"type": "text",
			"text": &Text{
				Content: PlainText,
				Link:    nil,
			},
			"plain_text": PlainText,
			"href":       nil,
		},
	}
}

func (rt *RichText) Type() string {
	if rt.JSON == nil {
		return ""
	}
	return rt.JSON.GetString("type")
}

func (rt *RichText) PlainText() string {
	if rt.JSON == nil {
		return ""
	}

	return rt.JSON.GetString("plain_text")
}

func (rt *RichText) Href() string {
	if rt.JSON == nil {
		return ""
	}

	return rt.JSON.GetString("href")
}

func (rt *RichText) GetAnnotations() (*Annotations, error) {
	if rt.JSON == nil {
		return nil, fmt.Errorf("Invalid RichText Object")
	}

	annotations := &Annotations{}
	j, ok := rt.JSON.GetJSON("annotations")
	if !ok {
		return nil, fmt.Errorf("'annotations' is nil")
	}

	if err := j.Unmarshal(annotations); err != nil {
		return nil, err
	}

	return annotations, nil
}

func (rt *RichText) SetAnnotations(annotations *Annotations) {
	if annotations != nil {
		rt.JSON["annotations"] = *annotations
	}
}

func (rt *RichText) getObject(name string, v interface{}) error {
	if rt.JSON == nil {
		return fmt.Errorf("Invalid RichText Object")
	}

	if rt.Type() != name {
		return fmt.Errorf("type of RichText is not '%s'", name)
	}

	j, ok := rt.JSON.GetJSON(name)
	if !ok {
		return fmt.Errorf("Invalid '%s' Object", name)
	}

	return j.Unmarshal(v)
}

func (rt *RichText) setObject(name string, v interface{}) error {
	if rt.JSON == nil {
		return fmt.Errorf("Invalid RichText Object")
	} else if v == nil {
		return fmt.Errorf("object is nil pointer")
	}

	for _, name := range []string{"text", "mention", "equation", "date"} {
		delete(rt.JSON, name)
	}

	rt.JSON["type"] = name
	rt.JSON.Set(name, v)

	return nil
}

func (rt *RichText) getMention(name string, v interface{}) error {
	mention := make(JSON)
	if err := rt.getObject("mention", &mention); err != nil {
		return err
	}

	j, ok := mention.GetJSON(name)
	if !ok {
		return fmt.Errorf("not found '%s' Object", name)
	}

	if err := j.Unmarshal(v); err != nil {
		return err
	}

	return nil
}

func (rt *RichText) setMention(name string, v interface{}) error {
	j := JSON{}
	if err := j.Marshal(v); err != nil {
		return err
	}
	mention := JSON{
		"type": name,
		name:   j,
	}

	return rt.setObject("mention", mention)
}

func (rt *RichText) GetText() (*Text, error) {
	text := &Text{}

	if err := rt.getObject("text", text); err != nil {
		return nil, err
	}

	return text, nil
}

func (rt *RichText) SetText(text *Text) {
	rt.setObject("text", *text)
}

func (rt *RichText) GetEquation() (*Equation, error) {
	equation := &Equation{}

	if err := rt.getObject("equation", equation); err != nil {
		return nil, err
	}

	return equation, nil
}

func (rt *RichText) SetEquation(equation *Equation) {
	rt.setObject("equation", equation)
}

func (rt *RichText) GetMentionUser() (*User, error) {
	user := &User{}
	if err := rt.getMention("user", user); err != nil {
		return nil, err
	}
	return user, nil
}

func (rt *RichText) SetMentionUser(user *User) error {
	return rt.setMention("user", user)
}

func (rt *RichText) GetMentionPage() (*Page, error) {
	page := &Page{}
	if err := rt.getMention("page", &page.JSON); err != nil {
		return nil, err
	}

	return page, nil
}

func (rt *RichText) SetMentionPage(page *Page) error {
	j := JSON{
		"id": page.ID(),
	}

	return rt.setMention("page", j)
}

func (rt *RichText) GetMentionDatabase() (*Database, error) {
	database := &Database{}
	if err := rt.getMention("database", &database.JSON); err != nil {
		return nil, err
	}

	return database, nil
}

func (rt *RichText) SetMentionDatabase(database *Database) error {
	j := JSON{
		"id": database.ID(),
	}

	return rt.setMention("database", j)
}

func (rt *RichText) GetMentionDate() (*Date, error) {
	date := &Date{}

	if err := rt.getMention("date", date); err != nil {
		return nil, err
	}

	return date, nil
}

func (rt *RichText) SetMentionDate(date *Date) error {
	return rt.setMention("date", date)
}

type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link"`
}

type Link struct {
	URL string `json:"url"`
}

type Equation struct {
	Expression string `json:"expression"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

//NewAnnotations Create Annotations Struct, set color with "default"
func NewAnnotations() *Annotations {
	return &Annotations{
		Bold:          false,
		Italic:        false,
		Strikethrough: false,
		Underline:     false,
		Code:          false,
		Color:         "default",
	}
}

type SelectOption struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

type File struct {
	Name string `json:"name"`
}

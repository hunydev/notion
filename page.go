package notion

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TypeParentDatabase  = "database_id"
	TypeParentPage      = "page_id"
	TypeParentWorkspace = "workspace"
)

const (
	TypePropertyTitle       = "title"
	TypePropertyRichText    = "rich_text"
	TypePropertyNumber      = "number"
	TypePropertySelect      = "select"
	TypePropertyMultiSelect = "multi_select"
	TypePropertyDate        = "date"
	TypePropertyFormula     = "formula"
	TypePropertyPeople      = "people"
	TypePropertyFiles       = "files"
	TypePropertyCheckbox    = "checkbox"
	TypePropertyURL         = "url"
	TypePropertyEmail       = "email"
	TypePropertyPhoneNumber = "phone_number"
)

type Page struct {
	JSON JSON
}

func NewPage(parent *Parent) *Page {
	if parent == nil {
		return nil
	}

	page := &Page{
		JSON: JSON{
			"parent": parent.JSON,
		},
	}

	return page
}

func (page *Page) Object() string {
	return "page"
}

func (page *Page) ID() string {
	return page.JSON.GetString("id")
}

func (page *Page) CreatedTime() string {
	return page.JSON.GetString("created_time")
}

func (page *Page) LastEditedTime() string {
	return page.JSON.GetString("last_edited_time")
}

func (page *Page) Archived() bool {
	return page.JSON.GetBool("archived")
}

func (page *Page) Parent() *Parent {
	j, ok := page.JSON.GetJSON("parent")
	if !ok {
		return nil
	}

	switch j.GetString("type") {
	case TypeParentDatabase:
		return NewParentDatabase(j.GetString(TypeParentDatabase))
	case TypeParentPage:
		return NewParentPage(j.GetString(TypeParentPage))
	case TypeParentWorkspace:
		return NewParentWorkspace()
	}

	return nil
}

func (page *Page) Properties() []Property {
	properties := []Property{}

	j, ok := page.JSON.GetJSON("properties")
	if !ok {
		return nil
	}

	for k, v := range j {
		jj := JSON{}
		if jj.Marshal(v) != nil {
			continue
		}

		property, err := AssignProperty(k, jj)
		if err != nil {
			continue
		}

		properties = append(properties, property)
	}

	return properties
}

type Parent struct {
	ID string

	JSON JSON
}

func NewParentPage(ID string) *Parent {
	return &Parent{ID: ID, JSON: JSON{"type": TypeParentPage, TypeParentPage: ID}}
}

func NewParentDatabase(ID string) *Parent {
	return &Parent{ID: ID, JSON: JSON{"type": TypeParentDatabase, TypeParentDatabase: ID}}
}

func NewParentWorkspace() *Parent {
	return &Parent{ID: "", JSON: JSON{"type": TypeParentWorkspace, TypeParentWorkspace: true}}
}

type Property interface {
	Name() string
	ID() string
	Type() string

	Json() JSON
	Interface() interface{}
}

func NewProperty(j JSON) Property {
	return &BaseProperty{
		JSON: j,
	}
}

func AssignProperty(name string, json JSON) (Property, error) {
	var property Property

	t := json.GetString("type")

	switch t {
	case TypePropertyTitle:
		property = &PropertyTitle{&RichTextProperty{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyRichText:
		property = &PropertyRichText{&RichTextProperty{&BaseProperty{name: name, JSON: json}}}
	case TypePropertyNumber:
		property = &PropertyNumber{&BaseProperty{name: name, JSON: json}}
	case TypePropertySelect:
		property = &PropertySelect{&BaseProperty{name: name, JSON: json}}
	case TypePropertyMultiSelect:
		property = &PropertyMultiSelect{&BaseProperty{name: name, JSON: json}}
	case TypePropertyDate:
		property = &PropertyDate{&BaseProperty{name: name, JSON: json}}
	case TypePropertyFormula:
		property = &PropertyFormula{&BaseProperty{name: name, JSON: json}}
	case TypePropertyPeople:
		property = &PropertyPeople{&BaseProperty{name: name, JSON: json}}
	case TypePropertyFiles:
		property = &PropertyFiles{&BaseProperty{name: name, JSON: json}}
	case TypePropertyCheckbox:
		property = &PropertyCheckbox{&BaseProperty{name: name, JSON: json}}
	case TypePropertyURL:
		property = &PropertyURL{&BaseProperty{name: name, JSON: json}}
	case TypePropertyEmail:
		property = &PropertyEmail{&BaseProperty{name: name, JSON: json}}
	case TypePropertyPhoneNumber:
		property = &PropertyPhoneNumber{&BaseProperty{name: name, JSON: json}}
	default:
		return nil, fmt.Errorf("invalid type: '%s'", t)
	}
	return property, nil
}

type BaseProperty struct {
	name string

	JSON JSON
}

func newBaseProperty(Name, ID, Type string, v interface{}) *BaseProperty {
	base := &BaseProperty{
		name: Name,
		JSON: JSON{
			"type": Type,
			Type:   v,
		},
	}

	if len(strings.TrimSpace(ID)) > 0 {
		base.JSON["id"] = ID
	}
	return base
}

func (property *BaseProperty) Name() string {
	return property.name
}

func (property *BaseProperty) ID() string {
	return property.JSON.GetString("id")
}

func (property *BaseProperty) Type() string {
	return property.JSON.GetString("type")
}

func (property *BaseProperty) Json() JSON {
	return property.JSON
}

func (property *BaseProperty) Interface() interface{} {
	return property
}

type RichTextProperty struct {
	*BaseProperty
}

func newRichTextProperty(Name, ID, Type string, Text []RichText) *RichTextProperty {
	property := &RichTextProperty{
		BaseProperty: newBaseProperty(Name, ID, Type, []JSON{}),
	}

	for _, t := range Text {
		property.JSON.Append(Type, t.JSON)
	}

	return property
}

func (property *RichTextProperty) Interface() interface{} {
	return property
}

func (property *RichTextProperty) RichText() []RichText {
	j, ok := property.JSON.GetJSONList(property.Type())
	if !ok {
		return nil
	}

	list := []RichText{}

	for _, jj := range j {
		rt := RichText{}

		if jj.Unmarshal(&rt.JSON) == nil {
			list = append(list, rt)
		}
	}

	return list
}

type PropertyTitle struct {
	*RichTextProperty
}

func NewPropertyTitle(Name string, text []RichText) Property {
	property := &PropertyTitle{
		RichTextProperty: newRichTextProperty(Name, "title", "title", text),
	}

	return property
}

func (property *PropertyTitle) Interface() interface{} {
	return property
}

type PropertyRichText struct {
	*RichTextProperty
}

func NewPropertyRichText(Name string, text []RichText) Property {
	property := &PropertyRichText{
		RichTextProperty: newRichTextProperty(Name, "", "rich_text", text),
	}

	return property
}

func (property *PropertyRichText) Interface() interface{} {
	return property
}

type PropertyNumber struct {
	*BaseProperty
}

func NewPropertyNumber(Name string, Number int) Property {
	property := &PropertyNumber{
		BaseProperty: newBaseProperty(Name, "", "number", Number),
	}
	return property
}

func (property *PropertyNumber) Interface() interface{} {
	return property
}

func (property *PropertyNumber) Number() (int, error) {
	return property.JSON.GetInt("number"), nil
}

type PropertySelect struct {
	*BaseProperty
}

func NewPropertySelect(Name string, Select *SelectOption) Property {
	property := &PropertySelect{
		BaseProperty: newBaseProperty(Name, "", "select", Select),
	}
	return property
}

func (property *PropertySelect) Interface() interface{} {
	return property
}

func (property *PropertySelect) Option() *SelectOption {
	j, ok := property.JSON.GetJSON(property.Type())
	if !ok {
		return nil
	}

	option := &SelectOption{}
	if j.Unmarshal(option) != nil {
		return nil
	}
	return option
}

type PropertyMultiSelect struct {
	*BaseProperty
}

func NewPropertyMultiSelect(Name string, Select ...SelectOption) Property {
	list := []JSON{}

	for _, s := range Select {
		j := JSON{}
		j.Marshal(s)
		list = append(list, j)
	}

	property := &PropertyMultiSelect{
		BaseProperty: newBaseProperty(Name, "", "multi_select", list),
	}

	return property
}

func (property *PropertyMultiSelect) Interface() interface{} {
	return property
}

func (property *PropertyMultiSelect) Options() []SelectOption {
	j, ok := property.JSON.GetJSONList(property.Type())
	if !ok {
		return nil
	}

	list := []SelectOption{}

	for _, jj := range j {
		s := SelectOption{}

		if jj.Unmarshal(&s) == nil {
			list = append(list, s)
		}
	}

	return list
}

type PropertyDate struct {
	*BaseProperty
}

func NewPropertyDate(Name string, Date *Date) Property {
	property := &PropertyDate{
		BaseProperty: newBaseProperty(Name, "", "date", Date),
	}

	return property
}

func (property *PropertyDate) Interface() interface{} {
	return property
}

func (property *PropertyDate) Date() *Date {
	j, ok := property.JSON.GetJSON(property.Type())
	if !ok {
		return nil
	}

	date := &Date{}
	if j.Unmarshal(date) != nil {
		return nil
	}
	return date
}

type PropertyFormula struct {
	*BaseProperty
}

func NewPropertyFormula(Name string, v interface{}) Property {
	j := JSON{}
	switch v.(type) {
	case int, int64, int32, int16, int8:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 64)
		j["type"] = "number"
		j["number"] = i
	case uint, uint64, uint32, uint16, uint8:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 64)
		j["type"] = "number"
		j["number"] = i
	case bool:
		s := fmt.Sprint(v)
		b, _ := strconv.ParseBool(s)
		j["type"] = "boolean"
		j["boolean"] = b
	case *Date, Date:
		j["type"] = "date"
		j["date"] = v
	default:
		j["type"] = "string"
		j["string"] = fmt.Sprint(v)
	}

	property := &PropertyFormula{
		BaseProperty: newBaseProperty(Name, "", "formula", j),
	}

	return property
}

func (property *PropertyFormula) Interface() interface{} {
	return property
}

func (property *PropertyFormula) Formula() (Type string, v interface{}) {
	j, ok := property.JSON.GetJSON(property.Type())
	if !ok {
		return "", nil
	}

	t := j.GetString("type")
	switch t {
	case "number":
		return "number", j.GetInt("number")
	case "boolean":
		return "boolean", j.GetBool("boolean")
	case "string":
		return "string", j.GetString("string")
	case "date":
		date := &Date{}

		if jj, ok := j.GetJSON("date"); ok {
			if jj.Unmarshal(date) == nil {
				return "date", date
			}
		}
	}

	return "", nil
}

type PropertyPeople struct {
	*BaseProperty
}

func NewPropertyPeople(Name string, User ...User) Property {
	property := &PropertyPeople{
		BaseProperty: newBaseProperty(Name, "", "people", []JSON{}),
	}

	for _, u := range User {
		property.BaseProperty.JSON.Append(property.Type(), u.JSON)
	}

	return property
}

func (property *PropertyPeople) Interface() interface{} {
	return property
}

func (property *PropertyPeople) Users() []User {
	j, ok := property.JSON.GetJSONList(property.Type())
	if !ok {
		return nil
	}

	users := []User{}

	for _, jj := range j {
		u := User{ID: jj.GetString("id")}

		if jj.Unmarshal(&u.JSON) == nil {
			users = append(users, u)
		}
	}

	return users
}

type PropertyFiles struct {
	*BaseProperty
}

func NewPropertyFiles(Name string, Files ...File) Property {
	list := []JSON{}

	for _, f := range Files {
		j := JSON{}
		j.Marshal(f)
		list = append(list, j)
	}

	property := &PropertyFiles{
		BaseProperty: newBaseProperty(Name, "", "files", list),
	}

	return property
}

func (property *PropertyFiles) Interface() interface{} {
	return property
}

func (property *PropertyFiles) Files() []File {
	j, ok := property.JSON.GetJSONList(property.Type())
	if !ok {
		return nil
	}

	list := []File{}

	for _, jj := range j {
		f := File{}

		if jj.Unmarshal(&f) == nil {
			list = append(list, f)
		}
	}

	return list
}

type PropertyCheckbox struct {
	*BaseProperty
}

func NewPropertyCheckbox(Name string, checked bool) Property {
	property := &PropertyCheckbox{
		BaseProperty: newBaseProperty(Name, "", "checkbox", checked),
	}

	return property
}

func (property *PropertyCheckbox) Interface() interface{} {
	return property
}

func (property *PropertyCheckbox) Checked() bool {
	return property.JSON.GetBool(property.Type())
}

type PropertyURL struct {
	*BaseProperty
}

func NewPropertyURL(Name string, URL string) Property {
	property := &PropertyURL{
		BaseProperty: newBaseProperty(Name, "", "url", URL),
	}

	return property
}

func (property *PropertyURL) Interface() interface{} {
	return property
}

func (property *PropertyURL) URL() string {
	return property.Json().GetString(property.Type())
}

type PropertyEmail struct {
	*BaseProperty
}

func NewPropertyEmail(Name string, Email string) Property {
	property := &PropertyEmail{
		BaseProperty: newBaseProperty(Name, "", "email", Email),
	}

	return property
}

func (property *PropertyEmail) Interface() interface{} {
	return property
}

func (property *PropertyEmail) Email() string {
	return property.JSON.GetString(property.Type())
}

type PropertyPhoneNumber struct {
	*BaseProperty
}

func NewPropertyPhoneNumber(Name string, PhoneNumber string) Property {
	property := &PropertyPhoneNumber{
		BaseProperty: newBaseProperty(Name, "", "phone_number", PhoneNumber),
	}

	return property
}

func (property *PropertyPhoneNumber) Interface() interface{} {
	return property
}

func (property *PropertyPhoneNumber) PhoneNumber() string {
	return property.JSON.GetString(property.Type())
}

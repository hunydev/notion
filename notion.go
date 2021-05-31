package notion

import (
	"fmt"
	"time"
)

type Color string

const (
	ColorDefault          = Color("default")
	ColorGray             = Color("gray")
	ColorBrown            = Color("brown")
	ColorOrange           = Color("orange")
	ColorYellow           = Color("yellow")
	ColorGreen            = Color("green")
	ColorBlue             = Color("blue")
	ColorPurple           = Color("purple")
	ColorPink             = Color("pink")
	ColorRed              = Color("red")
	ColorGrayBackground   = Color("gray_background")
	ColorBrownBackground  = Color("brown_background")
	ColorOrangeBackground = Color("orange_background")
	ColorYellowBackground = Color("yellow_background")
	ColorGreenBackground  = Color("green_background")
	ColorBlueBackground   = Color("blue_background")
	ColorPurpleBackground = Color("purple_background")
	ColorPinkBackground   = Color("pink_background")
	ColorRedBackground    = Color("red_background")
)

type Object string

const (
	ObjectUser     = Object("user")
	ObjectDatabase = Object("database")
	ObjectPage     = Object("page")
	ObjectBlock    = Object("block")
)

func (object Object) String() string {
	return string(object)
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

func ParseTime(t string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", t)
}

type Notion struct {
	api API
}

func New(API API) *Notion {
	notion := &Notion{
		api: API,
	}

	return notion
}

func (notion *Notion) invalid() bool {
	if notion.api == nil {
		return true
	}

	return false
}

func (notion *Notion) APIVersion() string {
	return notion.api.Version()
}

func (notion *Notion) ListAllUsers(pagination *PaginationRequest) (*PaginationResponse, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer API Implementation")
	}

	return notion.api.ListAllUsers(pagination)
}

func (notion *Notion) RetrieveUser(UserID string) (*User, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer API Implementation")
	}

	return notion.api.RetrieveUser(UserID)
}

func (notion *Notion) RetrieveBlockChildren(BlockID string, pagination *PaginationRequest) (*PaginationResponse, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer API Implementation")
	}

	return notion.api.RetrieveBlockChildren(BlockID, pagination)
}

func (notion *Notion) AppendBlockChildren(BlockID string, Children []Block) (Block, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer API Implementation")
	}

	return notion.api.AppendBlockChildren(BlockID, Children)
}

func (notion *Notion) RetrievePage(PageID string) (*Page, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.RetrievePage(PageID)
}

func (notion *Notion) CreatePage(Parent *Parent, Properties []Property, Children ...Block) (*Page, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.CreatePage(Parent, Properties, Children...)
}

func (notion *Notion) UpdatePageProperties(PageID string, Properties ...Property) (*Page, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.UpdatePageProperties(PageID, Properties...)
}

func (notion *Notion) RetrieveDatabase(DatabaseID string) (*Database, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.RetrieveDatabase(DatabaseID)
}

func (notion *Notion) QueryDatabase(DatabaseID string, Pagination *PaginationRequest, Filter Filter, Sorts []Sort) (*PaginationResponse, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.QueryDatabase(DatabaseID, Pagination, Filter, Sorts)
}

func (notion *Notion) ListDatabases(Pagination *PaginationRequest) (*PaginationResponse, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.ListDatabases(Pagination)
}

func (notion *Notion) Search(Query string, Pagination *PaginationRequest, Filter Object, Sort *Sort) (*PaginationResponse, error) {
	if notion.invalid() {
		return nil, fmt.Errorf("Nil pointer APi Implementation")
	}

	return notion.api.Search(Query, Pagination, Filter, Sort)
}

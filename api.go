package notion

type API interface {
	ListAllUsers(pagination *PaginationRequest) (*PaginationResponse, error)
	RetrieveUser(UserID string) (*User, error)

	RetrieveBlockChildren(BlockID string, pagination *PaginationRequest) (*PaginationResponse, error)
	AppendBlockChildren(BlockID string, Children []Block) (Block, error)

	RetrievePage(PageID string) (*Page, error)
	CreatePage(Parent *Parent, Properties []Property, Children ...Block) (*Page, error)
	UpdatePageProperties(PageID string, Properties ...Property) (*Page, error)

	RetrieveDatabase(DatabaseID string) (*Database, error)
	QueryDatabase(DatabaseID string, Pagination *PaginationRequest, Filter Filter, Sorts []Sort) (*PaginationResponse, error)
	ListDatabases(Pagination *PaginationRequest) (*PaginationResponse, error)

	Search(Query string, Pagination *PaginationRequest, Filter Object, Sort *Sort) (*PaginationResponse, error)

	Version() string
}

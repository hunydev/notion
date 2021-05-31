package api

import (
	"fmt"
	"net/http"

	"github.com/hunydev/notion"
)

func (api *API) RetrieveDatabase(DatabaseID string) (*notion.Database, error) {
	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/databases/%s",
			api.baseURL(), api.contextVersion(), DatabaseID),
		nil)
	if err != nil {
		return nil, err
	}

	database := &notion.Database{}
	if err := api.doRequest(req, &database.JSON); err != nil {
		return nil, err
	}

	return database, nil
}

func (api *API) QueryDatabase(DatabaseID string, Pagination *notion.PaginationRequest, Filter notion.Filter, Sorts []notion.Sort) (*notion.PaginationResponse, error) {
	body := notion.JSON{}
	if Pagination != nil {
		for k, v := range Pagination.Json() {
			body[k] = v
		}
	}

	if Filter != nil {
		body.Set("filter", Filter.Json())
	}
	if len(Sorts) > 0 {
		list := []notion.JSON{}
		for _, sort := range Sorts {
			list = append(list, notion.JSON{
				"property":  sort.Property,
				"timestamp": sort.Timestamp,
				"direction": sort.Direction,
			})
		}
		body.Set("sorts", list)
	}

	req, err := api.prepareRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/databases/%s/query",
			api.baseURL(), api.contextVersion(), DatabaseID),
		body)
	if err != nil {
		return nil, err
	}

	p := &notion.PaginationResponse{}
	if err := api.doRequest(req, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *API) ListDatabases(Pagination *notion.PaginationRequest) (*notion.PaginationResponse, error) {
	query := ""
	if Pagination != nil {
		query = Pagination.QueryString()
	}

	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/databases?%s",
			api.baseURL(), api.contextVersion(), query),
		nil)
	if err != nil {
		return nil, err
	}

	p := &notion.PaginationResponse{}
	if err := api.doRequest(req, &p); err != nil {
		return nil, err
	}

	return p, nil
}

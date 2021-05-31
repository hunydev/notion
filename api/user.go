package api

import (
	"fmt"
	"net/http"

	"github.com/hunydev/notion"
)

func (api *API) ListAllUsers(pagination *notion.PaginationRequest) (*notion.PaginationResponse, error) {
	query := ""
	if pagination != nil {
		query = pagination.QueryString()
	}

	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/users?%s",
			api.baseURL(), api.contextVersion(), query),
		nil)
	if err != nil {
		return nil, err
	}

	p := &notion.PaginationResponse{}
	if err := api.doRequest(req, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *API) RetrieveUser(UserID string) (*notion.User, error) {
	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/users/%s",
			api.baseURL(), api.contextVersion(), UserID),
		nil)
	if err != nil {
		return nil, err
	}

	m := notion.JSON{}
	if err := api.doRequest(req, &m); err != nil {
		return nil, err
	}

	return &notion.User{
		ID:   m.GetString("id"),
		JSON: m,
	}, nil
}

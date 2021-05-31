package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hunydev/notion"
)

func (api *API) Search(Query string, Pagination *notion.PaginationRequest, Filter notion.Object, Sort *notion.Sort) (*notion.PaginationResponse, error) {
	body := notion.JSON{}

	if len(strings.TrimSpace(Query)) > 0 {
		body["query"] = Query
	}

	if Pagination != nil {
		for k, v := range Pagination.Json() {
			body[k] = v
		}
	}

	if len(strings.TrimSpace(Filter.String())) > 0 {
		body.Set("filter", notion.JSON{
			"property": "object",
			"value":    Filter.String(),
		})
	}

	if Sort != nil {
		body.Set("sort", notion.JSON{
			"direction": Sort.Direction,
			"timestamp": Sort.Timestamp,
		})
	}

	req, err := api.prepareRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/search",
			api.baseURL(), api.contextVersion()),
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

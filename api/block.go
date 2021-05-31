package api

import (
	"fmt"
	"net/http"

	"github.com/hunydev/notion"
)

func (api *API) RetrieveBlockChildren(BlockID string, pagination *notion.PaginationRequest) (*notion.PaginationResponse, error) {
	query := ""
	if pagination != nil {
		query = pagination.QueryString()
	}

	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/blocks/%s/children?%s",
			api.baseURL(), api.contextVersion(), BlockID, query),
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
func (api *API) AppendBlockChildren(BlockID string, blocks []notion.Block) (notion.Block, error) {
	children := []notion.JSON{}

	for _, block := range blocks {
		children = append(children, block.Json())
	}

	req, err := api.prepareRequest(http.MethodPatch,
		fmt.Sprintf("%s/%s/blocks/%s/children",
			api.baseURL(), api.contextVersion(), BlockID),
		notion.JSON{"children": children})
	if err != nil {
		return nil, err
	}

	j := notion.JSON{}
	if err := api.doRequest(req, &j); err != nil {
		return nil, err
	}

	return notion.AssignBlock(j)
}

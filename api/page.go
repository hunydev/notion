package api

import (
	"fmt"
	"net/http"

	"huny.dev/notion"
)

func (api *API) RetrievePage(PageID string) (*notion.Page, error) {
	req, err := api.prepareRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/pages/%s",
			api.baseURL(), api.contextVersion(), PageID),
		nil)
	if err != nil {
		return nil, err
	}

	page := &notion.Page{}
	if err := api.doRequest(req, &page.JSON); err != nil {
		return nil, err
	}

	return page, nil
}

func (api *API) CreatePage(Parent *notion.Parent, Properties []notion.Property, Children ...notion.Block) (*notion.Page, error) {
	if Parent == nil {
		return nil, fmt.Errorf("Parent is Nil pointer")
	}

	body := notion.JSON{}
	body.Set("parent", Parent.JSON)

	properties := notion.JSON{}
	for _, property := range Properties {
		properties.Set(property.Name(), property.Json())
	}
	body.Set("properties", properties)

	if len(Children) > 0 {
		blocks := []notion.JSON{}
		for _, block := range Children {
			blocks = append(blocks, block.Json())
		}
		body.Set("children", blocks)
	}

	req, err := api.prepareRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/pages", api.baseURL(), api.contextVersion()),
		body)
	if err != nil {
		return nil, err
	}

	page := &notion.Page{}
	if err := api.doRequest(req, &page.JSON); err != nil {
		return nil, err
	}

	return page, nil
}

func (api *API) UpdatePageProperties(PageID string, Properties ...notion.Property) (*notion.Page, error) {
	body := notion.JSON{}

	properties := notion.JSON{}
	for _, property := range Properties {
		properties.Set(property.Name(), property.Json())
	}
	body.Set("properties", properties)

	req, err := api.prepareRequest(http.MethodPatch,
		fmt.Sprintf("%s/%s/pages/%s", api.baseURL(), api.contextVersion(), PageID),
		body)
	if err != nil {
		return nil, err
	}

	page := &notion.Page{}
	if err := api.doRequest(req, &page.JSON); err != nil {
		return nil, err
	}

	return page, nil
}

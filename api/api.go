package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type API struct {
	token  string
	client *http.Client
}

type Option struct {
	Timeout time.Duration
}

func New(Token string, Opt *Option) *API {
	api := &API{
		token:  Token,
		client: &http.Client{},
	}

	if Opt != nil {
		api.client.Timeout = Opt.Timeout
	}

	return api
}

func (api *API) Version() string {
	return ""
}

func (api *API) baseURL() string {
	return "https://api.notion.com"
}

func (api *API) contextVersion() string {
	return "v1"
}

func (api *API) prepareRequest(method string, url string, requestBody interface{}) (*http.Request, error) {
	var r io.Reader

	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)

		if err != nil {
			return nil, err
		}

		r = bytes.NewReader(jsonBody)

		//Debug
		fmt.Printf("[Json Body] ==> %s\n", jsonBody)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))
	if len(api.Version()) > 0 {
		req.Header.Set("Notion-Version", api.Version())
	}
	if r != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (api *API) doRequest(req *http.Request, responseBody interface{}) error {
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := ReadError(resp.Body)

		return err
	}

	return api.parseResponse(resp.Body, responseBody)
}

func (api *API) parseResponse(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	if d == nil {
		return fmt.Errorf("decoder is nil")
	}

	return d.Decode(v)
}

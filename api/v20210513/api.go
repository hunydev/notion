package v20210513

import (
	"time"

	"github.com/hunydev/notion/api"
)

type API struct {
	*api.API
}

type Option struct {
	Timeout time.Duration
}

func (api *API) Version() string {
	return "2021-05-13"
}

func New(Token string, Opt *Option) *API {
	api := &API{
		API: api.New(Token, &api.Option{
			Timeout: Opt.Timeout,
		}),
	}

	return api
}

package myapi

import (
	"io"
	"net/http"
)

type MYAPI struct {
	Client  *http.Client
	BaseURL string
}

func (api *MYAPI) DoStuff() ([]byte, error) {
	resp, err := api.Client.Get(api.BaseURL + "/some/path")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// handling error and doing stuff with body that needs to be unit tested
	return body, err
}

//

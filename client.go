// Copyright 2017 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	Version = "1.0"

	// UserAgent is the header string used to identify this package.
	UserAgent = "blockchain-api-v1-client-go/" + Version + " (go; github; +https://git.io/v5dN0)"

	// APIRootTor the root address in the tor network
	APIRootTor = "https://blockchainbdgpzk.onion"

	// APIRootNet the root address in the network
	APIRootNet = "https://blockchain.info"
)

var (
	RRE = errors.New("could not read answer response")
	RSE = errors.New("incorrect response status")
	RPE = errors.New("response parsing error")
	WAE = errors.New("address is wrong")
	PAE = errors.New("no address(es) provided")
	CDE = errors.New("cannot get data on url")
	THW = errors.New("transaction hash is wrong")
	BEW = errors.New("block height is wrong")
	BHW = errors.New("block hash is wrong")
)

// Client specifies the mechanism by which individual API requests are made.
type Client struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment
	error     *Error
}

func (c *Client) userAgent() string {
	if c.UserAgent == "" {
		return UserAgent
	}

	return UserAgent + " " + c.UserAgent
}

func (c *Client) Error() *Error {
	defer func(c *Client) {
		c.error = nil
	}(c)
	return c.error
}

type Error struct {
	// ErrorMain error information from the standard package error set,
	ErrorMain error
	// ErrorExec information about the error that occurred during
	// the operation of the standard library or external packages
	ErrorExec error
	// Response http response
	Response *http.Response
	// Address wrong address
	Address *string
}

func (e Error) Error() string {
	return e.ErrorMain.Error()
}

func (c *Client) setErrorOne(errorMain error) error {
	return c.setError(errorMain, nil, nil, nil)
}

func (c *Client) setErrorTwo(errorMain error, errorExec error) error {
	return c.setError(errorMain, errorExec, nil, nil)
}

func (c *Client) setError(errorMain error, errorExec error, response *http.Response, address *string) error {
	c.error = nil

	if errorMain == nil {
		return nil
	}

	c.error = &Error{
		ErrorMain: errorMain,
		ErrorExec: errorExec,
		Response:  response,
		Address:   address,
	}

	return errorMain
}

// DoRequest to send an client request, which is then converted to the passed type.
func (c *Client) DoRequest(path string, i interface{}, params map[string]string) error {
	params["format"] = "json"
	urlValues := url.Values{}
	for k, v := range params {
		urlValues.Set(k, v)
	}

	req, e := http.NewRequest("GET", c.BasePath+path+"?"+(urlValues.Encode()), nil)
	if e != nil {
		return c.setErrorTwo(CDE, e)
	}

	req.Header.Set("User-Agent", c.userAgent())

	resp, e := c.client.Do(req)
	if e != nil {
		return c.setError(CDE, e, resp, nil)
	}

	defer resp.Body.Close()

	bytes, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return c.setError(RRE, e, resp, nil)
	}

	if resp.Status[0] != '2' {
		return c.setError(RSE, e, resp, nil)
	}

	e = json.Unmarshal(bytes, &i)
	if e != nil {
		return c.setError(RPE, e, resp, nil)
	}

	return nil
}

// New specifies the mechanism by create new client the network internet
func New() *Client {
	return &Client{client: &http.Client{}, BasePath: APIRootNet}
}

// NewTor specifies the mechanism by create new client the network internet
func NewTor() *Client {
	return &Client{client: &http.Client{}, BasePath: APIRootTor}
}

// SetHTTP client client setter
func (c *Client) SetHTTP(cli *http.Client) {
	c.client = cli
}

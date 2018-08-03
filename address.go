// Copyright 2017-2018 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

import (
	"strings"
)

// Address description of the address structure returned from the API,
// Some fields in some cases may be empty or absent.
type Address struct {
	// Exist only in the case address
	Hash160 string `json:"hash160,omitempty"`

	Address       string `json:"address"`
	NTx           uint64 `json:"n_tx"`
	TotalReceived uint64 `json:"total_received"`
	TotalSent     uint64 `json:"total_sent"`
	FinalBalance  uint64 `json:"final_balance"`
	Txs           []Tx   `json:"txs,omitempty"`

	// Exist only in the case multiaddr
	ChangeIndex  uint64 `json:"change_index,omitempty"`
	AccountIndex uint64 `json:"account_index,omitempty"`
}

// MultiAddr structure of the result when querying multiple addresses
type MultiAddr struct {
	RecommendIncludeFee bool      `json:"recommend_include_fee,omitempty"`
	SharedcoinEndpoint  string    `json:"sharedcoin_endpoint,omitempty"`
	Wallet              Wallet    `json:"wallet"`
	Addresses           []Address `json:"addresses"`
	Txs                 []Tx      `json:"txs"`
	Info                Info      `json:"info"`
}

// Wallet summary data about the requested addresses
type Wallet struct {
	NTx           uint64 `json:"n_tx"`
	NTxFiltered   uint64 `json:"n_tx_filtered"`
	TotalReceived uint64 `json:"total_received"`
	TotalSent     uint64 `json:"total_sent"`
	FinalBalance  uint64 `json:"final_balance"`
}

// SymbolLocal ...
type SymbolLocal struct {
	Code               string  `json:"code"`
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Conversion         float64 `json:"conversion"`
	SymbolAppearsAfter bool    `json:"symbolAppearsAfter"`
	Local              bool    `json:"local"`
}

// SymbolBtc ...
type SymbolBtc struct {
	Code               string  `json:"code"`
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Conversion         float64 `json:"conversion"`
	SymbolAppearsAfter bool    `json:"symbolAppearsAfter"`
	Local              bool    `json:"local"`
}

// Info ...
type Info struct {
	NConnected  uint64      `json:"nconnected"`
	Conversion  float64     `json:"conversion"`
	SymbolLocal SymbolLocal `json:"symbol_local"`
	SymbolBtc   SymbolBtc   `json:"symbol_btc"`
	LatestBlock LatestBlock `json:"latest_block"`
}

func (c *Client) CheckAddress(address string) error {
	if res := ValidateBitcoinAddress(address); res == -1 {
		return c.setError(AIW, nil, nil, nil)
	}

	return nil
}

func (c *Client) CheckAddresses(addresses []string) (e error) {
	if len(addresses) == 0 {
		return c.setErrorOne(ANP)
	}

	for _, address := range addresses {
		if e = c.CheckAddress(address); e != nil {
			return e
		}
	}

	return nil
}

// GetAddress alias GetAddressAdv without additional parameters
func (c *Client) GetAddress(address string) (*Address, error) {
	return c.GetAddressAdv(address, nil)
}

// GetAddressAdv is a mechanism which is used to obtain information about the address
func (c *Client) GetAddressAdv(address string, options map[string]string) (resp *Address, e error) {
	if e = c.CheckAddress(address); e != nil {
		return
	}
	resp = &Address{}
	return resp, c.DoRequest("/address/"+address, resp, options)
}

// GetAddresses alias GetAddressesAdv without additional parameters
func (c *Client) GetAddresses(addresses []string) (*MultiAddr, error) {
	return c.GetAddressesAdv(addresses, nil)
}

// GetAddressesAdv is a mechanism which is used to obtain information about the addresses
func (c *Client) GetAddressesAdv(addresses []string, options map[string]string) (resp *MultiAddr, e error) {
	if e = c.CheckAddresses(addresses); e != nil {
		return
	}

	options = ApproveParams(options)
	options["active"] = strings.Join(addresses, "|")
	resp = &MultiAddr{}
	return resp, c.DoRequest("/multiaddr", resp, options)
}

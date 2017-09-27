// Copyright 2017 Vasilyuk Vasiliy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

import (
	"errors"
	"strings"
)

// Description of the address structure returned from the API,
// Some fields in some cases may be empty or absent.
type Address struct {
	// Exist only in the case address
	Hash160 string `json:"hash160,omitempty"`

	Address       string `json:"address"`
	NTx           uint64 `json:"n_tx"`
	TotalReceived uint64 `json:"total_received"`
	TotalSent     uint64 `json:"total_sent"`
	FinalBalance  uint64 `json:"final_balance"`
	Txs           []*Tx  `json:"txs,omitempty"`

	// Exist only in the case multiaddr
	ChangeIndex  uint64 `json:"change_index,omitempty"`
	AccountIndex uint64 `json:"account_index,omitempty"`
}

// The structure of the result when querying multiple addresses
type MultiAddr struct {
	RecommendIncludeFee bool       `json:"recommend_include_fee,omitempty"`
	SharedcoinEndpoint  string     `json:"sharedcoin_endpoint,omitempty"`
	Wallet              *Wallet    `json:"wallet"`
	Addresses           []*Address `json:"addresses"`
	Txs                 []*Tx      `json:"txs"`
	Info                *Info      `json:"info"`
}

// Summary data about the requested addresses
type Wallet struct {
	NTx           uint64 `json:"n_tx"`
	NTxFiltered   uint64 `json:"n_tx_filtered"`
	TotalReceived uint64 `json:"total_received"`
	TotalSent     uint64 `json:"total_sent"`
	FinalBalance  uint64 `json:"final_balance"`
}

type SymbolLocal struct {
	Code               string  `json:"code"`
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Conversion         float64 `json:"conversion"`
	SymbolAppearsAfter bool    `json:"symbolAppearsAfter"`
	Local              bool    `json:"local"`
}

type SymbolBtc struct {
	Code               string  `json:"code"`
	Symbol             string  `json:"symbol"`
	Name               string  `json:"name"`
	Conversion         float64 `json:"conversion"`
	SymbolAppearsAfter bool    `json:"symbolAppearsAfter"`
	Local              bool    `json:"local"`
}

type Info struct {
	NConnected  uint64       `json:"nconnected"`
	Conversion  float64      `json:"conversion"`
	SymbolLocal *SymbolLocal `json:"symbol_local"`
	SymbolBtc   *SymbolBtc   `json:"symbol_btc"`
	LatestBlock *LatestBlock `json:"latest_block"`
}

// Receiving data about one particular address
func (c *Client) GetAddress(address string, params ...map[string]string) (response *Address, e error) {
	if address == "" {
		return nil, errors.New("No Address Provided")
	}

	options := map[string]string{"format": "json"}
	if len(params) > 0 {
		for k, v := range params[0] {
			options[k] = v
		}
	}
	response = &Address{}
	e = c.DoRequest("/address/"+address, response, options)

	return
}

// Method for obtaining data about the set of addresses. No more than 80 addresses a time.
func (c *Client) GetAddresses(addresses []string, params ...map[string]string) (response *MultiAddr, e error) {
	if len(addresses) < 2 {
		return nil, errors.New("Invalid argument, you must pass an array with two or more addresses!")
	}

	options := map[string]string{"active": strings.Join(addresses, "|")}
	if len(params) > 0 {
		for k, v := range params[0] {
			options[k] = v
		}
	}

	response = &MultiAddr{}
	e = c.DoRequest("/multiaddr", response, options)

	return
}

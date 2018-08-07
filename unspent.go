// Copyright 2017-2018 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

import "strings"

// UnspentOutputs the set of unspent outputs
type UnspentOutputs struct {
	Notice         string          `json:"notice,omitempty"`
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

// UnspentOutput the basic structure unspent outputs
type UnspentOutput struct {
	TxAge     string `json:"tx_age"`
	TxHash    string `json:"tx_hash"`
	TxIndex   uint64 `json:"tx_index"`
	TxOutputN uint64 `json:"tx_output_n"`
	Script    string `json:"script"`
	Value     int64  `json:"value"`
}

// GetUnspent alias GetUnspentAdv without additional parameters
func (c *Client) GetUnspent(addresses []string) (*UnspentOutputs, error) {
	return c.GetUnspentAdv(addresses, nil)
}

// GetUnspentAdv specifies the mechanism by getting unspent outputs multiple addresses
func (c *Client) GetUnspentAdv(addresses []string, options map[string]string) (resp *UnspentOutputs, e error) {
	if e = c.checkAddresses(addresses); e != nil {
		return
	}

	options = ApproveOptions(options)
	options["active"] = strings.Join(addresses, "|")
	resp = &UnspentOutputs{}
	return resp, c.Do("/unspent", resp, options)
}

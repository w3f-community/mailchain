// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mailchain/mailchain/cmd/internal/http/params"
	"github.com/mailchain/mailchain/errs"
	"github.com/mailchain/mailchain/internal/address"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/pkg/errors"
)

// GetAddresses returns a handler get spec.
func GetAddresses(ks keystore.Store) func(w http.ResponseWriter, r *http.Request) {
	// Get swagger:route GET /addresses Addresses GetAddresses
	//
	// Get addresses.
	//
	// Get all address that this user has access to. The addresses can be used to send or receive messages.
	// Responses:
	//   200: GetAddressesResponse
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		protocol, err := params.QueryRequireProtocol(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		network, err := params.QueryRequireNetwork(r)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusUnprocessableEntity, errors.WithStack(err))
			return
		}

		rawAddresses, err := ks.GetAddresses(protocol, network)
		if err != nil {
			errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err))
			return
		}

		addresses := []GetAddressesItem{}
		for _, x := range rawAddresses {
			value, encoding, err := address.EncodeByProtocol(x, protocol)
			if err != nil {
				errs.JSONWriter(w, r, http.StatusInternalServerError, errors.WithStack(err))
				return
			}

			addresses = append(addresses, GetAddressesItem{
				Value:    value,
				Encoding: encoding,
			})
		}

		_ = json.NewEncoder(w).Encode(GetAddressesResponse{Addresses: addresses})
		w.Header().Set("Content-Type", "application/json")
	}
}

// GetAddressesResponse Holds the response messages
//
// swagger:response GetAddressesResponse
type GetAddressesResponse struct {
	// in: body
	Addresses []GetAddressesItem `json:"addresses"`
}

// swagger:response // in: body
type GetAddressesItem struct {
	// Address value
	//
	// Required: true
	// example: 0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761
	Value string `json:"value"`
	// Encoding method used for encoding the `address`
	//
	// Required: true
	// example: hex/0x-prefix
	Encoding string `json:"encoding"`
}

// GetAddressesRequest body
// swagger:parameters GetAddresses
type GetAddressesRequest struct {
	// Network to use when finding addresses.
	//
	// enum: mainnet,goerli,ropsten,rinkeby,local
	// in: query
	// required: true
	// example: goerli
	Network string `json:"network"`

	// Protocol to use when finding addresses.
	//
	// enum: ethereum, substrate
	// in: query
	// required: true
	// example: ethereum
	Protocol string `json:"protocol"`
}

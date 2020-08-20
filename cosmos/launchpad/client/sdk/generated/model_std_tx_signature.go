/*
 * Gaia-Lite for Cosmos
 *
 * A REST interface for state queries, transaction generation and broadcasting.
 *
 * API version: 3.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// StdTxSignature struct for StdTxSignature
type StdTxSignature struct {
	Signature     string               `json:"signature,omitempty"`
	PubKey        StdTxSignaturePubKey `json:"pub_key,omitempty"`
	AccountNumber string               `json:"account_number,omitempty"`
	Sequence      string               `json:"sequence,omitempty"`
}

/*
 * Gaia-Lite for Cosmos
 *
 * A REST interface for state queries, transaction generation and broadcasting.
 *
 * API version: 3.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// InlineResponse2005 struct for InlineResponse2005
type InlineResponse2005 struct {
	Height string `json:"height,omitempty"`
	Result []Coin `json:"result,omitempty"`
}

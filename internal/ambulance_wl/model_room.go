/*
 * Waiting List Api
 *
 * Ambulance Waiting List management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: test@test.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package ambulance_wl

// Room - Describes dimensions and equipment of ambulance rooms
type Room struct {

	Id string `json:"id,omitempty"`

	Width string `json:"width,omitempty"`

	Height string `json:"height,omitempty"`

	// Link to something
	Reference string `json:"reference,omitempty"`

	TipicalCostToOperate int32 `json:"tipicalCostToOperate,omitempty"`

	Equipment string `json:"equipment,omitempty"`

	Name string `json:"name,omitempty"`
}

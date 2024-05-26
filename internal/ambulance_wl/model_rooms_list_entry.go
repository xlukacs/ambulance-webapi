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

type RoomsListEntry struct {

	// Unique id of the entry in this waiting list
	Id string `json:"id"`

	// Width of room in waiting list
	Width string `json:"width"`

	// Height of the room known to Web-In-Cloud system
	Height string `json:"height"`

	// Timestamp since when the room entered the rooms list
	TipicalCostToOperate int32 `json:"tipicalCostToOperate"`

	Room Room `json:"Room,omitempty"`
}

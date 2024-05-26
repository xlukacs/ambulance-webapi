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

import (
   "net/http"

   "github.com/gin-gonic/gin"
)

type AmbulanceRoomsAPI interface {

   // internal registration of api routes
   addRoutes(routerGroup *gin.RouterGroup)

    // CreateRoom - Saves new entry into rooms list
   CreateRoom(ctx *gin.Context)

    // DeleteRoom - Deletes specific room
   DeleteRoom(ctx *gin.Context)

    // GetRooms - Provides the list of rooms associated with ambulance
   GetRooms(ctx *gin.Context)

}

// partial implementation of AmbulanceRoomsAPI - all functions must be implemented in add on files
type implAmbulanceRoomsAPI struct {

}

func newAmbulanceRoomsAPI() AmbulanceRoomsAPI {
  return &implAmbulanceRoomsAPI{}
}

func (this *implAmbulanceRoomsAPI) addRoutes(routerGroup *gin.RouterGroup) {
  routerGroup.Handle( http.MethodPost, "/rooms/:ambulanceId/entries", this.CreateRoom)
  routerGroup.Handle( http.MethodDelete, "/rooms/:ambulanceId/entries", this.DeleteRoom)
  routerGroup.Handle( http.MethodGet, "/rooms/:ambulanceId/entries", this.GetRooms)
}


// Copy following section to separate file, uncomment, and implement accordingly
// // CreateRoom - Saves new entry into rooms list
// func (this *implAmbulanceRoomsAPI) CreateRoom(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // DeleteRoom - Deletes specific room
// func (this *implAmbulanceRoomsAPI) DeleteRoom(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//
// // GetRooms - Provides the list of rooms associated with ambulance
// func (this *implAmbulanceRoomsAPI) GetRooms(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//


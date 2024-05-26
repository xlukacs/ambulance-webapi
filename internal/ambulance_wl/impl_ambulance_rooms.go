package ambulance_wl

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (this *implAmbulanceRoomsAPI) GetRooms(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.Rooms
		if result == nil {
			result = []Room{}
		}
		// return nil ambulance - no need to update it in db
		return nil, result, http.StatusOK
	})
}

func (this *implAmbulanceRoomsAPI) DeleteRoom(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		roomId := ctx.Param("roomId")

		if roomId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room ID is required",
			}, http.StatusBadRequest
		}

		roomIndx := slices.IndexFunc(ambulance.Rooms, func(current Room) bool {
			return roomId == current.Id
		})

		if roomIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Room not found",
			}, http.StatusNotFound
		}

		ambulance.Rooms = append(ambulance.Rooms[:roomIndx], ambulance.Rooms[roomIndx+1:]...)
		//ambulance.reconcileWaitingList()
		return ambulance, nil, http.StatusNoContent
	})
}

func (this *implAmbulanceRoomsAPI) CreateRoom(ctx *gin.Context) {
	// ctx.AbortWithStatus(http.StatusNotImplemented)
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry Room

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.Id == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room ID is required",
			}, http.StatusBadRequest
		}

		if entry.Width == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room width is required",
			}, http.StatusBadRequest
		}

		if entry.Height == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room height is required",
			}, http.StatusBadRequest
		}

		if entry.Equipment == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room equipment is required",
			}, http.StatusBadRequest
		}

		//if entry.TipicalCostToOperate == -1 {
		//	return nil, gin.H{
		//		"status":  http.StatusBadRequest,
		//		"message": "Room TipicalCostToOperate is required",
		//	}, http.StatusBadRequest
		//}

		if entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ambulance.Rooms, func(room_entry Room) bool {
			return entry.Id == room_entry.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.Rooms = append(ambulance.Rooms, entry)
		//ambulance.reconcileWaitingList() TODO: this is not needed here, since we dont need to update other room data
		//entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.Rooms, func(room_entry Room) bool {
			return entry.Id == room_entry.Id
		})
		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}

		return ambulance, ambulance.Rooms[entryIndx], http.StatusOK
	})
}

func (this *implAmbulanceRoomsAPI) UpdateRoom(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var room Room

		if err := c.ShouldBindJSON(&room); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		roomId := ctx.Param("roomId")

		if roomId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room ID is required",
			}, http.StatusBadRequest
		}

		roomIndx := slices.IndexFunc(ambulance.Rooms, func(current Room) bool {
			return roomId == current.Id
		})

		if roomIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Room not found",
			}, http.StatusNotFound
		}

		if room.Id != "" {
			ambulance.Rooms[roomIndx].Id = room.Id
		}

		if room.Width != "" {
			ambulance.Rooms[roomIndx].Width = room.Width
		}

		if room.Height != "" {
			ambulance.Rooms[roomIndx].Height = room.Height
		}

		if room.Equipment != "" {
			ambulance.Rooms[roomIndx].Equipment = room.Equipment
		}

		//ambulance.reconcileWaitingList()
		return ambulance, ambulance.Rooms[roomIndx], http.StatusOK
	})
}

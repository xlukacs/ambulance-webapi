package ambulance_wl

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (this *implSchedulesAPI) CreateSchedule(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry Schedule

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
				"message": "Schedule ID is required",
			}, http.StatusBadRequest
		}

		if entry.PatientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patiend ID is required",
			}, http.StatusBadRequest
		}

		if entry.RoomId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room ID is required",
			}, http.StatusBadRequest
		}

		if entry.Start.IsZero() {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Schedule start date is required",
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

		conflictIndx := slices.IndexFunc(ambulance.Schedules, func(schedule_entry Schedule) bool {
			return entry.Id == schedule_entry.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.Schedules = append(ambulance.Schedules, entry)
		// //ambulance.reconcileWaitingList() TODO: this is not needed here, since we dont need to update other room data
		// //entry was copied by value return reconciled value from the list
		entryIndx := slices.IndexFunc(ambulance.Schedules, func(schedule_entry Schedule) bool {
			return entry.Id == schedule_entry.Id
		})
		// if entryIndx < 0 {
		// 	return nil, gin.H{
		// 		"status":  http.StatusInternalServerError,
		// 		"message": "Failed to save entry",
		// 	}, http.StatusInternalServerError
		// }

		return ambulance, ambulance.Schedules[entryIndx], http.StatusOK

		// return ambulance, entry, http.StatusOK
	})
}

func (this *implSchedulesAPI) DeleteSchedule(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		scheduleId := ctx.Param("scheduleId")

		if scheduleId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Schedule iD is required",
			}, http.StatusBadRequest
		}

		scheduleIndx := slices.IndexFunc(ambulance.Schedules, func(current Schedule) bool {
			return scheduleId == current.Id
		})

		if scheduleIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Room not found",
			}, http.StatusNotFound
		}

		ambulance.Schedules = append(ambulance.Schedules[:scheduleIndx], ambulance.Schedules[scheduleIndx+1:]...)
		//ambulance.reconcileWaitingList()
		return ambulance, nil, http.StatusNoContent
	})
}

func (this *implSchedulesAPI) GetSchedule(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		result := ambulance.Schedules
		if result == nil {
			result = []Schedule{}
		}
		// return nil ambulance - no need to update it in db
		return ambulance, result, http.StatusOK
	})
}

func (this *implSchedulesAPI) UpdateSchedule(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var schedule Schedule

		if err := c.ShouldBindJSON(&schedule); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		scheduleId := ctx.Param("scheduleId")

		if scheduleId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Room ID is required",
			}, http.StatusBadRequest
		}

		scheduleIdx := slices.IndexFunc(ambulance.Schedules, func(current Schedule) bool {
			return scheduleId == current.Id
		})

		if scheduleIdx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Room not found",
			}, http.StatusNotFound
		}

		if schedule.Id != "" {
			ambulance.Schedules[scheduleIdx].Id = schedule.Id
		}

		if schedule.RoomId != "" {
			ambulance.Schedules[scheduleIdx].RoomId = schedule.RoomId
		}

		if schedule.PatientId != "" {
			ambulance.Schedules[scheduleIdx].PatientId = schedule.PatientId
		}

		if schedule.Note != "" {
			ambulance.Schedules[scheduleIdx].Note = schedule.Note
		}

		if !schedule.Start.IsZero() {
			ambulance.Schedules[scheduleIdx].Start = schedule.Start
		}

		if !schedule.End.IsZero() {
			ambulance.Schedules[scheduleIdx].End = schedule.End
		}

		//ambulance.reconcileWaitingList()
		return ambulance, ambulance.Schedules[scheduleIdx], http.StatusOK
	})
}

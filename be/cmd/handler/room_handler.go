package handler

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/types"
	"homestay-be/core/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoomHandler struct {
	roomLogic *logic.RoomLogic
}

func NewRoomHandler(roomLogic *logic.RoomLogic) *RoomHandler {
	return &RoomHandler{
		roomLogic: roomLogic,
	}
}

// CreateRoom - Create a new room
// @Summary Create a new room
// @Description Create a new room for a homestay
// @Tags Room
// @Accept json
// @Produce json
// @Param room body types.CreateRoomRequest true "Room information"
// @Success 200 {object} response.Response{data=types.Room}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms [post]
// @Security BearerAuth
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req types.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	room, err := h.roomLogic.CreateRoom(&req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, room)
}

// GetRoomByID - Get room by ID
// @Summary Get room by ID
// @Description Get detailed information of a room by ID
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} response.Response{data=types.RoomDetailResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms/{id} [get]
// @Security BearerAuth
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	roomDetail, err := h.roomLogic.GetRoomByID(roomID, hostID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgRoomNotFound)
		return
	}

	response.ResponseSuccess(c, roomDetail)
}

// GetRoomList - Get list of rooms
// @Summary Get list of rooms
// @Description Get paginated list of rooms for a homestay
// @Tags Room
// @Accept json
// @Produce json
// @Param homestay_id query int true "Homestay ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10, max: 100)"
// @Param search query string false "Search by room name"
// @Param status query string false "Filter by status (active, inactive)"
// @Param type query string false "Filter by room type"
// @Param sort_by query string false "Sort by field (name, price, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} response.Response{data=types.RoomListResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms [get]
// @Security BearerAuth
func (h *RoomHandler) GetRoomList(c *gin.Context) {
	var req types.RoomListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	roomList, err := h.roomLogic.GetRoomList(&req, hostID)
	if err != nil {
		logx.Error(err)
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, roomList)
}

// UpdateRoom - Update room
// @Summary Update room
// @Description Update room information
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body types.UpdateRoomRequest true "Updated room information"
// @Success 200 {object} response.Response{data=types.Room}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms/{id} [put]
// @Security BearerAuth
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	var req types.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	room, err := h.roomLogic.UpdateRoom(roomID, &req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, room)
}

// DeleteRoom - Delete room
// @Summary Delete room
// @Description Delete a room (only if no active bookings)
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms/{id} [delete]
// @Security BearerAuth
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	err = h.roomLogic.DeleteRoom(roomID, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, nil)
}

// CreateAvailability - Create room availability
// @Summary Create room availability
// @Description Create availability for a specific room and date
// @Tags Room Availability
// @Accept json
// @Produce json
// @Param availability body types.CreateAvailabilityRequest true "Availability information"
// @Success 200 {object} response.Response{data=types.RoomAvailability}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms/availability [post]
// @Security BearerAuth
func (h *RoomHandler) CreateAvailability(c *gin.Context) {
	var req types.CreateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	availability, err := h.roomLogic.CreateAvailability(&req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, availability)
}

// UpdateAvailability - Update room availability
// @Summary Update room availability
// @Description Update availability for a specific room and date
// @Tags Room Availability
// @Accept json
// @Produce json
// @Param id path int true "Availability ID"
// @Param availability body types.UpdateAvailabilityRequest true "Updated availability information"
// @Success 200 {object} response.Response{data=types.RoomAvailability}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/rooms/availability/{id} [put]
// @Security BearerAuth
func (h *RoomHandler) UpdateAvailability(c *gin.Context) {
	availabilityIDStr := c.Param("id")
	availabilityID, err := strconv.Atoi(availabilityIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	var req types.UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	availability, err := h.roomLogic.UpdateAvailability(availabilityID, &req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, availability)
}

// GetRoomStats - Get room statistics
// @Summary Get room statistics
// @Description Get statistics for rooms in a homestay
// @Tags Room
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Success 200 {object} response.Response{data=types.RoomStatsResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id}/rooms/stats [get]
// @Security BearerAuth
func (h *RoomHandler) GetRoomStats(c *gin.Context) {
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}
	hostID := userID.(int)

	stats, err := h.roomLogic.GetRoomStats(homestayID, hostID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgHomestayNotFound)
		return
	}

	response.ResponseSuccess(c, stats)
}

func (h *RoomHandler) GetPublicRoomList(c *gin.Context) {
	var req types.RoomListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	roomList, err := h.roomLogic.GetPublicRoomList(&req)
	if err != nil {
		logx.Error(err)
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, roomList)
}

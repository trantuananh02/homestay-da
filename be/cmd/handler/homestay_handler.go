package handler

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/types"
	"homestay-be/core/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomestayHandler struct {
	homestayLogic *logic.HomestayLogic
}

func NewHomestayHandler(homestayLogic *logic.HomestayLogic) *HomestayHandler {
	return &HomestayHandler{
		homestayLogic: homestayLogic,
	}
}

// CreateHomestay - Create a new homestay
// @Summary Create a new homestay
// @Description Create a new homestay for the authenticated host
// @Tags Homestay
// @Accept json
// @Produce json
// @Param homestay body types.CreateHomestayRequest true "Homestay information"
// @Success 200 {object} response.Response{data=types.Homestay}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays [post]
// @Security BearerAuth
func (h *HomestayHandler) CreateHomestay(c *gin.Context) {
	var req types.CreateHomestayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Create homestay
	homestay, err := h.homestayLogic.CreateHomestay(&req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, homestay)
}

// GetHomestayByID - Get homestay by ID
// @Summary Get homestay by ID
// @Description Get detailed information of a homestay by ID
// @Tags Homestay
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Success 200 {object} response.Response{data=types.HomestayDetailResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id} [get]
// @Security BearerAuth
func (h *HomestayHandler) GetHomestayByID(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Get homestay
	homestayDetail, err := h.homestayLogic.GetHomestayByID(homestayID, hostID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgHomestayNotFound)
		return
	}

	response.ResponseSuccess(c, homestayDetail)
}

// GetHomestayList - Get list of homestays
// @Summary Get list of homestays
// @Description Get paginated list of homestays for the authenticated host
// @Tags Homestay
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param status query string false "Filter by status (active, inactive)"
// @Param city query string false "Filter by city"
// @Param district query string false "Filter by district"
// @Success 200 {object} response.Response{data=types.HomestayListResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays [get]
// @Security BearerAuth
func (h *HomestayHandler) GetHomestayList(c *gin.Context) {
	var req types.HomestayListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Get homestay list
	homestayList, err := h.homestayLogic.GetHomestayList(&req, hostID)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, response.MsgDatabaseError)
		return
	}

	response.ResponseSuccess(c, homestayList)
}

// UpdateHomestay - Update homestay
// @Summary Update homestay
// @Description Update homestay information
// @Tags Homestay
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Param homestay body types.UpdateHomestayRequest true "Updated homestay information"
// @Success 200 {object} response.Response{data=types.Homestay}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id} [put]
// @Security BearerAuth
func (h *HomestayHandler) UpdateHomestay(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	var req types.UpdateHomestayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Update homestay
	homestay, err := h.homestayLogic.UpdateHomestay(homestayID, &req, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, homestay)
}

// DeleteHomestay - Delete homestay
// @Summary Delete homestay
// @Description Delete a homestay (only if no active bookings)
// @Tags Homestay
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id} [delete]
// @Security BearerAuth
func (h *HomestayHandler) DeleteHomestay(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Delete homestay
	err = h.homestayLogic.DeleteHomestay(homestayID, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, nil)
}

// GetHomestayStats - Get homestay statistics
// @Summary Get homestay statistics
// @Description Get statistics for all homestays of the authenticated host
// @Tags Homestay
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=types.HomestayStatsResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/stats [get]
// @Security BearerAuth
func (h *HomestayHandler) GetHomestayStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Get statistics
	stats, err := h.homestayLogic.GetHomestayStats(hostID)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, response.MsgDatabaseError)
		return
	}

	response.ResponseSuccess(c, stats)
}

// GetHomestayStatsByID - Get homestay statistics by ID
// @Summary Get homestay statistics by ID
// @Description Get statistics for a specific homestay
// @Tags Homestay
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Success 200 {object} response.Response{data=types.HomestayStatsResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id}/stats [get]
// @Security BearerAuth
func (h *HomestayHandler) GetHomestayStatsByID(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Get statistics
	stats, err := h.homestayLogic.GetHomestayStatsByID(homestayID, hostID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgHomestayNotFound)
		return
	}

	response.ResponseSuccess(c, stats)
}

// ToggleHomestayStatus - Toggle homestay status (active/inactive)
// @Summary Toggle homestay status
// @Description Toggle homestay status between active and inactive
// @Tags Homestay
// @Accept json
// @Produce json
// @Param id path int true "Homestay ID"
// @Success 200 {object} response.Response{data=types.Homestay}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/host/homestays/{id}/toggle-status [put]
// @Security BearerAuth
func (h *HomestayHandler) ToggleHomestayStatus(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		response.ResponseError(c, response.Unauthorized, response.MsgUnauthorized)
		return
	}

	hostID := userID.(int)

	// Toggle homestay status
	homestay, err := h.homestayLogic.ToggleHomestayStatus(homestayID, hostID)
	if err != nil {
		response.ResponseError(c, response.BadRequest, err.Error())
		return
	}

	response.ResponseSuccess(c, homestay)
}

// logic guest
func (h *HomestayHandler) GetPublicHomestayList(c *gin.Context) {
	var req types.HomestayListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidData+": "+err.Error())
		return
	}

	// Get homestay list
	homestayList, err := h.homestayLogic.GetPublicHomestayList(&req)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, response.MsgDatabaseError)
		return
	}

	response.ResponseSuccess(c, homestayList)
}

func (h *HomestayHandler) GetPublicHomestayByID(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get homestay
	homestayDetail, err := h.homestayLogic.GetPublicHomestayByID(homestayID)
	if err != nil {
		response.ResponseError(c, response.NotFound, response.MsgHomestayNotFound)
		return
	}

	response.ResponseSuccess(c, homestayDetail)
}

func (h *HomestayHandler) GetTopHomestays(c *gin.Context) {
	// Get top homestays
	topHomestays, err := h.homestayLogic.GetTopHomestays(8)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, response.MsgDatabaseError)
		return
	}

	response.ResponseSuccess(c, topHomestays)
}

// GetHomestayReviews - Get reviews for a homestay
// @Summary Get reviews for a homestay
func (h *HomestayHandler) GetHomestayReviews(c *gin.Context) {
	// Get homestay ID from URL
	homestayIDStr := c.Param("id")
	homestayID, err := strconv.Atoi(homestayIDStr)
	if err != nil {
		response.ResponseError(c, response.BadRequest, response.MsgInvalidID)
		return
	}

	// Get reviews
	reviews, err := h.homestayLogic.GetHomestayReviews(homestayID)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, response.MsgDatabaseError)
		return
	}

	response.ResponseSuccess(c, reviews)
}
package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type BroadcastHandler struct {
	service *service.BroadcastService
}

func NewBroadcastHandler(broadcastService *service.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{service: broadcastService}
}

func (h *BroadcastHandler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *BroadcastHandler) UpdateConfig(c *gin.Context) {
	var request service.BroadcastConfig
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	config, err := h.service.UpdateConfig(c.Request.Context(), request)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *BroadcastHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), false, 100)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

func (h *BroadcastHandler) Create(c *gin.Context) {
	var request service.ActivityBroadcastInput
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	adminID := int64(0)
	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		adminID = subject.UserID
	}
	item, err := h.service.Create(c.Request.Context(), request, adminID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *BroadcastHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid broadcast id")
		return
	}
	var request service.ActivityBroadcastInput
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	item, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *BroadcastHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid broadcast id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *BroadcastHandler) Cleanup(c *gin.Context) {
	var request struct {
		RetentionDays int `json:"retention_days"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	deleted, err := h.service.CleanupExpiredSystem(c.Request.Context(), request.RetentionDays)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": deleted})
}

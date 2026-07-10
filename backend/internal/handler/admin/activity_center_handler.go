package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityCenterHandler struct {
	activityCenterService *service.ActivityCenterService
}

func NewActivityCenterHandler(activityCenterService *service.ActivityCenterService) *ActivityCenterHandler {
	return &ActivityCenterHandler{activityCenterService: activityCenterService}
}

func (h *ActivityCenterHandler) GetConfig(c *gin.Context) {
	cfg, err := h.activityCenterService.GetConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, cfg)
}

func (h *ActivityCenterHandler) UpdateConfig(c *gin.Context) {
	var req service.ActivityCenterConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	cfg, err := h.activityCenterService.UpdateConfig(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, cfg)
}

func (h *ActivityCenterHandler) List(c *gin.Context) {
	items, err := h.activityCenterService.ListAdmin(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

func (h *ActivityCenterHandler) Create(c *gin.Context) {
	var req service.ActivityCenterItemInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	adminID := int64(0)
	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		adminID = subject.UserID
	}
	item, err := h.activityCenterService.Create(c.Request.Context(), req, adminID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *ActivityCenterHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid activity id")
		return
	}
	var req service.ActivityCenterItemInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	item, err := h.activityCenterService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *ActivityCenterHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid activity id")
		return
	}
	if err := h.activityCenterService.Delete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

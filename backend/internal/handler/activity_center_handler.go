package handler

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityCenterHandler struct {
	activityCenterService *service.ActivityCenterService
}

func NewActivityCenterHandler(activityCenterService *service.ActivityCenterService) *ActivityCenterHandler {
	return &ActivityCenterHandler{activityCenterService: activityCenterService}
}

func (h *ActivityCenterHandler) List(c *gin.Context) {
	result, err := h.activityCenterService.ListVisible(c.Request.Context(), time.Now())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

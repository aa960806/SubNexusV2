package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type BroadcastHandler struct {
	service *service.BroadcastService
}

func NewBroadcastHandler(broadcastService *service.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{service: broadcastService}
}

func (h *BroadcastHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.service.List(c.Request.Context(), true, limit)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

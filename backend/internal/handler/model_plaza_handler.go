package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *SettingHandler) GetModelPlaza(c *gin.Context) {
	if !h.settingService.IsModelPlazaEnabled(c.Request.Context()) {
		response.NotFound(c, "Model plaza is not available")
		return
	}
	config, err := h.settingService.GetModelPlazaConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

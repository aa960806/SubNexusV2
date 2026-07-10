package admin

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type modelPlazaConfigPayload struct {
	Enabled bool                     `json:"enabled"`
	Config  service.ModelPlazaConfig `json:"config"`
}

func (h *SettingHandler) GetModelPlazaConfig(c *gin.Context) {
	config, err := h.settingService.GetModelPlazaConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, modelPlazaConfigPayload{
		Enabled: h.settingService.IsModelPlazaEnabled(c.Request.Context()),
		Config:  *config,
	})
}

func (h *SettingHandler) UpdateModelPlazaConfig(c *gin.Context) {
	var request modelPlazaConfigPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.settingService.SaveModelPlazaConfig(c.Request.Context(), &request.Config); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if err := h.settingService.SetModelPlazaEnabled(c.Request.Context(), request.Enabled); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	updated, err := h.settingService.GetModelPlazaConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, modelPlazaConfigPayload{
		Enabled: h.settingService.IsModelPlazaEnabled(c.Request.Context()),
		Config:  *updated,
	})
}

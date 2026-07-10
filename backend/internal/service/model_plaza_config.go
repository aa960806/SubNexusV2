package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"golang.org/x/sync/singleflight"
)

// ModelPlazaPrice holds the four display prices for a single model.
// All fields are free-form strings entered by the admin (e.g. "3.00"); the
// backend never parses or computes with them — this feature is display-only.
type ModelPlazaPrice struct {
	Input      string `json:"input"`       // 输入价格 / per 1M tokens
	Output     string `json:"output"`      // 输出价格 / per 1M tokens
	CacheRead  string `json:"cache_read"`  // 缓存读取 / per 1M tokens
	CacheWrite string `json:"cache_write"` // 缓存创建 / per 1M tokens
}

// ModelPlazaModel is one row in a plaza group: a model name + its prices.
type ModelPlazaModel struct {
	Name  string          `json:"name"`
	Price ModelPlazaPrice `json:"price"`
	Note  string          `json:"note,omitempty"`
}

// ModelPlazaGroup is a display group (decoupled from real project groups).
type ModelPlazaGroup struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Currency    string            `json:"currency,omitempty"` // group-level default currency symbol, e.g. "$" / "¥"
	Models      []ModelPlazaModel `json:"models"`
}

// ModelPlazaCategory is one of the fixed top-level model types. Its display
// label lives in the frontend i18n (keyed by Key); the backend only stores the
// key and the admin-configured groups under it.
type ModelPlazaCategory struct {
	Key    string            `json:"key"`
	Groups []ModelPlazaGroup `json:"groups"`
}

// Fixed top-level category keys (order is the display order).
const (
	ModelPlazaCategoryClaude   = "claude"
	ModelPlazaCategoryOpenAI   = "openai"
	ModelPlazaCategoryGemini   = "gemini"
	ModelPlazaCategoryDomestic = "domestic" // 国产模型
	ModelPlazaCategoryImage    = "image"    // 生图模型
)

// ModelPlazaCategoryKeys is the canonical, ordered set of category keys.
var ModelPlazaCategoryKeys = []string{
	ModelPlazaCategoryClaude,
	ModelPlazaCategoryOpenAI,
	ModelPlazaCategoryGemini,
	ModelPlazaCategoryDomestic,
	ModelPlazaCategoryImage,
}

func isValidModelPlazaCategory(key string) bool {
	for _, k := range ModelPlazaCategoryKeys {
		if k == key {
			return true
		}
	}
	return false
}

// ModelPlazaConfig is the whole plaza display config (excludes the on/off
// switch, which lives in the public settings flag SettingKeyModelPlazaEnabled).
// Two levels: fixed Categories -> admin-configured Groups -> Models.
type ModelPlazaConfig struct {
	Categories []ModelPlazaCategory `json:"categories"`
}

// --- Validation: bounds only, to prevent abuse. No business semantics. ---

const (
	maxModelPlazaGroups        = 100 // per category
	maxModelPlazaModelsPerGrp  = 500
	maxModelPlazaStringLen     = 200
	maxModelPlazaLongStringLen = 2000
)

func validateModelPlazaConfig(cfg *ModelPlazaConfig) error {
	if cfg == nil {
		return nil
	}
	if len(cfg.Categories) > len(ModelPlazaCategoryKeys) {
		return fmt.Errorf("too many categories (max %d)", len(ModelPlazaCategoryKeys))
	}
	for ci, cat := range cfg.Categories {
		if !isValidModelPlazaCategory(cat.Key) {
			return fmt.Errorf("category[%d]: invalid key %q", ci, cat.Key)
		}
		if len(cat.Groups) > maxModelPlazaGroups {
			return fmt.Errorf("category %q: too many groups (max %d)", cat.Key, maxModelPlazaGroups)
		}
		for gi, g := range cat.Groups {
			if len(g.Name) > maxModelPlazaStringLen {
				return fmt.Errorf("category %q group[%d]: name too long (max %d)", cat.Key, gi, maxModelPlazaStringLen)
			}
			if len(g.Description) > maxModelPlazaLongStringLen {
				return fmt.Errorf("category %q group[%d]: description too long (max %d)", cat.Key, gi, maxModelPlazaLongStringLen)
			}
			if len(g.Currency) > maxModelPlazaStringLen {
				return fmt.Errorf("category %q group[%d]: currency too long (max %d)", cat.Key, gi, maxModelPlazaStringLen)
			}
			if len(g.Models) > maxModelPlazaModelsPerGrp {
				return fmt.Errorf("category %q group[%d]: too many models (max %d)", cat.Key, gi, maxModelPlazaModelsPerGrp)
			}
			for mi, m := range g.Models {
				if len(m.Name) > maxModelPlazaStringLen {
					return fmt.Errorf("category %q group[%d].model[%d]: name too long (max %d)", cat.Key, gi, mi, maxModelPlazaStringLen)
				}
				if len(m.Note) > maxModelPlazaLongStringLen {
					return fmt.Errorf("category %q group[%d].model[%d]: note too long (max %d)", cat.Key, gi, mi, maxModelPlazaLongStringLen)
				}
				for _, v := range []string{m.Price.Input, m.Price.Output, m.Price.CacheRead, m.Price.CacheWrite} {
					if len(v) > maxModelPlazaStringLen {
						return fmt.Errorf("category %q group[%d].model[%d]: price value too long (max %d)", cat.Key, gi, mi, maxModelPlazaStringLen)
					}
				}
			}
		}
	}
	return nil
}

// normalizeModelPlazaCategories returns exactly the fixed categories in fixed
// order, carrying over any stored groups by key and dropping unknown keys.
// This guarantees both the admin editor and the user view always see all five
// categories in a stable order, regardless of what was persisted.
func normalizeModelPlazaCategories(cfg *ModelPlazaConfig) *ModelPlazaConfig {
	byKey := make(map[string][]ModelPlazaGroup, len(ModelPlazaCategoryKeys))
	if cfg != nil {
		for _, c := range cfg.Categories {
			if isValidModelPlazaCategory(c.Key) {
				byKey[c.Key] = c.Groups
			}
		}
	}
	out := &ModelPlazaConfig{Categories: make([]ModelPlazaCategory, 0, len(ModelPlazaCategoryKeys))}
	for _, key := range ModelPlazaCategoryKeys {
		groups := byKey[key]
		if groups == nil {
			groups = []ModelPlazaGroup{}
		}
		out.Categories = append(out.Categories, ModelPlazaCategory{Key: key, Groups: groups})
	}
	return out
}

// --- In-process cache (same pattern as web search emulation config) ---

const sfKeyModelPlazaConfig = "model_plaza_config"

type cachedModelPlazaConfig struct {
	config    *ModelPlazaConfig
	expiresAt int64 // unix nano
}

var modelPlazaCache atomic.Value // *cachedModelPlazaConfig
var modelPlazaSF singleflight.Group

const (
	modelPlazaCacheTTL  = 60 * time.Second
	modelPlazaErrorTTL  = 5 * time.Second
	modelPlazaDBTimeout = 5 * time.Second
)

// GetModelPlazaConfig returns the display config with in-process cache + singleflight.
// Missing/unset returns an empty config (never an error to the caller's UI).
func (s *SettingService) GetModelPlazaConfig(ctx context.Context) (*ModelPlazaConfig, error) {
	if cached := modelPlazaCache.Load(); cached != nil {
		if c, ok := cached.(*cachedModelPlazaConfig); ok && time.Now().UnixNano() < c.expiresAt {
			return c.config, nil
		}
	}
	result, err, _ := modelPlazaSF.Do(sfKeyModelPlazaConfig, func() (any, error) {
		return s.loadModelPlazaConfigFromDB()
	})
	if err != nil {
		return normalizeModelPlazaCategories(nil), err
	}
	if cfg, ok := result.(*ModelPlazaConfig); ok {
		return cfg, nil
	}
	return normalizeModelPlazaCategories(nil), nil
}

func (s *SettingService) loadModelPlazaConfigFromDB() (*ModelPlazaConfig, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), modelPlazaDBTimeout)
	defer cancel()

	raw, err := s.settingRepo.GetValue(dbCtx, SettingKeyModelPlazaConfig)
	if err != nil {
		// Not found is expected when the feature was never configured; treat as empty.
		if isSettingNotFound(err) {
			cfg := normalizeModelPlazaCategories(nil)
			modelPlazaCache.Store(&cachedModelPlazaConfig{
				config:    cfg,
				expiresAt: time.Now().Add(modelPlazaCacheTTL).UnixNano(),
			})
			return cfg, nil
		}
		modelPlazaCache.Store(&cachedModelPlazaConfig{
			config:    normalizeModelPlazaCategories(nil),
			expiresAt: time.Now().Add(modelPlazaErrorTTL).UnixNano(),
		})
		return normalizeModelPlazaCategories(nil), err
	}
	cfg := parseModelPlazaConfigJSON(raw)
	modelPlazaCache.Store(&cachedModelPlazaConfig{
		config:    cfg,
		expiresAt: time.Now().Add(modelPlazaCacheTTL).UnixNano(),
	})
	return cfg, nil
}

// parseModelPlazaConfigJSON parses stored JSON and normalizes it to the fixed
// five categories in order (invalid JSON falls back to an empty-but-normalized config).
func parseModelPlazaConfigJSON(raw string) *ModelPlazaConfig {
	cfg := &ModelPlazaConfig{}
	if strings.TrimSpace(raw) != "" {
		if err := json.Unmarshal([]byte(raw), cfg); err != nil {
			slog.Warn("model plaza: failed to parse config JSON", "error", err)
			cfg = &ModelPlazaConfig{}
		}
	}
	return normalizeModelPlazaCategories(cfg)
}

// SaveModelPlazaConfig validates and persists the display config.
func (s *SettingService) SaveModelPlazaConfig(ctx context.Context, cfg *ModelPlazaConfig) error {
	if err := validateModelPlazaConfig(cfg); err != nil {
		return infraerrors.BadRequest("INVALID_MODEL_PLAZA_CONFIG", err.Error())
	}
	// Normalize to the fixed five categories in order before persisting so the
	// stored JSON and cache are always in canonical shape.
	cfg = normalizeModelPlazaCategories(cfg)
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("model plaza: marshal config: %w", err)
	}
	if err := s.settingRepo.Set(ctx, SettingKeyModelPlazaConfig, string(data)); err != nil {
		return fmt.Errorf("model plaza: save config: %w", err)
	}
	// Invalidate: forget singleflight first, then store new value.
	modelPlazaSF.Forget(sfKeyModelPlazaConfig)
	modelPlazaCache.Store(&cachedModelPlazaConfig{
		config:    cfg,
		expiresAt: time.Now().Add(modelPlazaCacheTTL).UnixNano(),
	})
	return nil
}

// IsModelPlazaEnabled reports whether the Model Plaza feature switch is on.
// Fail-closed: on error returns false, matching the opt-in default
// (unknown ↔ disabled). Reads the flag directly from the settings store.
func (s *SettingService) IsModelPlazaEnabled(ctx context.Context) bool {
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyModelPlazaEnabled)
	if err != nil {
		return false
	}
	return raw == "true"
}

// SetModelPlazaEnabled persists the Model Plaza feature switch.
func (s *SettingService) SetModelPlazaEnabled(ctx context.Context, enabled bool) error {
	if err := s.settingRepo.Set(ctx, SettingKeyModelPlazaEnabled, strconv.FormatBool(enabled)); err != nil {
		return fmt.Errorf("model plaza: save enabled flag: %w", err)
	}
	return nil
}

// isSettingNotFound reports whether err is the service's "setting not found" sentinel.
func isSettingNotFound(err error) bool {
	return err != nil && errors.Is(err, ErrSettingNotFound)
}

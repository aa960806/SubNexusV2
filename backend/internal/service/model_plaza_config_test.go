//go:build unit

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseModelPlazaConfigJSON(t *testing.T) {
	t.Run("empty returns the five fixed categories in order", func(t *testing.T) {
		cfg := parseModelPlazaConfigJSON("")
		require.NotNil(t, cfg)
		require.Len(t, cfg.Categories, len(ModelPlazaCategoryKeys))
		for i, key := range ModelPlazaCategoryKeys {
			require.Equal(t, key, cfg.Categories[i].Key)
			require.Empty(t, cfg.Categories[i].Groups)
		}
	})

	t.Run("invalid JSON falls back to normalized empty", func(t *testing.T) {
		cfg := parseModelPlazaConfigJSON("{not json")
		require.Len(t, cfg.Categories, len(ModelPlazaCategoryKeys))
	})

	t.Run("valid JSON round-trips and normalizes order", func(t *testing.T) {
		raw := `{"categories":[{"key":"claude","groups":[{"name":"标准池","currency":"$","models":[{"name":"opus","price":{"input":"15","output":"75","cache_read":"1.5","cache_write":"18.75"},"note":"flagship"}]}]}]}`
		cfg := parseModelPlazaConfigJSON(raw)
		require.Len(t, cfg.Categories, len(ModelPlazaCategoryKeys))
		// claude is first in the fixed order and carries the group
		require.Equal(t, ModelPlazaCategoryClaude, cfg.Categories[0].Key)
		require.Len(t, cfg.Categories[0].Groups, 1)
		g := cfg.Categories[0].Groups[0]
		require.Equal(t, "标准池", g.Name)
		require.Equal(t, "$", g.Currency)
		require.Len(t, g.Models, 1)
		require.Equal(t, "opus", g.Models[0].Name)
		require.Equal(t, "18.75", g.Models[0].Price.CacheWrite)
		require.Equal(t, "flagship", g.Models[0].Note)
		// other categories present but empty
		require.Empty(t, cfg.Categories[1].Groups)
	})

	t.Run("unknown category keys are dropped", func(t *testing.T) {
		raw := `{"categories":[{"key":"bogus","groups":[{"name":"x","models":[]}]}]}`
		cfg := parseModelPlazaConfigJSON(raw)
		require.Len(t, cfg.Categories, len(ModelPlazaCategoryKeys))
		for _, c := range cfg.Categories {
			require.True(t, isValidModelPlazaCategory(c.Key))
			require.Empty(t, c.Groups)
		}
	})
}

func TestNormalizeModelPlazaCategories(t *testing.T) {
	t.Run("nil yields five empty categories in order", func(t *testing.T) {
		cfg := normalizeModelPlazaCategories(nil)
		require.Len(t, cfg.Categories, len(ModelPlazaCategoryKeys))
		for i, key := range ModelPlazaCategoryKeys {
			require.Equal(t, key, cfg.Categories[i].Key)
			require.NotNil(t, cfg.Categories[i].Groups) // non-nil empty slice
		}
	})

	t.Run("image category is supported", func(t *testing.T) {
		require.True(t, isValidModelPlazaCategory(ModelPlazaCategoryImage))
		cfg := normalizeModelPlazaCategories(&ModelPlazaConfig{Categories: []ModelPlazaCategory{
			{Key: ModelPlazaCategoryImage, Groups: []ModelPlazaGroup{{Name: "g1"}}},
		}})
		// find image category
		var found bool
		for _, c := range cfg.Categories {
			if c.Key == ModelPlazaCategoryImage {
				found = true
				require.Len(t, c.Groups, 1)
			}
		}
		require.True(t, found)
	})
}

func TestValidateModelPlazaConfig(t *testing.T) {
	t.Run("nil is valid", func(t *testing.T) {
		require.NoError(t, validateModelPlazaConfig(nil))
	})

	t.Run("normal config is valid", func(t *testing.T) {
		cfg := &ModelPlazaConfig{Categories: []ModelPlazaCategory{
			{Key: ModelPlazaCategoryClaude, Groups: []ModelPlazaGroup{
				{Name: "A", Models: []ModelPlazaModel{{Name: "m1", Price: ModelPlazaPrice{Input: "1"}}}},
			}},
		}}
		require.NoError(t, validateModelPlazaConfig(cfg))
	})

	t.Run("invalid category key rejected", func(t *testing.T) {
		cfg := &ModelPlazaConfig{Categories: []ModelPlazaCategory{{Key: "nope"}}}
		require.Error(t, validateModelPlazaConfig(cfg))
	})

	t.Run("too many groups rejected", func(t *testing.T) {
		cfg := &ModelPlazaConfig{Categories: []ModelPlazaCategory{
			{Key: ModelPlazaCategoryOpenAI, Groups: make([]ModelPlazaGroup, maxModelPlazaGroups+1)},
		}}
		require.Error(t, validateModelPlazaConfig(cfg))
	})

	t.Run("too many models rejected", func(t *testing.T) {
		cfg := &ModelPlazaConfig{Categories: []ModelPlazaCategory{
			{Key: ModelPlazaCategoryGemini, Groups: []ModelPlazaGroup{
				{Name: "A", Models: make([]ModelPlazaModel, maxModelPlazaModelsPerGrp+1)},
			}},
		}}
		require.Error(t, validateModelPlazaConfig(cfg))
	})

	t.Run("overlong price value rejected", func(t *testing.T) {
		cfg := &ModelPlazaConfig{Categories: []ModelPlazaCategory{
			{Key: ModelPlazaCategoryDomestic, Groups: []ModelPlazaGroup{
				{Name: "A", Models: []ModelPlazaModel{
					{Name: "m1", Price: ModelPlazaPrice{Input: strings.Repeat("9", maxModelPlazaStringLen+1)}},
				}},
			}},
		}}
		require.Error(t, validateModelPlazaConfig(cfg))
	})
}

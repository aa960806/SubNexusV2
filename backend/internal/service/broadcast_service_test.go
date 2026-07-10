package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type broadcastSettingRepo struct {
	values map[string]string
}

func (r *broadcastSettingRepo) Get(_ context.Context, key string) (*Setting, error) {
	value, ok := r.values[key]
	if !ok {
		return nil, ErrSettingNotFound
	}
	return &Setting{Key: key, Value: value}, nil
}

func (r *broadcastSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	value, ok := r.values[key]
	if !ok {
		return "", ErrSettingNotFound
	}
	return value, nil
}

func (r *broadcastSettingRepo) Set(_ context.Context, key, value string) error {
	r.values[key] = value
	return nil
}

func (r *broadcastSettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	values := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			values[key] = value
		}
	}
	return values, nil
}

func (r *broadcastSettingRepo) SetMultiple(_ context.Context, settings map[string]string) error {
	for key, value := range settings {
		r.values[key] = value
	}
	return nil
}

func (r *broadcastSettingRepo) GetAll(_ context.Context) (map[string]string, error) {
	values := make(map[string]string, len(r.values))
	for key, value := range r.values {
		values[key] = value
	}
	return values, nil
}

func (r *broadcastSettingRepo) Delete(_ context.Context, key string) error {
	delete(r.values, key)
	return nil
}

func TestBroadcastServiceUpdateConfigPreservesUnknownActivityConfigFields(t *testing.T) {
	repo := &broadcastSettingRepo{values: map[string]string{
		SettingKeyActivityConfig: `{"daily_spin_enabled":true,"nested":{"prizes":[1,2,3]},"broadcast_enabled":true}`,
	}}
	service := NewBroadcastService(nil, repo)

	updated, err := service.UpdateConfig(context.Background(), BroadcastConfig{Enabled: false})
	require.NoError(t, err)
	require.False(t, updated.Enabled)

	var stored map[string]any
	require.NoError(t, json.Unmarshal([]byte(repo.values[SettingKeyActivityConfig]), &stored))
	require.Equal(t, false, stored["broadcast_enabled"])
	require.Equal(t, true, stored["daily_spin_enabled"])
	require.Equal(t, map[string]any{"prizes": []any{float64(1), float64(2), float64(3)}}, stored["nested"])
}

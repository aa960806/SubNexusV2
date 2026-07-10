package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	SettingKeyActivityConfig      = "ACTIVITY_CONFIG"
	ActivityBroadcastSourceAdmin  = "admin"
	ActivityBroadcastSourceSystem = "system"
)

type BroadcastConfig struct {
	Enabled bool `json:"enabled"`
}

type ActivityBroadcast struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Source    string     `json:"source"`
	Enabled   bool       `json:"enabled"`
	Priority  int        `json:"priority"`
	StartAt   *time.Time `json:"start_at,omitempty"`
	EndAt     *time.Time `json:"end_at,omitempty"`
	CreatedBy *int64     `json:"created_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ActivityBroadcastInput struct {
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Enabled  bool       `json:"enabled"`
	Priority int        `json:"priority"`
	StartAt  *time.Time `json:"start_at"`
	EndAt    *time.Time `json:"end_at"`
}

type BroadcastService struct {
	db          *sql.DB
	settingRepo SettingRepository
}

func NewBroadcastService(db *sql.DB, settingRepo SettingRepository) *BroadcastService {
	return &BroadcastService{db: db, settingRepo: settingRepo}
}

func (s *BroadcastService) GetConfig(ctx context.Context) (BroadcastConfig, error) {
	cfg := BroadcastConfig{Enabled: true}
	if s == nil || s.settingRepo == nil {
		return cfg, nil
	}
	setting, err := s.settingRepo.Get(ctx, SettingKeyActivityConfig)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return cfg, nil
		}
		return cfg, err
	}
	var stored struct {
		BroadcastEnabled *bool `json:"broadcast_enabled"`
	}
	if json.Unmarshal([]byte(setting.Value), &stored) == nil && stored.BroadcastEnabled != nil {
		cfg.Enabled = *stored.BroadcastEnabled
	}
	return cfg, nil
}

func (s *BroadcastService) UpdateConfig(ctx context.Context, cfg BroadcastConfig) (BroadcastConfig, error) {
	if s == nil || s.settingRepo == nil {
		return BroadcastConfig{}, infraerrors.InternalServer("BROADCAST_SETTINGS_UNAVAILABLE", "broadcast settings repository is unavailable")
	}
	stored := map[string]json.RawMessage{}
	if setting, err := s.settingRepo.Get(ctx, SettingKeyActivityConfig); err == nil {
		_ = json.Unmarshal([]byte(setting.Value), &stored)
	} else if !errors.Is(err, ErrSettingNotFound) {
		return BroadcastConfig{}, err
	}
	enabled, err := json.Marshal(cfg.Enabled)
	if err != nil {
		return BroadcastConfig{}, err
	}
	stored["broadcast_enabled"] = enabled
	raw, err := json.Marshal(stored)
	if err != nil {
		return BroadcastConfig{}, err
	}
	if err := s.settingRepo.Set(ctx, SettingKeyActivityConfig, string(raw)); err != nil {
		return BroadcastConfig{}, err
	}
	return cfg, nil
}

func (s *BroadcastService) List(ctx context.Context, activeOnly bool, limit int) ([]ActivityBroadcast, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("BROADCAST_DB_UNAVAILABLE", "broadcast database is unavailable")
	}
	if activeOnly {
		cfg, err := s.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		if !cfg.Enabled {
			return []ActivityBroadcast{}, nil
		}
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query := `
		SELECT id, title, content, source, enabled, priority, start_at, end_at, created_by, created_at, updated_at
		FROM activity_broadcasts`
	if activeOnly {
		query += `
		WHERE enabled = TRUE
		  AND (start_at IS NULL OR start_at <= NOW())
		  AND (end_at IS NULL OR end_at >= NOW())`
	}
	query += ` ORDER BY priority DESC, created_at DESC LIMIT $1`
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]ActivityBroadcast, 0, limit)
	for rows.Next() {
		item, err := scanActivityBroadcast(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *BroadcastService) Create(ctx context.Context, input ActivityBroadcastInput, adminID int64) (*ActivityBroadcast, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("BROADCAST_DB_UNAVAILABLE", "broadcast database is unavailable")
	}
	input = normalizeBroadcastInput(input)
	if input.Content == "" {
		return nil, infraerrors.BadRequest("BROADCAST_CONTENT_REQUIRED", "broadcast content is required")
	}
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO activity_broadcasts (title, content, source, enabled, priority, start_at, end_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, 0))
		RETURNING id, title, content, source, enabled, priority, start_at, end_at, created_by, created_at, updated_at
	`, input.Title, input.Content, ActivityBroadcastSourceAdmin, input.Enabled, input.Priority, input.StartAt, input.EndAt, adminID)
	item, err := scanActivityBroadcast(row)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *BroadcastService) Update(ctx context.Context, id int64, input ActivityBroadcastInput) (*ActivityBroadcast, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("BROADCAST_DB_UNAVAILABLE", "broadcast database is unavailable")
	}
	input = normalizeBroadcastInput(input)
	if input.Content == "" {
		return nil, infraerrors.BadRequest("BROADCAST_CONTENT_REQUIRED", "broadcast content is required")
	}
	row := s.db.QueryRowContext(ctx, `
		UPDATE activity_broadcasts
		SET title = $2, content = $3, enabled = $4, priority = $5, start_at = $6, end_at = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, title, content, source, enabled, priority, start_at, end_at, created_by, created_at, updated_at
	`, id, input.Title, input.Content, input.Enabled, input.Priority, input.StartAt, input.EndAt)
	item, err := scanActivityBroadcast(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, infraerrors.NotFound("BROADCAST_NOT_FOUND", "broadcast not found")
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *BroadcastService) Delete(ctx context.Context, id int64) error {
	if s == nil || s.db == nil {
		return infraerrors.InternalServer("BROADCAST_DB_UNAVAILABLE", "broadcast database is unavailable")
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM activity_broadcasts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return infraerrors.NotFound("BROADCAST_NOT_FOUND", "broadcast not found")
	}
	return nil
}

func (s *BroadcastService) CleanupExpiredSystem(ctx context.Context, retentionDays int) (int64, error) {
	if s == nil || s.db == nil {
		return 0, infraerrors.InternalServer("BROADCAST_DB_UNAVAILABLE", "broadcast database is unavailable")
	}
	if retentionDays < 0 {
		retentionDays = 0
	}
	if retentionDays > 3650 {
		retentionDays = 3650
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM activity_broadcasts
		WHERE source = $1 AND end_at IS NOT NULL AND end_at < $2
	`, ActivityBroadcastSourceSystem, cutoff)
	if err != nil {
		return 0, err
	}
	deleted, _ := result.RowsAffected()
	return deleted, nil
}

type activityBroadcastScanner interface {
	Scan(dest ...any) error
}

func scanActivityBroadcast(scanner activityBroadcastScanner) (ActivityBroadcast, error) {
	var item ActivityBroadcast
	var startAt sql.NullTime
	var endAt sql.NullTime
	var createdBy sql.NullInt64
	if err := scanner.Scan(&item.ID, &item.Title, &item.Content, &item.Source, &item.Enabled, &item.Priority, &startAt, &endAt, &createdBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
		return item, err
	}
	if startAt.Valid {
		item.StartAt = &startAt.Time
	}
	if endAt.Valid {
		item.EndAt = &endAt.Time
	}
	if createdBy.Valid {
		item.CreatedBy = &createdBy.Int64
	}
	return item, nil
}

func normalizeBroadcastInput(input ActivityBroadcastInput) ActivityBroadcastInput {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if input.Priority < 0 {
		input.Priority = 0
	}
	if input.Priority > 1000 {
		input.Priority = 1000
	}
	if input.StartAt != nil && input.EndAt != nil && input.EndAt.Before(*input.StartAt) {
		input.EndAt = nil
	}
	return input
}

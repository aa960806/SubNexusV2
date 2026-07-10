package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const SettingKeyActivityCenterConfig = "ACTIVITY_CENTER_CONFIG"

type ActivityCenterConfig struct {
	Enabled bool `json:"enabled"`
}

type ActivityCenterItem struct {
	ID          int64           `json:"id"`
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	Subtitle    string          `json:"subtitle"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	CoverURL    string          `json:"cover_url"`
	RoutePath   string          `json:"route_path"`
	ExternalURL string          `json:"external_url"`
	ActionLabel string          `json:"action_label"`
	Type        string          `json:"activity_type"`
	Enabled     bool            `json:"enabled"`
	SortOrder   int             `json:"sort_order"`
	StartAt     *time.Time      `json:"start_at,omitempty"`
	EndAt       *time.Time      `json:"end_at,omitempty"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedBy   *int64          `json:"created_by,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ActivityCenterItemInput struct {
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	Subtitle    string          `json:"subtitle"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	CoverURL    string          `json:"cover_url"`
	RoutePath   string          `json:"route_path"`
	ExternalURL string          `json:"external_url"`
	ActionLabel string          `json:"action_label"`
	Type        string          `json:"activity_type"`
	Enabled     bool            `json:"enabled"`
	SortOrder   int             `json:"sort_order"`
	StartAt     *time.Time      `json:"start_at"`
	EndAt       *time.Time      `json:"end_at"`
	Metadata    json.RawMessage `json:"metadata"`
}

type ActivityCenterListResponse struct {
	Enabled bool                 `json:"enabled"`
	Items   []ActivityCenterItem `json:"items"`
}

type ActivityCenterService struct {
	db          *sql.DB
	settingRepo SettingRepository
}

func NewActivityCenterService(db *sql.DB, settingRepo SettingRepository) *ActivityCenterService {
	return &ActivityCenterService{db: db, settingRepo: settingRepo}
}

func DefaultActivityCenterConfig() ActivityCenterConfig {
	return ActivityCenterConfig{Enabled: false}
}

func (s *ActivityCenterService) GetConfig(ctx context.Context) (ActivityCenterConfig, error) {
	cfg := DefaultActivityCenterConfig()
	if s == nil || s.settingRepo == nil {
		return cfg, nil
	}
	setting, err := s.settingRepo.Get(ctx, SettingKeyActivityCenterConfig)
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return cfg, nil
		}
		return cfg, err
	}
	if strings.TrimSpace(setting.Value) == "" {
		return cfg, nil
	}
	if err := json.Unmarshal([]byte(setting.Value), &cfg); err != nil {
		return DefaultActivityCenterConfig(), nil
	}
	return cfg, nil
}

func (s *ActivityCenterService) UpdateConfig(ctx context.Context, cfg ActivityCenterConfig) (ActivityCenterConfig, error) {
	if s == nil || s.settingRepo == nil {
		return ActivityCenterConfig{}, infraerrors.InternalServer("ACTIVITY_CENTER_SETTINGS_UNAVAILABLE", "activity center settings repository is unavailable")
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return cfg, err
	}
	if err := s.settingRepo.Set(ctx, SettingKeyActivityCenterConfig, string(raw)); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *ActivityCenterService) ListVisible(ctx context.Context, now time.Time) (*ActivityCenterListResponse, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	resp := &ActivityCenterListResponse{Enabled: cfg.Enabled, Items: []ActivityCenterItem{}}
	if !cfg.Enabled {
		return resp, nil
	}
	items, err := s.list(ctx, true, now)
	if err != nil {
		return nil, err
	}
	resp.Items = items
	return resp, nil
}

func (s *ActivityCenterService) ListAdmin(ctx context.Context) ([]ActivityCenterItem, error) {
	return s.list(ctx, false, time.Now())
}

func (s *ActivityCenterService) list(ctx context.Context, visibleOnly bool, now time.Time) ([]ActivityCenterItem, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("ACTIVITY_CENTER_DB_UNAVAILABLE", "activity center database is unavailable")
	}
	query := `
		SELECT id, slug, title, subtitle, description, icon, cover_url, route_path, external_url,
		       action_label, activity_type, enabled, sort_order, start_at, end_at, metadata,
		       created_by, created_at, updated_at
		FROM activity_center_items`
	args := []any{}
	if visibleOnly {
		query += `
		WHERE enabled = TRUE
		  AND (start_at IS NULL OR start_at <= $1)
		  AND (end_at IS NULL OR end_at >= $1)`
		args = append(args, now)
	}
	query += `
		ORDER BY sort_order ASC, created_at DESC, id DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := []ActivityCenterItem{}
	for rows.Next() {
		item, err := scanActivityCenterItem(rows)
		if err != nil {
			return nil, err
		}
		if visibleOnly {
			allowed, err := s.isActivityCenterItemEnabledByFeature(ctx, item)
			if err != nil {
				return nil, err
			}
			if !allowed {
				continue
			}
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ActivityCenterService) isActivityCenterItemEnabledByFeature(ctx context.Context, item ActivityCenterItem) (bool, error) {
	_ = ctx
	switch strings.ToLower(strings.TrimSpace(item.Type)) {
	case "daily_spin", "invite_lottery", "recharge_wheel", "invite_milestone":
		return false, nil
	default:
		return true, nil
	}
}

func (s *ActivityCenterService) Create(ctx context.Context, input ActivityCenterItemInput, adminID int64) (*ActivityCenterItem, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("ACTIVITY_CENTER_DB_UNAVAILABLE", "activity center database is unavailable")
	}
	input = normalizeActivityCenterItemInput(input)
	if input.Slug == "" || input.Title == "" {
		return nil, infraerrors.BadRequest("ACTIVITY_CENTER_ITEM_INVALID", "activity slug and title are required")
	}
	if err := validateActivityCenterDestination(input); err != nil {
		return nil, infraerrors.BadRequest("ACTIVITY_CENTER_DESTINATION_INVALID", err.Error())
	}
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO activity_center_items
			(slug, title, subtitle, description, icon, cover_url, route_path, external_url,
			 action_label, activity_type, enabled, sort_order, start_at, end_at, metadata, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15::jsonb,NULLIF($16,0))
		RETURNING id, slug, title, subtitle, description, icon, cover_url, route_path, external_url,
		          action_label, activity_type, enabled, sort_order, start_at, end_at, metadata,
		          created_by, created_at, updated_at
	`, input.Slug, input.Title, input.Subtitle, input.Description, input.Icon, input.CoverURL,
		input.RoutePath, input.ExternalURL, input.ActionLabel, input.Type, input.Enabled, input.SortOrder,
		input.StartAt, input.EndAt, string(input.Metadata), adminID)
	item, err := scanActivityCenterItem(row)
	if err != nil {
		if isActivityCenterSlugDuplicate(err) {
			return nil, infraerrors.BadRequest("ACTIVITY_CENTER_SLUG_DUPLICATE", "activity slug already exists")
		}
		return nil, err
	}
	return &item, nil
}

func (s *ActivityCenterService) Update(ctx context.Context, id int64, input ActivityCenterItemInput) (*ActivityCenterItem, error) {
	if s == nil || s.db == nil {
		return nil, infraerrors.InternalServer("ACTIVITY_CENTER_DB_UNAVAILABLE", "activity center database is unavailable")
	}
	input = normalizeActivityCenterItemInput(input)
	if id <= 0 || input.Slug == "" || input.Title == "" {
		return nil, infraerrors.BadRequest("ACTIVITY_CENTER_ITEM_INVALID", "activity id, slug and title are required")
	}
	if err := validateActivityCenterDestination(input); err != nil {
		return nil, infraerrors.BadRequest("ACTIVITY_CENTER_DESTINATION_INVALID", err.Error())
	}
	row := s.db.QueryRowContext(ctx, `
		UPDATE activity_center_items
		SET slug=$2, title=$3, subtitle=$4, description=$5, icon=$6, cover_url=$7,
		    route_path=$8, external_url=$9, action_label=$10, activity_type=$11,
		    enabled=$12, sort_order=$13, start_at=$14, end_at=$15, metadata=$16::jsonb, updated_at=NOW()
		WHERE id=$1
		RETURNING id, slug, title, subtitle, description, icon, cover_url, route_path, external_url,
		          action_label, activity_type, enabled, sort_order, start_at, end_at, metadata,
		          created_by, created_at, updated_at
	`, id, input.Slug, input.Title, input.Subtitle, input.Description, input.Icon, input.CoverURL,
		input.RoutePath, input.ExternalURL, input.ActionLabel, input.Type, input.Enabled, input.SortOrder,
		input.StartAt, input.EndAt, string(input.Metadata))
	item, err := scanActivityCenterItem(row)
	if err == sql.ErrNoRows {
		return nil, infraerrors.NotFound("ACTIVITY_CENTER_ITEM_NOT_FOUND", "activity center item not found")
	}
	if err != nil {
		if isActivityCenterSlugDuplicate(err) {
			return nil, infraerrors.BadRequest("ACTIVITY_CENTER_SLUG_DUPLICATE", "activity slug already exists")
		}
		return nil, err
	}
	return &item, nil
}

func (s *ActivityCenterService) Delete(ctx context.Context, id int64) error {
	if s == nil || s.db == nil {
		return infraerrors.InternalServer("ACTIVITY_CENTER_DB_UNAVAILABLE", "activity center database is unavailable")
	}
	res, err := s.db.ExecContext(ctx, `DELETE FROM activity_center_items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return infraerrors.NotFound("ACTIVITY_CENTER_ITEM_NOT_FOUND", "activity center item not found")
	}
	return nil
}

type activityCenterScanner interface {
	Scan(dest ...any) error
}

func scanActivityCenterItem(scanner activityCenterScanner) (ActivityCenterItem, error) {
	var item ActivityCenterItem
	var startAt sql.NullTime
	var endAt sql.NullTime
	var createdBy sql.NullInt64
	var metadata []byte
	if err := scanner.Scan(
		&item.ID, &item.Slug, &item.Title, &item.Subtitle, &item.Description, &item.Icon, &item.CoverURL,
		&item.RoutePath, &item.ExternalURL, &item.ActionLabel, &item.Type, &item.Enabled, &item.SortOrder,
		&startAt, &endAt, &metadata, &createdBy, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
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
	if len(metadata) == 0 {
		metadata = []byte("{}")
	}
	item.Metadata = json.RawMessage(metadata)
	return item, nil
}

func normalizeActivityCenterItemInput(input ActivityCenterItemInput) ActivityCenterItemInput {
	input.Slug = strings.ToLower(strings.TrimSpace(input.Slug))
	input.Title = strings.TrimSpace(input.Title)
	input.Subtitle = strings.TrimSpace(input.Subtitle)
	input.Description = strings.TrimSpace(input.Description)
	input.Icon = strings.TrimSpace(input.Icon)
	if input.Icon == "" {
		input.Icon = "gift"
	}
	input.CoverURL = strings.TrimSpace(input.CoverURL)
	input.RoutePath = strings.TrimSpace(input.RoutePath)
	input.ExternalURL = strings.TrimSpace(input.ExternalURL)
	input.ActionLabel = strings.TrimSpace(input.ActionLabel)
	if input.ActionLabel == "" {
		input.ActionLabel = "查看"
	}
	input.Type = strings.TrimSpace(input.Type)
	if input.Type == "" {
		input.Type = "custom"
	}
	if len(input.Metadata) == 0 || !json.Valid(input.Metadata) {
		input.Metadata = json.RawMessage(`{}`)
	}
	if input.StartAt != nil && input.EndAt != nil && input.EndAt.Before(*input.StartAt) {
		input.EndAt = nil
	}
	return input
}

func validateActivityCenterDestination(input ActivityCenterItemInput) error {
	if input.ExternalURL != "" {
		parsed, err := url.ParseRequestURI(input.ExternalURL)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return errors.New("external_url must be an absolute HTTP or HTTPS URL")
		}
	}
	if input.RoutePath != "" && (!strings.HasPrefix(input.RoutePath, "/") || strings.HasPrefix(input.RoutePath, "//")) {
		return errors.New("route_path must be a site-relative path beginning with a single slash")
	}
	return nil
}

func isActivityCenterSlugDuplicate(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "activity_center_items_slug_key") ||
		(strings.Contains(msg, "duplicate key") && strings.Contains(msg, "slug"))
}

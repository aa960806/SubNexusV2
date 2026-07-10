package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActivityCenterRemovedGameplayTypesStayHidden(t *testing.T) {
	service := &ActivityCenterService{}
	removedTypes := []string{"daily_spin", "invite_lottery", "recharge_wheel", "invite_milestone"}

	for _, activityType := range removedTypes {
		t.Run(activityType, func(t *testing.T) {
			allowed, err := service.isActivityCenterItemEnabledByFeature(
				context.Background(),
				ActivityCenterItem{Type: activityType},
			)
			require.NoError(t, err)
			require.False(t, allowed)
		})
	}

	allowed, err := service.isActivityCenterItemEnabledByFeature(
		context.Background(),
		ActivityCenterItem{Type: "custom"},
	)
	require.NoError(t, err)
	require.True(t, allowed)
}

func TestValidateActivityCenterDestination(t *testing.T) {
	tests := []struct {
		name    string
		input   ActivityCenterItemInput
		wantErr bool
	}{
		{name: "empty destination", input: ActivityCenterItemInput{}},
		{name: "internal route", input: ActivityCenterItemInput{RoutePath: "/activities/summer"}},
		{name: "https external", input: ActivityCenterItemInput{ExternalURL: "https://example.com/event"}},
		{name: "javascript external", input: ActivityCenterItemInput{ExternalURL: "javascript:alert(1)"}, wantErr: true},
		{name: "protocol relative external", input: ActivityCenterItemInput{ExternalURL: "//example.com/event"}, wantErr: true},
		{name: "external route", input: ActivityCenterItemInput{RoutePath: "https://example.com/event"}, wantErr: true},
		{name: "protocol relative route", input: ActivityCenterItemInput{RoutePath: "//example.com/event"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateActivityCenterDestination(test.input)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

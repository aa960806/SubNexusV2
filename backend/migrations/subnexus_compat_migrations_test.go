package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigration174OnlyDisablesRemovedSubNexusActivityEntries(t *testing.T) {
	content, err := FS.ReadFile("174_disable_removed_subnexus_activity_entries.sql")
	require.NoError(t, err)

	sql := strings.ToLower(string(content))
	require.Contains(t, sql, "update activity_center_items")
	require.Contains(t, sql, "set enabled = false")
	require.NotContains(t, sql, "delete from")
	for _, slug := range []string{"daily-spin", "invite-lottery", "recharge-wheel", "invite-milestone"} {
		require.Contains(t, sql, "'"+slug+"'")
	}
}

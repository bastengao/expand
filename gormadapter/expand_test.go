package gormadapter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var expandMap = map[string]any{
	"Avatar": nil,
	"User":   expandSlice,
	"UserMap": map[string]any{
		"Avatar": nil,
	},
}

var expandSlice = []string{"Avatar"}

func TestValidateExpand(t *testing.T) {
	require.NoError(t, ValidateExpand([]string{}, expandMap))
	require.NoError(t, ValidateExpand([]string{"Avatar", "User", "User.Avatar"}, expandMap))
	require.NoError(t, ValidateExpand([]string{"Avatar"}, expandSlice))

	require.Error(t, ValidateExpand([]string{"not_exists"}, expandMap))
	require.Error(t, ValidateExpand([]string{"not_exists"}, expandSlice))

	t.Run("WithCamelCase", func(t *testing.T) {
		require.NoError(t, ValidateExpand([]string{}, expandMap, WithCamelCase))
		require.NoError(t, ValidateExpand([]string{"avatar", "user", "user.avatar"}, expandMap, WithCamelCase))
		require.NoError(t, ValidateExpand([]string{"avatar"}, expandSlice, WithCamelCase))

		require.Error(t, ValidateExpand([]string{"not_exists"}, expandMap, WithCamelCase))
		require.Error(t, ValidateExpand([]string{"not_exists"}, expandSlice, WithCamelCase))
	})
}

func TestShrinkExpand(t *testing.T) {
	require.Equal(t, []string{"User.Avatar"}, shrinkExpand([]string{"User", "User.Avatar"}))
	require.Equal(t, []string{"User.Avatar.ID"}, shrinkExpand([]string{"User", "User.Avatar", "User.Avatar.ID"}))
	require.Equal(t, []string{"User.Avatar.ID"}, shrinkExpand([]string{"User.Avatar.ID", "User", "User.Avatar"}))
	require.Equal(t, []string{"User.Avatar", "User.Address"}, shrinkExpand([]string{"User", "User.Avatar", "User.Address"}))
}

func TestCamelCaseExpand(t *testing.T) {
	require.Equal(t, "User", camelCaseExpand("user"))
	require.Equal(t, "User.AvatarImg", camelCaseExpand("user.avatar_img"))
}

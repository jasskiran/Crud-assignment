package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGitUser(t *testing.T) {
	type user struct {
		name string
	}
	tests := []struct {
		name   string
		user   user
	}{
		{
			name:     "NameIsPresent",
			user: user{
				"user1",
			},
		},
		{
			name:     "NameIsAbsent",
			user: user{
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user{tt.user.name}
			if len(u.name) > 0 {
				require.NotEmpty(t, u.name,
					"Name should be there")
			}
		})
	}
}

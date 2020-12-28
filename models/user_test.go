package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	type user struct {
		name string
		password string
	}
	tests := []struct {
		name   string
		user   user
	}{
		{
			name:     "NameIsPresentAndPasswordIsAbsent",
			user: user{
				"user1",
				"",
			},
		},
		{
			name:     "NameIsAbsentAndPasswordIsPresent",
			user: user{
				"",
				"user1",
			},
		},
		{
			name:     "NameAndPasswordAreAbsent",
			user: user{
				"",
				"",
			},
		},
		{
			name:     "NameAndPasswordArePresent",
			user: user{
				"user1",
				"Password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user{tt.user.name, tt.user.password}
			if len(u.name) > 0 {
				require.NotEmpty(t, u.name,
					"Name should be there")
			}
			if len(u.password) > 0 {
				require.NotEmpty(t, u.password,
					"password should be there")
			}
		})
	}
}

func TestUser_Validate(t *testing.T) {
	type user struct {
		name string
		password string
	}
	tests := []struct {
		name   string
		user   user
	}{
		{
			name:     "NameIsPresentAndPasswordIsAbsent",
			user: user{
				"user1",
				"",
			},
		},
		{
			name:     "NameIsAbsentAndPasswordIsPresent",
			user: user{
				"",
				"user1",
			},
		},
		{
			name:     "NameAndPasswordAreAbsent",
			user: user{
				"",
				"",
			},
		},
		{
			name:     "NameAndPasswordArePresent",
			user: user{
				"user1",
				"Password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := user{tt.user.name, tt.user.password}
			if len(u.name) > 0 {
				require.NotEmpty(t, u.name,
					"Name should be there")
			}
			if len(u.password) > 0 {
				require.NotEmpty(t, u.password,
					"password should be there")
			}
		})
	}
}
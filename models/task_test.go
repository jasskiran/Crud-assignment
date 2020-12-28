package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	type task struct {
		name string
		description string
	}
	tests := []struct {
		name   string
		task   task
	}{
		{
			name:     "NameIsPresentAndDescriptionIsAbsent",
			task: task{
				"user1",
				"",
			},
		},
		{
			name:     "NameIsAbsentAndDescriptionIsPresent",
			task: task{
				"",
				"user1",
			},
		},
		{
			name:     "NameAndDescriptionAreAbsent",
			task: task{
				"",
				"",
			},
		},
		{
			name:     "NameAndDescriptionArePresent",
			task: task{
				"user1",
				"Password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task{tt.task.name, tt.task.description}
			if len(tsk.name) > 0 {
				require.NotEmpty(t, tsk.name,
					"Name should be there")
			}
			if len(tsk.description) > 0 {
				require.NotEmpty(t, tsk.description,
					"password should be there")
			}
		})
	}
}

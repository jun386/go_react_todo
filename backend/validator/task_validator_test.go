package validator_test

import (
	"backend/model"
	"backend/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidate(t *testing.T) {
	validator := validator.NewTaskValidator()

	tests := []struct {
		name    string
		task    model.Task
		wantErr bool
	}{
		{
			name:    "Empty Title",
			task:    model.Task{Title: ""},
			wantErr: true,
		},
		{
			name:    "Title Too Long",
			task:    model.Task{Title: "This title is definitely more than ten characters"},
			wantErr: true,
		},
		{
			name:    "Valid Title",
			task:    model.Task{Title: "Valid"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.TaskValidate(tt.task)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error for test: %v", tt.name)
			} else {
				assert.NoError(t, err, "Expected no error for test: %v", tt.name)
			}
		})
	}
}

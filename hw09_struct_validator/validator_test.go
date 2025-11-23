package hw09structvalidator

import (
	"encoding/json"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		desc        string
		expectedErr bool
	}{
		{
			User{
				ID:     "550e8400-e29b-41d4-a716-446655440000",
				Name:   "Alexey",
				Age:    36,
				Email:  "alexey@example.com",
				Role:   "admin",
				Phones: []string{"89199999999"},
			},
			"valid user Alexey", false,
		},
		{
			User{
				ID:     "1234567",
				Name:   "Andrey",
				Age:    38,
				Email:  "andrey@example.com",
				Role:   "guest",
				Phones: []string{"9876"},
			},
			"invalid user Andrey", true,
		},
		{
			App{Version: "0.1.1"},
			"valid app version",
			false,
		},
		{
			App{Version: "123456789"},
			"invalid app version",
			true,
		},
		{
			Response{
				Code: 404,
				Body: "Not Found",
			},
			"valid response 404",
			false,
		},
		{
			Response{
				Code: 418,
				Body: "I'm a teapot",
			},
			"invalid response code 418",
			true,
		},
		{
			Token{
				Header:    []byte{1, 1, 1},
				Payload:   []byte{2, 2, 2},
				Signature: []byte{3, 3, 3},
			},
			"valid token",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := Validate(tt.in)
			if (err != nil) != tt.expectedErr {
				t.Errorf("%s: expected error: %v, got: %v", tt.desc, tt.expectedErr, err)
			}
		})
	}
}

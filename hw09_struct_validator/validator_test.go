package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	InvalidValidationRules struct {
		ID     string   `validate:"len:short"`
		AgeMin int      `validate:"min:youngster"`
		AgeMax int      `validate:"max:adult"`
		Code   int      `validate:"in:one,two"`
		Phones []string `validate:"len:normal"`
		Email  string   `validate:"regexp:(@"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "not a struct",
			expectedErr: ErrUnsupportedType,
		},
		{
			in:          1234,
			expectedErr: ErrUnsupportedType,
		},
		{
			in: User{
				Name:   "all fields valid",
				ID:     "testID123456789101112131415161789012",
				Age:    45,
				Email:  "test@mail.test",
				Role:   "admin",
				Phones: []string{"89611111118", "89614111741"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{
				Name:   "all fields invalid",
				ID:     "invalid",
				Age:    450,
				Email:  "testmail.test",
				Role:   "unknown",
				Phones: []string{"8961", "821265"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrExactLen},
				{Field: "Age", Err: ErrLessOrEqual},
				{Field: "Email", Err: ErrMatchRegExp},
				{Field: "Role", Err: ErrNotInList},
				{Field: "Phones.0", Err: ErrExactLen},
				{Field: "Phones.1", Err: ErrExactLen},
			},
		},
		{
			in: User{
				Name:   "one phone is invalid age is less then needed",
				ID:     "testID123456789101112131415161789012",
				Age:    12,
				Email:  "test@mail.test",
				Role:   "admin",
				Phones: []string{"89611111118", "8961741"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrGreaterOrEqual},
				{Field: "Phones.1", Err: ErrExactLen},
			},
		},
		{
			in:          App{Version: "1.2.3"},
			expectedErr: nil,
		},
		{
			in: App{Version: "1.2"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrExactLen},
			},
		},
		{
			in: Token{
				Header:    []byte("test"),
				Payload:   []byte("without"),
				Signature: []byte("validation"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 401,
				Body: "Unauthorized",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrNotInList},
			},
		},
		{
			in: InvalidValidationRules{
				ID:     "testID123456789101112131415161789012",
				AgeMin: 12,
				AgeMax: 50,
				Code:   200,
				Phones: []string{"89611111118", "89614111741"},
				Email:  "test@mail.test",
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrInvalidRule},
				{Field: "AgeMin", Err: ErrInvalidRule},
				{Field: "AgeMax", Err: ErrInvalidRule},
				{Field: "Code", Err: ErrInvalidRule},
				{Field: "Phones", Err: ErrInvalidRule},
				{Field: "Email", Err: ErrInvalidRule},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)

			var wantValidateErrs ValidationErrors
			if errors.As(tt.expectedErr, &wantValidateErrs) {
				var gotValidateErrs ValidationErrors

				require.ErrorAs(t, err, &gotValidateErrs)
				require.Len(t, gotValidateErrs, len(wantValidateErrs))

				for j, gotE := range gotValidateErrs {
					wantE := wantValidateErrs[j]

					require.Equal(t, wantE.Field, gotE.Field)
					require.ErrorIs(t, gotE.Err, wantE.Err)
				}
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

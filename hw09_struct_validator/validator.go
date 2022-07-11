package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")

	ErrNotInList      = errors.New("value must exist in list")
	ErrExactLen       = errors.New("value must be exact length")
	ErrLessOrEqual    = errors.New("value must be less or equal")
	ErrGreaterOrEqual = errors.New("value must be greater or equal")
	ErrMatchRegExp    = errors.New("value must match regular expression")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder

	for _, err := range v {
		builder.WriteString(err.Field)
		builder.WriteString(": ")
		builder.WriteString(err.Err.Error())
		builder.WriteString("\n")
	}

	return builder.String()
}

func Validate(v interface{}) error { //nolint:gocognit
	var validationErr ValidationErrors

	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	filedCount := refValue.NumField()
	for i := 0; i < filedCount; i++ {
		fieldType := refValue.Type().Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldValue := refValue.Field(i)
		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int:
			if err := parseIntFieldValidationRules(tag).validate(fieldValue.Int()); err != nil {
				validationErr = append(validationErr, ValidationError{Field: fieldType.Name, Err: err})
			}
		case reflect.String:
			if err := parseStringFieldValidationRules(tag).validate(fieldValue.String()); err != nil {
				validationErr = append(validationErr, ValidationError{Field: fieldType.Name, Err: err})
			}
		case reflect.Slice:
			if fieldValue.Len() == 0 {
				continue
			}

			switch fieldValue.Index(0).Kind() { //nolint:exhaustive
			case reflect.Int:
				validator := parseIntFieldValidationRules(tag)

				for j := 0; j < fieldValue.Len(); j++ {
					if err := validator.validate(fieldValue.Index(j).Int()); err != nil {
						validationErr = append(validationErr, ValidationError{
							Err:   err,
							Field: fmt.Sprintf("%s.%d", fieldType.Name, j),
						})
					}
				}
			case reflect.String:
				validator := parseStringFieldValidationRules(tag)

				for j := 0; j < fieldValue.Len(); j++ {
					if err := validator.validate(fieldValue.Index(j).String()); err != nil {
						validationErr = append(validationErr, ValidationError{
							Err:   err,
							Field: fmt.Sprintf("%s.%d", fieldType.Name, j),
						})
					}
				}
			}
		}
	}

	if len(validationErr) != 0 {
		return validationErr
	}

	return nil
}

type (
	intFieldValidator struct {
		min, max int64
		in       []int64
	}

	stringFieldValidator struct {
		len int
		in  []string
		re  *regexp.Regexp
	}
)

func (v intFieldValidator) validate(item int64) error {
	if v.min != 0 && item < v.min {
		return fmt.Errorf("%w %d", ErrGreaterOrEqual, v.min)
	}

	if v.max != 0 && item > v.max {
		return fmt.Errorf("%w %d", ErrLessOrEqual, v.max)
	}

	if len(v.in) > 0 {
		var found bool
		for i := range v.in {
			if v.in[i] == item {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%w %q", ErrNotInList, v.in)
		}
	}

	return nil
}

func (v stringFieldValidator) validate(item string) error {
	if v.len != 0 && utf8.RuneCountInString(item) != v.len {
		return fmt.Errorf("%w %d", ErrExactLen, v.len)
	}

	if v.re != nil && v.re.FindString(item) != item {
		return fmt.Errorf("%w %s", ErrMatchRegExp, v.re.String())
	}

	if len(v.in) > 0 {
		var found bool
		for i := range v.in {
			if v.in[i] == item {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%w %q", ErrNotInList, v.in)
		}
	}

	return nil
}

func parseIntFieldValidationRules(tag string) intFieldValidator {
	var validator intFieldValidator

	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "min":
			validator.min, _ = strconv.ParseInt(parts[1], 10, 64)
		case "max":
			validator.max, _ = strconv.ParseInt(parts[1], 10, 64)
		case "in":
			numbers := strings.Split(parts[1], ",")
			validator.in = make([]int64, 0, len(numbers))

			for _, n := range numbers {
				if v, err := strconv.ParseInt(n, 10, 64); err == nil {
					validator.in = append(validator.in, v)
				}
			}
		}
	}

	return validator
}

func parseStringFieldValidationRules(tag string) stringFieldValidator {
	var validator stringFieldValidator

	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "len":
			validator.len, _ = strconv.Atoi(parts[1])
		case "regexp":
			validator.re, _ = regexp.Compile(parts[1])
		case "in":
			validator.in = strings.Split(parts[1], ",")
		}
	}

	return validator
}

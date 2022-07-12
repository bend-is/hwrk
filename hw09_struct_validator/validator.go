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
	ErrInvalidRule     = errors.New("invalid validation rule")

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
	var validationErrs ValidationErrors

	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	filedCount := refValue.NumField()
	for i := 0; i < filedCount; i++ {
		fieldValue := refValue.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}

		fieldType := refValue.Type().Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int:
			validator, err := parseIntFieldValidationRules(tag)
			if err != nil {
				validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
				continue
			}

			if err := validator.validate(fieldValue.Int()); err != nil {
				validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
			}
		case reflect.String:
			validator, err := parseStringFieldValidationRules(tag)
			if err != nil {
				validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
				continue
			}

			if err := validator.validate(fieldValue.String()); err != nil {
				validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
			}
		case reflect.Slice:
			switch val := fieldValue.Interface().(type) {
			case []int:
				validator, err := parseIntFieldValidationRules(tag)
				if err != nil {
					validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
					continue
				}

				if err := validator.validateSlice(fieldType.Name, val); err != nil {
					validationErrs = append(validationErrs, err...)
				}
			case []string:
				validator, err := parseStringFieldValidationRules(tag)
				if err != nil {
					validationErrs = append(validationErrs, ValidationError{Field: fieldType.Name, Err: err})
					continue
				}

				if err := validator.validateSlice(fieldType.Name, val); err != nil {
					validationErrs = append(validationErrs, err...)
				}
			}
		}
	}

	if len(validationErrs) != 0 {
		return validationErrs
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

func (v intFieldValidator) validateSlice(fieldName string, items []int) ValidationErrors {
	if len(items) == 0 {
		return nil
	}

	var errs ValidationErrors

	for i := range items {
		if err := v.validate(int64(items[i])); err != nil {
			errs = append(errs, ValidationError{
				Err:   err,
				Field: fmt.Sprintf("%s.%d", fieldName, i),
			})
		}
	}

	return errs
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

func (v stringFieldValidator) validateSlice(fieldName string, items []string) ValidationErrors {
	if len(items) == 0 {
		return nil
	}

	var errs ValidationErrors

	for i := range items {
		if err := v.validate(items[i]); err != nil {
			errs = append(errs, ValidationError{
				Err:   err,
				Field: fmt.Sprintf("%s.%d", fieldName, i),
			})
		}
	}

	return errs
}

func parseIntFieldValidationRules(tag string) (intFieldValidator, error) {
	var validator intFieldValidator

	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "min":
			v, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return validator, fmt.Errorf("%w: min must be a valid integer", ErrInvalidRule)
			}

			validator.min = v
		case "max":
			v, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return validator, fmt.Errorf("%w: max must be a valid integer", ErrInvalidRule)
			}

			validator.max = v
		case "in":
			numbers := strings.Split(parts[1], ",")
			validator.in = make([]int64, 0, len(numbers))

			for _, n := range numbers {
				v, err := strconv.ParseInt(n, 10, 64)
				if err != nil {
					return validator, fmt.Errorf("%w: in must contain valid intergers", ErrInvalidRule)
				}
				validator.in = append(validator.in, v)
			}
		}
	}

	return validator, nil
}

func parseStringFieldValidationRules(tag string) (stringFieldValidator, error) {
	var validator stringFieldValidator

	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "len":
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return validator, fmt.Errorf("%w: len must be a valid integer", ErrInvalidRule)
			}

			validator.len = v
		case "regexp":
			re, err := regexp.Compile(parts[1])
			if err != nil {
				return validator, fmt.Errorf("%w: regexp must contain valid regular expression", ErrInvalidRule)
			}

			validator.re = re
		case "in":
			validator.in = strings.Split(parts[1], ",")
		}
	}

	return validator, nil
}

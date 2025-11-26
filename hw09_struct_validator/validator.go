package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	result := make([]string, 0, len(v))
	for _, val := range v {
		result = append(result, fmt.Sprintf("%s: %s", val.Field, val.Err))
	}
	return strings.Join(result, "; ")
}

func Validate(v interface{}) error {
	refValue := reflect.ValueOf(v)
	refType := reflect.TypeOf(v)

	if refType.Kind() != reflect.Struct {
		return errors.New("expected a struct")
	}

	var vErrors ValidationErrors

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := refValue.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			name, arg, _ := strings.Cut(rule, ":")
			switch name {
			case "len":
				want, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}
				vErrors = append(vErrors, checkLen(field, value, want)...)
			case "min":
				minVal, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}
				vErrors = append(vErrors, checkMin(field, value, minVal)...)
			case "max":
				maxVal, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}
				vErrors = append(vErrors, checkMax(field, value, maxVal)...)
			case "in":
				options := strings.Split(arg, ",")
				vErrors = append(vErrors, checkIn(field, value, options)...)
			case "regexp":
				vErrorsElem, err := checkRegexp(field, value, arg)
				if err != nil {
					return err
				}
				vErrors = append(vErrors, vErrorsElem...)
			default:
				vErrors = append(vErrors, ValidationError{
					Field: field.Name,
					Err:   fmt.Errorf("unknown rule %s", name),
				})
			}
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func checkLen(field reflect.StructField, value reflect.Value, want int) ValidationErrors {
	var vErrors ValidationErrors

	switch value.Kind() {
	case reflect.String:
		if len(value.String()) != want {
			vErrors = append(vErrors, ValidationError{
				field.Name,
				fmt.Errorf("length must be %d", want),
			})
		}
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.String {
			for i := 0; i < value.Len(); i++ {
				s := value.Index(i).String()
				if len(s) != want {
					vErrors = append(vErrors, ValidationError{
						field.Name,
						fmt.Errorf("element %d length must be %d", i, want),
					})
				}
			}
		} else if value.Len() != want {
			vErrors = append(vErrors, ValidationError{
				field.Name,
				fmt.Errorf("length must be %d", want),
			})
		}
	case reflect.Array:
		if value.Len() != want {
			vErrors = append(vErrors, ValidationError{
				field.Name,
				fmt.Errorf("length must be %d", want),
			})
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Struct, reflect.UnsafePointer:
		vErrors = append(vErrors, ValidationError{field.Name, fmt.Errorf("len rule not supported for type %s", value.Kind())})
	}

	return vErrors
}

func checkMin(field reflect.StructField, value reflect.Value, minVal int) ValidationErrors {
	var vErrors ValidationErrors
	if isInt(value.Kind()) {
		if value.Int() < int64(minVal) {
			vErrors = append(vErrors, ValidationError{
				field.Name,
				fmt.Errorf("value must be >= %d", minVal),
			})
		}
	} else {
		vErrors = append(vErrors, ValidationError{
			field.Name,
			fmt.Errorf("min rule not supported for type %s", value.Kind()),
		})
	}
	return vErrors
}

func checkMax(field reflect.StructField, value reflect.Value, maxVal int) ValidationErrors {
	var vErrors ValidationErrors
	if isInt(value.Kind()) {
		if value.Int() > int64(maxVal) {
			vErrors = append(vErrors, ValidationError{
				field.Name,
				fmt.Errorf("value must be <= %d", maxVal),
			})
		}
	} else {
		vErrors = append(vErrors, ValidationError{
			field.Name,
			fmt.Errorf("max rule not supported for type %s", value.Kind()),
		})
	}
	return vErrors
}

func checkIn(field reflect.StructField, value reflect.Value, options []string) ValidationErrors {
	var vErrors ValidationErrors
	valStr := fmt.Sprintf("%v", value.Interface())
	foundFlag := false
	for _, option := range options {
		if valStr == option {
			foundFlag = true
			break
		}
	}
	if !foundFlag {
		vErrors = append(vErrors, ValidationError{
			field.Name,
			fmt.Errorf("value must be one of [%s]",
				strings.Join(options, ", ")),
		})
	}
	return vErrors
}

func checkRegexp(field reflect.StructField, value reflect.Value, pattern string) (ValidationErrors, error) {
	var vErrors ValidationErrors
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if !re.MatchString(fmt.Sprintf("%v", value.Interface())) {
		vErrors = append(vErrors, ValidationError{
			field.Name,
			fmt.Errorf("value does not match regexp %s", pattern),
		})
	}
	return vErrors, nil
}

func isInt(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

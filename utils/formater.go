package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"stock-tracker/constant"
	"strings"
	"time"
)

var regFormatHp *regexp.Regexp
var replacer *strings.Replacer

func init() {
	regFormatHp = regexp.MustCompile(`^(?:\+62|62|0)`)
	replacer = strings.NewReplacer(" ", "", "-", "")
}

func FormatPhoneTo62(phone string) string {
	cleaned := replacer.Replace(phone)
	formatted := regFormatHp.ReplaceAllString(cleaned, "62")
	return formatted
}

// toCamelCase converts PascalCase or UpperCamelCase to camelCase
func PascalcasetoCamelcase(str string) string {
	if str == "" {
		return str
	}

	// Handle the first character: lowercasing it
	str = strings.ToLower(string(str[0])) + str[1:]

	// Use regex to insert an underscore before consecutive uppercase letters followed by lowercase
	re := regexp.MustCompile("([A-Z])([A-Z]+)([a-z])")
	str = re.ReplaceAllStringFunc(str, func(s string) string {
		return string(s[0]) + strings.ToLower(s[1:len(s)-1]) + string(s[len(s)-1])
	})

	return str
}

// TrimWhitespace recursively trims whitespace from all string fields in a struct
func TrimWhitespace(v interface{}) {
	// Ensure the value is a pointer to a struct
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	timeType := reflect.TypeOf(time.Time{})

	// Iterate over all fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			// Trim string fields
			field.SetString(strings.TrimSpace(field.String()))

		case reflect.Ptr:
			// If it's a pointer, check if it points to a string
			if field.IsNil() {
				continue
			}

			// Skip *time.Time to avoid panic on unexported fields
			if field.Type().Elem() == timeType {
				continue
			}

			if field.Type().Elem().Kind() == reflect.String {
				// If it's a pointer to a string, trim its value
				trimmedStr := strings.TrimSpace(field.Elem().String())
				// Only update if the trimmed string has content
				if trimmedStr != "" {
					field.Elem().SetString(trimmedStr)
				}
			} else if field.Type().Elem().Kind() == reflect.Struct {
				// If it's a pointer to a struct, recursively call TrimWhitespace
				TrimWhitespace(field.Interface())
			}

		case reflect.Struct:
			// Skip time.Time struct
			if field.Type() == timeType {
				continue
			}

			// If it's a struct, recursively call TrimWhitespace
			TrimWhitespace(field.Addr().Interface())

		case reflect.Slice:
			// Handle slices of structs or pointers to structs
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)

				if elem.Kind() == reflect.Ptr {
					if elem.IsNil() {
						continue
					}
					// Skip *time.Time inside slice
					if elem.Type().Elem() == timeType {
						continue
					}
					TrimWhitespace(elem.Interface())

				} else if elem.Kind() == reflect.Struct {
					// Skip time.Time inside slice
					if elem.Type() == timeType {
						continue
					}
					TrimWhitespace(elem.Addr().Interface())
				}
			}
		}
	}
}

func DisplayDateLayout(date time.Time, layout string) string {
	return date.Format(layout)
}

func DisplayDate(date time.Time) string {
	return date.Format(constant.FormatDateLayout)
}

func DisplayDatetime(date time.Time) string {
	return date.Format(constant.FormatDatetimeLayout)
}

func DisplayTime(date time.Time) string {
	return date.Format(constant.FormatTimeLayout)
}

func DisplayBool(data bool, trueText string, falseText string) string {
	if data {
		return trueText
	}
	return falseText
}

func DisplayPhoneNumber(value string) string {
	// Remove all non-digit characters
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(value, "")

	// Match phone number format with optional '62' prefix
	re := regexp.MustCompile(`^(62|)?(\d{3})(\d{4})(\d{3,6})$`)
	match := re.FindStringSubmatch(cleaned)

	if match != nil {
		intlCode := ""
		if match[1] != "" {
			intlCode = "+62 "
		}
		return fmt.Sprintf("%s %s-%s-%s", intlCode, match[2], match[3], match[4])
	}

	return value
}

func DisplayNumber(value int64) string {
	return formatNumberWithSeparator(value, getThousandSeparator("id-ID"))
}

func DisplayMoney(value int64) string {
	return "Rp " + DisplayNumber(value)
}

func getThousandSeparator(locales string) string {
	if strings.HasPrefix(locales, "id") {
		return "."
	}
	return ","
}

func formatNumberWithSeparator(value int64, sep string) string {
	n := fmt.Sprintf("%d", value)
	if len(n) <= 3 {
		return n
	}

	// Insert the thousand separator
	var result strings.Builder
	remain := len(n) % 3
	if remain > 0 {
		result.WriteString(n[:remain])
		result.WriteString(sep)
	}
	for i := remain; i < len(n); i += 3 {
		if i > remain {
			result.WriteString(sep)
		}
		result.WriteString(n[i : i+3])
	}
	return result.String()
}

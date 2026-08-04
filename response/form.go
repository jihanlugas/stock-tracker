package response

import (
	"errors"
	"fmt"
	"reflect"
	"stock-tracker/config"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func getFieldError(str ...string) interface{} {
	switch str[1] {
	case ErrorFormSelectionInvalid:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormDataNotFound:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormLengthTooBig:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormLengthTooShort:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormInvalidEmail:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFieldRequired:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFieldNumeric:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFixedLength:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " ") + " " + str[3],
		}
	case ErrorFormAlreadyExists:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormFlexibleMsg:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormPhoto:
		return FieldError{
			Field: str[0],
			//Msg:   "Photo maks " + strconv.FormatInt(config.PhotoUploadMaxSizeByte/1000000, 10) + " mb and ext: jpg, jpeg, png",
			Msg: fmt.Sprintf("Photo max size %d mb and allowd ext: %s", config.PhotoUploadMaxSizeByte/1000000, strings.Join(config.PhotoUploadAllowedExtensions, ", ")),
		}
	case ErrorFormLowercase:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormUppercase:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	case ErrorFormNotMatch:
		return FieldError{
			Field: str[0],
			Msg:   strings.ReplaceAll(str[2], "_", " "),
		}
	default:
		return FieldError{
			Field: str[0],
			Msg:   "Error Message",
		}
	}
}

func getListError3(err error, obj interface{}) map[string]string {
	errors := make(map[string]string)

	// Use reflection to get the type of the object
	objType := reflect.TypeOf(obj)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var jsonTag string
			var fieldName string

			// Reflect on the struct field to get the JSON tag
			structField, found := objType.FieldByName(e.StructField())

			if found {
				jsonTag = structField.Tag.Get("json")
			}

			// Fallback to the field name if JSON tag is not set
			if jsonTag == "" {
				jsonTag = e.Field()
			}

			// Check if the error is on a slice element by looking for an index in the error path
			if e.Kind() == reflect.Slice || e.Kind() == reflect.Array {
				// Use StructNamespace to include the index in the error path
				fieldName = e.StructNamespace()
			} else {
				fieldName = jsonTag
			}

			// Generate a user-friendly error message based on the validation tag
			switch e.Tag() {
			case "required":
				errors[fieldName] = fmt.Sprintf("%s is required", jsonTag)
			case "email":
				errors[fieldName] = fmt.Sprintf("%s must be a valid email", jsonTag)
			case "gte":
				errors[fieldName] = fmt.Sprintf("%s must be greater than or equal to %s", jsonTag, e.Param())
			case "lte":
				errors[fieldName] = fmt.Sprintf("%s must be less than or equal to %s", jsonTag, e.Param())
			case "min":
				errors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", jsonTag, e.Param())
			case "len":
				errors[fieldName] = fmt.Sprintf("%s must be exactly %s characters long", jsonTag, e.Param())
			case "lowercase":
				errors[fieldName] = fmt.Sprintf("%s must be in lowercase", jsonTag)
			case "numeric":
				errors[fieldName] = fmt.Sprintf("%s must be numeric", jsonTag)
			// Add cases for additional tags as needed
			default:
				errors[fieldName] = fmt.Sprintf("%s is invalid", jsonTag)
			}
		}
	}
	return errors
}

func getFieldName(field string) string {
	arr := strings.Split(field, ".")
	field = strings.Join(arr[1:], ".")

	return field
}

//func getListError2(err error) map[string]string {
//	errors := make(map[string]string)
//
//	if validationErrors, ok := err.(validator.ValidationErrors); ok {
//		for _, e := range validationErrors {
//			//fmt.Println("==================")
//			//fmt.Println(e)
//			//fmt.Println("e.Kind() ", e.Kind())
//			//fmt.Println("e.Value() ", e.Value())
//			//fmt.Println("e.Namespace() ", e.Namespace())
//			//fmt.Println("e.Field() ", e.Field())
//			//fmt.Println("e.Tag() ", e.Tag())
//			//fmt.Println("e.ActualTag() ", e.ActualTag())
//			//fmt.Println("e.StructField() ", e.StructField())
//			//fmt.Println("e.StructNamespace() ", e.StructNamespace())
//			//fmt.Println("===========================================================================")
//
//			fieldName := getFieldName(e.Namespace()) // Includes nested struct names
//
//			fmt.Println("fieldName: ", fieldName)
//
//			//// Handle `dive` tag by constructing a more detailed error message for array items
//			//if e.Kind() == reflect.Slice || e.Kind() == reflect.Map {
//			//	fieldName = fmt.Sprintf("%s[%d]", e.StructNamespace(), e.Index())
//			//}
//
//			// Generate a user-friendly error message based on the validation tag
//			switch e.Tag() {
//			case "required":
//				errors[fieldName] = fmt.Sprintf("%s is required", e.Field())
//			case "email":
//				errors[fieldName] = fmt.Sprintf("%s must be a valid email", e.Field())
//			case "gte":
//				errors[fieldName] = fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
//			case "lte":
//				errors[fieldName] = fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
//			case "min":
//				errors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
//			case "max":
//				errors[fieldName] = fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
//			case "len":
//				errors[fieldName] = fmt.Sprintf("%s must be exactly %s characters long", e.Field(), e.Param())
//			case "lowercase":
//				errors[fieldName] = fmt.Sprintf("%s must be in lowercase", e.Field())
//			case "uppercase":
//				errors[fieldName] = fmt.Sprintf("%s must be in uppercase", e.Field())
//			case "alphanum":
//				errors[fieldName] = fmt.Sprintf("%s must be alphanumeric", e.Field())
//			case "numeric":
//				errors[fieldName] = fmt.Sprintf("%s must be numeric", e.Field())
//			case "url":
//				errors[fieldName] = fmt.Sprintf("%s must be a valid URL", e.Field())
//			case "uuid":
//				errors[fieldName] = fmt.Sprintf("%s must be a valid UUID", e.Field())
//			case "oneof":
//				errors[fieldName] = fmt.Sprintf("%s must be one of the following: %s", e.Field(), e.Param())
//			case "unique":
//				errors[fieldName] = fmt.Sprintf("%s must be unique", e.Field())
//			// Add cases for additional tags as needed
//			default:
//				errors[fieldName] = fmt.Sprintf("%s is invalid", e.Field())
//			}
//		}
//	}
//	return errors
//}

func getListError(err error) Payload {
	listError := Payload{}
	var fieldsError validator.ValidationErrors
	errors.As(err, &fieldsError)

	for _, fieldError := range fieldsError {
		fieldName := getFieldName(fieldError.Namespace())

		switch fieldError.ActualTag() {
		case "required":
			listError[fieldName] = getFieldError(fieldName, ErrorFormFieldRequired, ErrorFormFieldRequired)
		case "notexists":
			listError[fieldName] = getFieldError(fieldName, ErrorFormAlreadyExists, ErrorFormAlreadyExists)
		case "existsdata":
			listError[fieldName] = getFieldError(fieldName, ErrorFormDataNotFound, ErrorFormDataNotFound)
		case "no_hp":
			listError[fieldName] = getFieldError(fieldName, ErrorFormFlexibleMsg, "Format nomor HP tidak benar")
		case "oneof", "exists", "weekday":
			listError[fieldName] = getFieldError(fieldName, ErrorFormSelectionInvalid, ErrorFormSelectionInvalid)
		case "lte", "max":
			listError[fieldName] = getFieldError(fieldName, ErrorFormLengthTooBig, ErrorFormLengthTooBig, fieldError.Param())
		case "gte", "min":
			listError[fieldName] = getFieldError(fieldName, ErrorFormLengthTooShort, ErrorFormLengthTooShort, fieldError.Param())
		case "email":
			listError[fieldName] = getFieldError(fieldName, ErrorFormInvalidEmail, ErrorFormInvalidEmail)
		case "numeric":
			listError[fieldName] = getFieldError(fieldName, ErrorFormFieldNumeric, ErrorFormFieldNumeric)
		case "len":
			listError[fieldName] = getFieldError(fieldName, ErrorFormFixedLength, ErrorFormFixedLength, fieldError.Param())
		case "passwdComplex":
			listError[fieldName] = getFieldError(fieldName, ErrorFormFlexibleMsg, "Password harus 1 lowercase 1 uppercase 1 numberic")
		case "base64PhotoCheck":
			listError[fieldName] = getFieldError(fieldName, ErrorFormPhoto)
		case "lowercase":
			listError[fieldName] = getFieldError(fieldName, ErrorFormLowercase, ErrorFormLowercase)
		case "uppercase":
			listError[fieldName] = getFieldError(fieldName, ErrorFormUppercase, ErrorFormUppercase)
		case "eqfield":
			listError[fieldName] = getFieldError(fieldName, ErrorFormNotMatch, ErrorFormNotMatch, fieldError.Param())
		default:
			listError[fieldName] = getFieldError(fieldName, ErrorFormInvalidValue, ErrorFormInvalidValue)
		}
	}

	return listError
}

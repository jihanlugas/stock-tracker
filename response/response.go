package response

import (
	"encoding/json"
	"stock-tracker/config"

	"github.com/labstack/echo/v4"
)

const (
	// Error General
	ErrorInternalServer = "internal server error"
	ErrorUnauthorized   = "unauthorized"

	// Error Middleware
	ErrorMiddlewareUserNotFound = "Token Expired!"
	ErrorMiddlewarePassVersion  = "Token Expired~"

	// Error Handler
	ErrorHandlerFailedValidation = "error validation"
	ErrorHandlerIDOR             = "unable to perform this action"
	ErrorHandlerGetUserInfo      = "bad request"
	ErrorHandlerGetParam         = "invalid param"
	ErrorHandlerBind             = "invalid request"

	// Error Usecase
	ErrorUsecaseBadRequest = "bad request"

	// Error Form
	ErrorFormSelectionInvalid = "error_form_selection_invalid"
	ErrorFormDataNotFound     = "error_form_data_not_found"
	ErrorFormLengthTooBig     = "error_form_length_too_big"
	ErrorFormLengthTooShort   = "error_form_length_too_short"
	ErrorFormInvalidEmail     = "error_form_invalid_email"
	ErrorFormFieldRequired    = "error_form_field_required"
	ErrorFormFieldNumeric     = "error_form_field_numeric"
	ErrorFormFixedLength      = "error_form_fix_length"
	ErrorFormAlreadyExists    = "error_form_already_exists"
	ErrorFormFlexibleMsg      = "error_form_flexible_msg"
	ErrorFormPhoto            = "error_form_photo"
	ErrorFormLowercase        = "error_form_lowercase"
	ErrorFormUppercase        = "error_form_uppercase"
	ErrorFormNotMatch         = "error_form_not_match"
	ErrorFormInvalidValue     = "error_form_invalid_value"
)

// SuccessResponse type for Success Response
type Response struct {
	Code    int         `json:"code"`
	Err     string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Payload interface{} `json:"payload,omitempty" swaggertype:"object"`
}

type Payload map[string]interface{}

func (e *Response) Error() string {
	return e.Message
}

func Success(code int, msg string, payload interface{}) *Response {
	return &Response{
		Code:    code,
		Message: msg,
		Status:  true,
		Payload: payload,
	}
}

func Error(code int, msg string, err error, payload interface{}) *Response {
	if config.Debug {
		return &Response{
			Code:    code,
			Err:     err.Error(),
			Message: msg,
			Status:  false,
			Payload: payload,
		}
	} else {
		return &Response{
			Code:    code,
			Message: msg,
			Status:  false,
			Payload: payload,
		}
	}
}

func ErrorForce(code int, msg string) *Response {
	payload := Payload{
		"forceLogout": true,
	}
	return &Response{
		Code:    code,
		Message: msg,
		Status:  false,
		Payload: payload,
	}
}

func (r *Response) SendJSON(c echo.Context) error {
	js, err := json.Marshal(r)
	if err != nil {
		panic(err)
	} else {
		//c.Set(constant.Response, js)
		return c.Blob(r.Code, echo.MIMEApplicationJSON, js)
	}
}

func ValidationError(err error) *Payload {
	return &Payload{
		"listError": getListError(err),
	}
}

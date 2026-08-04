package validator

import (
	"encoding/base64"
	"errors"
	"reflect"
	"regexp"
	"slices"
	"stock-tracker/app/user"
	"stock-tracker/config"
	"stock-tracker/db"
	"stock-tracker/utils"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	Validate        *CustomValidator
	regxPhoneNumber *regexp.Regexp
)

type CustomValidator struct {
	validator *validator.Validate
}

func init() {
	Validate = NewValidator()
	regxPhoneNumber = regexp.MustCompile(`((^\+?628\d{8,14}$)|(^0?8\d{8,14}$)){1}`)
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// ValidateVar for validate field against tag. Expl: ValidateVar("abc@gmail.com", "required,email")
func (v *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

func NewValidator() *CustomValidator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = validate.RegisterValidation("notexists", notExistsOnDbTable)
	_ = validate.RegisterValidation("existsdata", existsDataOnDbTable)
	_ = validate.RegisterValidation("no_hp", validPhoneNumber)
	_ = validate.RegisterValidation("passwdComplex", checkPasswordComplexity)
	_ = validate.RegisterValidation("base64PhotoCheck", base64PhotoCheck, true)

	return &CustomValidator{
		validator: validate,
	}
}

func notExistsOnDbTable(fl validator.FieldLevel) bool {
	var err error
	params := strings.Fields(fl.Param())

	val := strings.TrimSpace(fl.Field().String())
	if val == "" {
		return true
	}

	userRepo := user.NewRepository()

	conn := db.GetPostgresConnection()

	switch params[0] {
	case "username":
		_, err = userRepo.GetByUsername(conn, val)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false

	case "email":
		_, err = userRepo.GetByEmail(conn, val)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false

	case "no_hp":
		_, err = userRepo.GetByPhoneNumber(conn, utils.FormatPhoneTo62(val))
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		}
		return false

	}

	return false
}

func existsDataOnDbTable(fl validator.FieldLevel) bool {
	var err error
	params := strings.Fields(fl.Param())

	if fl.Field().String() == "" {
		return true
	}

	// companyRepo := company.NewRepository()
	userRepo := user.NewRepository()

	conn := db.GetPostgresConnection()

	switch params[0] {
	case "user_id":
		ID := fl.Field().String()
		_, err = userRepo.GetTableById(conn, ID)
		if err != nil {
			return false
		}
		return true
		// case "company_id":
		// 	ID := fl.Field().String()
		// 	_, err = companyRepo.GetTableById(conn, ID)
		// 	if err != nil {
		// 		return false
		// 	}
		// 	return true
	}
	return false
}

func IsSameDate(date1, date2 *time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func checkPasswordComplexity(fl validator.FieldLevel) bool {
	passwd := fl.Field().String()

	var capitalFlag, lowerCaseFlag, numberFlag bool
	for _, c := range passwd {
		if unicode.IsUpper(c) {
			capitalFlag = true
		} else if unicode.IsLower(c) {
			lowerCaseFlag = true
		} else if unicode.IsDigit(c) {
			numberFlag = true
		}
		if capitalFlag && lowerCaseFlag && numberFlag {
			return true
		}
	}
	return false
}

func validPhoneNumber(fl validator.FieldLevel) bool {
	return regxPhoneNumber.MatchString(fl.Field().String())
}

func base64PhotoCheck(fl validator.FieldLevel) bool {
	base64String := fl.Field().String()

	if base64String == "" {
		return true
	}

	// Check if Base64 data contains image type metadata
	if !strings.HasPrefix(base64String, "data:image/") {
		return false
	}

	// Extract image type (e.g., png, jpg)
	imageType := strings.TrimPrefix(base64String[:strings.Index(base64String, ";")], "data:image/")

	//// Validate the image extension
	//if !regExt.MatchString(imageType) {
	//	return false
	//}

	// Validate the image extension
	if !slices.Contains(config.PhotoUploadAllowedExtensions, imageType) {
		return false
	}

	// Remove metadata prefix (like "data:image/png;base64,") and validate base64 data
	base64Data := base64String[strings.Index(base64String, ",")+1:]

	// Check if the string is valid base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return false
	}

	// Check image size (limit is 10 MB)
	if int64(len(imageData)) > config.PhotoUploadMaxSizeByte {
		return false
	}

	return true
}

//func photoCheck(fl validator.FieldLevel) bool {
//	if len(params) == 0 {
//		return true
//	}
//	parentVal := fl.Parent()
//	if parentVal.Kind() == reflect.Ptr {
//		parentVal = reflect.Indirect(parentVal)
//	}
//	// field photo harus dengan tipe data: *multipart.FileHeader ( pointer )
//	photoVal := parentVal.FieldByName(params[0])
//	if photoVal.Kind() != reflect.Ptr {
//		return false
//	}
//	if photoVal.IsZero() {
//		return true
//	}
//	if f, ok := photoVal.Interface().(*multipart.FileHeader); !ok {
//		return false
//	} else {
//		if !regExt.MatchString(filepath.Ext(f.Filename)) {
//			return false
//		}
//		if f.Size > config.PhotoUploadMaxSizeByte {
//			return false
//		}
//		return true
//	}
//}

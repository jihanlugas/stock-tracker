package base

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// ApplyGlobalSearch menerapkan global search berdasarkan tag `search`
func ApplyGlobalSearch(db *gorm.DB, keyword string, filter any) *gorm.DB {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return db
	}

	fields := GetSearchFields(filter)
	if len(fields) == 0 {
		return db
	}

	keyword = "%" + keyword + "%"

	conditions := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))

	for _, field := range fields {
		conditions = append(conditions, field+" ILIKE ?")
		args = append(args, keyword)
	}

	return db.Where(
		"("+strings.Join(conditions, " OR ")+")",
		args...,
	)
}

// GetSearchFields mengambil semua field yang memiliki tag `search`
func GetSearchFields(filter any) []string {
	t := reflect.TypeOf(filter)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if search := field.Tag.Get("search"); search != "" {
			fields = append(fields, search)
		}
	}

	return fields
}

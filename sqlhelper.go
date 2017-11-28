package sqlhelper

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var (
	// tables holds the parsed fields for each table given
	tables []table
	// IgnoreFields are the fields to skip when parsing the struct. Overwrite
	// if you want
	IgnoreFields = []string{"ID", "Created"}
)

type table struct {
	Name      string
	Fields    []string
	sqlFields []string
}

// Reset the storage array
func Reset() {
	tables = []table{}
}

func (t *table) generateSQLFields() {
	t.sqlFields = []string{}
	for _, f := range t.Fields {
		for _, r := range IgnoreFields {
			if f == r {
				goto SKIP
			}
		}
		t.sqlFields = append(t.sqlFields, fmt.Sprintf("%s.%s", t.Name, f))

	SKIP:
	}

	// Sort
	sort.Strings(t.sqlFields)
}

// StructFields returns the structs fields in a SQL friendly way
func StructFields(tableName string, in interface{}) {

	t := table{
		Name:   tableName,
		Fields: []string{},
	}

	val := reflect.ValueOf(in).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag.Get("db") != "-" {
			t.Fields = append(t.Fields, typeField.Name)
		}
	}

	t.generateSQLFields()

	tables = append(tables, t)
}

// InsertFields returns a tables fields except for ID and Created
func InsertFields(tableName string) (string, string) {
	for _, t := range tables {
		if t.Name == tableName {
			return strings.Join(t.sqlFields, ", "), strings.TrimRight(strings.Repeat("?, ", len(t.sqlFields)), ", ")
		}
	}

	return "", ""
}

// UpdateFields returns a tables fields ready for updating
func UpdateFields(tableName string) string {
	for _, t := range tables {
		if t.Name == tableName {
			return strings.Join(t.sqlFields, " = ?, ") + " = ?"
		}
	}

	return ""
}

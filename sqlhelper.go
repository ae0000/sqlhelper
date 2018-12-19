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
	IgnoreFields = []string{"ID", "Created"}
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

// SelectFields returns a tables fields including those excluded in
// update/insert
func SelectFields(tableName string) string {
	for _, t := range tables {
		if t.Name == tableName {
			sql := ""
			for _, f := range t.Fields {
				sql += fmt.Sprintf("%s.%s, ", tableName, f)
			}
			return strings.TrimRight(sql, ", ")
		}
	}

	return ""
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

// ShowInsert shows what a insert would most likely look like (for quickly
// getting inserts setup)
func ShowInsert(tableName string, instanceName string) {

	for _, t := range tables {
		if t.Name == tableName {

			fmt.Printf("// %sInsert inserts a %s\n", strings.Title(instanceName), instanceName)
			fmt.Printf("func %sInsert(%s *%s) error {\n", strings.Title(instanceName), instanceName, strings.Title(instanceName))
			fmt.Printf("fields, params := sqlhelper.InsertFields(\"%s\")\n", tableName)
			fmt.Printf("sql := fmt.Sprintf(\"INSERT INTO %s (%%s) VALUES (%%s)\", fields, params)\n", tableName)

			fmt.Println("r, err := d.Exec(sql,")
			sort.Strings(t.Fields)
			for _, f := range t.Fields {
				for _, r := range IgnoreFields {
					if f == r {
						goto SKIP
					}
				}
				fmt.Printf("  %s.%s,\n", instanceName, f)

			SKIP:
			}

			fmt.Println(")")
			fmt.Println("")
			fmt.Println("if err == nil {")
			fmt.Printf("		%s.ID, err = r.LastInsertId()\n", instanceName)
			fmt.Println("}")
			fmt.Println("")
			fmt.Println("return err")
			fmt.Println("}")
		}
	}

}

// ShowUpdate shows what a update would most likely look like (for quickly
// getting setup)
func ShowUpdate(tableName string, instanceName string) {

	for _, t := range tables {
		if t.Name == tableName {

			fmt.Printf("// %sUpdate updates a %s\n", strings.Title(instanceName), instanceName)
			fmt.Printf("func %sUpdate(%s *%s) error {\n", strings.Title(instanceName), instanceName, strings.Title(instanceName))
			fmt.Printf("fields := sqlhelper.UpdateFields(\"%s\")\n", tableName)
			fmt.Printf("sql := fmt.Sprintf(\"UPDATE %s SET %%s WHERE ID = ? LIMIT 1\", fields)\n", tableName)

			fmt.Println("r, err := d.Exec(sql,")
			sort.Strings(t.Fields)
			for _, f := range t.Fields {
				for _, r := range IgnoreFields {
					if f == r {
						goto SKIP
					}
				}
				fmt.Printf("  %s.%s,\n", instanceName, f)

			SKIP:
			}

			fmt.Printf("  %s.ID,\n", instanceName)
			fmt.Println(")")
			fmt.Println("")
			fmt.Println("return err")
			fmt.Println("}")
		}
	}

}

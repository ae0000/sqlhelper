SQLhelper for Go
================

This is just a little helper that saves mucking around with inserts and updates when using non-orm SQL drivers.

Usage:
------

```
package db

import (
	"fmt"
	"github.com/ae0000/sqlhelper
)

type user struct {
	ID      int64
	Name    string
	Email   string
	NotMe   string `db:"-"`
	Created time.Time
}

func init() {
	sqlhelper.StructFields("Users", &User{})
}

// UserInsert inserts a user
func UserInsert(user *user) error {
	fields, params := sqlhelper.InsertFields("Users")

	// We get: "INSERT INTO Users (Users.Email, Users.Name) VALUES (?, ?)"
	sql := fmt.Sprintf("INSERT INTO Users (%s) VALUES (%s)", fields, params)

	// Note that we don't have to include ID, NotMe or Created as they are 
	// excluded (ID, Created are excluded by default, NotMe via the tag db:"-")
	r, err := d.Exec(sql,
	user.Email,
	user.Name)

	if err == nil {
		user.ID, err = r.LastInsertId()
	}

	return err
}

```


You can edit the fields you want to exclude via:
```
sqlhelper.IgnoreFields = []string{"LastUpdated"}
```
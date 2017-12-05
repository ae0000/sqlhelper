package sqlhelper

import (
	"fmt"
	"testing"
	"time"
)

type user struct {
	ID      int64
	Name    string
	Email   string
	NotMe   string `db:"-"`
	Created time.Time
}

func TestInsert(t *testing.T) {
	// Generate the fields (do this on init normally)
	Reset()
	StructFields("Users", &user{})

	fields, params := InsertFields("Users")

	sql := fmt.Sprintf("INSERT INTO Users (%s) VALUES (%s)", fields, params)

	if sql != "INSERT INTO Users (Users.Email, Users.Name) VALUES (?, ?)" {
		t.Error("incorrect insertFields, ", sql)
	}
}

func TestInsertEditingIgnoreFields(t *testing.T) {
	Reset()
	IgnoreFields = []string{}

	// Generate the fields (do this on init normally)
	StructFields("Users", &user{})

	fields, params := InsertFields("Users")

	sql := fmt.Sprintf("INSERT INTO Users (%s) VALUES (%s)", fields, params)

	if sql != "INSERT INTO Users (Users.Created, Users.Email, Users.ID, Users.Name) VALUES (?, ?, ?, ?)" {
		t.Error("incorrect insertFields, ", sql)
	}
}
func TestUpdate(t *testing.T) {
	// Generate the fields (do this on init normally)
	Reset()
	StructFields("Users", &user{})

	fields := UpdateFields("Users")

	sql := fmt.Sprintf("UPDATE Users SET %s", fields)

	if sql != "UPDATE Users SET Users.Created = ?, Users.Email = ?, Users.ID = ?, Users.Name = ?" {
		t.Error("incorrect updateFields, ", sql)
	}
}

func TestNoTable(t *testing.T) {
	Reset()
	StructFields("Users", &user{})

	fields, params := InsertFields("NotUsers")
	if fields != "" {
		t.Error("was expecting fields to be empty, got:", fields)
	}
	if params != "" {
		t.Error("was expecting params to be empty, got:", params)
	}

	fields = UpdateFields("NotUsers")
	if fields != "" {
		t.Error("was expecting fields to be empty, got:", fields)
	}
}

func TestShowInsert(t *testing.T) {
	Reset()
	StructFields("Users", &user{})

	ShowInsert("Users", "user")
}

func TestShowUpdate(t *testing.T) {
	Reset()
	StructFields("Users", &user{})

	ShowUpdate("Users", "user")
	t.Error("!")
}

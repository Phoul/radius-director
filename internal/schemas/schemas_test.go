package schemas

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadMySQLSchema(t *testing.T) {
	loader := NewLoader(fstest.MapFS{
		"freeradius/3.2.10/mysql/schema.sql": {
			Data: []byte("CREATE TABLE radcheck (...);"),
		},
	})

	schema, err := loader.LoadMySQLSchema("3.2.10")
	if err != nil {
		t.Fatalf("LoadMySQLSchema() error = %v", err)
	}

	if got, want := schema, "CREATE TABLE radcheck (...);"; got != want {
		t.Errorf("LoadMySQLSchema() = %q, want %q", got, want)
	}
}

func TestLoadMySQLSchemaMissingVersion(t *testing.T) {
	loader := NewLoader(fstest.MapFS{})

	_, err := loader.LoadMySQLSchema("3.2.10")
	if err == nil {
		t.Fatal("LoadMySQLSchema() expected error for missing version")
	}

	if !strings.Contains(err.Error(), `FreeRADIUS version "3.2.10"`) {
		t.Errorf("LoadMySQLSchema() error = %q, want version in error", err)
	}
}

func TestLoadMySQLSchemaInvalidVersion(t *testing.T) {
	loader := NewLoader(fstest.MapFS{})

	tests := []string{
		"",
		".",
		"../3.2.10",
		"3.2.10/../../other",
	}

	for _, version := range tests {
		t.Run(version, func(t *testing.T) {
			_, err := loader.LoadMySQLSchema(version)
			if err == nil {
				t.Fatalf("LoadMySQLSchema(%q) expected error", version)
			}
		})
	}
}

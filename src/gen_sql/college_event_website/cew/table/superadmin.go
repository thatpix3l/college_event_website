//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Superadmin = newSuperadminTable("cew", "superadmin", "")

type superadminTable struct {
	postgres.Table

	// Columns
	ID postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type SuperadminTable struct {
	superadminTable

	EXCLUDED superadminTable
}

// AS creates new SuperadminTable with assigned alias
func (a SuperadminTable) AS(alias string) *SuperadminTable {
	return newSuperadminTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new SuperadminTable with assigned schema name
func (a SuperadminTable) FromSchema(schemaName string) *SuperadminTable {
	return newSuperadminTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new SuperadminTable with assigned table prefix
func (a SuperadminTable) WithPrefix(prefix string) *SuperadminTable {
	return newSuperadminTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new SuperadminTable with assigned table suffix
func (a SuperadminTable) WithSuffix(suffix string) *SuperadminTable {
	return newSuperadminTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newSuperadminTable(schemaName, tableName, alias string) *SuperadminTable {
	return &SuperadminTable{
		superadminTable: newSuperadminTableImpl(schemaName, tableName, alias),
		EXCLUDED:        newSuperadminTableImpl("", "excluded", ""),
	}
}

func newSuperadminTableImpl(schemaName, tableName, alias string) superadminTable {
	var (
		IDColumn       = postgres.StringColumn("id")
		allColumns     = postgres.ColumnList{IDColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return superadminTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID: IDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}

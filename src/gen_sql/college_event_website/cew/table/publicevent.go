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

var Publicevent = newPubliceventTable("cew", "publicevent", "")

type publiceventTable struct {
	postgres.Table

	// Columns
	ID       postgres.ColumnString
	Approved postgres.ColumnBool

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PubliceventTable struct {
	publiceventTable

	EXCLUDED publiceventTable
}

// AS creates new PubliceventTable with assigned alias
func (a PubliceventTable) AS(alias string) *PubliceventTable {
	return newPubliceventTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PubliceventTable with assigned schema name
func (a PubliceventTable) FromSchema(schemaName string) *PubliceventTable {
	return newPubliceventTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PubliceventTable with assigned table prefix
func (a PubliceventTable) WithPrefix(prefix string) *PubliceventTable {
	return newPubliceventTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PubliceventTable with assigned table suffix
func (a PubliceventTable) WithSuffix(suffix string) *PubliceventTable {
	return newPubliceventTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPubliceventTable(schemaName, tableName, alias string) *PubliceventTable {
	return &PubliceventTable{
		publiceventTable: newPubliceventTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newPubliceventTableImpl("", "excluded", ""),
	}
}

func newPubliceventTableImpl(schemaName, tableName, alias string) publiceventTable {
	var (
		IDColumn       = postgres.StringColumn("id")
		ApprovedColumn = postgres.BoolColumn("approved")
		allColumns     = postgres.ColumnList{IDColumn, ApprovedColumn}
		mutableColumns = postgres.ColumnList{ApprovedColumn}
	)

	return publiceventTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:       IDColumn,
		Approved: ApprovedColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}

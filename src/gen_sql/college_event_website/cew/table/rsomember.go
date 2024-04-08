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

var Rsomember = newRsomemberTable("cew", "rsomember", "")

type rsomemberTable struct {
	postgres.Table

	// Columns
	ID      postgres.ColumnString
	RsoID   postgres.ColumnString
	IsAdmin postgres.ColumnBool

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type RsomemberTable struct {
	rsomemberTable

	EXCLUDED rsomemberTable
}

// AS creates new RsomemberTable with assigned alias
func (a RsomemberTable) AS(alias string) *RsomemberTable {
	return newRsomemberTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RsomemberTable with assigned schema name
func (a RsomemberTable) FromSchema(schemaName string) *RsomemberTable {
	return newRsomemberTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RsomemberTable with assigned table prefix
func (a RsomemberTable) WithPrefix(prefix string) *RsomemberTable {
	return newRsomemberTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RsomemberTable with assigned table suffix
func (a RsomemberTable) WithSuffix(suffix string) *RsomemberTable {
	return newRsomemberTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRsomemberTable(schemaName, tableName, alias string) *RsomemberTable {
	return &RsomemberTable{
		rsomemberTable: newRsomemberTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newRsomemberTableImpl("", "excluded", ""),
	}
}

func newRsomemberTableImpl(schemaName, tableName, alias string) rsomemberTable {
	var (
		IDColumn       = postgres.StringColumn("id")
		RsoIDColumn    = postgres.StringColumn("rso_id")
		IsAdminColumn  = postgres.BoolColumn("is_admin")
		allColumns     = postgres.ColumnList{IDColumn, RsoIDColumn, IsAdminColumn}
		mutableColumns = postgres.ColumnList{IsAdminColumn}
	)

	return rsomemberTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:      IDColumn,
		RsoID:   RsoIDColumn,
		IsAdmin: IsAdminColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}

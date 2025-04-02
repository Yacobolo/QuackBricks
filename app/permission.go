package main

type ObjectType string

const (
	Catalog ObjectType = "catalog"
	Schema  ObjectType = "schema"
	Table   ObjectType = "table"
)

type Privilege string

const (
	Select Privilege = "SELECT"
	Modify Privilege = "MODIFY"
	Usage  Privilege = "USAGE"
	All    Privilege = "ALL_PRIVILEGES"
)

type ObjectReference struct {
	Type ObjectType
	Name string // e.g., "main.sales.customers"
}

type Principal struct {
	Name string // e.g., "alice@databricks.com" or "data_analysts"
}

type PermissionEntry struct {
	Principal  Principal
	Object     ObjectReference
	Privileges map[Privilege]bool
}

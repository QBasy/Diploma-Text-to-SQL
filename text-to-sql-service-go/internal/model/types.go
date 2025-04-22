package model

type Request struct {
	Query  string  `json:"query"`
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	Name       string   `json:"name"`
	Columns    []Column `json:"columns"`
	PrimaryKey string   `json:"primaryKey"`
}

type Column struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	IsForeignKey     bool   `json:"isForeignKey"`
	ReferencedTable  string `json:"referencedTable,omitempty"`
	ReferencedColumn string `json:"referencedColumn,omitempty"`
}

type Response struct {
	SqlQuery string `json:"sql"`
	Error    string `json:"error,omitempty"`
}

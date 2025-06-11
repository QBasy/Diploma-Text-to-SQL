package parsers

import "strings"

type SQLQueryType int

const (
	QueryTypeSelect SQLQueryType = iota
	QueryTypeCreate
	QueryTypeDrop
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
	QueryTypeAlter
	QueryTypeUnknown
)

type SQLParser struct{}

func NewSQLParser() *SQLParser {
	return &SQLParser{}
}

func (p *SQLParser) ParseQueryType(query string) SQLQueryType {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))

	switch {
	case strings.HasPrefix(upperQuery, "SELECT"):
		return QueryTypeSelect
	case strings.HasPrefix(upperQuery, "CREATE"):
		return QueryTypeCreate
	case strings.HasPrefix(upperQuery, "DROP"):
		return QueryTypeDrop
	case strings.HasPrefix(upperQuery, "INSERT"):
		return QueryTypeInsert
	case strings.HasPrefix(upperQuery, "UPDATE"):
		return QueryTypeUpdate
	case strings.HasPrefix(upperQuery, "DELETE"):
		return QueryTypeDelete
	case strings.HasPrefix(upperQuery, "ALTER"):
		return QueryTypeAlter
	default:
		return QueryTypeUnknown
	}
}

func (p *SQLParser) IsReadOnlyQuery(queryType SQLQueryType) bool {
	return queryType == QueryTypeSelect
}

func (p *SQLParser) IsWriteQuery(queryType SQLQueryType) bool {
	return queryType == QueryTypeCreate ||
		queryType == QueryTypeDrop ||
		queryType == QueryTypeInsert ||
		queryType == QueryTypeUpdate ||
		queryType == QueryTypeDelete ||
		queryType == QueryTypeAlter
}

package services

import (
	pb "database-service/internal/proto/generated/visualisationpb"
	"database-service/pkg/adapters"
	"database-service/pkg/parsers"
	"database-service/pkg/utils" // добавим сюда наш utils
	"fmt"
	"strings"
)

type QueryExecutionService struct {
	parser *parsers.SQLParser
}

func NewQueryExecutionService() *QueryExecutionService {
	return &QueryExecutionService{
		parser: parsers.NewSQLParser(),
	}
}

type QueryResult struct {
	Columns  []string
	Rows     []map[string]interface{}
	RowCount int
}

func (s *QueryExecutionService) ExecuteSelectQuery(conn adapters.Queryable, query string) (*QueryResult, error) {
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL query: %v", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	var result []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))
		for i := range columns {
			columnPtrs[i] = &columns[i]
		}

		if err := rows.Scan(columnPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[col] = columns[i]
		}
		result = append(result, rowMap)
	}

	return &QueryResult{
		Columns:  cols,
		Rows:     result,
		RowCount: len(result),
	}, nil
}

func (s *QueryExecutionService) ExecuteModifyQuery(conn adapters.Queryable, query string) error {
	statements := utils.SplitSQLStatements(query)
	for _, stmt := range statements {
		queryType := s.parser.ParseQueryType(stmt)
		if queryType == parsers.QueryTypeUnknown {
			return fmt.Errorf("query type not allowed: %s", stmt)
		}
		_, err := conn.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement: %v", err)
		}
	}
	return nil
}

func (s *QueryExecutionService) PrepareForVisualization(conn adapters.Queryable, query string) (*pb.QueryResult, error) {
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("invalid SQL query: %v", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch column info: %v", err)
	}

	var queryResult pb.QueryResult
	queryResult.SqlQuery = query

	for rows.Next() {
		columnVals := make([]interface{}, len(cols))
		columnPtrs := make([]interface{}, len(cols))
		for i := range columnVals {
			columnPtrs[i] = &columnVals[i]
		}

		if err := rows.Scan(columnPtrs...); err != nil {
			return nil, fmt.Errorf("failed to process row data: %v", err)
		}

		row := &pb.Row{}
		for _, val := range columnVals {
			row.Values = append(row.Values, fmt.Sprintf("%v", val))
		}
		queryResult.Result = append(queryResult.Result, row)
	}

	return &queryResult, nil
}

func cleanSQLStatement(stmt string) string {
	stmt = strings.TrimSpace(stmt)
	if strings.HasPrefix(stmt, "--") || strings.HasPrefix(stmt, "//") {
		return ""
	}
	return stmt
}

func (s *QueryExecutionService) ValidateQuery(query string) (parsers.SQLQueryType, error) {
	statements := utils.SplitSQLStatements(query)
	if len(statements) == 0 {
		return parsers.QueryTypeUnknown, fmt.Errorf("empty query")
	}

	for _, stmt := range statements {
		stmt = cleanSQLStatement(stmt)
		if stmt == "" {
			continue
		}
		queryType := s.parser.ParseQueryType(stmt)
		if queryType == parsers.QueryTypeUnknown {
			return queryType, fmt.Errorf("query type not allowed: %s", stmt)
		}
	}
	return s.parser.ParseQueryType(cleanSQLStatement(statements[0])), nil
}

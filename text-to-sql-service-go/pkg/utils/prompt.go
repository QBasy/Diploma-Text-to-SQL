package utils

import (
	"fmt"
	"text-to-sql/internal/model"
)

func CreatePrompt(req *model.Request) (string, string) {
	prompt := "Convert this natural language query to SQL: " + req.Query
	systemMessage := "You are a SQL expert. Convert natural language queries to valid SQL statements. Return only the SQL code without any explanation."
	return prompt, systemMessage
}

func CreateComplexPrompt(req *model.Request) (string, string) {

	schemaText := ""
	for _, table := range req.Schema.Tables {
		tableDesc := fmt.Sprintf("Table %s: ", table.Name)
		for _, column := range table.Columns {
			colDesc := fmt.Sprintf("%s (%s)", column.Name, column.Type)
			if column.IsForeignKey {
				colDesc += fmt.Sprintf(" [FK to %s.%s]", column.ReferencedTable, column.ReferencedColumn)
			}
			tableDesc += colDesc + ", "
		}
		schemaText += tableDesc + "\n"
	}

	prompt := fmt.Sprintf(`
Given the following database schema:

%s

Convert this natural language query to a valid SQL statement:
"%s"

Return only the SQL code without any explanation.
`, schemaText, req.Query)

	systemMessage := "You are a SQL expert. Convert natural language queries to valid SQL based on the provided schema."
	return prompt, systemMessage
}

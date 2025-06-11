package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	pb "visualisation-service/generated/visualisationpb"
	"visualisation-service/pkg/logger"
)

func GenerateChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	if len(data.Result) < 2 {
		return nil, fmt.Errorf("not enough data to draw graph")
	}

	chartType := determineChartType(data.SqlQuery, data.Result)
	logger.InfoLogger.Printf("Определен тип графика: %s", chartType)

	switch chartType {
	case ChartTypeBar:
		return generateBarChart(data)
	case ChartTypeLine:
		return generateLineChart(data)
	case ChartTypePie:
		return generatePieChart(data)
	case ChartTypeScatter:
		return generateScatterChart(data)
	default:
		logger.WarnLogger.Printf("Не удалось определить тип графика. Используется Bar по умолчанию.")
		return generateBarChart(data)
	}
}

func determineChartType(query string, results []*pb.Row) string {
	query = strings.ToLower(query)

	if strings.Contains(query, "percent") || strings.Contains(query, "ratio") {
		return ChartTypePie
	}

	if strings.Contains(query, "count") || strings.Contains(query, "sum") ||
		strings.Contains(query, "group by") {
		return ChartTypeBar
	}

	if strings.Contains(query, "date") || strings.Contains(query, "time") ||
		strings.Contains(query, "strftime") || strings.Contains(query, "month") {
		return ChartTypeLine
	}

	if len(results) > 0 && len(results[0].Values) >= 2 {
		numericColumns := 0
		for i := 0; i < len(results[0].Values); i++ {
			if _, err := strconv.ParseFloat(results[0].Values[i], 64); err == nil {
				numericColumns++
			}
		}
		if numericColumns >= 2 {
			return ChartTypeScatter
		}
	}

	return ChartTypeBar
}

func extractTitle(query string) string {
	query = strings.ToLower(query)

	if strings.Contains(query, "select") && strings.Contains(query, "from") {
		parts := strings.Split(query, "from")
		if len(parts) > 1 {
			tableParts := strings.Fields(parts[1])
			if len(tableParts) > 0 {
				return "Data from " + strings.TrimSpace(tableParts[0])
			}
		}
	}

	return "SQL Chart Visualization"
}

func formatNumber(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%d", int(value))
	}
	return fmt.Sprintf("%.2f", value)
}

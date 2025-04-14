package service

import (
	"bytes"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"math"
	"strconv"
	"strings"
	"time"
	pb "visualisation-service/generated/visualisationpb"
)

func generateLineChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	rows := data.Result[1:]

	width, height := 800, 500
	margin := struct {
		top, right, bottom, left int
	}{60, 60, 80, 80}

	chartWidth := width - margin.left - margin.right
	chartHeight := height - margin.top - margin.bottom

	maxValue := 0.0
	minValue := math.MaxFloat64
	values := make([]float64, len(rows))
	labels := make([]string, len(rows))

	for i, row := range rows {
		if len(row.Values) < 2 {
			continue
		}

		labels[i] = row.Values[0]
		v, err := strconv.ParseFloat(row.Values[1], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования данных: %v", err)
		}
		values[i] = v
		if v > maxValue {
			maxValue = v
		}
		if v < minValue {
			minValue = v
		}
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	canvas.Style("text { font-family: 'Arial', sans-serif; }")

	canvas.Rect(0, 0, width, height, "fill:#f8f9fa")

	title := extractTitle(data.SqlQuery)
	canvas.Text(width/2, 30, title, "text-anchor:middle;font-size:18px;font-weight:bold;fill:#212529")

	gridLineStyle := "stroke:#dee2e6;stroke-width:1;stroke-dasharray:5,5"
	numGridLines := 5
	for i := 0; i <= numGridLines; i++ {
		y := margin.top + chartHeight - (i * chartHeight / numGridLines)
		canvas.Line(margin.left, y, width-margin.right, y, gridLineStyle)

		valueAtLine := minValue + (maxValue-minValue)*float64(i)/float64(numGridLines)
		valueLabel := formatNumber(valueAtLine)
		canvas.Text(margin.left-10, y+5, valueLabel, "text-anchor:end;font-size:12px;fill:#6c757d")
	}

	for i := 0; i < len(rows); i++ {
		x := margin.left + i*(chartWidth/(len(rows)-1))
		canvas.Line(x, margin.top, x, height-margin.bottom, gridLineStyle)
	}

	axisStyle := "stroke:#343a40;stroke-width:2"
	canvas.Line(margin.left, height-margin.bottom, width-margin.right, height-margin.bottom, axisStyle) // Ось X
	canvas.Line(margin.left, margin.top, margin.left, height-margin.bottom, axisStyle)                  // Ось Y

	var path strings.Builder
	var areaPath strings.Builder

	path.WriteString("M ")
	areaPath.WriteString("M ")

	mainColor := "#4895ef"

	for i, v := range values {
		if i >= len(rows) || i >= len(values) {
			continue
		}

		x := margin.left + i*(chartWidth/(len(rows)-1))
		y := margin.top + chartHeight - int((v/maxValue)*float64(chartHeight))

		if i == 0 {
			path.WriteString(fmt.Sprintf("%d %d ", x, y))
			areaPath.WriteString(fmt.Sprintf("%d %d ", x, y))
		} else {
			path.WriteString(fmt.Sprintf("L %d %d ", x, y))
			areaPath.WriteString(fmt.Sprintf("L %d %d ", x, y))
		}

		canvas.Circle(x, y, 5, fmt.Sprintf("fill:%s;stroke:#fff;stroke-width:2", mainColor))

		if len(rows) <= 20 || i%(len(rows)/10+1) == 0 {
			canvas.Text(x, height-margin.bottom+20, labels[i],
				"text-anchor:middle;font-size:12px;fill:#495057")
		}

		canvas.Title(fmt.Sprintf("%s: %s", labels[i], formatNumber(v)))
	}

	if len(rows) > 0 {
		areaPath.WriteString(fmt.Sprintf("L %d %d L %d %d Z",
			margin.left+(len(rows)-1)*(chartWidth/(len(rows)-1)), height-margin.bottom,
			margin.left, height-margin.bottom))

		canvas.Path(areaPath.String(), fmt.Sprintf("fill:%s;fill-opacity:0.1", mainColor))
	}

	canvas.Path(path.String(), fmt.Sprintf("fill:none;stroke:%s;stroke-width:3;stroke-linecap:round;stroke-linejoin:round", mainColor))

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	canvas.Text(width-margin.right, height-10, "Generated: "+timestamp,
		"text-anchor:end;font-size:10px;fill:#6c757d")

	canvas.End()
	return buf, nil
}

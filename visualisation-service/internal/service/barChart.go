package service

import (
	"bytes"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"log"
	"math"
	"strconv"
	"time"
	pb "visualisation-service/generated/visualisationpb"
)

func generateBarChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	log.Println("Generating Bar Chart")
	rows := data.Result[1:]

	width, height := 800, 500
	margin := struct {
		top, right, bottom, left int
	}{60, 60, 80, 80}

	chartWidth := width - margin.left - margin.right
	chartHeight := height - margin.top - margin.bottom

	barWidth := int(math.Max(5, float64(chartWidth/len(rows))-10))

	maxValue := 0.0
	minValue := math.MaxFloat64
	values := make([]float64, len(rows))
	labels := make([]string, len(rows))

	for i, row := range rows {
		if len(row.Values) < 3 {
			continue
		}

		v, err := strconv.ParseFloat(row.Values[0], 64)
		if err != nil {
			log.Printf("Skipping non-numeric count value: %v", row.Values[0])
			continue
		}
		values[i] = v

		labels[i] = row.Values[1] + " (ID: " + row.Values[2] + ")"

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

	canvas.Rect(0, 0, width, height, "fill:#f8f9fa")

	title := extractTitle(data.SqlQuery)
	canvas.Text(width/2, 30, title, "text-anchor:middle;font-size:18px;font-weight:bold;fill:#212529")

	colors := []string{
		"#4361ee", "#3a0ca3", "#7209b7", "#f72585", "#4cc9f0",
		"#4895ef", "#560bad", "#b5179e", "#480ca8", "#3f37c9",
	}

	gridLineStyle := "stroke:#dee2e6;stroke-width:1;stroke-dasharray:5,5"
	numGridLines := 5
	for i := 0; i <= numGridLines; i++ {
		y := margin.top + chartHeight - (i * chartHeight / numGridLines)
		canvas.Line(margin.left, y, width-margin.right, y, gridLineStyle)

		valueAtLine := minValue + (maxValue-minValue)*float64(i)/float64(numGridLines)
		valueLabel := formatNumber(valueAtLine)
		canvas.Text(margin.left-10, y+5, valueLabel, "text-anchor:end;font-size:12px;fill:#6c757d")
	}

	axisStyle := "stroke:#343a40;stroke-width:2"
	canvas.Line(margin.left, height-margin.bottom, width-margin.right, height-margin.bottom, axisStyle) // Ось X
	canvas.Line(margin.left, margin.top, margin.left, height-margin.bottom, axisStyle)                  // Ось Y

	for i, v := range values {
		if i >= len(rows) {
			continue
		}

		x := margin.left + i*(chartWidth/len(rows)) + (chartWidth/len(rows)-barWidth)/2
		h := int((v / maxValue) * float64(chartHeight))
		y := margin.top + chartHeight - h

		color := colors[i%len(colors)]

		canvas.Rect(x, y, barWidth, h, fmt.Sprintf("fill:%s;rx:3;ry:3", color))

		canvas.Text(x+barWidth/2, y-10, formatNumber(v),
			"text-anchor:middle;font-size:12px;font-weight:bold;fill:#212529")

		labelStyle := "text-anchor:middle;font-size:12px;fill:#495057"
		if len(rows) > 10 {
			labelStyle += ";transform:rotate(-45deg);transform-origin:center"
		}
		canvas.Text(x+barWidth/2, height-margin.bottom+20, labels[i], labelStyle)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	canvas.Text(width-margin.right, height-10, "Generated: "+timestamp,
		"text-anchor:end;font-size:10px;fill:#6c757d")

	canvas.End()
	return buf, nil
}

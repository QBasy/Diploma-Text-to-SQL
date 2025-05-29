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

func generateScatterChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	log.Println("Generate Scatter Chart")
	rows := data.Result[1:]

	width, height := 800, 500
	margin := struct {
		top, right, bottom, left int
	}{60, 60, 80, 80}

	chartWidth := width - margin.left - margin.right
	chartHeight := height - margin.top - margin.bottom

	xValues := make([]float64, len(rows))
	yValues := make([]float64, len(rows))
	labels := make([]string, len(rows))

	var maxX, maxY float64
	var minX, minY float64 = math.MaxFloat64, math.MaxFloat64

	for i, row := range rows {
		if len(row.Values) < 3 {
			continue
		}

		labels[i] = row.Values[0]

		x, err := strconv.ParseFloat(row.Values[1], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования X-координаты: %v", err)
		}
		xValues[i] = x

		y, err := strconv.ParseFloat(row.Values[2], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования Y-координаты: %v", err)
		}
		yValues[i] = y

		if x > maxX {
			maxX = x
		}
		if x < minX {
			minX = x
		}
		if y > maxY {
			maxY = y
		}
		if y < minY {
			minY = y
		}
	}

	if maxX == minX {
		maxX += 1.0
	}
	if maxY == minY {
		maxY += 1.0
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	canvas.Style("text { font-family: 'Arial', sans-serif; }")

	canvas.Rect(0, 0, width, height, "fill:#ffffff")

	title := extractTitle(data.SqlQuery)
	canvas.Text(width/2, 30, title, "text-anchor:middle;font-size:18px;font-weight:bold;fill:#212529")

	gridLineStyle := "stroke:#dee2e6;stroke-width:1;stroke-dasharray:5,5"
	numGridLines := 5

	for i := 0; i <= numGridLines; i++ {
		y := margin.top + chartHeight - (i * chartHeight / numGridLines)
		canvas.Line(margin.left, y, width-margin.right, y, gridLineStyle)

		valueAtLine := minY + (maxY-minY)*float64(i)/float64(numGridLines)
		valueLabel := formatNumber(valueAtLine)
		canvas.Text(margin.left-10, y+5, valueLabel, "text-anchor:end;font-size:12px;fill:#6c757d")
	}

	for i := 0; i <= numGridLines; i++ {
		x := margin.left + (i * chartWidth / numGridLines)
		canvas.Line(x, margin.top, x, height-margin.bottom, gridLineStyle)

		valueAtLine := minX + (maxX-minX)*float64(i)/float64(numGridLines)
		valueLabel := formatNumber(valueAtLine)
		canvas.Text(x, height-margin.bottom+20, valueLabel, "text-anchor:middle;font-size:12px;fill:#6c757d")
	}

	axisStyle := "stroke:#343a40;stroke-width:2"
	canvas.Line(margin.left, height-margin.bottom, width-margin.right, height-margin.bottom, axisStyle)
	canvas.Line(margin.left, margin.top, margin.left, height-margin.bottom, axisStyle)

	pointColor := "#4895ef"
	for i := range rows {
		if i >= len(xValues) || i >= len(yValues) {
			continue
		}

		x := margin.left + int(((xValues[i]-minX)/(maxX-minX))*float64(chartWidth))
		y := margin.top + chartHeight - int(((yValues[i]-minY)/(maxY-minY))*float64(chartHeight))

		canvas.Circle(x, y, 6, fmt.Sprintf("fill:%s;stroke:#ffffff;stroke-width:2", pointColor))

		canvas.Title(fmt.Sprintf("%s: (%s, %s)", labels[i], formatNumber(xValues[i]), formatNumber(yValues[i])))
	}

	canvas.Text(width/2, height-10, "X", "text-anchor:middle;font-size:14px;fill:#343a40;font-weight:bold")
	canvas.Text(10, height/2, "Y", "text-anchor:middle;font-size:14px;fill:#343a40;font-weight:bold;writing-mode:tb")

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	canvas.Text(width-margin.right, height-10, "Generated: "+timestamp,
		"text-anchor:end;font-size:10px;fill:#6c757d")

	canvas.End()
	return buf, nil
}

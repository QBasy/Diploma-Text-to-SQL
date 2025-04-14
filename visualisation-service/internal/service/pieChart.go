package service

import (
	"bytes"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"math"
	"strconv"
	"time"
	pb "visualisation-service/generated/visualisationpb"
)

func generatePieChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	rows := data.Result[1:]

	width, height := 800, 500
	margin := 50
	centerX, centerY := width/2, height/2
	var radius int
	radius = int(math.Min(float64(width-2*margin), float64(height-2*margin)) / 2)

	values := make([]float64, len(rows))
	labels := make([]string, len(rows))
	total := 0.0

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
		total += v
	}

	colors := []string{
		"#4361ee", "#3a0ca3", "#7209b7", "#f72585", "#4cc9f0",
		"#4895ef", "#560bad", "#b5179e", "#480ca8", "#3f37c9",
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	canvas.Style("text { font-family: 'Arial', sans-serif; }")

	canvas.Rect(0, 0, width, height, "fill:#f8f9fa")

	title := extractTitle(data.SqlQuery)
	canvas.Text(width/2, 30, title, "text-anchor:middle;font-size:18px;font-weight:bold;fill:#212529")

	deg2rad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	startAngle := 0.0
	for i, value := range values {
		if value <= 0 {
			continue
		}

		percentage := value / total
		endAngle := startAngle + 360.0*percentage

		x1 := centerX + int(math.Cos(deg2rad(startAngle))*float64(radius))
		y1 := centerY + int(math.Sin(deg2rad(startAngle))*float64(radius))

		x2 := centerX + int(math.Cos(deg2rad(endAngle))*float64(radius))
		y2 := centerY + int(math.Sin(deg2rad(endAngle))*float64(radius))

		largeArcFlag := 0
		if endAngle-startAngle > 180 {
			largeArcFlag = 1
		}

		color := colors[i%len(colors)]

		pathData := fmt.Sprintf("M %d,%d L %d,%d A %d,%d 0 %d,1 %d,%d Z",
			centerX, centerY,
			x1, y1,
			radius, radius, largeArcFlag,
			x2, y2)

		canvas.Path(pathData, fmt.Sprintf("fill:%s", color))

		textAngle := startAngle + (endAngle-startAngle)/2
		textRadius := int(float64(radius) * 0.7)
		textX := centerX + int(math.Cos(deg2rad(textAngle))*float64(textRadius))
		textY := centerY + int(math.Sin(deg2rad(textAngle))*float64(textRadius))

		if percentage > 0.03 {
			canvas.Text(textX, textY, fmt.Sprintf("%.1f%%", percentage*100),
				"text-anchor:middle;font-size:12px;font-weight:bold;fill:white")
		}

		startAngle = endAngle
	}

	legendX := width - margin - 150
	legendY := margin + 50

	for i, label := range labels {
		if i >= len(values) || values[i] <= 0 {
			continue
		}

		color := colors[i%len(colors)]
		canvas.Rect(legendX, legendY+i*25, 15, 15, "fill:"+color)

		percentage := values[i] / total * 100
		canvas.Text(legendX+25, legendY+i*25+12,
			fmt.Sprintf("%s: %.1f%% (%.2f)", label, percentage, values[i]),
			"font-size:12px;fill:#212529")
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	canvas.Text(width-margin, height-10, "Generated: "+timestamp,
		"text-anchor:end;font-size:10px;fill:#6c757d")

	canvas.End()
	return buf, nil
}

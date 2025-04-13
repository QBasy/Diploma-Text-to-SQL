package main

import (
	"bytes"
	"context"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	pb "visualisation-service/generated/visualisationpb"
)

type server struct {
	pb.UnimplementedVisualisationServiceServer
}

func (s *server) GenerateChart(ctx context.Context, req *pb.QueryResult) (*pb.SVGResponse, error) {
	svgBuffer, err := generateChart(req)
	if err != nil {
		return nil, err
	}
	return &pb.SVGResponse{Svg: svgBuffer.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5007")
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterVisualisationServiceServer(s, &server{})

	log.Println("gRPC-сервер Visualisation запущен на :5007")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}

func generateChart(data *pb.QueryResult) (*bytes.Buffer, error) {
	if len(data.Result) < 2 {
		return nil, fmt.Errorf("недостаточно данных для построения графика")
	}

	rows := data.Result[1:]

	width, height := 600, 400
	margin := 50
	barWidth := (width - 2*margin) / len(rows)

	maxValue := 0.0
	values := make([]float64, len(rows))
	for i, row := range rows {
		v, err := strconv.ParseFloat(row.Values[1], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка преобразования данных: %v", err)
		}
		values[i] = v
		if v > maxValue {
			maxValue = v
		}
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")

	canvas.Line(margin, height-margin, width-margin, height-margin, "stroke:black;stroke-width:2")
	canvas.Line(margin, height-margin, margin, margin, "stroke:black;stroke-width:2")

	for i, v := range values {
		x := margin + i*barWidth
		h := int((v / maxValue) * float64(height-2*margin))
		y := height - margin - h

		canvas.Rect(x, y, barWidth-5, h, "fill:blue")
		canvas.Text(x+barWidth/2, height-margin+20, rows[i].Values[0], "text-anchor:middle;font-size:12px;fill:black")
	}

	canvas.Text(width/2, 20, "SQL Query Result", "text-anchor:middle;font-size:16px;fill:black")

	canvas.End()

	return buf, nil
}

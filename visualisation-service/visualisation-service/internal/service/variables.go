package service

const (
	ChartTypeBar     = "bar"
	ChartTypeLine    = "line"
	ChartTypePie     = "pie"
	ChartTypeScatter = "scatter"
)

var (
	modernPalette = []string{
		"#4361ee", "#3a0ca3", "#7209b7", "#f72585", "#4cc9f0",
		"#4895ef", "#560bad", "#b5179e", "#480ca8", "#3f37c9",
	}

	pastelPalette = []string{
		"#ffadad", "#ffd6a5", "#fdffb6", "#caffbf", "#9bf6ff",
		"#a0c4ff", "#bdb2ff", "#ffc6ff", "#fffffc", "#ffd6ff",
	}

	vibrantPalette = []string{
		"#ff595e", "#ffca3a", "#8ac926", "#1982c4", "#6a4c93",
		"#f94144", "#f3722c", "#f8961e", "#f9c74f", "#90be6d",
	}
)

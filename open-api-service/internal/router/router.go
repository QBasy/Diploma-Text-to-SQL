package router

import (
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/config"
	"github.com/QBasy/Diploma-text-to-SQL/open-api-service/internal/handler"
	"net/http"
)

func New(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	textHandler := handler.NewTextSQLHandler(cfg)
	visualHandler := handler.NewVisualHandler(cfg)

	mux.HandleFunc("/api/simple", textHandler.HandleSimple)
	mux.HandleFunc("/api/complex", textHandler.HandleComplex)
	mux.HandleFunc("/api/visualise", visualHandler.HandleVisual)

	return mux
}

package handler

import (
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
	Config "text-to-sql/internal/config"
	"text-to-sql/internal/model"
	"text-to-sql/internal/service"
	"text-to-sql/pkg/utils"
)

func TextToSQLHandler(w http.ResponseWriter, r *http.Request) {
	chRequest := make(chan *model.Request)
	go func() {
		var req *model.Request
		req, err := utils.ParseRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			close(chRequest)
			return
		}
		chRequest <- req
	}()

	chConfig := make(chan *Config.Config)
	go func() {
		config, err := Config.LoadConfig("config/local.yaml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			close(chConfig)
			return
		}
		chConfig <- config
	}()

	req := <-chRequest
	config := <-chConfig

	if req == nil || config == nil {
		return
	}

	client := openai.NewClient(option.WithAPIKey(config.TTSQL.APIKey))

	prompt, systemMessage := utils.CreatePrompt(req)

	chSQL := make(chan string)
	go func() {
		sqlQuery, err := service.GenerateSQL(client, prompt, systemMessage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating SQL: %v", err), http.StatusInternalServerError)
			close(chSQL)
			return
		}
		chSQL <- sqlQuery
	}()

	sqlQuery := <-chSQL
	if sqlQuery == "" {
		return
	}

	response := model.Response{SqlQuery: sqlQuery}
	utils.ParseResponse(w, response)
}

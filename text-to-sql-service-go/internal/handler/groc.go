package handler

import (
	"fmt"
	"net/http"
	Config "text-to-sql/internal/config"
	"text-to-sql/internal/model"
	"text-to-sql/internal/service"
	"text-to-sql/pkg/utils"
)

func TextToSQLHandlerWithGroc(w http.ResponseWriter, r *http.Request) {
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

	client := &http.Client{}

	var prompt, systemMessage string
	if req.Schema == nil || len(req.Schema.Tables) == 0 {
		prompt, systemMessage = utils.CreateSimplePrompt(req)
	} else {
		prompt, systemMessage = utils.CreateComplexPrompt(req)
	}

	chSQL := make(chan string)
	go func() {
		sqlQuery, err := service.GenerateSQLWithGroc(client, prompt, systemMessage, config.TTSQL.GROC.MODEL, config.TTSQL.GROC.APIKey, config.TTSQL.GROC.BASEURL)
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

	sqlQuery = utils.StripSQLMarkdown(sqlQuery)

	response := model.Response{SqlQuery: sqlQuery}
	utils.ParseResponse(w, response)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ServiceConfig struct {
	Name     string `json:"name"`
	Target   string `json:"target"`
	BasePath string `json:"basePath"`
}

type ReverseProxy struct {
	services map[string]*httputil.ReverseProxy
	paths    map[string]string
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		services: make(map[string]*httputil.ReverseProxy),
		paths:    make(map[string]string),
	}
}

func (rp *ReverseProxy) AddService(config ServiceConfig) error {
	target, err := url.Parse(config.Target)
	if err != nil {
		return fmt.Errorf("invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Настраиваем модификацию запроса перед проксированием
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", target.Host)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, config.BasePath)
	}

	// Добавляем обработку ошибок
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Error proxying %s: %v", config.Name, err)
		http.Error(w, fmt.Sprintf("Service %s unavailable", config.Name), http.StatusServiceUnavailable)
	}

	rp.services[config.BasePath] = proxy
	rp.paths[config.BasePath] = config.Name
	return nil
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Находим подходящий сервис по префиксу пути
	var targetPath string
	var targetProxy *httputil.ReverseProxy

	for path, proxy := range rp.services {
		if strings.HasPrefix(r.URL.Path, path) {
			if len(path) > len(targetPath) {
				targetPath = path
				targetProxy = proxy
			}
		}
	}

	if targetProxy == nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	log.Printf("Proxying request to %s: %s", rp.paths[targetPath], r.URL.Path)
	targetProxy.ServeHTTP(w, r)
}

func main() {
	proxy := NewReverseProxy()

	var services []ServiceConfig

	err := json.Unmarshal([]byte("services"), &services)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	for _, service := range services {
		newService := ServiceConfig{service.Name, service.Target, service.BasePath}
		services = append(services, newService)
	}

	for _, service := range services {
		if err := proxy.AddService(service); err != nil {
			log.Fatalf("Failed to add service %s: %v", service.Name, err)
		}
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Proxy is healthy")
	})

	http.Handle("/", proxy)

	log.Printf("Starting reverse proxy on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

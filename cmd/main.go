package main

import (
	config "gateway/config"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	serviceConfig, err := config.ReadConfigYAML(os.Getenv("CONFIG_FILE_PATH"))
	if err != nil {
		log.Fatal("Config not readable")
	}

	servicesConfig := config.NewServiceConfig(serviceConfig.Services)

	proxy := &httputil.ReverseProxy{Director: servicesConfig.DirectorFunc}

	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Send request to " + r.URL.Path)
		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)
	log.Printf("Gateway listening on port 8000")
	http.ListenAndServe(":8000", nil)
}

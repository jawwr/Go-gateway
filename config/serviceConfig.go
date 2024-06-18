package config

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type ServicesConfig struct {
	configmap map[string]*url.URL
}

func NewServiceConfig(config map[string]ServiceConfig) *ServicesConfig {
	services := make(map[string]*url.URL)
	for serviceName, serviceConfig := range config {
		service, err := url.Parse(fmt.Sprintf(serviceConfig.Url))
		if err != nil {
			panic(err)
		}
		if conn, err := net.Dial("tcp", service.Host); err != nil {
			log.Printf("Service [ %s ] not available", serviceName)
			continue
		} else {
			conn.Close()
		}
		services[serviceConfig.Prefix] = service
		log.Printf("Service [ %s ] has been configured successfully", serviceName)
	}
	return &ServicesConfig{configmap: services}
}

func (config *ServicesConfig) DirectorFunc(req *http.Request) {
	prefix := getServicePrefixFromUrl(req.URL.Path)
	serviceUrl := config.configmap[prefix]
	if serviceUrl == nil {
		log.Printf("Prefix '%s' not found, try to get service", prefix)
		origin := req.Header["Referer"]
		if origin == nil {
			origin = req.Header["Origin"]
			if origin == nil {
				log.Printf("Origin not exists")
				return
			}
		}
		from := getServicePrefixFromUrl(origin[0])
		serviceUrl = config.configmap[from]
		if serviceUrl == nil {
			log.Printf("Service %s not fount", from)
			return
		}
	} else {
		req.URL.Path = strings.Split(req.URL.Path, fmt.Sprintf("/%s", prefix))[1]
	}
	req.URL.Host = serviceUrl.Host
	req.URL.Scheme = serviceUrl.Scheme
}

func getServicePrefixFromUrl(url string) string {
	regex, err := regexp.Compile("(https?://[A-Za-z-.]+)(:\\d{1,5})?")
	if err != nil {
		log.Fatal("pupupu expression")
	}
	if regex.MatchString(url) {
		servicePrefix := regex.Split(url, -1)[1]
		return strings.Split(servicePrefix, "/")[1]
	}
	return strings.Split(url, "/")[1]
}

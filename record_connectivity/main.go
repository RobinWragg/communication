package main

import (
	"os"
	"strings"
	"net/http"
	"time"
	"waylon_socket"
)

func GetHostname() string {
	hostname, _ := os.Hostname()
	hostname = strings.ToLower(hostname)
	return strings.TrimSuffix(hostname, ".local")
}

func RecordInternetConnectivity() {
	const timeout = time.Second * 2
	
	hostname := GetHostname()
	
	response_time_name := hostname + "_internet_response_time"
	access_name := hostname + "_internet_access"
	
	client := http.Client{Timeout: timeout}
	
	time.Sleep(time.Second)
	
	request_start_t := time.Now()
	response, e := client.Get("https://google.com")
	delta_secs := float32(time.Since(request_start_t).Seconds())
	
	if e != nil {
		delta_secs = float32(timeout.Seconds())
	} else {
		response.Body.Close()
	}
	
	waylon_socket.SendMeasurement(response_time_name, delta_secs)
	
	var access float32
	if delta_secs < float32(timeout.Seconds() * 0.999) { 
		access = 1
	} else {
		access = 0
	}
	
	waylon_socket.SendMeasurement(access_name, access)
}

func RecordNetworkConnectivity() {
	RecordPlatformSpecificConnectivity()
	// TODO: check router connectivity & modem connectivity
}

func RecordAll() {
	RecordInternetConnectivity()
	RecordNetworkConnectivity()
}

func main() {
	arguments := os.Args[1:]
	
	if len(arguments) > 0 && arguments[0] == "-repeat" {
		for {
			RecordAll()
			time.Sleep(time.Second)
		}
	} else {
		RecordAll()
	}
}





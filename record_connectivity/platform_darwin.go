package main

// TODO: maxRate, RSSI and noise: /System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport --getinfo

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"waylon_socket"
)

func RecordPlatformSpecificConnectivity() {
	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport", "--getinfo")
	
	response, e := cmd.Output()
	
	if e != nil {
		fmt.Println("airport command failed")
		return
	}
	
	response_str := string(response)
	
	entries := []string{
		"agrCtlRSSI",
		"agrCtlNoise",
		"maxRate",
	}
	
	for _, entry := range entries {
		reg, _ := regexp.Compile(entry + ":\\s+?([-\\d]+)")
		
		results := reg.FindStringSubmatch(response_str)
		
		if len(results) > 0 {
			if value, e := strconv.ParseFloat(results[1], 32); e == nil {
				waylon_socket.SendMeasurement("bmo_wifi_" + entry, float32(value))
			}
		}
	}
}





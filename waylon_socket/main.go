package waylon_socket

import (
	"os"
	"io/ioutil"
	"net"
	"time"
	"fmt"
	"strings"
	"encoding/json"
)

const serverAddress = "192.168.0.51:8080"
const timeout = time.Second

type Measurement struct {
	RelativeSeconds int64 `json:"relative_seconds"`
	Value float32 `json:"value"`
}

type Submission struct {
	MeasurementName string `json:"measurement_name"`
	RelativeSecondsNow int64 `json:"relative_seconds_now"`
	Measurements []Measurement `json:"measurements"`
}

func GetFilePath(measurementName string) string {
	return getTempDir() + measurementName + ".json"
}

func MakeJsonData(measurementName string, newValue float32) []byte {
	var submission Submission
	
	if _, e := os.Stat(getTempDir()); os.IsNotExist(e) {
		os.MkdirAll(getTempDir(), 0777)
	}
	
	filePath := GetFilePath(measurementName)
	
	_, e := os.Stat(filePath)
	
	if e == nil {
		jsonData, e := ioutil.ReadFile(filePath)
		if e != nil { panic(e) }
		
		e = json.Unmarshal(jsonData, &submission)
		if e != nil { panic(e) }
		
		submission.RelativeSecondsNow = time.Now().UTC().Unix()
		newMeasurement := Measurement{submission.RelativeSecondsNow, newValue}
		submission.Measurements = append(submission.Measurements, newMeasurement)
		fmt.Println("Added to existing submission")
	} else {
		if os.IsExist(e) {
			fmt.Println("FILE EXISTS BUT CANNOT BE OPENED.", filePath)
		} else {
			// File doesn't exist, so create it.
			submission.MeasurementName = measurementName
			submission.RelativeSecondsNow = time.Now().UTC().Unix()
			
			submission.Measurements = []Measurement{
				{submission.RelativeSecondsNow, newValue},
			}
			
			fmt.Println("Created submission")
		}
	}
	
	jsonData, e := json.MarshalIndent(submission, "", "\t")
	if e != nil { panic(e) }
	
	e = ioutil.WriteFile(filePath, jsonData, 0777)
	if e != nil { panic(e) }
	
	return jsonData
}

func SendMeasurement(name string, value float32) {
	fmt.Println("SendMeasurement", name, value)
	
	jsonData := MakeJsonData(name, value)
	
	conn, e := net.DialTimeout("tcp", serverAddress, timeout)
	
	if e == nil {
		defer conn.Close()
		
		conn.SetDeadline(time.Now().Add(timeout))
		
		fmt.Println("JSON", string(jsonData))
		_, e = conn.Write(jsonData)
		
		if e == nil {
			response := make([]byte, 32)
			_, e = conn.Read(response)
			
			if e == nil {
				response_str := string(response)
				
				if strings.Contains(response_str, "success") {
					fmt.Println("Sent")
					
					e = os.Remove(GetFilePath(name))
					if e != nil { panic(e) }
					
					fmt.Println("Removed", GetFilePath(name))
					return
				}
			}
		}
	}
	
	fmt.Println("Failed to send")
	fmt.Println("Unsent measurement(s) stored at", GetFilePath(name))
}




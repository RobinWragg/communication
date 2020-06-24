package waylon_socket

import "os"

func getTempDir() string {
	return os.ExpandEnv("${TEMP}\\waylon_socket\\")
}
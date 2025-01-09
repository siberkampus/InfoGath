package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir(".\\files"))


	http.Handle("/", fs)
	port := ":9000"
	fmt.Printf("File server başlatıldı, %s portunda dinleniyor...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

package main

import (
	"net/http"
	"os"
	"log"
	"fmt"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)

	port := "8080"
	if envPort, exists := os.LookupEnv("PORT"); exists {
		port = envPort
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}


func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	builder := strings.Builder{}
	_, _ = builder.WriteString("headers for this request:\n\n")
	for key, val := range r.Header {
		vals := strings.Join(val, " ")
		_, _ = builder.WriteString(fmt.Sprintf("%s : %s\n", key, vals))
	}
	output := builder.String()
	_, _ = fmt.Print(output)
	_, _ = w.Write([]byte(output))
}
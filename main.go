package main

import (
	"net/http"
	"os"
	"log"
	"fmt"
	"strings"
	"github.com/coreos/go-oidc"
	"context"
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

func getBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	tok := strings.ReplaceAll(auth, "Bearer", "")
	return strings.TrimSpace(tok)
}

func handler(w http.ResponseWriter, r *http.Request) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("request to %s\n", r.URL.Path))
	if val := r.URL.Query().Get("googleoidc"); len(val) > 0 {
		ctx := context.Background()
		provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		verifier := provider.Verifier(&oidc.Config{
			ClientID: r.URL.Path,
			SkipClientIDCheck: false,
			SkipExpiryCheck: false,
			SkipIssuerCheck: false,
		})
		var token string
		if val := r.URL.Query().Get("tok"); len(val) > 0 {
			token = val
		} else {
			token = getBearerToken(r)
		}
		tokVer, err := verifier.Verify(ctx, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		builder.WriteString("token valid\n")
		builder.WriteString(fmt.Sprintf("token detail: %+v \n\n", tokVer))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = builder.WriteString("headers for this request:\n\n")
	for key, val := range r.Header {
		vals := strings.Join(val, " ")
		_, _ = builder.WriteString(fmt.Sprintf("%s : %s\n", key, vals))
	}
	output := builder.String()
	_, _ = fmt.Print(strings.ReplaceAll(output, "\n", " "))
	_, _ = w.Write([]byte(output))
}
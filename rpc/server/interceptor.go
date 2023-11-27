package server

import (
	"net/http"
	"strings"
)

func Auth(handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := strings.Split(request.Header.Get("Authorization"), " ")
		if len(data) != 2 || data[1] != "test_auth" {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("unauth"))
			return
		}
		handler.ServeHTTP(writer, request)
	}
}

func Auth2() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data := strings.Split(request.Header.Get("Authorization"), " ")
		if len(data) != 2 || data[1] != "test_auth" {
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("unauth"))
			return
		}
	}
}

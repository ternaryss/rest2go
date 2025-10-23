package rest2go

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Middleware func(http.Handler) http.HandlerFunc

func Middlewares(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next.ServeHTTP
	}
}

func LogRequestAndResponseMiddleware(next http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var body []byte
		uid := uuid.New().String()
		reqContentType := request.Header.Get("Content-Type")

		if strings.Contains(reqContentType, "application/json") {
			body, _ = io.ReadAll(request.Body)
			request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		slog.Info("HTTP Request", "uid", uid, "method", request.Method, "path", request.RequestURI, "body", body)
		writer := newLogResponseWriter(response)
		next.ServeHTTP(writer, request)
		resContentType := writer.Header().Get("Content-Type")

		switch resContentType {
		case "application/json":
			slog.Info("HTTP Response", "uid", uid, "status", writer.status, "body", writer.body.String())
		default:
			slog.Info("HTTP Response", "uid", uid, "status", writer.status)
		}
	}
}

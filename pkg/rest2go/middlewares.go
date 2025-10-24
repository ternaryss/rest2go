package rest2go

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/ternaryss/rest2go/pkg/rest2go/settings"
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

func ApiKeyAuthMiddleware(authorization settings.Authorization) Middleware {
	conf := authorization.Header
	patterns := make([]*regexp.Regexp, 0, len(conf.Public))

	for _, pattern := range conf.Public {
		patterns = append(patterns, regexp.MustCompile("^"+antToRegex(pattern)+"$"))
	}

	return func(next http.Handler) http.HandlerFunc {
		return func(response http.ResponseWriter, request *http.Request) {
			if !conf.Enabled {
				next.ServeHTTP(response, request)
				return
			}

			url := path.Clean(request.URL.Path)

			for _, pattern := range patterns {
				if pattern.MatchString(url) {
					next.ServeHTTP(response, request)
					return
				}
			}

			key := request.Header.Get("Api-Key")

			if key != conf.Key {
				HandleError(NewApiError(401, "unauthorized"), response)
				return
			}

			next.ServeHTTP(response, request)
		}
	}
}

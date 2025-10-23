package rest2go

import (
	"bytes"
	"net/http"
)

type logResponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newLogResponseWriter(writer http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: writer,
		status:         -1,
		body:           bytes.NewBuffer(nil),
	}
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logResponseWriter) Write(body []byte) (int, error) {
	w.body.Write(body)

	return w.ResponseWriter.Write(body)
}

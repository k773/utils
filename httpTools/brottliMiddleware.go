package httpTools

import (
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	"strings"
)

func BrotliMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode incoming Brotli request if "Content-Encoding" is "br"
		if r.Header.Get("Content-Encoding") == "br" {
			brotliReader := brotli.NewReader(r.Body)
			r.Body = io.NopCloser(brotliReader)
			r.Header.Del("Content-Encoding")
		}

		// Encode outgoing response with Brotli if "Accept-Encoding" contains "br"
		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header().Set("Content-Encoding", "br")
			w.Header().Add("Vary", "Accept-Encoding")
			w.Header().Set("Accept-Encoding", "br")

			bw := brotli.NewWriterOptions(w, brotli.WriterOptions{Quality: 8})
			defer bw.Close()

			bwResponseWriter := &brotliResponseWriter{bw, w}
			next.ServeHTTP(bwResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
}

type brotliResponseWriter struct {
	bw *brotli.Writer
	http.ResponseWriter
}

func (brw *brotliResponseWriter) Write(b []byte) (int, error) {
	return brw.bw.Write(b)
}

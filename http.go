package utils

import (
	"net"
	"net/http"
)

// RunHttpOnRandomFreePortAsync will start new http server running on a separate goroutine serving the provided handler.
// If an error has occurred while setting-up the listener, it will be returned and the returned port will be 0.
// After http.Serve(...) is finished, onFinish() will be called with the error returned by http.Serve.
func RunHttpOnRandomFreePortAsync(handler http.Handler, onFinish func(e error)) (port int, listenerError error) {
	listener, listenerError := net.Listen("tcp", ":0")
	if listenerError != nil {
		return 0, listenerError
	}

	go onFinish(http.Serve(listener, handler))
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// RunHttpOnRandomFreePortSync will start new http server running on the same goroutine
// onPort will be called (synchronously) if the listener has successfully bound to a free port
func RunHttpOnRandomFreePortSync(handler http.Handler, onPort func(port int)) (e error) {
	listener, e := net.Listen("tcp", ":0")
	if e != nil {
		return e
	}

	onPort(listener.Addr().(*net.TCPAddr).Port)

	return http.Serve(listener, handler)
}

// RunHttpsTlsOnRandomFreePortSync will start new https+TLS server running on the same goroutine
// onPort will be called (synchronously) if the listener has successfully bound to a free port
func RunHttpsTlsOnRandomFreePortSync(handler http.Handler, certFile, keyFile string, onPort func(port int)) (e error) {
	listener, e := net.Listen("tcp", ":0")
	if e != nil {
		return e
	}

	onPort(listener.Addr().(*net.TCPAddr).Port)

	return http.ServeTLS(listener, handler, certFile, keyFile)
}

// HttpWrapWithResponseCodeFunc sets the specified status code to the http response;
func HttpWriteResponseCodeFunc(next http.HandlerFunc, code int) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(code)
		if next != nil {
			next(writer, request)
		}
	}
}

// HttpServeBytesFunc writes content of 'a' slice into the response. No extra headers are applied.
func HttpServeBytesFunc(next http.HandlerFunc, a []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(a)
		if next != nil {
			next(writer, request)
		}
	}
}

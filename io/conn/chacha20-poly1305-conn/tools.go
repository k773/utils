package chacha20_poly1305_conn

import "crypto/rand"

func generateNonce(size int) []byte {
	var buf = make([]byte, size)
	if _, e := rand.Read(buf); e != nil {
		panic(e)
	}
	return buf
}

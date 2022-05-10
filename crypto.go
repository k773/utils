package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func EnsureKeysExistence(pubKeyPath, prvKeyPath string, minKeySizeBits int, log io.Writer) (privateKey *rsa.PrivateKey, e error) {
	var pubKeyFound = FileExist(pubKeyPath)
	var prvKeyFound = FileExist(prvKeyPath)

	if !prvKeyFound {
		_, _ = log.Write([]byte("keys: Generating keys (this may take a while)..."))
		pubKeyFound = true

		if privateKey, e = rsa.GenerateKey(rand.Reader, minKeySizeBits); e != nil {
			panic(e)
		}

		// Saving keys
		if e = JoinErrors(ioutil.WriteFile(pubKeyPath, []byte(base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))), 0600), ioutil.WriteFile(prvKeyPath, []byte(base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))), 0600)); e == nil {
			_, _ = log.Write([]byte("keys: Done! Private key length: " + strconv.Itoa(privateKey.Size()*8) + ". Keys saved at locations: public key: " + pubKeyPath + "; private key: " + prvKeyPath))
		}
	} else {
		println("keys: Loading keys...")
		var data []byte
		data, e = ioutil.ReadFile(prvKeyPath)
		if data, e = base64.StdEncoding.DecodeString(string(data)); e == nil {
			if e == nil {
				if privateKey, e = x509.ParsePKCS1PrivateKey(data); e == nil {
					if privateKey.Size()*8 < minKeySizeBits {
						_ = os.RemoveAll(prvKeyPath)
						_ = os.RemoveAll(pubKeyPath)
						return EnsureKeysExistence(pubKeyPath, prvKeyPath, minKeySizeBits, log)
					}

					if !pubKeyFound {
						_, _ = log.Write([]byte("keys: Public key not found, regenerating it at location: " + pubKeyPath))
						if e = ioutil.WriteFile(pubKeyPath, []byte(base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))), 0600); e == nil {
							_, _ = log.Write([]byte("keys: Done! Private key length: " + strconv.Itoa(privateKey.Size()*8)))
						}
					}
				}
			}
		}
	}
	return privateKey, errors.Wrap(e, "EnsureKeysExistence")
}

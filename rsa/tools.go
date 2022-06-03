package rsaTools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/k773/utils"
	"io/ioutil"
	"os"
)

/*
	Public key serializing
*/

func ExportPublicKey(key *rsa.PublicKey) []byte {
	keyBytes := x509.MarshalPKCS1PublicKey(key)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})
}

func ImportPublicKey(key []byte) (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode(key)
	if publicKeyBlock == nil {
		return nil, errors.New("public key's decoded block is null")
	}
	return x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
}

/*
	Private key serializing
*/

func ExportPrivateKey(key *rsa.PrivateKey) []byte {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)
}

func ImportPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(key)
	if privateKeyBlock == nil {
		return nil, errors.New("private key's PEM decoded block is null")
	}

	return x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
}

/*
	Sign
*/

func SignRsa(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	return rsa.SignPSS(rand.Reader, key, crypto.SHA512, utils.Sha512B2B(data), &opts)
}

func VerifySign(pubKey *rsa.PublicKey, data, sign []byte) bool {
	var opts = rsa.PSSOptions{SaltLength: 20}
	return rsa.VerifyPSS(pubKey, crypto.SHA512, utils.Sha512B2B(data), sign, &opts) == nil
}

/*
	Keys generation
*/

func InitKeys(privateKeyPath, publicKeyPath string, keyLenBits int) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, e error) {
	if privateKey, e = LoadPrivateKey(privateKeyPath); e != nil {
		if privateKey, e = rsa.GenerateKey(rand.Reader, keyLenBits); e == nil {
			publicKey = &privateKey.PublicKey
			e = SavePublicAndPrivateKey(privateKeyPath, publicKeyPath, privateKey, publicKey)
		}
	} else {
		if privateKey.Size()*8 < keyLenBits {
			if e = utils.JoinErrors(os.Remove(privateKeyPath), os.Remove(publicKeyPath)); e == nil {
				return InitKeys(privateKeyPath, publicKeyPath, keyLenBits)
			}
		} else {
			if publicKey, e = LoadPublicKey(publicKeyPath); e != nil {
				publicKey = &privateKey.PublicKey
				e = SavePublicKey(publicKeyPath, publicKey)
			}
		}
	}
	return
}

// Private key

func SavePrivateKey(filePath string, pk *rsa.PrivateKey) error {
	return ioutil.WriteFile(filePath, ExportPrivateKey(pk), 0600)
}

func LoadPrivateKey(filePath string) (pk *rsa.PrivateKey, e error) {
	privateKeyBytes, e := ioutil.ReadFile(filePath)
	if e == nil {
		pk, e = ImportPrivateKey(privateKeyBytes)
	}
	return
}

// Public key

func SavePublicKey(filePath string, pk *rsa.PublicKey) error {
	return ioutil.WriteFile(filePath, ExportPublicKey(pk), 0600)
}

func LoadPublicKey(filePath string) (pk *rsa.PublicKey, e error) {
	publicKeyBytes, e := ioutil.ReadFile(filePath)
	if e == nil {
		pk, e = ImportPublicKey(publicKeyBytes)
	}
	return
}

// Public && private key

func SavePublicAndPrivateKey(privateFilePath, publicFilePath string, private *rsa.PrivateKey, public *rsa.PublicKey) error {
	return utils.JoinErrors(SavePrivateKey(privateFilePath, private), SavePublicKey(publicFilePath, public))
}

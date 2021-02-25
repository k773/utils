package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RSA struct {
}

func (RSA) ExportPublicKey(key *rsa.PublicKey) []byte {
	keyBytes := x509.MarshalPKCS1PublicKey(key)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})
}

func (RSA) ImportPublicKey(key []byte) (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode(key)
	if publicKeyBlock == nil {
		return nil, errors.New("public key's decoded block is null")
	}
	return x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
}

func (RSA) ExportPrivateKey(key *rsa.PrivateKey) []byte {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)
}

func (RSA) ImportPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(key)
	if privateKeyBlock == nil {
		return nil, errors.New("private key's PEM decoded block is null")
	}

	return x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
}

func (RSA) EncryptRsa(key *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, message, []byte(""))
}

func (RSA) DecryptRsa(key *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key, message, []byte(""))
}

func (RSA) SignRsa(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	return rsa.SignPSS(rand.Reader, key, crypto.SHA256, Sha256B2B(data), &opts)
}

func (RSA) VerifySign(pubKey *rsa.PublicKey, data, sign []byte) bool {
	var opts rsa.PSSOptions = rsa.PSSOptions{SaltLength: 20}
	return rsa.VerifyPSS(pubKey, crypto.SHA256, Sha256B2B(data), sign, &opts) == nil
}

//F

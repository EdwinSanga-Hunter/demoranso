package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// Encrypt content with a public key
func Encrypt(publicKey []byte, text []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return []byte{}, fmt.Errorf("failed to decode PEM public key")
	}
	pubkeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}
	pubkey, ok := pubkeyInterface.(*rsa.PublicKey)
	if !ok {
		return []byte{}, fmt.Errorf("the PEM data does not contain a RSA public key")
	}

	out, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, text, []byte(""))
	if err != nil {
		return []byte{}, err
	}

	return out, nil
}

// Decrypt content using a Private Key
func Decrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	// Extract the PEM-encoded data block
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return []byte{}, fmt.Errorf("bad key data: %s", "not PEM-encoded")
	}
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return []byte{}, fmt.Errorf(
			"Failed to read private key: password protected keys are\n" +
				"not supported. Please decrypt the key prior to use.")
	}

	// Accept both PKCS#1 ("RSA PRIVATE KEY", old openssl) and
	// PKCS#8 ("PRIVATE KEY", modern openssl) formats
	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return []byte{}, fmt.Errorf("Unknown key type %q, want %q or %q", block.Type, "RSA PRIVATE KEY", "PRIVATE KEY")
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err8 != nil {
			return []byte{}, fmt.Errorf("Bad private key: %s", err)
		}
		var ok bool
		priv, ok = key.(*rsa.PrivateKey)
		if !ok {
			return []byte{}, fmt.Errorf("Bad private key: not a RSA private key")
		}
	}

	var out []byte

	// Decrypt the data
	out, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, []byte(""))
	if err != nil {
		return []byte{}, err
	}

	return out, nil
}

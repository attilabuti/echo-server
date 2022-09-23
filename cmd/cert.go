package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func generateCert() (certFilePath string, keyFilePath string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = fmt.Errorf("failed to generate private key: %v", err)
		return
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		err = fmt.Errorf("failed to generate serial number: %v", err)
		return
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if certFilePath, err = createCert(&template, priv); err != nil {
		return
	}

	if keyFilePath, err = createKey(priv); err != nil {
		return
	}

	return
}

func createCert(template *x509.Certificate, priv *rsa.PrivateKey) (certFilePath string, err error) {
	certOut, err := os.CreateTemp(os.TempDir(), "server_crt_")
	if err != nil {
		return
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		err = fmt.Errorf("failed to create certificate: %v", err)
		return
	}

	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		err = fmt.Errorf("failed to write data to %s: %v", certOut.Name(), err)
		return
	}

	if err = certOut.Close(); err != nil {
		err = fmt.Errorf("error closing %s: %v", certOut.Name(), err)
		return
	}

	certFilePath = certOut.Name()

	return
}

func createKey(priv *rsa.PrivateKey) (keyFilePath string, err error) {
	keyOut, err := os.CreateTemp(os.TempDir(), "server_key_")
	if err != nil {
		return
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		err = fmt.Errorf("unable to marshal private key: %v", err)
		return
	}

	if err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		err = fmt.Errorf("failed to write data to %s: %v", keyOut.Name(), err)
		return
	}

	if err = keyOut.Close(); err != nil {
		err = fmt.Errorf("error closing %s: %v", keyOut.Name(), err)
		return
	}

	keyFilePath = keyOut.Name()

	return
}

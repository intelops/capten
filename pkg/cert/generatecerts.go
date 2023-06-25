package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func GenerateCerts() error {

	// Create directories
	err := os.MkdirAll("certs/root", 0755)
	if err != nil {
		logrus.Printf("failed to create root directory. Error - %vs", err.Error())
		return err
	}
	err = os.MkdirAll("certs/intermediate", 0755)
	if err != nil {
		logrus.Printf("failed to create intermediate directory. Error - %vs", err.Error())
		return err
	}
	err = os.MkdirAll("certs/server", 0755)
	if err != nil {
		logrus.Printf("failed to create server directory. Error - %vs", err.Error())
		return err
	}
	err = os.MkdirAll("certs/client", 0755)
	if err != nil {
		logrus.Printf("failed to create client directory. Error - %vs", err.Error())
		return err
	}

	// Root certificate creation
	rootKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.Printf("failed to generate RSA key for root certificate. Error - %vs", err.Error())
		return err
	}

	rootCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Intelops"},
			CommonName:   "Root CA",
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	rootCert, err := x509.CreateCertificate(rand.Reader, &rootCertTemplate, &rootCertTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		logrus.Printf("failed to create root certificate. Error - %vs", err.Error())
		return err
	}

	rootCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCert})
	err = ioutil.WriteFile("certs/root/root-cert.pem", rootCertPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from root cert to certs/root/root-cert.pem. Error - %vs", err.Error())
		return err
	}

	rootKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey)})
	err = ioutil.WriteFile("certs/root/root-key.pem", rootKeyPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from root key to certs/root/root-key.pem. Error - %vs", err.Error())
		return err
	}

	// Intermediate certificate creation
	interKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.Printf("failed to generate RSA key for intermediate certificate. Error - %vs", err.Error())
		return err
	}

	interCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Intelops"},
			CommonName:   "Optimizor CA",
			Locality:     []string{"agent"},
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	interCert, err := x509.CreateCertificate(rand.Reader, &interCertTemplate, &rootCertTemplate, &interKey.PublicKey, rootKey)
	if err != nil {
		logrus.Printf("failed to create intermediate certificate. Error - %vs", err.Error())
		return err
	}

	interCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: interCert})
	err = ioutil.WriteFile("certs/intermediate/ca-cert.pem", interCertPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from intermediate cert to certs/intermediate/ca-cert.pem. Error - %vs", err.Error())
		return err
	}

	interKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(interKey)})
	err = ioutil.WriteFile("certs/intermediate/ca-key.pem", interKeyPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from intermediate key to certs/intermediate/ca-key.pem. Error - %vs", err.Error())
		return err
	}

	catCertPEM, err := ioutil.ReadFile("certs/root/root-cert.pem")
	if err != nil {
		logrus.Printf("error while reading certs/root/root-cert.pem file . Error - %vs", err.Error())
		return err
	}

	err = ioutil.WriteFile("certs/intermediate/cert-chain.pem", append(interCertPEM, catCertPEM...), 0644)
	if err != nil {
		logrus.Printf("error while writing to certs/intermediate/cert-chain.pem file. Error - %vs", err.Error())
		return err
	}

	err = ioutil.WriteFile("certs/intermediate/root-cert.pem", catCertPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing to certs/intermediate/root-cert.pem file. Error - %vs", err.Error())
		return err
	}

	// Server certificate creation
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Printf("failed to generate RSA key for server certificate. Error - %vs", err.Error())
		return err
	}

	serverCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Intelops Inc."},
			CommonName:   "*.dev.optimizor.app",
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 1, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{"captenagent.dev.optimizor.app", "dev.optimizor.app", "localhost"},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	serverCert, err := x509.CreateCertificate(rand.Reader, &serverCertTemplate, &interCertTemplate, &serverKey.PublicKey, interKey)
	if err != nil {
		logrus.Printf("failed to create server certificate. Error - %vs", err.Error())
		return err
	}

	serverCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCert})
	err = ioutil.WriteFile("certs/server/server.crt", serverCertPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from server cert to certs/server/server.crt. Error - %vs", err.Error())
		return err
	}

	serverKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})
	err = ioutil.WriteFile("certs/server/server.key", serverKeyPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from server key to certs/server/server.key. Error - %vs", err.Error())
		return err
	}

	// Client certificate creation
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Printf("failed to generate RSA key for client certificate. Error - %vs", err.Error())
		return err
	}

	clientCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Intelops Inc."},
			CommonName:   "*.dev.optimizor.app",
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 1, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		DNSNames:              []string{"client.dev.optimizor.app", "localhost"},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &clientCertTemplate, &interCertTemplate, &clientKey.PublicKey, interKey)
	if err != nil {
		logrus.Printf("failed to create client certificate. Error - %vs", err.Error())
		return err
	}

	clientCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	err = ioutil.WriteFile("certs/client/client.crt", clientCertPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from client cert to certs/client/client.crt. Error - %vs", err.Error())
		return err
	}

	clientKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	err = ioutil.WriteFile("certs/client/client.key", clientKeyPEM, 0644)
	if err != nil {
		logrus.Printf("error while writing from client key to certs/client/client.key. Error - %vs", err.Error())
		return err
	}

	return nil
}

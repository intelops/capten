package cert

import (
	"archive/zip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const FILE_PERMISSION os.FileMode = 0755
const WRITE_FILE_PERMISSION os.FileMode = 0644
const CERTIFICATE_BITS_4096 = 4096
const CERTIFICATE_BITS_2048 = 2048

type Config struct {
	Domain     string     `yaml:"Domain"`
	CertConfig CertConfig `yaml:"CertConfig"`
}

type CertConfig struct {
	OrgName                    string   `yaml:"OrgName"`
	RootCertCommonName         string   `yaml:"RootCertCommonName"`
	IntermediateCertCommonName string   `yaml:"IntermediateCertCommonName"`
	ServerCertCommonName       string   `yaml:"ServerCertCommonName"`
	ServerDNSNames             []string `yaml:"ServerDNSNames"`
	ClientCertCommonName       string   `yaml:"ClientCertCommonName"`
}

/*
certDirPath - root path of 'which ends with 'certs'
configYAMLPath -  captem config yaml path
*/
func GenerateCerts(certDirPath, configYAMLPath string) error {

	// Read the YAML file
	data, err := ioutil.ReadFile(configYAMLPath)
	if err != nil {
		errStr := fmt.Sprintf("failed to read YAML capten config file. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Unmarshal YAML into struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		errStr := fmt.Sprintf("failed to unmarshal YAML config. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Create directories
	err = os.MkdirAll(certDirPath, FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("failed to create root directory. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Root certificate creation
	rootKey, err := rsa.GenerateKey(rand.Reader, CERTIFICATE_BITS_4096)
	if err != nil {
		errStr := fmt.Sprintf("failed to generate RSA key for root certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	rootCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{config.CertConfig.OrgName},
			CommonName:   config.CertConfig.RootCertCommonName,
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
		errStr := fmt.Sprintf("failed to create root certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	rootCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCert})
	err = ioutil.WriteFile(certDirPath+"/root-cert.pem", rootCertPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from root cert to certs/root-cert.pem. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	rootKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey)})
	err = ioutil.WriteFile(certDirPath+"/root-key.pem", rootKeyPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from root key to certs/root-key.pem. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Intermediate certificate creation
	interKey, err := rsa.GenerateKey(rand.Reader, CERTIFICATE_BITS_4096)
	if err != nil {
		errStr := fmt.Sprintf("failed to generate RSA key for intermediate certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	interCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{config.CertConfig.OrgName},
			CommonName:   config.CertConfig.IntermediateCertCommonName,
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
		errStr := fmt.Sprintf("failed to create intermediate certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	interCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: interCert})
	err = ioutil.WriteFile(certDirPath+"/ca-cert.pem", interCertPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from intermediate cert to certs/ca-cert.pem. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	interKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(interKey)})
	err = ioutil.WriteFile(certDirPath+"/ca-key.pem", interKeyPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from intermediate key to certs/ca-key.pem. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	catCertPEM, err := ioutil.ReadFile(certDirPath + "/root-cert.pem")
	if err != nil {
		errStr := fmt.Sprintf("error while reading certs/root-cert.pem file . Error - %vs", err.Error())
		return errors.New(errStr)
	}

	err = ioutil.WriteFile(certDirPath+"/ca-cert-chain.pem", append(interCertPEM, catCertPEM...), WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing to certs/ca-cert-chain.pem file. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	err = ioutil.WriteFile(certDirPath+"/root-cert.pem", catCertPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing to certs/root-cert.pem file. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Server certificate creation
	serverKey, err := rsa.GenerateKey(rand.Reader, CERTIFICATE_BITS_2048)
	if err != nil {
		errStr := fmt.Sprintf("failed to generate RSA key for server certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	serverCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{config.CertConfig.OrgName},
			CommonName:   config.CertConfig.ServerCertCommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 1, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              config.CertConfig.ServerDNSNames,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	serverCert, err := x509.CreateCertificate(rand.Reader, &serverCertTemplate, &interCertTemplate, &serverKey.PublicKey, interKey)
	if err != nil {
		errStr := fmt.Sprintf("failed to create server certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	serverCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCert})
	err = ioutil.WriteFile(certDirPath+"/server.crt", serverCertPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from server cert to certs/server.crt. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	serverKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})
	err = ioutil.WriteFile(certDirPath+"/server.key", serverKeyPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from server key to certs/server.key. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Client certificate creation
	clientKey, err := rsa.GenerateKey(rand.Reader, CERTIFICATE_BITS_2048)
	if err != nil {
		errStr := fmt.Sprintf("failed to generate RSA key for client certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	clientCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{config.CertConfig.OrgName},
			CommonName:   config.CertConfig.ClientCertCommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 1, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &clientCertTemplate, &interCertTemplate, &clientKey.PublicKey, interKey)
	if err != nil {
		errStr := fmt.Sprintf("failed to create client certificate. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	clientCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	err = ioutil.WriteFile(certDirPath+"/client.crt", clientCertPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from client cert to certs/client.crt. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	clientKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	err = ioutil.WriteFile(certDirPath+"/client.key", clientKeyPEM, WRITE_FILE_PERMISSION)
	if err != nil {
		errStr := fmt.Sprintf("error while writing from client key to certs/client.key. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	if err := GenerateClientCertZipFile(certDirPath); err != nil {
		return err
	}

	return nil
}

func GenerateClientCertZipFile(certDirPath string) error {
	// Create a new zip file
	zipFile, err := os.Create(certDirPath + "/server-client-auth-certs.zip")
	if err != nil {
		errStr := fmt.Sprintf("error while creating server-client-auth-certs.zip zip file. Error - %vs", err.Error())
		return errors.New(errStr)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add the client cert file to the zip
	err = addFileToZip(zipWriter, "client.crt", certDirPath+"/client.crt")
	if err != nil {
		errStr := fmt.Sprintf("error while adding client cert file to zip folder. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Add the client key file to the zip
	err = addFileToZip(zipWriter, "client.key", certDirPath+"/client.key")
	if err != nil {
		errStr := fmt.Sprintf("error while adding client key file to zip folder. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	// Add the ca cert chain cert file to the zip
	err = addFileToZip(zipWriter, "ca-cert-chain.pem", certDirPath+"/ca-cert-chain.pem")
	if err != nil {
		errStr := fmt.Sprintf("error while adding ca-cert-chain.pem file to zip folder. Error - %vs", err.Error())
		return errors.New(errStr)
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, fileName, filePath string) error {
	// Open the file to be added to the zip
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new file header
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = fileName

	// Create a new zip file entry
	entry, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write the file data to the zip entry
	_, err = io.Copy(entry, file)
	if err != nil {
		return err
	}

	return nil
}

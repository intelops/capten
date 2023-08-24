package cert

import (
	"archive/zip"
	"capten/pkg/clog"
	"capten/pkg/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	folderPrmission     os.FileMode = 0755
	filePrmission       os.FileMode = 0644
	caBitSize                       = 4096
	certBitSize                     = 2048
	rootCAKeyFileName               = "root.key"
	rootCACertFileName              = "root.crt"
	interCAKeyFileName              = "inter-ca.key"
	interCACertFileName             = "inter-ca.crt"
)

func PrepareCerts(captenConfig config.CaptenConfig) error {
	if !checkCertsExist(captenConfig) || captenConfig.ForceGenerateCerts {
		return generateCerts(captenConfig)
	}
	clog.Logger.Debug("Cert files exist, skipped generating certs")
	return nil
}

func checkCertsExist(captenConfig config.CaptenConfig) bool {
	certFiles := []string{
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, rootCACertFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, rootCAKeyFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, interCACertFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, interCAKeyFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentCertFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentKeyFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName),
	}

	for _, certFile := range certFiles {
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			clog.Logger.Debugf("Cert file %s does not exist", certFile)
			return false
		}
	}
	return true
}

func generateCerts(captenConfig config.CaptenConfig) error {

	err := os.MkdirAll(captenConfig.PrepareDirPath(captenConfig.CertDirPath), folderPrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to create directory %s", captenConfig.CertDirPath)
	}

	rootKey, rootCertTemplate, err := generateCACert(captenConfig)
	if err != nil {
		return err
	}

	interKey, interCACertTemplate, err := generateIntermediateCACert(captenConfig, rootKey, rootCertTemplate)
	if err != nil {
		return err
	}

	err = generateAgentCert(captenConfig, interKey, interCACertTemplate)
	if err != nil {
		return err
	}

	err = generateClientCert(captenConfig, interKey, interCACertTemplate)
	if err != nil {
		return err
	}

	if err := generateCACertChain(captenConfig); err != nil {
		return err
	}

	if err := generateCaptenClientCertZip(captenConfig); err != nil {
		return err
	}
	return nil
}

func generateCACert(captenConfig config.CaptenConfig) (rootKey *rsa.PrivateKey,
	rootCertTemplate *x509.Certificate, err error) {
	rootKey, err = rsa.GenerateKey(rand.Reader, caBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for root certificate")
		return
	}

	rootCertTemplate = &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{captenConfig.OrgName},
			CommonName:   captenConfig.RootCACommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	rootCert, err := x509.CreateCertificate(rand.Reader, rootCertTemplate, rootCertTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		err = errors.WithMessage(err, "failed to create root CA certificate")
		return
	}

	rootCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCert})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, rootCACertFileName), rootCertPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from root CA cert")
		return
	}

	rootKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey)})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, rootCAKeyFileName), rootKeyPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from root CA key")
		return
	}
	return
}

func generateIntermediateCACert(captenConfig config.CaptenConfig, rootKey *rsa.PrivateKey,
	rootCertTemplate *x509.Certificate) (interKey *rsa.PrivateKey, interCACertTemplate *x509.Certificate, err error) {
	interKey, err = rsa.GenerateKey(rand.Reader, caBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for intermediate certificate")
		return
	}

	interCACertTemplate = &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{captenConfig.OrgName},
			CommonName:   captenConfig.IntermediateCACommonName,
			Locality:     []string{"agent"},
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	interCert, err := x509.CreateCertificate(rand.Reader, interCACertTemplate, rootCertTemplate, &interKey.PublicKey, rootKey)
	if err != nil {
		err = errors.WithMessage(err, "failed to create intermediate CA certificate")
		return
	}

	interCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: interCert})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, interCACertFileName), interCertPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from intermediate CA cert")
		return
	}

	interKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(interKey)})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, interCAKeyFileName), interKeyPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from intermediate CA key")
		return
	}
	return
}

func generateAgentCert(captenConfig config.CaptenConfig, interKey *rsa.PrivateKey,
	interCACertTemplate *x509.Certificate) (err error) {
	agentKey, err := rsa.GenerateKey(rand.Reader, certBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for agent certificate")
		return
	}

	agentCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{captenConfig.OrgName},
			CommonName:   captenConfig.AgentCertCommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              captenConfig.AgentDNSNames,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	agentCert, err := x509.CreateCertificate(rand.Reader, &agentCertTemplate, interCACertTemplate, &agentKey.PublicKey, interKey)
	if err != nil {
		err = errors.WithMessage(err, "failed to create server certificate")
		return
	}

	agentCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: agentCert})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentCertFileName), agentCertPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from agent cert to certs/server.crt")
		return
	}

	agentKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(agentKey)})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.AgentKeyFileName), agentKeyPEM, filePrmission)
	if err != nil {
		err = errors.WithMessage(err, "error while writing from agent key to certs/server.key")
		return
	}
	return
}

func generateClientCert(captenConfig config.CaptenConfig, interKey *rsa.PrivateKey,
	interCACertTemplate *x509.Certificate) (err error) {
	clientKey, err := rsa.GenerateKey(rand.Reader, certBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for capten client certificate")
		return
	}

	clientCertTemplate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{captenConfig.OrgName},
			CommonName:   captenConfig.CaptenClientCertCommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &clientCertTemplate, interCACertTemplate, &clientKey.PublicKey, interKey)
	if err != nil {
		return errors.WithMessage(err, "failed to create client certificate")
	}

	clientCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName), clientCertPEM, filePrmission)
	if err != nil {
		return errors.WithMessage(err, "error while writing from client cert to certs/client.crt")
	}

	clientKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})
	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName), clientKeyPEM, filePrmission)
	if err != nil {
		return errors.WithMessage(err, "error while writing from client key to certs/client.key")
	}
	return
}

func generateCACertChain(captenConfig config.CaptenConfig) error {
	caCertPEMFromFile, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, rootCACertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading root CA cert file")
	}

	interCACertPEMFromFile, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, interCACertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while reading intermediate CA cert file")
	}

	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName), append(caCertPEMFromFile, interCACertPEMFromFile...), filePrmission)
	if err != nil {
		return errors.WithMessage(err, "error while writing to ca cert file")
	}
	return nil
}

func generateCaptenClientCertZip(captenConfig config.CaptenConfig) error {
	zipFile, err := os.Create(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertExportFileName))
	if err != nil {
		return errors.WithMessage(err, "error while creating client cert export file")
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = addFileToZip(zipWriter, captenConfig.ClientCertFileName,
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName))
	if err != nil {
		return errors.WithMessage(err, "error while adding client cert file to zip")
	}

	err = addFileToZip(zipWriter, captenConfig.ClientKeyFileName,
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName))
	if err != nil {
		return errors.WithMessage(err, "error while adding client key file to zip")
	}

	err = addFileToZip(zipWriter, captenConfig.CAFileName,
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return errors.WithMessage(err, "error while adding ca cert file to zip")
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, fileName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = fileName

	entry, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(entry, file)
	if err != nil {
		return err
	}
	return nil
}

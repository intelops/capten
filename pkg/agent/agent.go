package agent

import (
	"capten/pkg/config"
	"crypto/tls"
	"crypto/x509"
	"os"

	"capten/pkg/agent/agentpb"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetAgentClient(config config.CaptenConfig) (agentpb.AgentClient, error) {
	tlsCredentials, err := loadTLSCredentials(config)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to load capten agent client certs")
	}

	conn, err := grpc.Dial(config.GetCaptenAgentEndpoint(), grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to connect to capten agent")
	}

	return agentpb.NewAgentClient(conn), nil
}

func loadTLSCredentials(captenConfig config.CaptenConfig) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientCertFileName),
		captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.ClientKeyFileName))
	if err != nil {
		return nil, err
	}

	caCertChain, err := os.ReadFile(captenConfig.PrepareFilePath(captenConfig.CertDirPath, captenConfig.CAFileName))
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertChain); !ok {
		return nil, errors.New("failed to add server CA's certificate")
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientAuth:   tls.RequireAnyClientCert,
		ClientCAs:    caCertPool,
	}), nil
}

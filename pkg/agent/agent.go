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
	"google.golang.org/grpc/credentials/insecure"
)

func GetAgentClient(config config.CaptenConfig) (agentpb.AgentClient, error) {
	agentEndpoint := config.GetCaptenAgentEndpoint()

	var conn *grpc.ClientConn
	var err error
	if config.AgentSecure {
		tlsCredentials, err := loadTLSCredentials(config)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to load capten agent client certs")
		}

		conn, err = grpc.Dial(agentEndpoint, grpc.WithTransportCredentials(tlsCredentials))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to connect to capten agent")
		}
	} else {
		conn, err = grpc.Dial(agentEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to connect to capten agent")
		}
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

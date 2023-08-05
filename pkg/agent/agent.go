package agent

import (
	"bytes"
	"capten/pkg/clog"
	"capten/pkg/config"
	"capten/pkg/types"
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"

	"capten/pkg/agent/agentpb"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
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

func SaveAppConfigsOnAgent(client agentpb.AgentClient, configDir string) error {
	appConfigs, err := readAppConfigs(configDir)
	if err != nil {
		return err
	}

	for _, appConfig := range appConfigs {
		data, err := appConfig.ToSyncAppData()
		if err != nil {
			clog.Logger.Errorf("Err while making SyncAppRequest: %v for release: %v", err, appConfig.ReleaseName)
		}
		res, err := client.SyncApp(context.TODO(), &agentpb.SyncAppRequest{Data: &data})
		if err != nil {
			clog.Logger.Errorf("Err while receiving SyncAppResponse: %v for release: %v", err, appConfig.ReleaseName)
		}
		if res.Status != agentpb.StatusCode_OK {
			clog.Logger.Errorf("Response message: %v for release: %v", res.GetStatusMessage(), appConfig.ReleaseName)
		}
	}
	return nil
}

func readAppConfigs(configDir string) (ret []types.AppConfig, err error) {

	err = filepath.Walk(configDir, func(path string, info os.FileInfo, er error) error {
		if er != nil {
			return errors.Wrapf(er, "in file: %v", path)
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		byt, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "in file: %v", path)
		}
		var appConfig types.AppConfig
		if err := yaml.NewDecoder(bytes.NewBuffer(byt)).Decode(&appConfig); err != nil {
			return errors.Wrapf(err, "in file: %v", path)
		}
		ret = append(ret, appConfig)
		return nil
	})

	return
}

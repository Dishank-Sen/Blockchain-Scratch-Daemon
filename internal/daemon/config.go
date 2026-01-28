package daemon

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/quic-go/quic-go"
)

func getConfig(addr string) (node.Config, error){
	tlsCfg, err := getTlsConfig()
	if err != nil{
		return node.Config{}, err
	}
	quicCfg := &quic.Config{
		MaxIdleTimeout: constants.MaxIdleTimeout,
		KeepAlivePeriod: constants.KeepAlivePeriod,
	}

	return node.Config{
		ListenAddr: addr,
		TlsConfig: tlsCfg,
		QuicConfig: quicCfg,
	}, nil
}

func getTlsConfig() (*tls.Config, error){
	certFilePath := os.Getenv("TLS_CERT_PATH")
	keyFilePath  := os.Getenv("TLS_KEY_PATH")

	if certFilePath == "" || keyFilePath == "" {
		logger.Error("TLS_CERT_PATH or TLS_KEY_PATH not set")
		return &tls.Config{}, fmt.Errorf("TLS_CERT_PATH and TLS_KEY_PATH are not set in environment")
	}

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil{
		logger.Debug("error in tls")
		logger.Error(err.Error())
		return &tls.Config{}, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quicnode"},
	}
	return tlsConfig, nil
}
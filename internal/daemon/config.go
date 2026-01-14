package daemon

import (
	"crypto/tls"
	"os"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/utils/logger"
	"github.com/Dishank-Sen/quicnode/node"
	"github.com/quic-go/quic-go"
)

func getConfig(addr string) node.Config{
	tlsCfg := getTlsConfig()
	
	quicCfg := &quic.Config{
		MaxIdleTimeout: constants.QuicTimeout,
	}

	return node.Config{
		ListenAddr: addr,
		TlsConfig: tlsCfg,
		QuicConfig: quicCfg,
	}
}

func getTlsConfig() *tls.Config{
	certFilePath := os.Getenv("TLS_CERT_PATH")
	keyFilePath  := os.Getenv("TLS_KEY_PATH")

	if certFilePath == "" || keyFilePath == "" {
		logger.Error("TLS_CERT_PATH or TLS_KEY_PATH not set")
		return nil
	}

	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil{
		logger.Debug("error in tls")
		logger.Error(err.Error())
		return nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quicnode"},
	}
	return tlsConfig
}
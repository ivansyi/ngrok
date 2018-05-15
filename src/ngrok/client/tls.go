package client

import (
	_ "crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"ngrok/client/assets"
)

func LoadTLSConfig(rootCertPaths []string) (*tls.Config, error) {
	pool := x509.NewCertPool()

	for _, certPath := range rootCertPaths {
		rootCrt, err := assets.Asset(certPath)
		if err != nil {
			return nil, err
		}

		pemBlock, _ := pem.Decode(rootCrt)
		if pemBlock == nil {
			return nil, fmt.Errorf("Bad PEM data")
		}

		certs, err := x509.ParseCertificates(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		pool.AddCert(certs[0])
	}

	//add self-signed CA:
	caCrt, err := ioutil.ReadFile("/etc/ssl/cert.pem")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)

	return &tls.Config{RootCAs: pool}, nil
}

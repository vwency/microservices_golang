package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

func NewTLSCredentials(logger *zap.Logger) (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile("tls/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read root CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	serverCert, err := tls.LoadX509KeyPair("tls/db_server.crt", "tls/db_server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate and key: %w", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		ServerName:   "localhost",

		ClientSessionCache: tls.NewLRUClientSessionCache(128),

		SessionTicketsDisabled: false,
	}

	return credentials.NewTLS(config), nil
}

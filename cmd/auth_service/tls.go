package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func NewTLSCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("tls/ca.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read root CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	serverCert, err := tls.LoadX509KeyPair("tls/auth_client.crt", "tls/auth_client.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate and key: %w", err)
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          certPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(128),
		MinVersion:         tls.VersionTLS13, // TLS 1.3 быстрее, если сервер поддерживает
	}

	return credentials.NewTLS(config), nil
}

package shared

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/credentials"
)

func LoadClientTLSCredentials() (credentials.TransportCredentials, error) {
	caPath, err := resolveCertPath("ca-cert.pem")
	if err != nil {
		return nil, err
	}

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("reading CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("invalid CA certificate at %s", caPath)
	}

	return credentials.NewTLS(&tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}), nil
}

func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	serviceName := resolveServiceName()
	certName := serviceName + "-cert.pem"
	keyName := serviceName + "-key.pem"

	certPath, err := resolveCertPath(certName)
	if err != nil {
		return nil, err
	}

	keyPath, err := resolveCertPath(keyName)
	if err != nil {
		return nil, err
	}

	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("loading server certificate: %w", err)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS12,
	}), nil
}

func resolveServiceName() string {
	if serviceName := strings.TrimSpace(os.Getenv("SERVICE_NAME")); serviceName != "" {
		return serviceName
	}

	wd, err := os.Getwd()
	if err == nil {
		name := filepath.Base(wd)
		if name != "" && name != "." {
			return name
		}
	}

	return "user-service"
}

func resolveCertPath(filename string) (string, error) {
	serviceName := resolveServiceName()
	candidates := make([]string, 0, 16)

	for _, envName := range []string{"TLS_CERT_PATH", "TLS_KEY_PATH", "TLS_CA_CERT_PATH"} {
		value := strings.TrimSpace(os.Getenv(envName))
		if value == "" {
			continue
		}
		if filepath.Base(value) == filename {
			candidates = append(candidates, value)
		}
	}

	candidates = append(candidates,
		filepath.Clean(filename),
		filepath.Join(".", "cert", filename),
		filepath.Join(".", filename),
		filepath.Join("..", "certs", filename),
		filepath.Join("..", "..", "certs", filename),
		filepath.Join("..", "..", "..", "services", "certs", filename),
		filepath.Join("..", "services", "certs", filename),
		filepath.Join("..", serviceName, "cert", filename),
		filepath.Join("..", "services", serviceName, "cert", filename),
		filepath.Join("..", "..", "services", serviceName, "cert", filename),
		filepath.Join("..", "..", "..", "services", serviceName, "cert", filename),
	)

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		path, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		if _, statErr := os.Stat(path); statErr == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("certificate %s not found in the expected locations", filename)
}

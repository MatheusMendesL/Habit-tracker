package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	pb "shared/pb/user"
	"user-service/db"
	"user-service/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.Println("Starting server...")

	if err := startServer(); err != nil {
		panic(err)
	}
}

func startServer() (err error) {
	list, err := net.Listen("tcp", ":8080")

	if err != nil {
		return err
	}
	dbConn, queries, err := db.Conn()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tlsCredentials, err := loadTLScreadentials()

	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(
			grpcZap.UnaryServerInterceptor(logger),
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	if err := grpcServer.Serve(list); err != nil {
		return err
	}

	return nil
}

func loadTLScreadentials() (credentials.TransportCredentials, error) {
	pemClientCA, err := os.ReadFile("./cert/ca-cert.pem")

	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("error to add the certificate")
	}

	serverCert, err := tls.LoadX509KeyPair(
		"./cert/user-service-cert.pem",
		"./cert/user-service-key.pem",
	)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		/*ClientAuth:   tls.RequireAndVerifyClientCert,*/
		ClientCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

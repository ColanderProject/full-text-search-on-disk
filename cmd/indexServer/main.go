package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/ColanderProject/full-text-search-on-disk/endpoint"
	"github.com/ColanderProject/full-text-search-on-disk/search"
	"github.com/ColanderProject/full-text-search-on-disk/transport"
	"github.com/ColanderProject/search-protobuf/pb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const cACertificateIndex int = 1
const serverCertificateIndex int = 2
const serverKeyIndex int = 3
const listenerAddrIndex int = 4

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// create grpc three-tier service
	searchservice := search.NewService(logger)
	searchpoint := endpoint.MakeEndpoints(searchservice)
	grpcServer := transport.NewGRPCServer(searchpoint, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", os.Args[listenerAddrIndex])
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			logger.Log("Cannot load TLS credentials: ", err)
		}

		abs, err := filepath.Abs("./main.go")
		if err == nil {
			fmt.Println("Absolute:", abs)
		}

		// create a base server with TLS authentication configured
		baseServer := grpc.NewServer(
			grpc.Creds(tlsCredentials),
		)
		// register defined grpc three-tier service on base server
		pb.RegisterSearchServiceServer(baseServer, grpcServer)
		level.Info(logger).Log("msg", "Server started successfully!!!")
		// add defined listener to grpc server
		baseServer.Serve(grpcListener)
		//need to execute grpc handler to launch DB as well.
		//or creating a new exe to launch this grpc handler for DB.
	}()

	level.Error(logger).Log("exit", <-errs)
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemClientCA, readErr := ioutil.ReadFile(os.Args[cACertificateIndex])
	if readErr != nil {
		return nil, readErr
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(os.Args[serverCertificateIndex], os.Args[serverKeyIndex])
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

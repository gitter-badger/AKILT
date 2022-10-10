package consumer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/Xart3mis/AKILT/Client/lib/bundles"
	"github.com/Xart3mis/AKILTC/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Init(address string) (pb.ConsumerClient, error) {
	// certFile, err := filepath.Abs("./../../../cert/cert.pem")
	// if err != nil {
	// 	return nil, err
	// }

	creds, err := loadTLSCredentials()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return pb.NewConsumerClient(conn), nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	bundles.WriteCaCertPem()
	pemServerCA, err := ioutil.ReadFile("./ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func SubscribeOnScreenText(ctx context.Context, client pb.ConsumerClient, cid string) (pb.Consumer_SubscribeOnScreenTextClient, error) {
	receiver, err := client.SubscribeOnScreenText(ctx, &pb.ClientDataRequest{ClientId: cid})
	if err != nil {
		return nil, err
	}

	return receiver, nil
}

package utils

import (
	"log"
	"os"

	"google.golang.org/grpc"
)

func ConnectGRPCEndPoint(addr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func Getenv(envName, defaultValue string) string {
	value := os.Getenv(envName)

	if value == "" {
		value = defaultValue
	}

	return value
}

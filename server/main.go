package main

import (
	"context"

	"flag"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/YAWAL/GetMeConf/api"

	"fmt"

	"net"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "localhost", "Server host")
	port = flag.Int("port", 8081, "Server port")
)

type configServer struct {
	configCache *cache.Cache
}

type configExample struct {
	BField bool   `json:"b_field"`
	IField int    `json:"i_field"`
	SField string `json:"s_field"`
}

func getFromFile(info *pb.ConfigInfo) ([]byte, error) {
	raw, err := ioutil.ReadFile(info.ConfigPath)
	if err != nil {
		return nil, err
	}
	//var c []byte
	//err = json.Unmarshal(raw, &c)
	//if err != nil {
	//	return nil, err
	//}
	return raw, nil
}

func (s *configServer) SearchConfig(ctx context.Context, configInfo *pb.ConfigInfo) (*pb.ConfigInfo, error) {
	return &pb.ConfigInfo{}, nil
}

func (s *configServer) GetConfig(ctx context.Context, configInfo *pb.ConfigInfo) (*pb.Config, error) {
	if config, found := s.configCache.Get(configInfo.ConfigId); found {
		return config.(*pb.Config), nil
	}
	log.Println(configInfo.ConfigPath)
	raw, err := getFromFile(configInfo)
	if err != nil {
		return nil, err
	}
	config := &pb.Config{raw}
	s.configCache.Set(configInfo.ConfigId, config, cache.DefaultExpiration)

	return config, nil

}

func main() {
	flag.Parse()

	//Secure
	//cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	//if err != nil {
	//	log.Fatal("filed to load key pair: ", err)
	//}
	//
	//serverConf := &tls.Config{Certificates: []tls.Certificate{cer}}
	//lis, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port), serverConf)
	//if err != nil {
	//	log.Fatal("filed to listen: ", err)
	//}
	//Insecure

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started at %s:%d", *host, *port)

	grpcServer := grpc.NewServer()

	cache := cache.New(5*time.Minute, 10*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: cache})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("filed to serve: %v", err)
	}
}

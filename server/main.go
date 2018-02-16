package main

import (
	"flag"
	"fmt"

	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	"sync"

	"bytes"
	"encoding/gob"

	pb "github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

var (
	host = flag.String("host", "localhost", "Server host")
	port = flag.Int("port", 8081, "Server port")
)

type configServer struct {
	configCache *cache.Cache
	mut         *sync.Mutex
}

func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {
	s.mut.Lock()
	configResponse, found := s.configCache.Get(nameRequest.ConfigName)
	s.mut.Unlock()
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}

	res, err := database.GetConfigByNameFromDB(nameRequest.ConfigName, nameRequest.ConfigType)
	if err != nil {
		return nil, err
	}
	byteRes, err := getBytes(res)
	if err != nil {
		return nil, err
	}
	configResponse = &pb.GetConfigResponce{byteRes}
	s.mut.Lock()
	s.configCache.Set(nameRequest.ConfigName, configResponse, cache.DefaultExpiration)
	s.mut.Unlock()
	return configResponse.(*pb.GetConfigResponce), nil
}
func (s *configServer) GetConfigsByType(typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeServer) error {

	res, err := database.GetConfigsByTypeFromDB(typeRequest.ConfigType)

	if err != nil {
		return err
	}
	for _, v := range res {
		byteRes, _ := getBytes(v)
		if err := stream.Send(&pb.GetConfigResponce{byteRes}); err != nil {
			return err
		}
	}
	return nil
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func main() {
	flag.Parse()

	cfg, err := database.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read config from file with error : %v", err)
	}
	if err = database.InitPostgresDB(*cfg); err != nil {
		log.Fatal(err)
	}

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started at %s:%d", *host, *port)

	grpcServer := grpc.NewServer()

	cache := cache.New(5*time.Minute, 10*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: cache, mut: &sync.Mutex{}})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("filed to serve: %v", err)
	}
}

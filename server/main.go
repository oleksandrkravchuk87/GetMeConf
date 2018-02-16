package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"errors"

	pb "github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
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
	byteRes, err := json.Marshal(res)
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

	switch typeRequest.ConfigType {
	case "mongodb":
		res, err := database.GetMongoDBConfigs()
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	case "tempconfig":
		res, err := database.GetTempConfigs()
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	case "tsconfig":
		res, err := database.GetTsconfigs()
		if err != nil {
			return err
		}
		for _, v := range res {
			if err = marshalAndSend(v, stream); err != nil {
				return err
			}
		}
	default:
		log.Print("unexpacted type")
		return errors.New("unexpacted type")
	}
	return nil
}

func marshalAndSend(results interface{}, stream pb.ConfigService_GetConfigsByTypeServer) error {
	byteRes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return stream.Send(&pb.GetConfigResponce{byteRes})
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"errors"

	pb "github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type configServer struct {
	configCache *cache.Cache
	mut         *sync.Mutex
	db          *gorm.DB
}

var databaseGetConfigByNameFromDB = database.GetConfigByNameFromDB
var databaseGetMongoDBConfigs = database.GetMongoDBConfigs
var databaseGetTempConfigs = database.GetTempConfigs
var databaseGetTsconfigs = database.GetTsconfigs
var databaseSaveConfigToDB = database.SaveConfigToDB
var databaseDeleteConfigFromDB = database.DeleteConfigFromDB

//GetConfigByName returns one config in GetConfigResponce message
func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {
	s.mut.Lock()
	configResponse, found := s.configCache.Get(nameRequest.ConfigName)
	s.mut.Unlock()
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}

	res, err := databaseGetConfigByNameFromDB(nameRequest.ConfigName, nameRequest.ConfigType, s.db)
	if err != nil {
		return nil, err
	}
	byteRes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	configResponse = &pb.GetConfigResponce{Config: byteRes}
	s.mut.Lock()
	s.configCache.Set(nameRequest.ConfigName, configResponse, cache.DefaultExpiration)
	s.mut.Unlock()
	return configResponse.(*pb.GetConfigResponce), nil
}

//GetConfigByName streams configs as GetConfigResponce messages
func (s *configServer) GetConfigsByType(typeRequest *pb.GetConfigsByTypeRequest, stream pb.ConfigService_GetConfigsByTypeServer) error {
	switch typeRequest.ConfigType {
	case "mongodb":
		res, err := databaseGetMongoDBConfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			stream.Send(&pb.GetConfigResponce{Config: byteRes})
		}
	case "tempconfig":
		res, err := databaseGetTempConfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			stream.Send(&pb.GetConfigResponce{Config: byteRes})
		}
	case "tsconfig":
		res, err := databaseGetTsconfigs(s.db)
		if err != nil {
			return err
		}
		for _, v := range res {
			byteRes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			stream.Send(&pb.GetConfigResponce{Config: byteRes})
		}
	default:
		log.Print("unexpacted type")
		return errors.New("unexpacted type")
	}
	return nil
}

func (s *configServer) CreateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	s.mut.Lock()
	s.configCache.Flush()
	s.mut.Unlock()
	fmt.Println("main", config)
	response, err := databaseSaveConfigToDB(config.ConfigType, config.Config, s.db)
	if err != nil {
		return nil, err
	}
	return &pb.Responce{response}, nil
}

func (s *configServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	s.mut.Lock()
	s.configCache.Flush()
	s.mut.Unlock()
	response, err := databaseDeleteConfigFromDB(delConfigRequest.ConfigName, delConfigRequest.ConfigType, s.db)
	if err != nil {
		return nil, err
	}
	return &pb.Responce{response}, nil
}

func marshalAndSend(results interface{}, stream pb.ConfigService_GetConfigsByTypeServer) error {
	byteRes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return stream.Send(&pb.GetConfigResponce{Config: byteRes})
}

func main() {

	port := "3000"
	//port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("port is not set")
	}

	cfg, err := database.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read config from file with error : %v", err)
	}
	db, err := database.InitPostgresDB(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started at :%s", port)

	grpcServer := grpc.NewServer()

	cache := cache.New(5*time.Minute, 10*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: cache, mut: &sync.Mutex{}, db: db})
	defer grpcServer.GracefulStop()
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("filed to serve: %v", err)
	}
}

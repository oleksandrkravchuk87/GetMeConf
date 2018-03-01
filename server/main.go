package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"errors"

	"os"

	pb "github.com/YAWAL/GetMeConf/api"
	"github.com/YAWAL/GetMeConf/database"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type configServer struct {
	configCache *cache.Cache
	db          *gorm.DB
}

var databaseGetConfigByNameFromDB = database.GetConfigByNameFromDB
var databaseGetMongoDBConfigs = database.GetMongoDBConfigs
var databaseGetTempConfigs = database.GetTempConfigs
var databaseGetTsconfigs = database.GetTsconfigs
var databaseSaveConfigToDB = database.SaveConfigToDB
var databaseDeleteConfigFromDB = database.DeleteConfigFromDB
var databaseUpdateMongoDBConfigInDB = database.UpdateMongoDBConfigInDB
var databaseUpdateTempconfigInDB = database.UpdateTempConfigInDB
var databaseUpdateTsConfigInDB = database.UpdateTsConfigInDB

//GetConfigByName returns one config in GetConfigResponce message
func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {

	configResponse, found := s.configCache.Get(nameRequest.ConfigName)

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

	s.configCache.Set(nameRequest.ConfigName, configResponse, cache.DefaultExpiration)

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
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
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
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
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
			if err = stream.Send(&pb.GetConfigResponce{Config: byteRes}); err != nil {
				return err
			}
		}
	default:
		log.Print("unexpacted type")
		return errors.New("unexpacted type")
	}
	return nil
}

//CreateConfig calls the function from database package to add a new config record to the database, returns response structure containing a status message
func (s *configServer) CreateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	response, err := databaseSaveConfigToDB(config.ConfigType, config.Config, s.db)
	if err != nil {
		return nil, err
	}

	s.configCache.Flush()

	return &pb.Responce{Status: response}, nil
}

//DeleteConfig removes config records from the database. If successful, returns the amount of deleted records in a status message of the response structure
func (s *configServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	response, err := databaseDeleteConfigFromDB(delConfigRequest.ConfigName, delConfigRequest.ConfigType, s.db)
	if err != nil {
		return nil, err
	}

	s.configCache.Flush()

	return &pb.Responce{Status: response}, nil
}

func marshalAndSend(results interface{}, stream pb.ConfigService_GetConfigsByTypeServer) error {
	byteRes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return stream.Send(&pb.GetConfigResponce{Config: byteRes})
}

//UpdateConfig
func (s *configServer) UpdateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	var status string
	var err error
	switch config.ConfigType {
	case "mongodb":
		status, err = databaseUpdateMongoDBConfigInDB(config.Config, s.db)
		if err != nil {
			return nil, err
		}
	case "tempconfig":
		status, err = databaseUpdateTempconfigInDB(config.Config, s.db)
		if err != nil {
			return nil, err
		}
	case "tsconfig":
		status, err = databaseUpdateTsConfigInDB(config.Config, s.db)
		if err != nil {
			return nil, err
		}
	default:
		log.Print("unexpacted type")
		return nil, errors.New("unexpacted type")
	}

	s.configCache.Flush()

	return &pb.Responce{Status: status}, nil
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	cfg, err := database.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read config from file with error : %v", err)
	}
	db, err := database.InitPostgresDB(*cfg)
	if err != nil {
		log.Fatalf("failed to init postgres db: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started at :%s", port)

	grpcServer := grpc.NewServer()

	configCache := cache.New(5*time.Minute, 10*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: configCache, db: db})
	defer grpcServer.GracefulStop()
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("filed to serve: %v", err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"errors"

	"os"

	pb "github.com/YAWAL/GetMeConfAPI/api"

	"strconv"

	"os/signal"
	"syscall"

	"github.com/YAWAL/GetMeConf/entities"
	"github.com/YAWAL/GetMeConf/repository"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	defaultPort                 = "3000"
	defaultCacheExpirationTime  = 5
	defaultCacheCleanupInterval = 10
)

const (
	mongodb    = "mongodb"
	tempconfig = "tempconfig"
	tsconfig   = "tsconfig"
)

type configServer struct {
	configCache       *cache.Cache
	mongoDBConfigRepo repository.MongoDBConfigRepo
	tempConfigRepo    repository.TempConfigRepo
	tsConfigRepo      repository.TsConfigRepo
}

//GetConfigByName returns one config in GetConfigResponce message
func (s *configServer) GetConfigByName(ctx context.Context, nameRequest *pb.GetConfigByNameRequest) (*pb.GetConfigResponce, error) {

	configResponse, found := s.configCache.Get(nameRequest.ConfigName)
	if found {
		return configResponse.(*pb.GetConfigResponce), nil
	}
	var err error
	var res entities.ConfigInterface

	switch nameRequest.ConfigType {
	case mongodb:
		res, err = s.mongoDBConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		res, err = s.tempConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		res, err = s.tsConfigRepo.Find(nameRequest.ConfigName)
		if err != nil {
			return nil, err
		}
	default:
		log.Print("unexpected type")
		return nil, errors.New("unexpected type")
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
	case mongodb:
		res, err := s.mongoDBConfigRepo.FindAll()
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
	case tempconfig:
		res, err := s.tempConfigRepo.FindAll()
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
	case tsconfig:
		res, err := s.tsConfigRepo.FindAll()
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
		log.Print("unexpected type")
		return errors.New("unexpected type")
	}
	return nil
}

//CreateConfig calls the function from database package to add a new config record to the database, returns response structure containing a status message
func (s *configServer) CreateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	switch config.ConfigType {
	case mongodb:
		configStr := entities.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		response, err := s.mongoDBConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tempconfig:
		configStr := entities.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		response, err := s.tempConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil

	case tsconfig:
		configStr := entities.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		response, err := s.tsConfigRepo.Save(&configStr)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		log.Print("unexpected type")
		return nil, errors.New("unexpected type")
	}
}

//DeleteConfig removes config records from the database. If successful, returns the amount of deleted records in a status message of the response structure
func (s *configServer) DeleteConfig(ctx context.Context, delConfigRequest *pb.DeleteConfigRequest) (*pb.Responce, error) {
	switch delConfigRequest.ConfigType {
	case mongodb:
		response, err := s.mongoDBConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tempconfig:
		response, err := s.tempConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	case tsconfig:
		response, err := s.tsConfigRepo.Delete(delConfigRequest.ConfigName)
		if err != nil {
			return nil, err
		}
		s.configCache.Flush()
		return &pb.Responce{Status: response}, nil
	default:
		log.Print("unexpected type")
		return nil, errors.New("unexpected type")
	}
}

//UpdateConfig
func (s *configServer) UpdateConfig(ctx context.Context, config *pb.Config) (*pb.Responce, error) {
	var status string
	switch config.ConfigType {
	case mongodb:
		configStr := entities.Mongodb{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		status, err = s.mongoDBConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	case tempconfig:
		configStr := entities.Tempconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		status, err = s.tempConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	case tsconfig:
		configStr := entities.Tsconfig{}
		err := json.Unmarshal(config.Config, &configStr)
		if err != nil {
			log.Printf("unmarshal config err: %v", err)
			return nil, err
		}
		status, err = s.tsConfigRepo.Update(&configStr)
		if err != nil {
			return nil, err
		}
	default:
		log.Print("unexpected type")
		return nil, errors.New("unexpected type")
	}
	s.configCache.Flush()
	return &pb.Responce{Status: status}, nil
}

func main() {

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		log.Println("error during reading env. variable, default value is used")
		port = defaultPort
	}
	cacheExpirationTime, err := strconv.Atoi(os.Getenv("CACHE_EXPIRATION_TIME"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		cacheExpirationTime = defaultCacheExpirationTime
	}
	cacheCleanupInterval, err := strconv.Atoi(os.Getenv("CACHE_CLEANUP_INTERVAL"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		cacheCleanupInterval = defaultCacheCleanupInterval
	}

	dbConn, err := repository.InitPostgresDB()
	if err != nil {
		log.Fatalf("failed to init postgres db: %v", err)
	}
	mongoDBRepo := repository.MongoDBConfigRepoImpl{DB: dbConn}
	tsConfigRepo := repository.TsConfigRepoImpl{DB: dbConn}
	tempConfigRepo := repository.TempConfigRepoImpl{DB: dbConn}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started at :%s", port)

	grpcServer := grpc.NewServer()

	configCache := cache.New(time.Duration(cacheExpirationTime)*time.Minute, time.Duration(cacheCleanupInterval)*time.Minute)

	pb.RegisterConfigServiceServer(grpcServer, &configServer{configCache: configCache, mongoDBConfigRepo: &mongoDBRepo, tsConfigRepo: &tsConfigRepo, tempConfigRepo: &tempConfigRepo})

	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("shotdown signal received, exiting")
	grpcServer.GracefulStop()
}

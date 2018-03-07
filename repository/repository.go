package repository

import (
	"github.com/YAWAL/GetMeConf/entities"
)

//MongoDBConfigRepo is a repository interface for MongoDB configs
type MongoDBConfigRepo interface {
	Find(configName string) (*entities.Mongodb, error)
	FindAll() ([]entities.Mongodb, error)
	Update(config *entities.Mongodb) (string, error)
	Save(config *entities.Mongodb) (string, error)
	Delete(configName string) (string, error)
}

//TempConfigRepo is a repository interface for Tempconfigs
type TempConfigRepo interface {
	Find(configName string) (*entities.Tempconfig, error)
	FindAll() ([]entities.Tempconfig, error)
	Update(config *entities.Tempconfig) (string, error)
	Save(config *entities.Tempconfig) (string, error)
	Delete(configName string) (string, error)
}

//TsConfigRepo is a repository interface for Tsconfigs
type TsConfigRepo interface {
	Find(configName string) (*entities.Tsconfig, error)
	FindAll() ([]entities.Tsconfig, error)
	Update(config *entities.Tsconfig) (string, error)
	Save(config *entities.Tsconfig) (string, error)
	Delete(configName string) (string, error)
}

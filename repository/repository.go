package repository

import (
	"github.com/YAWAL/GetMeConf/entitys"
)

//MongoDBConfigRepo is a repository interface for MongoDB configs
type MongoDBConfigRepo interface {
	Find(configName string) (*entitys.Mongodb, error)
	FindAll() ([]entitys.Mongodb, error)
	Update(config *entitys.Mongodb) (string, error)
	Save(config *entitys.Mongodb) (string, error)
	Delete(configName string) (string, error)
}

//TempConfigRepo is a repository interface for Tempconfigs
type TempConfigRepo interface {
	Find(configName string) (*entitys.Tempconfig, error)
	FindAll() ([]entitys.Tempconfig, error)
	Update(config *entitys.Tempconfig) (string, error)
	Save(config *entitys.Tempconfig) (string, error)
	Delete(configName string) (string, error)
}

//TsConfigRepo is a repository interface for Tsconfigs
type TsConfigRepo interface {
	Find(configName string) (*entitys.Tsconfig, error)
	FindAll() ([]entitys.Tsconfig, error)
	Update(config *entitys.Tsconfig) (string, error)
	Save(config *entitys.Tsconfig) (string, error)
	Delete(configName string) (string, error)
}

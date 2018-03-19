// Package repository contains repository interfaces as well as their implementations for given databases
package repository

import (
	"github.com/YAWAL/GetMeConf/entitie"
)

//MongoDBConfigRepo is a repository interface for MongoDB configs
type MongoDBConfigRepo interface {
	Find(configName string) (*entitie.Mongodb, error)
	FindAll() ([]entitie.Mongodb, error)
	Update(config *entitie.Mongodb) (string, error)
	Save(config *entitie.Mongodb) (string, error)
	Delete(configName string) (string, error)
}

//TempConfigRepo is a repository interface for Tempconfigs
type TempConfigRepo interface {
	Find(configName string) (*entitie.Tempconfig, error)
	FindAll() ([]entitie.Tempconfig, error)
	Update(config *entitie.Tempconfig) (string, error)
	Save(config *entitie.Tempconfig) (string, error)
	Delete(configName string) (string, error)
}

//TsConfigRepo is a repository interface for Tsconfigs
type TsConfigRepo interface {
	Find(configName string) (*entitie.Tsconfig, error)
	FindAll() ([]entitie.Tsconfig, error)
	Update(config *entitie.Tsconfig) (string, error)
	Save(config *entitie.Tsconfig) (string, error)
	Delete(configName string) (string, error)
}

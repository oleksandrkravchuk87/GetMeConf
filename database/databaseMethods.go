package database

import (
	"fmt"
	"time"

	"log"

	"errors"
	"strings"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/gormigrate.v1"
)

var factory map[string]PersistedData

func initConfigDataMap() {
	if factory != nil {
		return
	}
	factory = map[string]PersistedData{
		"mongodb":    {ConfigType: new(Mongodb), IDField: "domain"},
		"tempconfig": {ConfigType: new(Tempconfig), IDField: "host"},
		"tsconfig":   {ConfigType: new(Tsconfig), IDField: "module"},
	}
}

//InitPostgresDB initiates database connection using configuration file
func InitPostgresDB(cfg PostgresConfig) (db *gorm.DB, err error) {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Dbhost, cfg.Dbport, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err = gorm.Open("postgres", dbInfo)

	if err != nil {
		log.Printf("error during connection to postgres database has occurred: %v", err)
		return nil, err
	}
	db.DB().SetMaxOpenConns(cfg.MaxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(cfg.MbConnMaxLifetimeMinutes))
	log.Printf("connection to postgres database has been established")

	initConfigDataMap()

	if err = gormMigrate(db); err != nil {
		log.Printf("error during migration: %v", err)
		return nil, err
	}

	return db, nil
}

func gormMigrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "Initial",
			Migrate: func(tx *gorm.DB) error {
				type Mongodb struct {
					//gorm.Model
					Domain  string `gorm:"primary_key"`
					Mongodb bool
					Host    string
					Port    string
				}
				type Tsconfig struct {
					//gorm.Model
					Module    string `gorm:"primary_key"`
					Target    string
					SourceMap bool
					Excluding int
				}
				type Tempconfig struct {
					//gorm.Model
					RestApiRoot    string `gorm:"primary_key"`
					Host           string
					Port           string
					Remoting       string
					LegasyExplorer bool
				}
				return tx.AutoMigrate(&Mongodb{}, &Tsconfig{}, &Tempconfig{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("mongodbs", "tsconfigs", "tempconfigs").Error
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		log.Fatalf("could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
	return err
}

//GetConfigByNameFromDB searches a config in database using the type of the config and a unique name
func GetConfigByNameFromDB(confName string, confType string, db *gorm.DB) (ConfigInterface, error) {
	cType := strings.ToLower(confType)
	configStruct, ok := factory[cType]
	if !ok {
		return nil, errors.New("unexpected config type")
	}
	result := configStruct.ConfigType
	err := db.Where(configStruct.IDField+" = ?", confName).Find(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

//GetMongoDBConfigs searches for all Mongodb configs in database
func GetMongoDBConfigs(db *gorm.DB) ([]Mongodb, error) {
	var confSlice []Mongodb
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//GetTempConfigs searches for all TempConfig in database
func GetTempConfigs(db *gorm.DB) ([]Tempconfig, error) {
	var confSlice []Tempconfig
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//GetTsconfigs searches for all Tsconfigs in database
func GetTsconfigs(db *gorm.DB) ([]Tsconfig, error) {
	var confSlice []Tsconfig
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//SaveConfigToDB saves new config record to the database
func SaveConfigToDB(confType string, config []byte, db *gorm.DB) (string, error) {
	cType := strings.ToLower(confType)
	configStruct, ok := factory[cType]
	if !ok {
		return "", errors.New("unexpected config type")
	}
	configTypeStr := configStruct.ConfigType
	err := json.Unmarshal(config, configTypeStr)
	if err != nil {
		log.Printf("unmarshal config err: %v", err)
		return "", err
	}
	err = db.Create(configTypeStr).Error
	if err != nil {
		log.Printf("error during saving to database: %v", err)
		return "", err
	}
	return "OK", nil
}

//DeleteConfigFromDB removes config record from database
func DeleteConfigFromDB(confName, confType string, db *gorm.DB) (string, error) {
	cType := strings.ToLower(confType)
	configStruct, ok := factory[cType]
	if !ok {
		return "", errors.New("unexpected config type")
	}
	result := configStruct.ConfigType
	rowsAffected := db.Delete(result, configStruct.IDField+" = ?", confName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//UpdateMongoDBConfigInDB updates a record in database, rewriting the fields if string fields are not empty
func UpdateMongoDBConfigInDB(configBytes []byte, db *gorm.DB) (string, error) {
	var newConfig, persistedConfig Mongodb
	err := json.Unmarshal(configBytes, &newConfig)
	if err != nil {
		log.Printf("unmarshal config err: %v", err)
		return "", err
	}
	err = db.Where("domain = ?", newConfig.Domain).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" {
		err = db.Model(&persistedConfig).Where("domain = ?", newConfig.Domain).Update(Mongodb{Mongodb: newConfig.Mongodb, Port: newConfig.Port, Host: newConfig.Host}).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")

}

//UpdateTempConfigInDB updates a record in database, rewriting the fields if string fields are not empty
func UpdateTempConfigInDB(configBytes []byte, db *gorm.DB) (string, error) {
	var newConfig, persistedConfig Tempconfig
	err := json.Unmarshal(configBytes, &newConfig)
	if err != nil {
		log.Printf("unmarshal config err: %v", err)
		return "", err
	}
	err = db.Where("rest_api_root = ?", newConfig.RestApiRoot).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" && newConfig.Remoting != "" {
		err = db.Model(&persistedConfig).Where("rest_api_root = ?", newConfig.RestApiRoot).Update(Tempconfig{Host: newConfig.Host, Port: newConfig.Port, Remoting: newConfig.Remoting, LegasyExplorer: newConfig.LegasyExplorer}).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

//UpdateTsConfigInDB updates a record in database, rewriting the fields if string fields are not empty
func UpdateTsConfigInDB(configBytes []byte, db *gorm.DB) (string, error) {
	var newConfig, persistedConfig Tsconfig
	err := json.Unmarshal(configBytes, &newConfig)
	if err != nil {
		log.Printf("unmarshal config err: %v", err)
		return "", err
	}
	err = db.Where("module = ?", newConfig.Module).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Target != "" {
		err = db.Model(&persistedConfig).Where("module = ?", newConfig.Module).Update(Tsconfig{Target: newConfig.Target, SourceMap: newConfig.SourceMap, Excluding: newConfig.Excluding}).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

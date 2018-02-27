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
		"mongodb":    PersistedData{ConfigType: new(Mongodb), IDField: "domain"},
		"tempconfig": PersistedData{ConfigType: new(Tempconfig), IDField: "host"},
		"tsconfig":   PersistedData{ConfigType: new(Tsconfig), IDField: "module"},
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

	db.LogMode(true)
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "Initial",
			Migrate: func(tx *gorm.DB) error {
				type Mongodb struct {
					gorm.Model
					Domain  string
					Mongodb bool
					Host    string
					Port    string
				}
				type Tsconfig struct {
					gorm.Model
					Module    string
					Target    string
					SourseMap bool
					Exclude   int
				}
				type Tempconfig struct {
					gorm.Model
					RestApiRoot    string
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

//GetConfigByNameFromDB(confName string, confType string) searches a config in database using the type of the config and a unique name
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

//GetMongoDBConfigs(db *gorm.DB) searches for all Mongodb configs in database
func GetMongoDBConfigs(db *gorm.DB) ([]Mongodb, error) {
	var confSlice []Mongodb
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//GetTempConfigs(db *gorm.DB) searches for all TempConfig in database
func GetTempConfigs(db *gorm.DB) ([]Tempconfig, error) {
	var confSlice []Tempconfig
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//GetTsconfigs(db *gorm.DB) searches for all Tsconfigs in database
func GetTsconfigs(db *gorm.DB) ([]Tsconfig, error) {
	var confSlice []Tsconfig
	err := db.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//SaveConfigToDB(confType string, config []byte, db *gorm.DB) saves new config record to the database
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

//DeleteConfigFromDB
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

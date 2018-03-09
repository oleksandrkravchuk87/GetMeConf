package repository

import (
	"fmt"
	"log"
	"strconv"

	"time"

	"os"

	"errors"

	"github.com/YAWAL/GetMeConf/entities"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/gormigrate.v1"
)

var (
	defaultDbHost                   = "horton.elephantsql.com"
	defaultDbPort                   = "5432"
	defaultDbUser                   = "dlxifkbx"
	defaultDbPassword               = "L7Cey-ucPY4L3T6VFlFdNykNE4jO0VjV"
	defaultDbName                   = "dlxifkbx"
	defaultMaxOpenedConnectionsToDb = 5
	defaultMaxIdleConnectionsToDb   = 0
	defaultmbConnMaxLifetimeMinutes = 30
)

//MongoDBConfigRepoImpl represents an implementation of a MongoDB configs repository
type MongoDBConfigRepoImpl struct {
	DB *gorm.DB
}

//TsConfigRepoImpl represents an implementation of a Tsconfigs repository
type TsConfigRepoImpl struct {
	DB *gorm.DB
}

//TempConfigRepoImpl represents an implementation of a Tempconfigs repository
type TempConfigRepoImpl struct {
	DB *gorm.DB
}

//NewMongoDBConfigRepo returns a new MongoDB configs repository
func NewMongoDBConfigRepo(db *gorm.DB) MongoDBConfigRepo {
	return &MongoDBConfigRepoImpl{
		DB: db,
	}
}

//NewTempConfigRepo returns a new Tempconfigs repository
func NewTempConfigRepo(db *gorm.DB) TempConfigRepo {
	return &TempConfigRepoImpl{
		DB: db,
	}
}

//NewTsConfigRepo returns a new TsConfig repository
func NewTsConfigRepo(db *gorm.DB) TsConfigRepo {
	return &TsConfigRepoImpl{
		DB: db,
	}
}

//InitPostgresDB initiates database connection using environmental variables
func InitPostgresDB() (db *gorm.DB, err error) {
	dbHost := os.Getenv("PDB_HOST")
	if dbHost == "" {
		log.Println("error during reading env. variable, default value is used")
		dbHost = defaultDbHost
	}
	dbPort := os.Getenv("PDB_PORT")
	if dbPort == "" {
		log.Println("error during reading env. variable, default value is used")
		dbPort = defaultDbPort
	}
	dbUser := os.Getenv("PDB_USER")
	if dbUser == "" {
		log.Println("error during reading env. variable, default value is used")
		dbUser = defaultDbUser
	}
	dbPassword := os.Getenv("PDB_PASSWORD")
	if dbPassword == "" {
		log.Println("error during reading env. variable, default value is used")
		dbPassword = defaultDbPassword
	}
	dbName := os.Getenv("PDB_NAME")
	if dbName == "" {
		log.Println("error during reading env. variable, default value is used")
		dbName = defaultDbName
	}
	maxOpenedConnectionsToDb, err := strconv.Atoi(os.Getenv("MAX_OPENED_CONNECTIONS_TO_DB"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		maxOpenedConnectionsToDb = defaultMaxOpenedConnectionsToDb
	}
	maxIdleConnectionsToDb, err := strconv.Atoi(os.Getenv("MAX_IDLE_CONNECTIONS_TO_DB"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		maxIdleConnectionsToDb = defaultMaxIdleConnectionsToDb
	}
	mbConnMaxLifetimeMinutes, err := strconv.Atoi(os.Getenv("MB_CONN_MAX_LIFETIME_MINUTES"))
	if err != nil {
		log.Printf("error during reading env. variable: %v, default value is used", err)
		mbConnMaxLifetimeMinutes = defaultmbConnMaxLifetimeMinutes
	}
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = gorm.Open("postgres", dbInfo)

	if err != nil {
		log.Printf("error during connection to postgres database has occurred: %v", err)
		return nil, err
	}

	db.DB().SetMaxOpenConns(maxOpenedConnectionsToDb)
	db.DB().SetMaxIdleConns(maxIdleConnectionsToDb)
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(mbConnMaxLifetimeMinutes))
	log.Printf("connection to postgres database has been established")

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

//Find returns a config record from database using the unique name
func (r *MongoDBConfigRepoImpl) Find(configName string) (*entities.Mongodb, error) {
	result := entities.Mongodb{}
	err := r.DB.Where("domain = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *MongoDBConfigRepoImpl) FindAll() ([]entities.Mongodb, error) {
	var confSlice []entities.Mongodb
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *MongoDBConfigRepoImpl) Save(config *entities.Mongodb) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		log.Printf("error during saving to database: %v", err)
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *MongoDBConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entities.Mongodb{}, "domain = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *MongoDBConfigRepoImpl) Update(newConfig *entities.Mongodb) (string, error) {
	var persistedConfig entities.Mongodb
	err := r.DB.Where("domain = ?", newConfig.Domain).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" {
		err = r.DB.Exec("UPDATE mongodbs SET mongodb = ?, port = ?, host = ? WHERE domain = ?", strconv.FormatBool(newConfig.Mongodb), newConfig.Port, newConfig.Host, persistedConfig.Domain).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

//Find returns a config record from database using the unique name
func (r *TempConfigRepoImpl) Find(configName string) (*entities.Tempconfig, error) {
	result := entities.Tempconfig{}
	err := r.DB.Where("rest_api_root = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *TempConfigRepoImpl) FindAll() ([]entities.Tempconfig, error) {
	var confSlice []entities.Tempconfig
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *TempConfigRepoImpl) Save(config *entities.Tempconfig) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		log.Printf("error during saving to database: %v", err)
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *TempConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entities.Tempconfig{}, "rest_api_root = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *TempConfigRepoImpl) Update(newConfig *entities.Tempconfig) (string, error) {
	var persistedConfig entities.Tempconfig
	err := r.DB.Where("rest_api_root = ?", newConfig.RestApiRoot).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Host != "" && newConfig.Port != "" && newConfig.Remoting != "" {
		err = r.DB.Exec("UPDATE tempconfigs SET remoting = ?, port = ?, host = ?, legasy_explorer = ? WHERE rest_api_root = ?", newConfig.Remoting, newConfig.Port, newConfig.Host, strconv.FormatBool(newConfig.LegasyExplorer), persistedConfig.RestApiRoot).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

//Find returns a config record from database using the unique name
func (r *TsConfigRepoImpl) Find(configName string) (*entities.Tsconfig, error) {
	result := entities.Tsconfig{}
	err := r.DB.Where("module = ?", configName).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//FindAll returns all config record of one type from database
func (r *TsConfigRepoImpl) FindAll() ([]entities.Tsconfig, error) {
	var confSlice []entities.Tsconfig
	err := r.DB.Find(&confSlice).Error
	if err != nil {
		return nil, err
	}
	return confSlice, nil
}

//Save saves new config record to the database
func (r *TsConfigRepoImpl) Save(config *entities.Tsconfig) (string, error) {
	err := r.DB.Create(config).Error
	if err != nil {
		log.Printf("error during saving to database: %v", err)
		return "", err
	}
	return "OK", nil
}

//Delete removes config record from database
func (r *TsConfigRepoImpl) Delete(configName string) (string, error) {
	rowsAffected := r.DB.Delete(entities.Tsconfig{}, "module = ?", configName).RowsAffected
	if rowsAffected < 1 {
		return "", errors.New("could not delete from database")
	}
	return fmt.Sprintf("deleted %d row(s)", rowsAffected), nil
}

//Update updates a record in database, rewriting the fields if string fields are not empty
func (r *TsConfigRepoImpl) Update(newConfig *entities.Tsconfig) (string, error) {
	var persistedConfig entities.Tsconfig
	err := r.DB.Where("module = ?", newConfig.Module).Find(&persistedConfig).Error
	if err != nil {
		return "", err
	}
	if newConfig.Target != "" {
		err = r.DB.Exec("UPDATE tsconfigs SET target = ?, source_map = ?, excluding = ? WHERE module = ?", newConfig.Target, strconv.FormatBool(newConfig.SourceMap), strconv.Itoa(newConfig.Excluding), persistedConfig.Module).Error
		if err != nil {
			log.Printf("error during saving to database: %v", err)
			return "", err
		}
		return "OK", nil
	}
	return "", errors.New("fields are empty")
}

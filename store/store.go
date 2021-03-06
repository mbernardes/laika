package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MEDIGO/laika/store/schema"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

var ErrNoRows = sql.ErrNoRows
var ErrTxDone = sql.ErrTxDone

type Store interface {
	GetFeatureByName(name string) (*Feature, error)
	ListFeatures() ([]*Feature, error)
	CreateFeature(feature *Feature) error
	UpdateFeature(feature *Feature) error

	GetEnvironmentByName(name string) (*Environment, error)
	ListEnvironments() ([]*Environment, error)
	CreateEnvironment(environment *Environment) error
	UpdateEnvironment(environment *Environment) error

	GetFeatureStatus(featureId int64, environmentId int64) (*FeatureStatus, error)
	ListFeatureStatus(featureId *int64, environmentId *int64) ([]*FeatureStatus, error)
	CreateFeatureStatus(featureStatus *FeatureStatus) error
	UpdateFeatureStatus(featureStatus *FeatureStatus) error

	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) error

	ListFeatureStatusHistory(featureId *int64, environmentId *int64, featureStatusId *int64) ([]*FeatureStatusHistory, error)

	Ping() error

	// Migrate migrates the database schema to the latest available version.
	Migrate() error

	// Reset removes all stored data.
	Reset() error
}

type store struct {
	db *sql.DB
}

// NewStore creates a new Store.
func NewStore(username, password, host, port, dbname string) (Store, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return &store{db}, nil
}

func (s *store) Ping() error {
	var err error

	for i := 0; i < 10; i++ {
		err = s.db.Ping()
		if err == nil {
			return nil
		}

		log.Warn("Failed to ping the database. Retry in 1s.")
		time.Sleep(time.Second)
	}

	return err
}

func (s *store) Migrate() error {
	migrations := &migrate.AssetMigrationSource{
		Asset:    schema.Asset,
		AssetDir: schema.AssetDir,
		Dir:      "store/schema",
	}

	_, err := migrate.Exec(s.db, "mysql", migrations, migrate.Up)
	return err
}

func (s *store) Reset() error {
	tables := []string{
		"feature_status_history",
		"feature_status",
		"environment",
		"feature",
		"user",
	}

	if _, err := s.db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return err
	}

	for _, table := range tables {
		if _, err := s.db.Exec("TRUNCATE TABLE " + table); err != nil {
			return err
		}
	}

	if _, err := s.db.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
		return err
	}

	return nil
}

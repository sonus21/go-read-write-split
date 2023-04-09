package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// MySqlConfig this is used to store the different MySQL configuration
type MySqlConfig struct {
	Host                 string
	Port                 int
	Database             string
	Username             string
	Password             string
	MaxConnectionRetries int
	MaxOpenConnection    int
	MaxIdleConnection    int
	ConnectionLifetime   int //in minutes
}

// MySqlDataBase opens the MySQL DB and pings the database for any unavailability like connection error, permission denied etc
func MySqlDataBase(config MySqlConfig) (*sql.DB, error) {
	dsn := getDatasourceString(config)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * time.Duration(config.ConnectionLifetime))
	db.SetMaxOpenConns(config.MaxOpenConnection)
	db.SetMaxIdleConns(config.MaxIdleConnection)

	//if MaxConnectionRetries is zero then we don't check the connection at server start
	if config.MaxConnectionRetries == 0 {
		return db, err
	}

	retries := config.MaxConnectionRetries

	// Open doesn't open a connection. Validate DSN data:
	for retries > 0 {
		err = db.Ping()
		if err == nil {
			return db, err
		}
		log.Println(fmt.Sprintf("failed pinging database server: %s, waiting 2 seconds before trying %d more times", err, retries))
		time.Sleep(time.Second * 2)
		retries--
	}
	err = fmt.Errorf("failed to connect to database even after %d attempt(s)", config.MaxConnectionRetries)
	return db, err
}

// getDatasourceString creates database datasource string
func getDatasourceString(config MySqlConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=True&rejectReadOnly=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
}

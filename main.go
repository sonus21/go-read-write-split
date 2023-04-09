package main

import (
	"database/sql"
	"github.com/sonus21/db-read-write/controller"
	"github.com/sonus21/db-read-write/pkg/database"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func getPrimaryDbConfig() database.MySqlConfig {
	// use any other mechanism like viper/toml/env file to create the cpnfig
	return database.MySqlConfig{
		Host:                 "localhost",
		Port:                 3306,
		Database:             "my_app_db",
		Username:             "my_app",
		Password:             "my_app_pass",
		MaxConnectionRetries: 3,
		MaxOpenConnection:    30,
		MaxIdleConnection:    5,
		ConnectionLifetime:   5,
	}
}

func getSecondaryDbConfig() database.MySqlConfig {
	// use any other mechanism like viper/toml/env file to create the config
	return database.MySqlConfig{
		Host:                 "localhost",
		Port:                 3306,
		Database:             "my_app_db",
		Username:             "my_app_reader",
		Password:             "my_app_reader_pass",
		MaxConnectionRetries: 3,
		MaxOpenConnection:    30,
		MaxIdleConnection:    5,
		ConnectionLifetime:   5,
	}
}

func main() {
	// get primary database configurations and create primary database
	primary, err := database.MySqlDataBase(getPrimaryDbConfig())
	defer primary.Close()

	if err != nil {
		panic("Primary database could not be created" + err.Error())
	}

	// get secondary database configurations and create secondary database
	secondary, err := database.MySqlDataBase(getSecondaryDbConfig())
	defer secondary.Close()
	if err != nil {
		panic("Secondary database could not be created" + err.Error())
	}
	// initialize database map
	databases := map[string]*sql.DB{
		database.Primary:   primary,
		database.Secondary: secondary,
	}
	// initialize database
	database.Init(databases, database.Primary)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// add database middleware
	r.Use(database.Middleware())
	r.Post("/api/v1/orders", controller.HandleCreateOrder)
	r.Get("/api/v1/orders/{orderId:[0-9-]+}", controller.OrderDetails)
	err = http.ListenAndServe(":3000", r)
	if err == nil {
		log.Fatal(err)
	}
}

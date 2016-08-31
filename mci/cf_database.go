package mci

import (
	"github.com/cloudfoundry-community/go-cfenv"
	toml "github.com/robfig/config"
	"net/url"
	"errors"
	"fmt"
)

func cfDatabase(appEnv *cfenv.App, leanoteConfig *toml.Config) error {
	var dbService *cfenv.Service
	var dbUrl *url.URL
	var err error
	dbService, err = getMongoDb(appEnv)
	if err != nil {
		return errors.New("Cannot find database service from cloudfoundry")
	}
	dbUrl, err = parseBasedUrlService(dbService)
	if err != nil {
		return err
	}
	dbHost := getStringValueFromService(dbService, "hostname")
	dbPort := getStringValueFromService(dbService, "port")
	dbUser := getStringValueFromService(dbService, "username")
	dbPwd  := getStringValueFromService(dbService, "password")
	dbName := getStringValueFromService(dbService, "name")

	c := toml.NewDefault()
	c.AddOption("", "db.host", dbHost)
	c.AddOption("", "db.port", dbPort)
	c.AddOption("", "db.dbname", dbName)
	c.AddOption("", "db.username", dbUser)
	c.AddOption("", "db.password", dbPwd)

	c.AddOption("", "db.url", dbUrl.String())
	leanoteConfig.Merge(c)

	return nil
}
func getMongoDb(appEnv *cfenv.App) (*cfenv.Service, error) {
	var service *cfenv.Service
	var services []cfenv.Service
	var err error
	services, err = appEnv.Services.WithTag("mongodb")
	if err == nil {
		return &services[0], err
	}
	service, err = appEnv.Services.WithName("mongodb-leanote")
	if err == nil {
		return service, err
	}
	return nil, err
}

func getPostgresDb(appEnv *cfenv.App) (*cfenv.Service, error) {
	var service *cfenv.Service
	var services []cfenv.Service
	var err error
	services, err = appEnv.Services.WithTag("postgres")
	if err == nil {
		return &services[0], err
	}
	services, err = appEnv.Services.WithTag("postgresql")
	if err == nil {
		return &services[0], err
	}
	service, err = appEnv.Services.WithName("postgres-mattermost")
	if err == nil {
		return service, err
	}
	return nil, err
}
func getMysqlDb(appEnv *cfenv.App) (*cfenv.Service, error) {
	var service *cfenv.Service
	var services []cfenv.Service
	var err error
	services, err = appEnv.Services.WithTag("mysql")
	if err == nil {
		return &services[0], err
	}
	service, err = appEnv.Services.WithName("mysql-mattermost")
	if err == nil {
		return service, err
	}
	return nil, err
}
func formatDataSource(dbUrl *url.URL, isPostgres bool) string {
	var postgresString = "postgres://%s:%s@%s%s?sslmode=disable&connect_timeout=10"
	var mysqlString = "%s:%s@tcp(%s)%s?charset=utf8mb4,utf8"
	var dataSource string
	if (isPostgres) {
		dataSource = postgresString
	} else {
		dataSource = mysqlString
	}
	password, _ := dbUrl.User.Password()
	return fmt.Sprintf(
		dataSource,
		dbUrl.User.Username(),
		password,
		dbUrl.Host,
		dbUrl.Path)
}
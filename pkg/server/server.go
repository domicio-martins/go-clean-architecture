package server

import (
	"context"
	"time"

	"github.com/domicio-martins/go-clean-architecture/pkg/newrelic"

	"github.com/domicio-martins/go-clean-architecture/pkg/http"
	"github.com/domicio-martins/go-clean-architecture/pkg/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Debug              int    `envconfig:"DEBUG" default:"0"`
	AppEnv             string `envconfig:"APP_ENV" default:"prod"`
	Host               string `default:"http://localhost:9000"`
	HttpAddress        string `envconfig:"HTTP_ADDRESS" default:"0.0.0.0:9001"`
	MongoAddress       string `envconfig:"MONGO_ADDRESS" default:"mongodb://localhost:27019"`
	MongoTimeout       int    `envconfig:"MONGO_TIMEOUT" default:"10"`
	MongoDatabase      string `envconfig:"MONGO_DATABASE" default:"clean_architecture"`
	NewRelicLicenseKey string `envconfig:"NEW_RELIC_LICENSE_KEY"`
	PrometheusName     string `envconfig:"PROMETHEUS_NAME" default:"clean_architecture_gin"`
	PrometheusPath     string `envconfig:"PROMETHEUS_PATH" default:"/metrics"`
}

type Server struct {
	Config     *Config
	DBConn     *mongo.Client
	DB         *mongo.Database
	HttpServer *http.Server
	NewRelic   newrelic.NewRelic

	quit chan struct{}
}

type RegisterableService interface {
	Load(*gin.RouterGroup, *Server) error
}

func New(conf *Config) (*Server, error) {
	nr := conf.NewRelicLicenseKey != ""
	conn, err := mongodb.Open(conf.MongoAddress, nr)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to mongodb")
	}

	timeout := time.Duration(conf.MongoTimeout) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err = conn.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "timeout: could not connect to mongo")
	}

	db := conn.Database(conf.MongoDatabase)
	httpServer := http.NewServer(conf.HttpAddress, conf.PrometheusName, conf.PrometheusPath)
	if conf.Debug == 1 {
		httpServer.Debug()
	}

	s := &Server{
		Config:     conf,
		DBConn:     conn,
		DB:         db,
		HttpServer: httpServer,
	}

	return s, err

}

func (s *Server) Load(group string, services ...RegisterableService) error {
	router := s.HttpServer.Router.Group(group)
	for _, service := range services {
		err := service.Load(router, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Start() {
	go s.HttpServer.Start()
	<-s.quit
}

func (s *Server) Stop() {
	s.HttpServer.Close()
	s.DBConn.Disconnect(context.TODO())
	s.quit <- struct{}{}
}

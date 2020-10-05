package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/domicio-martins/go-clean-architecture/pkg/password"

	"github.com/domicio-martins/go-clean-architecture/domain/usecase/transfer"

	"github.com/domicio-martins/go-clean-architecture/domain/entity/user"

	"github.com/domicio-martins/go-clean-architecture/domain/entity/wallet"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/codegangsta/negroni"
	"github.com/domicio-martins/go-clean-architecture/api/handler"
	"github.com/domicio-martins/go-clean-architecture/api/middleware"
	"github.com/domicio-martins/go-clean-architecture/config"
	"github.com/domicio-martins/go-clean-architecture/pkg/metric"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_DATABASE)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	walletRepo := wallet.NewMySQLRepository(db)
	walletManager := wallet.NewManager(walletRepo)

	userRepo := user.NewMySQLRepoRepository(db)
	userManager := user.NewManager(userRepo, password.NewService())

	transferUseCase := transfer.NewUseCase(userManager, walletManager)

	metricService, err := metric.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(middleware.Metrics(metricService)),
		negroni.NewLogger(),
	)
	//wallet
	handler.MakeWalletHandlers(r, *n, walletManager)

	//user
	handler.MakeUserHandlers(r, *n, userManager)

	//transfer
	handler.MakeTransferHandlers(r, *n, walletManager, userManager, transferUseCase)

	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(config.API_PORT),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

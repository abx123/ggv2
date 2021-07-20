package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"ggv2/logger"
)

func main() {
	logger := logger.NewLogger()
	zap.ReplaceGlobals(logger)

	dsn := getDSN()
	port := getPort()
	conn := initDb(*dsn)
	defer conn.Close()

	router := NewRouter(*port, conn)
	router.InitRouter()
}

func getDSN() *string {
	envdsn := os.Getenv("DSN")

	dsn := flag.String("dsn", "", "mysql datasource string")

	flag.Parse()
	if *dsn == "" {
		dsn = &envdsn
		fmt.Printf("-dsn flag not set, defaulting to %s \n", envdsn)
	}

	return dsn
}

func getPort() *int {
	envport := os.Getenv("PORT")
	port := flag.Int("p", 0, "port on which the application should run on")
	flag.Parse()
	if *port == 0 {
		p, err := strconv.Atoi(envport)
		if err != nil {
			zap.L().Fatal(err.Error(), zap.Error(err))
		}
		port = &p
		fmt.Printf("-p flag not set, defaulting to port %s \n", envport)
	}
	return port
}

func initDb(dsn string) *sqlx.DB {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		zap.L().Fatal(err.Error(), zap.Error(err))
		return nil
	}
	return db
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/siyu-nuro/go-example/config"
	"github.com/siyu-nuro/go-example/pkg/endpoint"
	"github.com/siyu-nuro/go-example/pkg/service/order"
	"github.com/siyu-nuro/go-example/pkg/store"
	db2 "github.com/siyu-nuro/go-example/pkg/store/db"
	"github.com/siyu-nuro/go-example/pkg/transport"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/oklog/pkg/group"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type flags struct {
	configFile *string
}

func parseFlags() *flags {
	// Define flags here
	f := &flags{
		configFile: flag.String("config", "", "Path to configuration yaml"),
	}
	// Define all flags above here
	flag.Parse()
	return f
}

func main() {
	flags := parseFlags()

	// load config from file and panic if fail to load
	cfg, err := config.GetConfig(flags.configFile)
	if err != nil {
		log.Panicf("failed to load config from file, error %v", err)
	}

	// initialize logger
	zapLogger, err := newLogger(cfg)
	if err != nil {
		log.Panicf("failed to initialize logger, error %v", err)
	}
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	// initialize storage manager
	mysqlUserName, mysqlPassWord := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")
	if mysqlUserName == "" {
		mysqlUserName = "root"
	}
	mySQLDataSourceName := mysqlUserName + ":" + mysqlPassWord
	mySQLDataSourceName += fmt.Sprintf("@%v/%v?parseTime=true", cfg.Database.Addr, cfg.Database.Name)
	db := sqlx.MustConnect("mysql", mySQLDataSourceName)
	defer db.Close()

	storageManager := makeStorageManager(logger, db)

	// initialize services
	orderService := order.NewOrderService(logger, storageManager)

	// initialize request handler
	requestHandler := endpoint.NewRequestHandler(orderService)

	// initialize http handler
	httpHandler := transport.NewHTTPHandler(requestHandler)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpAddr := fmt.Sprintf(":%v", cfg.HTTPServer.Port)
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			// Seems like we should be using panic instead of os.Exit to get deferred function execution? https://stackoverflow.com/a/28473339
			logger.Panic("err", err)
		}
		g.Add(func() error {
			logger.Info("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Info("exit", g.Run())
}

func makeStorageManager(logger *zap.SugaredLogger, db *sqlx.DB) store.StorageManager {
	storageManager, err := db2.NewMySQLStorageManager(db, logger)
	if err != nil {
		logger.Panicw("failed to NewMySQLStorageManager", "error", err.Error())
	}
	return storageManager
}

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerCfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   false,
		Encoding:      "json",
		EncoderConfig: encoderCfg,
		OutputPaths: []string{
			"stderr",
			"log.path"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return loggerCfg.Build()
}

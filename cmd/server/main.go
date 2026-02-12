package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-warehouse-ms/internal/api"
	"go-warehouse-ms/internal/infra"
	"go-warehouse-ms/internal/model"
	"go-warehouse-ms/internal/repository"
	"go-warehouse-ms/internal/service"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database repository.DatabaseConfig `mapstructure:"database"`
}

func loadConfig() (AppConfig, error) {
	var cfg AppConfig
	viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	db, err := repository.GetDB(cfg.Database)
	if err != nil {
		logger.Fatal("db init failed", zap.Error(err))
	}
	if err := db.AutoMigrate(&model.User{}, &model.Goods{}); err != nil {
		logger.Fatal("db migrate failed", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	goodsRepo := repository.NewGoodsRepository(db)
	goodsProxy := repository.NewGoodsRepositoryProxy(goodsRepo, logger)
	bus := service.NewEventBus()
	bus.Subscribe(service.NewAuditObserver(logger))
	caretaker := service.NewGoodsCaretaker(20)
	authService := service.NewAuthService(userRepo)
	goodsService := service.NewGoodsService(goodsProxy, bus, caretaker)

	authHandler := api.NewAuthHandler(authService)
	goodsHandler := api.NewGoodsHandler(goodsService)

	router := api.NewRouter(authHandler, goodsHandler, logger)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	logger.Info("server started", zap.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
}

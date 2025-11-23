package main

import (
	"log"
	"wallet-backend/internal/config"
	"wallet-backend/internal/eth"
	"wallet-backend/internal/http"
	"wallet-backend/internal/repository"
	"wallet-backend/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化以太坊客户端
	ethClient, err := eth.NewClient(cfg.EthNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	defer ethClient.Close()

	// 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	txRepo := repository.NewTxRepository(db)

	// 初始化 Service
	walletService := service.NewWalletService(walletRepo, ethClient)
	userService := service.NewUserService(userRepo, walletService)
	txService := service.NewTxService(txRepo, walletRepo, ethClient)

	// 启动区块监听器
	listener := eth.NewListener(ethClient, walletRepo, txRepo)
	go listener.Start(cfg.StartBlock)

	// 启动 HTTP 服务
	router := http.SetupRouter(userService, walletService, txService)
	log.Printf("Server starting on %s", cfg.ServerAddr)
	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

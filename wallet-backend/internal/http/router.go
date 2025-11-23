package http

import (
	"wallet-backend/internal/http/handlers"
	"wallet-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userService *service.UserService, walletService *service.WalletService, txService *service.TxService) *gin.Engine {
	router := gin.Default()

	// 初始化 handlers
	userHandler := handlers.NewUserHandler(userService)
	txHandler := handlers.NewTxHandler(txService)

	// API 路由组
	api := router.Group("/api")
	{
		// 用户相关（无需认证）
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
		}

		// 需要认证的路由
		auth := api.Group("")
		// auth.Use(authMiddleware()) // 可添加 JWT 中间件
		{
			// 用户信息
			auth.GET("/profile", userHandler.GetProfile)

			// 钱包相关
			wallet := auth.Group("/wallet")
			{
				wallet.GET("/balance", func(c *gin.Context) {
					userID := c.GetUint("userID")
					balance, err := walletService.GetBalance(userID)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"balance": balance})
				})
			}

			// 交易相关
			tx := auth.Group("/tx")
			{
				tx.POST("/transfer", txHandler.Transfer)
				tx.GET("/list", txHandler.GetUserTransactions)
			}
		}

		// 公共查询（无需认证）
		api.GET("/tx/:hash", txHandler.GetTransaction)
		api.GET("/tx/list/:address", txHandler.GetTransactionList)
	}

	return router
}

// authMiddleware JWT 认证中间件（示例）
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现 JWT 验证逻辑
		// token := c.GetHeader("Authorization")
		// 验证 token，解析 userID
		// c.Set("userID", userID)
		c.Set("userID", uint(1)) // 临时设置
		c.Next()
	}
}

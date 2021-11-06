package controllers

import (
	"dreampay/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) initializeRoutes() {
	docs.SwaggerInfo.Title = "Backend DREAMPAY "
	docs.SwaggerInfo.Description = "Backend untuk kebutuhan frontend DREAMPAY"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	s.Router.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	dreampay := s.Router.Group("/api")
	{
		dreampay.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "dreamland://sekolahimpian.com",
			})
		})
		dreampay.POST("/register", s.CreateAccountController)
		dreampay.POST("/update/:id", s.UpdateAccountController)
		dreampay.POST("/login", s.LoginAccountController)
		dreampay.POST("/verification", s.VerificationAccountController)
		dreampay.POST("/delete/:mobile", s.DeleteAccountController)
		dreampay.GET("/account", s.GetAllAccountController)
		dreampay.GET("/transactions", s.GetAllTransactionController)
		dreampay.GET("/transaction", s.GetTransactionByIDController)
		dreampay.POST("/transaction", s.CreateTransactionController)
		dreampay.GET("/money-status", s.GetMoneyStatusController)
		dreampay.POST("/withdraw", s.CreateWithdrawController)
		dreampay.GET("/withdraw", s.GetWithdrawController)
		dreampay.POST("/transaction/delete/multiple", s.DeleteMultipleTransaction)
		dreampay.POST("/withdraw/delete/multiple", s.DeleteMultipleWithdraw)
	}

}

package initRouter

import (
	"Lottery/handler"
	"Lottery/middleware"
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	if mode := gin.Mode(); mode == gin.TestMode {
		router.LoadHTMLGlob("./../templates/*")
	} else {
		router.LoadHTMLGlob("templates/*")
	}
	router.Static("/statics","./statics")
	router.StaticFile("/lottery_icon.svg","./lottery_icon.svg")
	indexRouter := router.Group("/")
	{
		indexRouter.GET("", handler.Index)
	}
	customerRouter := router.Group("/customer")
	{
		customerRouter.GET("", middleware.CustomerAuth(), handler.Customer)
		customerRouter.POST("/home", handler.CustomerLogin)
		customerRouter.GET("/home", middleware.CustomerHomeAuth(), handler.CustomerHome)
		customerRouter.POST("/logout", handler.CustomerLogout)
		customerRouter.GET("/register", handler.CustomerRegisterIndex)
		customerRouter.POST("/register", handler.CustomerRegister)
		customerRouter.GET("/lottery/:lotteryId", middleware.CustomerLotteryAuth(), handler.CustomerLottery)
		customerRouter.POST("/lottery/:lotteryId", middleware.CustomerLotteryAuth(), handler.CustomerLotteryJoin)
	}
	businessRouter := router.Group("/business")
	{
		businessRouter.GET("", middleware.BusinessAuth(), handler.Business)
		businessRouter.POST("/home", handler.BusinessLogin)
		businessRouter.GET("/home", middleware.BusinessHomeAuth(), handler.BusinessHome)
		businessRouter.POST("/logout", handler.BusinessLogout)
		businessRouter.GET("/register", handler.BusinessRegisterIndex)
		businessRouter.POST("/register", handler.BusinessRegister)
		businessRouter.GET("/home/create-lottery", middleware.BusinessHomeAuth(), handler.BusinessCreateLotteryIndex)
		businessRouter.POST("/home/create-lottery", middleware.BusinessHomeAuth(), handler.BusinessCreateLottery)
		businessRouter.GET("/lottery/:lotteryId", middleware.BusinessLotteryAuth(), handler.BusinessLottery)
	}
	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("", middleware.AdminAuth(), handler.Admin)
		adminRouter.POST("/home", handler.AdminLogin)
		adminRouter.GET("/home", middleware.AdminHomeAuth(), handler.AdminHome)
		adminRouter.POST("/logout", handler.AdminLogout)
		adminRouter.GET("/lottery/:lotteryId", middleware.AdminLotteryAuth(), handler.AdminLottery)
	}
	return router
}

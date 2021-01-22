package handler

import (
	"Lottery/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Admin(context *gin.Context) {
	//用code来表示页面状态
	//1：默认登录页面
	//2：用户名或密码错误
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"usernameError": false,
		})
	}else{
		context.HTML(http.StatusOK, "admin.tmpl", gin.H{
			"usernameError": true,
		})
	}
}

func AdminLogin(context *gin.Context) {
	admin := &model.Admin{}
	if err := context.ShouldBind(&admin); err != nil {
		context.Redirect(http.StatusMovedPermanently, "/admin?code=2")
	}
	adminQuery := admin.QueryByUsername()
	if admin.AdminPassword == adminQuery.AdminPassword&&adminQuery.AdminPassword!=""{
		expiresTime := time.Now().Unix() + int64(60 * 60 * 2)
		claims := jwt.StandardClaims{
			Audience:  admin.AdminUsername,         // 受众
			ExpiresAt: expiresTime,                 // 失效时间
			Id:        strconv.Itoa(adminQuery.AdminId), // 编号
			IssuedAt:  time.Now().Unix(),           // 签发时间
			Issuer:    "Lottery",                   // 签发人
			NotBefore: time.Now().Unix(),           // 生效时间
			Subject:   "adminLogin",                // 主题
		}
		var jwtSecret = []byte("lotteryJwtSecret")
		tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
			context.SetCookie("adminToken", token, 60 * 60 * 2, "/", "localhost", false, true)
			context.Redirect(http.StatusMovedPermanently, "/admin/home")
		} else {
			context.Redirect(http.StatusMovedPermanently, "/admin?code=2")
		}
	} else {
		context.Redirect(http.StatusMovedPermanently, "/admin?code=2")
	}
}

func AdminHome(context *gin.Context) {
	adminToken, _ := context.Cookie("adminToken")
	claim, _ := parseToken(adminToken)
	adminId, _ := strconv.Atoi(claim.Id)
	admin := model.Admin{
		AdminId: adminId,
	}
	lotteries := admin.GetAllLotteries()
	context.HTML(http.StatusOK, "adminHome.tmpl", gin.H{
		"lotteries": lotteries,
	})
}

func AdminLogout(context *gin.Context) {
	context.SetCookie("adminToken", "", -1, "/", "localhost", false, true)
	context.Redirect(http.StatusMovedPermanently, "/admin")
}

func AdminLottery(context *gin.Context) {
	lottery := model.Lottery{
		LotteryId : 1,
	}
	lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
	if err!=nil{
		context.Redirect(http.StatusMovedPermanently, "/admin/home")
	}
	lottery.LotteryId = lotteryId
	lottery, code := lottery.QueryById()
	business := model.Business{
		BusinessId: lottery.LotteryBusiness,
	}
	business = business.QueryByUsername()
	lotteryBusiness := business.BusinessName
	if code==0{
		context.Redirect(http.StatusMovedPermanently, "/admin/home")
	}else{
		context.HTML(http.StatusOK, "adminLottery.tmpl", gin.H{
			"lotteryBusiness": lotteryBusiness,
			"lotteryName": lottery.LotteryName,
			"lotteryGift": lottery.LotteryGift,
			"lotteryCreateTime": lottery.LotteryCreateTime,
			"lotteryStartTime": lottery.LotteryStartTime,
			"lotteryTime": lottery.LotteryTime,
			"lotteryQueueLength": lottery.LotteryQueueLength,
			"lotteryState": lottery.LotteryState,
			"lotteryTotalCustomer": lottery.LotteryTotalCustomer,
			"lotteryResult": lottery.LotteryResult,
			"lotteryLabel": lottery.LotteryLabel,
		})
	}
}
package handler

import (
	"Lottery/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)


func Business(context *gin.Context) {
	//用code来表示页面状态
	//1：默认登录页面
	//2：用户名或密码错误
	//3：注册成功
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "business.tmpl", gin.H{
			"usernameError": false,
			"registerSuccess": false,
		})
	}else if code=="2"{
		context.HTML(http.StatusOK, "business.tmpl", gin.H{
			"usernameError": true,
			"registerSuccess": false,
		})
	}else{
		context.HTML(http.StatusOK, "business.tmpl", gin.H{
			"usernameError": false,
			"registerSuccess": true,
		})
	}
}

func BusinessLogin(context *gin.Context) {
	business := &model.Business{}
	if err := context.ShouldBind(&business); err != nil {
		context.Redirect(http.StatusMovedPermanently, "/business?code=2")
	}
	businessQuery := business.QueryByUsername()
	if business.BusinessPassword == businessQuery.BusinessPassword&&businessQuery.BusinessPassword!=""{
		expiresTime := time.Now().Unix() + int64(60 * 60 * 2)
		claims := jwt.StandardClaims{
			Audience:  business.BusinessUsername,         // 受众
			ExpiresAt: expiresTime,                       // 失效时间
			Id:        strconv.Itoa(businessQuery.BusinessId), // 编号
			IssuedAt:  time.Now().Unix(),                 // 签发时间
			Issuer:    "Lottery",                         // 签发人
			NotBefore: time.Now().Unix(),                 // 生效时间
			Subject:   "businessLogin",                   // 主题
		}
		var jwtSecret = []byte("lotteryJwtSecret")
		tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
			context.SetCookie("businessToken", token, 60 * 60 * 2, "/", "localhost", false, true)
			context.Redirect(http.StatusMovedPermanently, "/business/home")
		} else {
			context.Redirect(http.StatusMovedPermanently, "/business?code=2")
		}
	} else {
		context.Redirect(http.StatusMovedPermanently, "/business?code=2")
	}
}

func BusinessHome(context *gin.Context) {
	businessToken, _ := context.Cookie("businessToken")
	claim, _ := parseToken(businessToken)
	businessId, _ := strconv.Atoi(claim.Id)
	business := model.Business{
		BusinessId: businessId,
	}
	lotteries := business.QueryLotteryByBusinessId()
	context.HTML(http.StatusOK, "businessHome.tmpl", gin.H{
		"lotteries": lotteries,
	})
}

func BusinessRegisterIndex(context *gin.Context) {
	//用code来表示页面状态
	//1:默认注册页面
	//2:非法注册，请重新填写
	//3:用户名已存在
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "businessRegister.tmpl", gin.H{
			"msg": "",
		})
	}else if code=="2"{
		context.HTML(http.StatusOK, "businessRegister.tmpl", gin.H{
			"msg": "非法注册，请重新填写",
		})
	}else if code=="3" {
		context.HTML(http.StatusOK, "businessRegister.tmpl", gin.H{
			"msg": "用户名已存在",
		})
	}else{
		context.HTML(http.StatusOK, "businessRegister.tmpl", gin.H{
			"msg": "",
		})
	}
}

func BusinessRegister(context *gin.Context) {
	var business model.Business
	if err := context.ShouldBind(&business); err != nil {
		context.Redirect(http.StatusMovedPermanently, "/register?code=2")
	}
	_, code := business.Create()
	if code==0{
		context.Redirect(http.StatusMovedPermanently, "/business?code=3")
	}else if code==1{
		context.Redirect(http.StatusMovedPermanently, "/business/register?code=2")
	}else if code==2{
		context.Redirect(http.StatusMovedPermanently, "/business/register?code=3")
	}
}

func BusinessLogout(context *gin.Context) {
	context.SetCookie("businessToken", "", -1, "/", "localhost", false, true)
	context.Redirect(http.StatusMovedPermanently, "/business")
}

func BusinessCreateLotteryIndex(context *gin.Context) {
	//用code来表示页面状态
	//1:默认创建页面
	//2:创建新抽奖失败
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "businessCreateLottery.tmpl", gin.H{
			"msg": "",
		})
	}else if code=="2"{
		context.HTML(http.StatusOK, "businessCreateLottery.tmpl", gin.H{
			"msg": "创建失败，请按照要求重新填写",
		})
	}else{
		context.HTML(http.StatusOK, "businessCreateLottery.tmpl", gin.H{
			"msg": "",
		})
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("lotteryJwtSecret"), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}

func BusinessCreateLottery(context *gin.Context) {
	lottery := &model.Lottery{}
	if err := context.ShouldBind(&lottery); err != nil {
		//fmt.Println("err1")
		context.Redirect(http.StatusMovedPermanently, "/business/home/create-lottery?code=2")
		return
	}
	lotteryQueueLength := context.PostForm("lotteryQueueLength")
	if lotteryQueueLength==""{
		//fmt.Println("err2")
		context.Redirect(http.StatusMovedPermanently, "/business/home/create-lottery?code=2")
		return
	}
	lotteryLotteryQueueLength, err := strconv.Atoi(lotteryQueueLength)
	if err!=nil{
		//fmt.Println("err3")
		context.Redirect(http.StatusMovedPermanently, "/business/home/create-lottery?code=2")
		return
	}
	lottery.LotteryQueueLength = lotteryLotteryQueueLength
	businessToken, err := context.Cookie("businessToken")
	if err!=nil{
		context.Abort()
		context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		return
	}
	// 校验businessToken
	jwtToken, err := jwt.ParseWithClaims(businessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("lotteryJwtSecret"), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			business := model.Business{
				BusinessUsername: claim.Audience,
			}
			business = business.QueryByUsername()
			lottery.LotteryBusiness = business.BusinessId
			if business.BusinessPassword==""{
				//fmt.Println("err4")
				context.Redirect(http.StatusMovedPermanently, "/business/home/create-lottery?code=2")
				return
			}
			lottery, code := lottery.Create()
			if code==0{
				//fmt.Println(lottery)
				//fmt.Println(lottery.LotteryId)
				//fmt.Println("/business/lottery/"+strconv.Itoa(lottery.LotteryId))
				customers := lottery.QueryCustomersByLottery()
				//fmt.Println(customers)
				for _, customer:= range customers{
					customer.UpdateCustomerLottery(lottery.LotteryId)
				}
				business.UpdateBusinessLottery(lottery.LotteryId)
				context.Redirect(http.StatusMovedPermanently, "/business/lottery/"+strconv.Itoa(lottery.LotteryId))
				return
			}else{
				//fmt.Println("err5")
				context.Redirect(http.StatusMovedPermanently, "/business/home/create-lottery?code=2")
				return
			}
		}else{
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business?code=1")
			return
		}
	}else{
		context.Abort()
		context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		return
	}
}

func BusinessLottery(context *gin.Context) {
	lottery := model.Lottery{
		LotteryId : 1,
	}
	lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
	if err!=nil{
		context.Redirect(http.StatusMovedPermanently, "/business/home")
	}
	lottery.LotteryId = lotteryId
	lottery, code := lottery.QueryById()
	if code==0{
		context.Redirect(http.StatusMovedPermanently, "/business/home")
	}else{
		context.HTML(http.StatusOK, "businessLottery.tmpl", gin.H{
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
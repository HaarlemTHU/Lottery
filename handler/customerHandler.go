package handler

import (
	"Lottery/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Customer(context *gin.Context) {
	//用code来表示页面状态
	//1：默认登录页面
	//2：用户名或密码错误
	//3：注册成功
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "customer.tmpl", gin.H{
			"usernameError": false,
			"registerSuccess": false,
		})
	}else if code=="2"{
		context.HTML(http.StatusOK, "customer.tmpl", gin.H{
			"usernameError": true,
			"registerSuccess": false,
		})
	}else{
		context.HTML(http.StatusOK, "customer.tmpl", gin.H{
			"usernameError": false,
			"registerSuccess": true,
		})
	}
}

func CustomerLogin(context *gin.Context) {
	customer := &model.Customer{}
	if err := context.ShouldBind(&customer); err != nil {
		context.Redirect(http.StatusMovedPermanently, "/customer?code=2")
	}
	customerQuery := customer.QueryByUsername()
	if customer.CustomerPassword == customerQuery.CustomerPassword&&customerQuery.CustomerPassword!=""{
		expiresTime := time.Now().Unix() + int64(60 * 60 * 2)
		claims := jwt.StandardClaims{
			Audience:  customer.CustomerUsername,         // 受众
			ExpiresAt: expiresTime,                 // 失效时间
			Id:        strconv.Itoa(customerQuery.CustomerId), // 编号
			IssuedAt:  time.Now().Unix(),           // 签发时间
			Issuer:    "Lottery",                   // 签发人
			NotBefore: time.Now().Unix(),           // 生效时间
			Subject:   "customerLogin",                // 主题
		}
		var jwtSecret = []byte("lotteryJwtSecret")
		tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
			context.SetCookie("customerToken", token, 60 * 60 * 2, "/", "localhost", false, true)
			context.Redirect(http.StatusMovedPermanently, "/customer/home")
		} else {
			context.Redirect(http.StatusMovedPermanently, "/customer?code=2")
		}
	} else {
		context.Redirect(http.StatusMovedPermanently, "/customer?code=2")
	}
}

func CustomerHome(context *gin.Context) {
	customerToken, _ := context.Cookie("customerToken")
	claim, _ := parseToken(customerToken)
	customerId, _ := strconv.Atoi(claim.Id)
	customer := model.Customer{
		CustomerId: customerId,
	}
	lotteries := customer.QueryLotteryByCustomerId()
	context.HTML(http.StatusOK, "customerHome.tmpl", gin.H{
		"lotteries": lotteries,
	})
}

func CustomerRegisterIndex(context *gin.Context) {
	//用code来表示页面状态
	//1:默认注册页面
	//2:非法注册，请重新填写
	//3:用户名非法，请重新填写
	//4:用户名已存在
	code := context.DefaultQuery("code", "1")
	if code=="1"{
		context.HTML(http.StatusOK, "customerRegister.tmpl", gin.H{
			"msg": "",
		})
	}else if code=="2"{
		context.HTML(http.StatusOK, "customerRegister.tmpl", gin.H{
			"msg": "非法注册，请重新填写",
		})
	}else if code=="3"{
		context.HTML(http.StatusOK, "customerRegister.tmpl", gin.H{
			"msg": "用户名非法，请重新填写",
		})
	}else if code=="4"{
		context.HTML(http.StatusOK, "customerRegister.tmpl", gin.H{
			"msg": "用户名已存在",
		})
	}else{
		context.HTML(http.StatusOK, "customerRegister.tmpl", gin.H{
			"msg": "",
		})
	}
}

func CustomerRegister(context *gin.Context) {
	var customer model.Customer
	if err := context.ShouldBind(&customer); err != nil {
		context.Redirect(http.StatusMovedPermanently, "/customer/register?code=2")
	}
	_, code := customer.Create()
	if code==0{
		context.Redirect(http.StatusMovedPermanently, "/customer?code=3")
	}else if code==1{
		context.Redirect(http.StatusMovedPermanently, "/customer/register?code=3")
	}else if code==2{
		context.Redirect(http.StatusMovedPermanently, "/customer/register?code=4")
	}else if code==3{
		context.Redirect(http.StatusMovedPermanently, "/customer/register?code=2")
	}
}

func CustomerLogout(context *gin.Context) {
	context.SetCookie("customerToken", "", -1, "/", "localhost", false, true)
	context.Redirect(http.StatusMovedPermanently, "/customer")
}

func CustomerLottery(context *gin.Context) {
	//access=0:抽奖仍在进行，未参与过此次抽奖，默认页面
	//access=1:抽奖仍在进行，但已参与过此次抽奖，显示中奖信息
	customerToken, _ := context.Cookie("customerToken")
	claim, _ := parseToken(customerToken)
	customerId, _ := strconv.Atoi(claim.Id)
	access := context.DefaultQuery("access", "0")
	lottery := model.Lottery{
		LotteryId : 1,
	}
	lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
	if err!=nil{
		fmt.Println("err1")
		context.Redirect(http.StatusMovedPermanently, "/customer/home")
	}
	lottery.LotteryId = lotteryId
	lottery, _ = lottery.QueryById()
	gift := "很遗憾您未中奖"
	business := model.Business{
		BusinessId: lottery.LotteryBusiness,
	}
	business = business.QueryByUsername()
	lotteryBusiness := business.BusinessName
	if lottery.LotteryResult!=""{
		lotteryResultList := strings.Split(lottery.LotteryResult, ";")
		for _, lotteryResult:=range lotteryResultList{
			customerIdQuery, _ :=  strconv.Atoi(strings.Split(lotteryResult, ":")[1])
			if customerIdQuery==customerId{
				gift = strings.Split(lotteryResult, ":")[0]
				break
			}
		}
	}
	context.HTML(http.StatusOK, "customerLottery.tmpl", gin.H{
		"lotteryBusiness": lotteryBusiness,
		"lotteryId": lotteryId,
		"lotteryAccess": access=="0",
		"lotteryName": lottery.LotteryName,
		"lotteryGift": lottery.LotteryGift,
		"lotteryStartTime": lottery.LotteryStartTime,
		"lotteryTime": lottery.LotteryTime,
		"lotteryState": lottery.LotteryState,
		"lotteryTotalCustomer": lottery.LotteryTotalCustomer,
		"gift": gift,
	})
}

func CustomerLotteryJoin(context *gin.Context) {
	lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
	if err!=nil{
		//fmt.Println("err1")
		context.Redirect(http.StatusMovedPermanently, "/customer/home")
	}
	customerToken, _ := context.Cookie("customerToken")
	claim, _ := parseToken(customerToken)
	customerId, _ := strconv.Atoi(claim.Id)
	customer := model.Customer{
		CustomerId: customerId,
	}
	access, gift :=customer.QueryLotteryAccessByCustomerId(lotteryId)
	if access{
		customer.JoinLottery(lotteryId)
		context.Redirect(http.StatusMovedPermanently, "/customer/lottery/"+strconv.Itoa(lotteryId)+"?access=1")
	}else{
		if gift==""{
			context.Redirect(http.StatusMovedPermanently, "/customer/home")
		}else{
			context.Redirect(http.StatusMovedPermanently, "/customer/lottery/"+strconv.Itoa(lotteryId)+"?access=1")
		}
	}
}
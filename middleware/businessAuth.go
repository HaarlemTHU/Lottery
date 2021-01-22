package middleware

import (
	"Lottery/initDB"
	"Lottery/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

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

func BusinessAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		businessToken, err := context.Cookie("businessToken")
		if err!=nil{
			context.Next()
		}
		// 校验businessToken
		_, err = parseToken(businessToken)
		if err != nil {
			context.Next()
		}
		context.Abort()
		context.Redirect(http.StatusMovedPermanently, "/business/home")
	}
}

func BusinessHomeAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		businessToken, err := context.Cookie("businessToken")
		if err!=nil{
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		}
		// 校验businessToken
		_, err = parseToken(businessToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		}
		context.Next()
	}
}

func BusinessLotteryAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		businessToken, err := context.Cookie("businessToken")
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		}
		// 校验businessToken
		claim, err := parseToken(businessToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business?code=1")
		}
		lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
		if err != nil {
			//fmt.Println("err1")
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/business/home")
		}
		lottery := model.Lottery{}
		lottery.LotteryId = lotteryId
		lottery, code := lottery.QueryById()
		if code == 0 {
			//fmt.Println("err2")
			context.Redirect(http.StatusMovedPermanently, "/business/home")
		} else{
			claimId, err := strconv.Atoi(claim.Id)
			if err!=nil{
				//fmt.Println("err3")
				context.Redirect(http.StatusMovedPermanently, "/business/home")
			}
			if claimId==lottery.LotteryBusiness{
				if lottery.LotteryState=="2"&&parseTimeStrToTimestamp(lottery.LotteryStartTime,5)<=time.Now().Unix(){
					//fmt.Println(lottery.LotteryStartTime)
					//fmt.Println(parseTimeStrToTimestamp(lottery.LotteryStartTime,5))
					//fmt.Println(time.Now().Unix())
					//fmt.Println("change1")
					lottery = lottery.UpdateLotteryState("1")
				}
				lotteryTime, _ := strconv.ParseInt(lottery.LotteryTime, 10, 64)
				if lottery.LotteryState=="1"&&parseTimeStrToTimestamp(lottery.LotteryStartTime,5)+lotteryTime<=time.Now().Unix(){
					//fmt.Println("change2")
					lottery = lottery.UpdateLotteryState("0")
					initDB.Db.DropTable("lottery_queue"+strconv.Itoa(lotteryId))
				}
				context.Next()
			}else{
				//fmt.Println("err4")
				//fmt.Println(claim, " | ", lottery.LotteryBusiness)
				context.Redirect(http.StatusMovedPermanently, "/business/home")
			}
		}
	}
}
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

func parseAdminToken(token string) (*jwt.StandardClaims, error) {
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

func AdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		adminToken, err := context.Cookie("adminToken")
		if err!=nil{
			context.Next()
		}
		// 校验adminToken
		_, err = parseAdminToken(adminToken)
		if err != nil {
			context.Next()
		}
		context.Abort()
		context.Redirect(http.StatusMovedPermanently, "/admin/home")
	}
}

func AdminHomeAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		adminToken, err := context.Cookie("adminToken")
		if err!=nil{
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/admin?code=1")
		}
		// 校验adminToken
		_, err = parseAdminToken(adminToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/admin?code=1")
		}
		context.Next()
	}
}

func AdminLotteryAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		adminToken, err := context.Cookie("adminToken")
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/admin?code=1")
		}
		// 校验adminToken
		_, err = parseToken(adminToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/admin?code=1")
		}
		lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
		if err != nil {
			//fmt.Println("err1")
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/admin/home")
		}
		lottery := model.Lottery{}
		lottery.LotteryId = lotteryId
		lottery, code := lottery.QueryById()
		if code == 0 {
			//fmt.Println("err2")
			context.Redirect(http.StatusMovedPermanently, "/admin/home")
		} else{
			if lottery.LotteryState=="2"&&parseTimeStrToTimestamp(lottery.LotteryStartTime,5)<=time.Now().Unix(){
				lottery = lottery.UpdateLotteryState("1")
			}
			lotteryTime, _ := strconv.ParseInt(lottery.LotteryTime, 10, 64)
			if lottery.LotteryState=="1"&&parseTimeStrToTimestamp(lottery.LotteryStartTime,5)+lotteryTime<=time.Now().Unix(){
				lottery = lottery.UpdateLotteryState("0")
				initDB.Db.DropTable("lottery_queue"+strconv.Itoa(lotteryId))
			}
			context.Next()
		}
	}
}

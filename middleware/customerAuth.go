package middleware

import (
	"Lottery/initDB"
	"Lottery/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseCustomerToken(token string) (*jwt.StandardClaims, error) {
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

func parseTimeStrToTimestamp(timeStr string, flag int) int64 {
	var t int64
	loc, _ := time.LoadLocation("Local")
	if flag == 1 {
		t1, _ := time.ParseInLocation("2006.01.02 15:04:05", timeStr, loc)
		t = t1.Unix()
	} else if flag == 2 {
		t1, _ := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
		t = t1.Unix()
	} else if flag == 3 {
		t1, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
		t = t1.Unix()
	} else if flag == 4 {
		t1, _ := time.ParseInLocation("2006.01.02", timeStr, loc)
		t = t1.Unix()
	} else {
		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
		t = t1.Unix()
	}
	return t
}

func CustomerAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		customerToken, err := context.Cookie("customerToken")
		if err!=nil{
			context.Next()
		}
		// 校验customerToken
		_, err = parseCustomerToken(customerToken)
		if err != nil {
			context.Next()
		}
		context.Abort()
		context.Redirect(http.StatusMovedPermanently, "/customer/home")
	}
}

func CustomerHomeAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		customerToken, err := context.Cookie("customerToken")
		if err!=nil{
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/customer?code=1")
		}
		// 校验customerToken
		_, err = parseCustomerToken(customerToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/customer?code=1")
		}
		context.Next()
	}
}

func CustomerLotteryAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		customerToken, err := context.Cookie("customerToken")
		if err != nil {
			context.Abort()
			//fmt.Println("err1")
			context.Redirect(http.StatusMovedPermanently, "/customer?code=1")
		}
		// 校验customerToken
		claim, err := parseToken(customerToken)
		if err != nil {
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/customer?code=1")
		}
		lotteryId, err := strconv.Atoi(context.Param("lotteryId"))
		if err != nil {
			//fmt.Println("err2")
			context.Abort()
			context.Redirect(http.StatusMovedPermanently, "/customer/home")
		}
		lottery := model.Lottery{}
		lottery.LotteryId = lotteryId
		lottery, code := lottery.QueryById()
		if code == 0 {
			//fmt.Println("err3")
			context.Redirect(http.StatusMovedPermanently, "/customer/home")
		} else{
			claimId, err := strconv.Atoi(claim.Id)
			if err!=nil{
				//fmt.Println("err4")
				context.Redirect(http.StatusMovedPermanently, "/customer/home")
			}else{
				customer := model.Customer{
					CustomerId: claimId,
				}
				customer = customer.QueryByUsername()
				lotteryIds := strings.Split(customer.CustomerLottery,";")
				//fmt.Println(lotteryIds)
				//fmt.Println(lottery.LotteryId)
				flag:=0
				for _,lotteryIdStr:=range lotteryIds{
					lotteryId,_:=strconv.Atoi(strings.Split(lotteryIdStr,":")[0])
					if lottery.LotteryId==lotteryId{
						flag = 1
						break
					}
				}
				if flag==0{
					//fmt.Println("err5")
					context.Redirect(http.StatusMovedPermanently, "/customer/home")
				}else {
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
	}
}

package model

import (
	"Lottery/initDB"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var tabName string

type Lottery struct {
	LotteryId int `gorm:"primary_key"`
	LotteryBusiness int
	LotteryName string `form:"lotteryName"`
	LotteryGift string `form:"lotteryGift"`
	LotteryCreateTime string
	LotteryStartTime string `form:"lotteryStartTime"`
	LotteryTime string `form:"lotteryTime"`
	LotteryLabel string `form:"lotteryLabel"`
	LotteryState string
	LotteryResult string
	LotteryTotalCustomer int
	LotteryQueueLength int
}

type LotteryQueueObject struct {
	LotteryQueueObjectId int `gorm:"primary_key"`
	LotteryQueueObjectGift string
}

func (LotteryQueueObject) TableName() string {
	return tabName
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

func (lottery Lottery) Create() (Lottery, int) {
	//code=0:创建成功
	//code=1:创建失败
	labels := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	flag := 0
	for _, label:=range labels{
		if label==lottery.LotteryLabel{
			flag = 1
			break
		}
	}
	if flag==0{
		//fmt.Println("err1")
		return lottery, 1
	}
	lottery.LotteryCreateTime = time.Now().Format("2006-01-02 15:04:05")
	flag = 0
	lotteryTimes := []string{"30", "60", "120", "180", "300"}
	for _, lotteryTime:=range lotteryTimes{
		if lotteryTime==lottery.LotteryTime{
			flag = 1
			break
		}
	}
	if flag==0{
		//fmt.Println("err2")
		return lottery, 1
	}
	if parseTimeStrToTimestamp(lottery.LotteryStartTime,5)<parseTimeStrToTimestamp(lottery.LotteryCreateTime,5){
		//fmt.Println("err3")
		return lottery, 1
	}
	if lottery.LotteryName==""||lottery.LotteryGift==""{
		//fmt.Println("err4")
		return lottery, 1
	}
	if lottery.LotteryQueueLength<=0{
		//fmt.Println("err5")
		return lottery, 1
	}
	lottery.LotteryState = "2"
	lottery.LotteryResult = ""
	lottery.LotteryTotalCustomer = 0
	initDB.Db.Create(&lottery)
	lottery, _ = lottery.QueryById()
	lottery.CreateLotteryQueue()
	return lottery, 0
}

func (lottery Lottery) QueryById() (Lottery,int) {
	// 0:未查询到
	// 1：查询到
	if errs := initDB.Db.Where(&lottery).First(&lottery).GetErrors(); len(errs)!=0{
		return lottery,0
	}
	return lottery,1
}

//查询所有对应label的customer
func (lottery Lottery) QueryCustomersByLottery() []Customer {
	var customers []Customer
	if errs := initDB.Db.Where("customer_label = ?", lottery.LotteryLabel).Find(&customers).GetErrors(); len(errs)!=0{
		return customers
	}
	return customers
}

//更改state
func (lottery Lottery) UpdateLotteryState(newState string) Lottery {
	initDB.Db.Model(&lottery).Where("lottery_id = ?", lottery.LotteryId).Update("lottery_state", newState)
	return lottery
}

func (lottery Lottery) CreateLotteryQueue() Lottery {
	tabName = "lottery_queue"+strconv.Itoa(lottery.LotteryId)
	initDB.Db.AutoMigrate(&LotteryQueueObject{})
	giftsStrList := strings.Split(lottery.LotteryGift,";")
	var gifts []string
	for _, giftsStr:=range giftsStrList{
		num,_ := strconv.Atoi(strings.Split(giftsStr,":")[1])
		gift := strings.Split(giftsStr,":")[0]
		for i := 0; i < num; i++{
			gifts = append(gifts, gift)
		}
	}
	nullLength := lottery.LotteryQueueLength - len(gifts)
	for i := 0; i < nullLength; i++ {
		//" "代表没抽中
		gifts = append(gifts, " ")
	}
	gifts = Random(gifts)
	for _, gift := range gifts{
		lotteryQueueObject := LotteryQueueObject{
			LotteryQueueObjectGift: gift,
		}
		initDB.Db.Create(&lotteryQueueObject)
	}
	return lottery
}

//洗牌算法，随机打乱
func Random(strings []string) []string {
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}
	return strings
}

//参与抽奖，修改lottery LotteryTotalCustomer
func (lottery Lottery) UpdateLotteryTotalCustomer(newTotalCustomer int) Lottery {
	initDB.Db.Model(&lottery).Where("lottery_id = ?", lottery.LotteryId).Update("lottery_total_customer", newTotalCustomer)
	lottery,_ = lottery.QueryById()
	return lottery
}

//抽取对应lottery抽奖队列中的值
func (lottery Lottery) Gift() (Lottery, string) {
	tabName = "lottery_queue"+strconv.Itoa(lottery.LotteryId)
	lottery, _ = lottery.QueryById()
	if lottery.LotteryTotalCustomer==lottery.LotteryQueueLength{
		return lottery, " "
	}
	lottery = lottery.UpdateLotteryTotalCustomer(lottery.LotteryTotalCustomer + 1)
	lotteryQueueObject := LotteryQueueObject{
		LotteryQueueObjectId: lottery.LotteryTotalCustomer,
	}
	initDB.Db.Where(&lotteryQueueObject).First(&lotteryQueueObject)
	return lottery,lotteryQueueObject.LotteryQueueObjectGift
}

//参与抽奖，修改lottery LotteryResult
func (lottery Lottery) UpdateLotteryResult(newResult string) Lottery {
	initDB.Db.Model(&lottery).Where("lottery_id = ?", lottery.LotteryId).Update("lottery_result", newResult)
	lottery,_ = lottery.QueryById()
	return lottery
}
package model

import (
	"Lottery/initDB"
	"strconv"
	"strings"
)

type Business struct {
	BusinessId int
	BusinessName string `form:"businessName"`
	BusinessUsername string `form:"businessUsername"`
	BusinessPassword string `form:"businessPassword"`
	BusinessLottery string
}

// 登录查询
func (business Business) QueryByUsername() Business {
	if errs := initDB.Db.Where(&business).First(&business).GetErrors(); len(errs)!=0{
		business.BusinessPassword = ""
		return business
	}
	return business
}

// 注册
func (business Business) Create() (Business, int) {
	//code=0:创建成功
	//code=1:创建失败，用户名非法
	//code=2:创建失败，用户名已存在
	businessByUsername := Business{
		BusinessUsername: business.BusinessUsername,
	}
	if business.BusinessUsername==""||business.BusinessPassword==""||business.BusinessName==""{
		return business, 1
	}
	if errs := initDB.Db.Where(&businessByUsername).First(&businessByUsername).GetErrors(); len(errs)==0{
		if businessByUsername.BusinessPassword!=""{
			return business, 2
		}
		return business, 1
	}
	initDB.Db.Create(&business)
	return business, 0
}

// 查询所有lottery
func (business Business) QueryLotteryByBusinessId() []Lottery {
	var lotteries []Lottery
	business = business.QueryByUsername()
	if business.BusinessLottery==""{
		return lotteries
	}
	lotteryIds := strings.Split(business.BusinessLottery,";")
	for _, lotteryIdStr:=range lotteryIds{
		lotteryId, _ := strconv.Atoi(lotteryIdStr)
		lottery := Lottery{LotteryId: lotteryId}
		lottery, _ = lottery.QueryById()
		lotteries = append(lotteries, lottery)
	}
	return lotteries
}

func (business Business) UpdateBusinessLottery(lotteryId int) Business {
	var newBusinessLottery string
	if business.BusinessLottery==""{
		newBusinessLottery = business.BusinessLottery + strconv.Itoa(lotteryId)
	}else{
		newBusinessLottery = business.BusinessLottery + ";" + strconv.Itoa(lotteryId)
	}
	initDB.Db.Model(&business).Where("business_id = ?", business.BusinessId).Update("business_lottery", newBusinessLottery)
	return business
}

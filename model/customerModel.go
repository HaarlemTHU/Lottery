package model

import (
	"Lottery/initDB"
	"strconv"
	"strings"
)

type Customer struct {
	CustomerId int
	CustomerName string `form:"customerName"`
	CustomerUsername string `form:"customerUsername"`
	CustomerPassword string `form:"customerPassword"`
	CustomerLabel string `form:"customerLabel"`
	CustomerLottery string
}

func (customer Customer) QueryByUsername() Customer {
	if errs := initDB.Db.Where(&customer).First(&customer).GetErrors(); len(errs)!=0{
		customer.CustomerPassword = ""
		return customer
	}
	return customer
}

func (customer Customer) Create() (Customer, int) {
	//code=0:创建成功
	//code=1:创建失败，用户名非法
	//code=2:创建失败，用户名已存在
	//code=3:创建失败，label非法
	customerByUsername := Customer{
		CustomerUsername: customer.CustomerUsername,
	}
	if customer.CustomerUsername==""||customer.CustomerPassword==""||customer.CustomerName==""{
		return customer, 1
	}
	if errs := initDB.Db.Where(&customerByUsername).First(&customerByUsername).GetErrors(); len(errs)==0{
		if customerByUsername.CustomerPassword!=""{
			return customer, 2
		}
		return customer, 1
	}
	labels := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	flag := 0
	for _, label:=range labels{
		if label==customer.CustomerLabel{
			flag = 1
			break
		}
	}
	if flag==0{
		return customer, 3
	}
	initDB.Db.Create(&customer)
	return customer, 0
}

// 查询所有lottery
func (customer Customer) QueryLotteryByCustomerId() []Lottery {
	var lotteries []Lottery
	customer = customer.QueryByUsername()
	if customer.CustomerLottery==""{
		return lotteries
	}
	lotteryIds := strings.Split(customer.CustomerLottery,";")
	//fmt.Println(len(lotteries))
	//fmt.Println(lotteryIds)
	//fmt.Println(len(lotteryIds))
	for _, lotteryIdStr:=range lotteryIds{
		//fmt.Println("loop1")
		lotteryId, _ := strconv.Atoi(strings.Split(lotteryIdStr,":")[0])
		lottery := Lottery{LotteryId: lotteryId}
		lottery, _ = lottery.QueryById()
		lotteries = append(lotteries, lottery)
	}
	//fmt.Println(len(lotteries))
	return lotteries
}

//查询对应lottery的能否抽奖access
func (customer Customer) QueryLotteryAccessByCustomerId(lotteryIdQuery int) (bool, string){
	customer = customer.QueryByUsername()
	if customer.CustomerLottery==""{
		return false, ""
	}
	lotteryIds := strings.Split(customer.CustomerLottery,";")
	//fmt.Println(len(lotteries))
	//fmt.Println(lotteryIds)
	//fmt.Println(len(lotteryIds))
	for _, lotteryIdStr:=range lotteryIds{
		//fmt.Println("loop1")
		lotteryId, _ := strconv.Atoi(strings.Split(lotteryIdStr,":")[0])
		lotteryAccess := strings.Split(lotteryIdStr,":")[1]
		if lotteryId==lotteryIdQuery{
			return lotteryAccess=="0", lotteryAccess
		}
	}
	return false, ""
}

//customer新增lottery
func (customer Customer) UpdateCustomerLottery(lotteryId int) Customer {
	var newCustomerLottery string
	if customer.CustomerLottery==""{
		newCustomerLottery = customer.CustomerLottery + strconv.Itoa(lotteryId) + ":0"
	}else{
		newCustomerLottery = customer.CustomerLottery + ";" + strconv.Itoa(lotteryId) + ":0"
	}
	initDB.Db.Model(&customer).Where("customer_id = ?", customer.CustomerId).Update("customer_lottery", newCustomerLottery)
	return customer
}

//customer参与对应lottery的抽奖
func (customer Customer) JoinLottery(lotteryId int) Customer {
	lottery := Lottery{
		LotteryId : lotteryId,
	}
	lottery, gift := lottery.Gift()
	var result string
	customer = customer.QueryByUsername()
	//fmt.Println("old",customer.CustomerLottery)
	newCustomerLottery := strings.ReplaceAll(customer.CustomerLottery, strconv.Itoa(lotteryId) + ":0", strconv.Itoa(lotteryId) + ":" + gift)
	initDB.Db.Model(&customer).Where("customer_id = ?", customer.CustomerId).Update("customer_lottery", newCustomerLottery)
	//fmt.Println("new",customer.CustomerLottery)
	if gift!=" "{
		if lottery.LotteryResult==""{
			result = gift+":"+ strconv.Itoa(customer.CustomerId)
		}else{
			result = lottery.LotteryResult+";"+gift+":"+ strconv.Itoa(customer.CustomerId)
		}
		lottery.UpdateLotteryResult(result)
	}
	return customer
}
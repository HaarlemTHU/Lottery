package model

import (
	"Lottery/initDB"
)

type Admin struct {
	AdminId int
	AdminName string
	AdminUsername string `form:"adminUsername"`
	AdminPassword string `form:"adminPassword"`
}

func (admin Admin) QueryByUsername() Admin {
	if errs := initDB.Db.Where(&admin).First(&admin).GetErrors(); len(errs)!=0{
		admin.AdminPassword = ""
		return admin
	}
	return admin
}

func (admin Admin) GetAllLotteries() []Lottery {
	var lotteries []Lottery
	initDB.Db.Find(&lotteries)
	return lotteries
}

package initDB

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

//请自行填写
func init() {
	var err error
	Db, err = gorm.Open()
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
}

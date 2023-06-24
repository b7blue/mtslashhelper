package models

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// *gorm.DB是多线程安全的，您可以在多个例程中使用一个*gorm.DB。您可以将其初始化一次，并在需要时获取它。
var db *gorm.DB = initdb()

// create user 'mtslashhelper'@'localhost' identified by 'mt51a5Hservice';
var dsn string = "mtslashhelper:mt51a5Hservice@tcp(127.0.0.1:3306)/mtslashhelper?charset=utf8mb4&parseTime=True"

func initdb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("数据库连接成功")

		// &Article{}, &User{}, &Subscription{}
		// if err := db.AutoMigrate(&Article{}, &User{}, &Subscription{}, &Cookie{}); err != nil {
		// 	fmt.Println(err)
		// 	return nil
		// }

		// fmt.Println("初始化表完成")
	}

	return db
}

func GetConnection() *gorm.DB {
	return db
}

func Insert(a []Article) {
	if err := db.Create(&a).Error; err != nil {
		f, _ := os.OpenFile("err.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		defer f.Close()
		f.WriteString(err.Error() + "\n")
	} else {
		fmt.Println("插入成功！")
	}
}

func Getarticles(a []Article) {
	db.Find(&a, a)
}

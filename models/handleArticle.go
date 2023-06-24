package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AddArticle(newa Article) {
	db = GetConnection()
	db.Create(newa)
}

// 当没有人订阅这个文章，就删除
func DelArticle(tid int) error {
	db.Delete(&Article{}, tid)
	return db.Error
}

func SetUpdateTime(list []Article) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	for i := range list {
		db.Save(&list[i])
	}
}

func TidExist(tid int) bool {
	result := db.First(&Article{}, tid)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func GetAllArticle() (allA []Article) {
	db := GetConnection()
	db.Find(&allA)
	return
}

func GetArticlesbyTids(tids []int) (a []Article) {
	db.Find(&a, tids)
	return
}

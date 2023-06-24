package models

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 登陆的时候用 1根据email检查用户是否存在 2在1通过后检查密码是否正确
func CheckUser(email, pw string) (int, error) {
	thisUser := User{
		Email: email,
	}
	result := db.Where(&thisUser, "Email").First(&thisUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return -1, gorm.ErrRecordNotFound
	}
	if fmt.Sprintf("%x", sha256.Sum256([]byte(pw))) != thisUser.PW {
		return -1, nil
	}
	return thisUser.Uid, nil
}

// 注册发邮件之前
func IsNewUser(email string) bool {
	thisUser := User{
		Email: email,
	}
	result := db.Where(&thisUser, "Email").First(&thisUser)
	// fmt.Println(thisUser.Uid, result.Error)
	return errors.Is(result.Error, gorm.ErrRecordNotFound)

}

func NewUser(email, pw string) (uid int, err error) {
	newUser := User{
		Email: email,
		PW:    fmt.Sprintf("%x", sha256.Sum256([]byte(pw))),
	}
	result := db.Create(&newUser)
	db.Find(&newUser)
	return newUser.Uid, result.Error
}

// 根据tid查找订阅列表
func GetUidByEmail(email string) (uid int) {
	u := User{Email: email}
	db.Where(&u).Select("Uid").First(&u)
	return u.Uid
}

func DelUser() {

}

func ChangePW() {

}

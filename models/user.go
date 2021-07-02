package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string
	Nickname string
	Password string
}

// 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func Query(N string, P string) (user User, ok bool) {
	DB.Where("username=?", N).Find(&user)
	ok = user.CheckPassword(P)
	return
}

func QueryUserById(id interface{}) (user User) {
	DB.Where("id=?", id).Find(&user)
	return
}

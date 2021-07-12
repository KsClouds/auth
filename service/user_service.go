package service

import (
	"auth/models"
	"auth/util"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DEFAULT_PASSWORD = "12345678"

type UserRegisterService struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=2,max=10"`
	Username string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password string `form:"password" json:"password" binding:"max=40"`
	Superior uint   `form:"superior" json:"superior"`
}

func (service *UserRegisterService) Register() error {
	count := 0
	models.DB.Model(&models.User{}).Where("username = ?", service.Username).Count(&count)
	if count > 0 {
		return errors.New("账号已存在")
	}
	if len(service.Password) == 0 {
		service.Password = DEFAULT_PASSWORD
	}
	var digest, err = util.Digest(service.Password)
	if err != nil {
		return err
	}
	user := &models.User{
		Nickname: service.Nickname,
		Username: service.Username,
		Password: digest,
	}

	tx := models.DB.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if service.Superior != 0 {
		superior := models.Superior{
			UserId:     user.ID,
			SuperiorId: service.Superior,
		}
		if err := tx.Create(&superior).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

type UserPasswordService struct {
	Password    string `form:"password" json:"password" binding:"required,min=8,max=40"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=8,max=40"`
}

// 修改密码
func (service *UserPasswordService) ChangePassword(c *gin.Context, user models.User) error {
	if !util.CheckPassword(service.Password, user.Password) {
		return errors.New("原密码错误")
	}

	var digest, err = util.Digest(service.NewPassword)
	if err != nil {
		return err
	}
	user.Password = digest
	if err := models.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func ResetPassword(c *gin.Context) error {
	var digest, err = util.Digest(DEFAULT_PASSWORD)
	if err != nil {
		return err
	}
	id := c.Param("id")
	return models.DB.Model(models.User{}).Where("id = ?", id).Update("password", digest).Error
}

func ChangeSuperior(c *gin.Context) error {
	id, superId := c.Param("id"), c.Param("superiorId")

	err := models.DB.Where("user_id = ?", id).Delete(&models.Superior{}).Error
	if err != nil {
		return err
	}
	userId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	superiorId, err := strconv.ParseUint(superId, 10, 32)
	if err != nil {
		return err
	}
	superior := &models.Superior{
		UserId:     uint(userId),
		SuperiorId: uint(superiorId),
	}
	return models.DB.Save(superior).Error
}

func DeleteUser(c *gin.Context) error {
	return models.DB.Where("id = ?", c.Param("id")).Delete(&models.User{}).Error
}

type Staff struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Username   string `json:"username"`
	SuperiorId uint   `json:"superiorId"`
	Superior   string `json:"superior"`
	Self       bool   `json:"self"`
}

func GetUserList(uid uint) []Staff {
	var staffList []Staff
	sql := "select a.id, a.nickname, a.username, b.superior_id, c.nickname as superior " +
		"from users a " +
		"left join superiors b on b.user_id = a.id " +
		"left join users c on c.id = b.superior_id " +
		"where a.deleted_at is null"
	models.DB.Raw(sql).Scan(&staffList)
	for i, staff := range staffList {
		if uid == staff.ID {
			staffList[i].Self = true
		}
	}
	return staffList
}

func GetUserListAll() []Staff {
	var staffList []Staff
	sql := "select a.id, a.nickname, a.username, b.superior_id, c.nickname as superior " +
		"from users a " +
		"left join superiors b on b.user_id = a.id " +
		"left join users c on c.id = b.superior_id"
	models.DB.Raw(sql).Scan(&staffList)
	return staffList
}

func UserMe(uid uint) models.User {
	user := &models.User{}
	models.DB.Where("id=?", uid).First(user)
	fmt.Println(user)
	return *user
}

func GetSubList(uid uint) []uint {
	subList := make([]uint, 0)
	superiors := make([]models.Superior, 0)
	models.DB.Where("superior_id=?", uid).Find(&superiors)
	for _, superior := range superiors {
		subList = append(subList, superior.UserId)
	}
	return subList
}

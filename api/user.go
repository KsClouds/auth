package api

import (
	"auth/models"
	"auth/serializer"
	"auth/service"
	"auth/util/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) uint {
	return c.GetUint("uid")
}

// 修改密码
func ChangePassword(c *gin.Context) {
	var service service.UserPasswordService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusOK, serializer.Err(-1, err))
	}
	err := service.ChangePassword(c, models.QueryUserById(CurrentUser(c)))
	if err != nil {
		c.JSON(http.StatusOK, serializer.Err(-1, err))
	}
	c.JSON(http.StatusOK, serializer.Success())
}

func ResetPassword(c *gin.Context) {
	err := service.ResetPassword(c)
	if err != nil {
		c.JSON(http.StatusOK, serializer.Err(-1, err))
	}
	c.JSON(http.StatusOK, serializer.Success())
}

func ChangeSuperior(c *gin.Context) {
	err := service.ChangeSuperior(c)
	if err != nil {
		c.JSON(http.StatusOK, serializer.Err(-1, err))
	} else {
		c.JSON(http.StatusOK, serializer.Success())
	}
}

func DeleteUser(c *gin.Context) {
	err := service.DeleteUser(c)
	if err != nil {
		c.JSON(http.StatusOK, serializer.Err(-1, err))
		return
	}
	c.JSON(http.StatusOK, serializer.Success())
}

func GetList(c *gin.Context) {
	uid := CurrentUser(c)
	res := service.GetUserList(uid)
	c.JSON(http.StatusOK, serializer.Suc(res))
}

func GetListAll(c *gin.Context) {
	res := service.GetUserListAll()
	c.JSON(http.StatusOK, serializer.Suc(res))
}

func UserMe(c *gin.Context) {
	uid := CurrentUser(c)
	user := service.UserMe(uid)
	c.JSON(http.StatusOK, serializer.UserRsp(user))
}

func GetSubList(c *gin.Context) {
	uid := CurrentUser(c)
	subList := service.GetSubList(uid)
	c.JSON(http.StatusOK, serializer.Suc(subList))
}

func CheckStatus(c *gin.Context) {
	userInfo, err := jwt.Get(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
	} else {
		c.JSON(http.StatusOK, userInfo)
	}
}

func Welcome(c *gin.Context) {
	_, err := jwt.Get(c)
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, c.Query("url"))
	} else {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}

func Login(c *gin.Context) {
	user, ok := models.Query(c.PostForm("usr"), c.PostForm("pas"))
	if ok {
		jwt.Set(c, user.ID)
		c.Redirect(http.StatusMovedPermanently, c.PostForm("url"))
	} else {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Name":   c.PostForm("usr"),
			"Status": "失败",
			"Pasw":   c.PostForm("pas"),
		})
	}
}

func Register(c *gin.Context) {
	var registerService service.UserRegisterService
	if err := c.ShouldBind(&registerService); err == nil {
		err := registerService.Register()
		if err != nil {
			c.JSON(http.StatusOK, serializer.Err(-1, err))
		} else {
			c.JSON(http.StatusOK, serializer.Success())
		}
	} else {
		c.JSON(http.StatusOK, err)
	}
}

func Logout(c *gin.Context) {
	jwt.Exp(c)
	c.JSON(http.StatusOK, serializer.Success())
}

package api

import (
	"fmt"
	"m/models"
	"net/http"

	"m/utils"

	"github.com/gin-gonic/gin"
)

type RegistStuct struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	Department uint8  `json:"department"`
}

type LoginStuct struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func GetUsers(c *gin.Context) {
	filterParams := make(map[string]interface{}, 8)
	for k, v := range c.Request.URL.Query() {
		filterParams[k] = v[0]
	}

	users, count, _ := models.GetUsers(&filterParams)
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
		"count":  count,
		"data":   users,
	})

}

func UpdateUser(c *gin.Context) {
	var delInfo []map[string]interface{}
	err := utils.GetRequestJson(c, &delInfo)
	if err != nil {
		c.JSON(500, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("%v", err),
		})
	}
	count := models.UpdateUsers(&delInfo)

	//	var u models.User
	//	_ = c.BindJSON(&u)
	//	count := models.UpdateUsers(&u)

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
		"count":  count,
	})

}

func AddUsers(c *gin.Context) {
	var delInfo []map[string]interface{}
	err := utils.GetRequestJson(c, &delInfo)
	if err != nil {
		c.JSON(500, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("%v", err),
		})
	}

	if err = models.AddUsers(&delInfo); err != nil {
		c.JSON(500, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("%v", err),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
	})

}

func Regist(c *gin.Context) {
	var R RegistStuct
	_ = c.BindJSON(&R)
	user := &models.User{Name: R.Name, Password: R.Password, Department: R.Department}
	userid, err := user.Regist()
	if err != nil {
		c.JSON(400, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("%v", err),
			"data": gin.H{
				"u_id": userid,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "success",
			"data": gin.H{
				"u_id": userid,
			},
		})
	}
}

func Login(c *gin.Context) {
	var L LoginStuct
	_ = c.BindJSON(&L)
	user := &models.User{Name: L.Name, Password: L.Password}
	token, uid, err := user.Login()
	if err != nil {
		c.JSON(400, gin.H{
			"status": -1,
			"msg":    fmt.Sprintf("%v", err),
			"data": gin.H{
				"u_id":  uid,
				"token": token,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "success",
			"data": gin.H{
				"u_id":  uid,
				"token": token,
			},
		})
	}
}

func Test(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "success",
		"data":   claims,
	})
}

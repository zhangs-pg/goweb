package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"m/models/db"
	"m/router/middleware"
	//"github.com/jinzhu/gorm"
)

const (
	PermDataRead   = 0x01
	PerDataWrite   = 0x02
	PermDataDelete = 0x04
	PermDataAdmin  = 0x08
)

type User struct {
	//gorm.Model
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"-"`
	Permission uint64 `json:"permission"`
	Status     uint8  `json:"status"`
	Department uint8  `json:"department"`
	//CreateTime time.Time `json:"create_time"`
	//ModifyTime time.Time `json:"modify_time"`
}

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//获取users
func GetUsers(f *map[string]interface{}) (*[]map[string]interface{}, int, error) {
	result, count, err := GetModels(User{}, f)
	return result, count, err
}

func AddUsers(data *[]map[string]interface{}) error {
	err := AddModels(User{}, data)
	return err
}

//更改users
func UpdateUsers(data *[]map[string]interface{}) int64 {
	count := UpdateModels(User{}, data)
	return count
}

//注册
func (u *User) Regist() (uint, error) {
	var user User
	var err error
	//判断用户名是否注册
	notResigt := db.DEFAULTDB.Where("name = ?", u.Name).First(&user).RecordNotFound()
	//notResigt为false表明读取到了 不能注册
	if !notResigt {
		return 0, errors.New("用户名已注册")
	} else {
		// 否则 附加uuid 密码md5简单加密 注册
		u.Password = MD5V([]byte(u.Password))
		err = db.DEFAULTDB.Create(u).Error
	}
	return user.ID, err
}

//登录
func (u *User) Login() (string, uint, error) {
	var user User
	u.Password = MD5V([]byte(u.Password))

	err := db.DEFAULTDB.Where("name = ?", u.Name).First(&user).Error
	if err != nil {
		return "", 0, errors.New("未注册")
	} else {
		if user.Password != u.Password {
			return "", 0, errors.New("密码错误")
		} else {
			c := middleware.CustomClaims{ID: user.ID, Name: user.Name, Permission: user.Permission}
			token, err := middleware.CreateToken(c)
			return token, user.ID, err
		}
	}
}

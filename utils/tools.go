package utils

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/goinggo/mapstructure"

	"bytes"
	"encoding/json"
	"io/ioutil"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MapToStruct(s interface{}, m map[string]interface{}) {
	if err := mapstructure.Decode(m, &s); err != nil {
		fmt.Println(err)
	}
}

func StructToMap(obj interface{}, data *map[string]interface{}, tag bool) {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	if tag {
		for i := 0; i < obj1.NumField(); i++ {
			if string(obj1.Field(i).Tag.Get("json")) != "-" {
				(*data)[string(obj1.Field(i).Tag.Get("json"))] = obj2.Field(i).Interface()
			}
		}
	} else {
		for i := 0; i < obj1.NumField(); i++ {
			(*data)[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
}

func getRequestBody(context *gin.Context, s interface{}) error {
	body, _ := context.Get("json")
	reqBody, _ := body.(string)

	decoder := json.NewDecoder(bytes.NewReader([]byte(reqBody)))
	//jsonnumber转int: data.(json.Number).Int64()
	//decoder.UseNumber()

	err := decoder.Decode(&s)
	return err
}

func SetBodyJson(context *gin.Context, key string) {

	body, _ := ioutil.ReadAll(context.Request.Body)
	context.Set(key, string(body))

	//body := make([]byte, l)
	//n, _ := context.Request.Body.Read(body)
	//fmt.Println("request body:", n)
	//context.Set(key, string(body[0:n]))
}

func GetRequestJson(c *gin.Context, info *[]map[string]interface{}) error {
	SetBodyJson(c, "json")
	err := getRequestBody(c, info)
	return err
}

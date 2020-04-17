package models

import (
	"errors"
	"fmt"
	"m/models/db"
	"strings"

	"github.com/goinggo/mapstructure"
	//"reflect"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

// sql build where
func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				switch v.(type) {
				case string:
					v := v.(string)
					tmp := strings.Split(v, ",")
					vals = append(vals, tmp)
				default:
					vals = append(vals, v)
				}
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}

func AddModels(model interface{}, data *[]map[string]interface{}) error {
	var err error
	tx := db.DEFAULTDB.Begin()

	switch model.(type) {
	case User:
		var user User
		for _, d := range *data {
			if err = mapstructure.Decode(d, &user); err != nil {
				tx.Rollback()
				return err
			}
			if err = db.DEFAULTDB.Create(&user).Error; err != nil {
				tx.Rollback()
				return err
			}
			user = User{}
		}
	default:
		return errors.New("no model add")
	}

	tx.Commit()
	return nil
}

func UpdateModels(model interface{}, data *[]map[string]interface{}) int64 {
	var count int64 = 0
	var i_d1 interface{}
	var i_d2 float64
	var id int64
	var err bool

	tx := db.DEFAULTDB.Begin()
	switch model.(type) {
	case User:
		var user User
		for _, d := range *data {
			i_d1, _ = (d)["id"]
			i_d2, _ = i_d1.(float64)
			id = int64(i_d2)
			user = User{}

			err = tx.First(&user, id).RecordNotFound()

			if err {
				//fmt.Println("find", err)
				continue
			}
			count += tx.Model(&user).Omit("id").Updates(d).RowsAffected
		}
	}
	tx.Commit()

	return count
}

func DeleteModels(model interface{}, ids *[]int32) int64 {
	var count int64 = 0

	tx := db.DEFAULTDB.Begin()
	switch model.(type) {
	case User:
		count = tx.Table("user").Where("id IN (?)", *ids).Update("is_del", 1).RowsAffected
	}
	tx.Commit()

	return count

}

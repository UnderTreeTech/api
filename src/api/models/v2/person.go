package v2

import (
	"api/dao"
	"api/util"
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
)

type Person struct {
	Id       int32  `json:"-"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Gender   int8   `json:"gender"`
}

func List() (persons []Person) {
	// persons = make([]Person, 0)
	// persons = append(persons, Person{UserId: "helloworld", Age: 18, Username: "johnsun"})
	return nil
}

func Get(userId string) *Person {
	//Go strings are immutable. Concatenating two strings generates a third.
	//Avoid string concatenation by appending into a []bytebu
	buffer := make([]byte, 27)
	buffer = append(buffer, util.USER_INFO...)
	buffer = append(buffer, userId...)
	cacheKey := string(buffer)
	cacheVal, err := dao.GetCache(cacheKey)

	person := new(Person)
	if cacheVal == nil {
		var (
			id              int32
			phone, nickname sql.NullString
			gender          int8
		)
		err = dao.GetDB().QueryRow("select id,phone,nickname,gender from t_user where userId = ?", userId).Scan(&id, &phone, &nickname, &gender)
		if err != nil {
			beego.Error(err.Error())
			return person
		}

		person.Id = id
		person.Gender = gender
		if phone.Valid {
			person.Phone = phone.String
		}

		if nickname.Valid {
			person.Nickname = nickname.String
		}

		cacheVal, err = json.Marshal(person)
		if err != nil {
			beego.Error("marshal error,", err.Error())
		} else {
			err = dao.SetCache(cacheKey, cacheVal)
			if err != nil {
				beego.Error("set cache error,", err.Error())
			}
		}
	}

	err = json.Unmarshal(cacheVal.([]byte), &person)
	if err != nil {
		beego.Error("unmarshal failed,", err.Error())
	}

	return person
}

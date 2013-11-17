package db

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

const (
	COLLECTION_NAME = "IDENTITY"
	INT32_MAX_VALUE = 2147483648
)

type Identity struct {
	Name  string
	Value int64
}

// 获取自增长列数值
// name: 自增长列唯一标识，一般使用集合名称+列名称格式
func IdOf(name string) int32 {
	c := Collection(name)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"Value": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	identity := &Identity{}
	info, err := c.Find(bson.M{"Name": name}).Apply(change, &identity)
	if err != nil {
		log.Println(info, err)
		return -1
	}

	return int32(identity.Value % INT32_MAX_VALUE)
}

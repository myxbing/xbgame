package db

import (
	"configs"
	"labix.org/v2/mgo"
)

var Session *mgo.Session
var Default *mgo.Database

func init() {
	session, err := mgo.Dial(configs.String("db.host"))
	if err == nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	Session = session
	Default = session.DB(configs.String("db.name"))
}

func Collection(name string) *mgo.Collection {
	return Default.C(name)
}

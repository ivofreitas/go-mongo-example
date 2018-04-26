package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type State struct {
	CustomerId  string
	CreatedOn   time.Time
	CreatedBy   string
	Description string
}

type DataStore struct {
	session *mgo.Session
	err     error
}

func (ds *DataStore) getCol(collectionName string) *mgo.Collection {

	ds.session, ds.err = mgo.Dial("localhost:27017")

	if ds.err != nil {
		panic(ds.err)
	}

	return ds.session.DB("c3po_db").C(collectionName)
}

func (ds *DataStore) GetAll() []State {

	var states []State
	ds.err = ds.getCol("state").Find(bson.M{}).All(&states)
	if ds.err != nil {
		panic(ds.err)
	}

	ds.session.Close()
	return states
}

func (ds *DataStore) GetById(customerId string) State {

	var state State
	ds.err = ds.getCol("state").Find(bson.M{"customerId": customerId}).One(&state)
	if ds.err != nil {
		panic(ds.err)
	}

	ds.session.Close()
	return state

}

func (ds *DataStore) CreateOrUpdate(state State) bool {

	_, ds.err = ds.getCol("state").Upsert(
		bson.M{"customerId": state.CustomerId},
		bson.M{"$set": state})
	if ds.err != nil {
		panic(ds.err)
	} else {
		ds.session.Close()
		return true
	}

	ds.session.Close()
	return false
}

func main() {

	ds := DataStore{}

	state := State{CustomerId: "3", CreatedOn: time.Now(), CreatedBy: "Ivo", Description: "6"}

	ds.CreateOrUpdate(state)

	fmt.Println(ds.GetById("3"))

	fmt.Println(ds.GetAll())

}

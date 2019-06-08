package repository

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/kennykarnama/go-mfcc/helper"
)

//KeyValueRepository represents interface
//to interact with key value db
type KeyValueRepository interface {
	Save(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

//BadgerRepo provide interaction with badgerdb
type repository struct {
	DB *badger.DB
}

//NewRepository constructs new repository that
//use badgerdb
func NewRepository(db *badger.DB) KeyValueRepository {
	return &repository{DB: db}
}

//Save value based on the key
func (br *repository) Save(key string, value interface{}) error {
	samples, err := helper.ConformToArrayFloat32(value)
	if err != nil {
		return err
	}
	delim := ","
	val := strings.Trim(strings.Join(strings.Split(fmt.Sprint(samples), " "), delim), "[]")
	err = br.DB.Update(func(txn *badger.Txn) error {
		//log.Printf("Now Processing val %s", val)
		err := txn.Set([]byte(key), []byte(val))
		return err
	})
	return err
}

//Get value from key specified
func (br *repository) Get(key string) (interface{}, error) {
	var res interface{}
	err := br.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valcopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		res = valcopy
		return err
	})
	return res, err
}

package main

import (
	"log"

	"github.com/kennykarnama/go-mfcc/preemphasis"

	"github.com/dgraph-io/badger"

	"github.com/kennykarnama/go-mfcc/mfcc"
)

func main() {

	//Init the DB
	conn := initBadgerDB()
	defer conn.Close()
	//Init repo
	repo := initRepository(conn)
	//First we create the preprocessing step
	//here we introduce the pre-emphasis process
	preEmphasis := preemphasis.NewPreEmphasis(preemphasis.WithAlfa(0.97), preemphasis.WithRepository(repo))
	//We then construct new mfcc object to do the processing
	mfcc := mfcc.NewMFCC(mfcc.WithPreProcessing(preEmphasis), mfcc.WithFilepath("sample_sounds/bird.wav"), mfcc.WithRepository(repo))
	mfcc.Run()
}

func initRepository(db *badger.DB) mfcc.KeyValueRepository {
	repo := mfcc.NewBadgerRepo(db)
	return repo
}

func initBadgerDB() *badger.DB {
	opts := badger.DefaultOptions
	opts.Dir = "db/mfcc"
	opts.ValueDir = "db/mfcc"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//failOnError
func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

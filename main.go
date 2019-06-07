package main

import (
	"fmt"
	"log"

	"github.com/kennykarnama/experiment-sound-proc/preemphasis"

	"github.com/dgraph-io/badger"

	"github.com/kennykarnama/experiment-sound-proc/mfcc"
)

func main() {
	//First we create the preprocessing step
	//here we introduce the pre-emphasis process
	preEmphasis := preemphasis.NewPreEmphasis(preemphasis.WithAlfa(0.97))
	//Init the DB
	conn := initBadgerDB()
	defer conn.Close()
	//Init repo
	repo := initRepository(conn)
	// val, err := repo.Get("pre-processing-result")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(val)

	//We then construct new mfcc object to do the processing
	mfcc := mfcc.NewMFCC(mfcc.WithPreProcessing(preEmphasis), mfcc.WithFilepath("sample_sounds/bird.wav"), mfcc.WithRepository(repo))
	samples := mfcc.Run()
	fmt.Println(mfcc.Processor.SampleRate)

	fmt.Println(samples[0:5])
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

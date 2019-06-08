package preemphasis

import (
	"log"

	"github.com/kennykarnama/go-mfcc/mfcc/plot"
	"github.com/kennykarnama/go-mfcc/mfcc/repository"
	"github.com/kennykarnama/go-mfcc/preprocessing"
)

const (
	//StatusSuccess saved to db
	StatusSuccess = "SUCCESS"
	//Failed means db failed
	Failed = "FAILED"
)

//Options for preemphasis
type Options struct {
	Alfa       float32
	KeyDB      string
	Repository repository.KeyValueRepository
	Plot       plot.Plot
}

//WithAlfa is a function to set the alfa
func WithAlfa(a float32) Option {
	return func(o *Options) {
		o.Alfa = a
	}
}

//SetKeyDB sets the key will be used
//when saving the value inside db
func SetKeyDB(keydb string) Option {
	return func(o *Options) {
		o.KeyDB = keydb
	}
}

//WithRepository sets the repository used
func WithRepository(keydb repository.KeyValueRepository) Option {
	return func(o *Options) {
		o.Repository = keydb
	}
}

//WithPlot sets the plotting process
//Will be called after each process
func WithPlot(plot plot.Plot) Option {
	return func(o *Options) {
		o.Plot = plot
	}
}

//Option is complier for optional pattern
type Option func(*Options)

//PreEmphasis represents the class
type PreEmphasis struct {
	Alfa       float32
	KeyDB      string
	Repository repository.KeyValueRepository
	Plot       plot.Plot
}

//NewPreEmphasis will construct new preemphasis signal
func NewPreEmphasis(options ...Option) preprocessing.PreProcessing {
	args := Options{
		Alfa:       1.0,
		KeyDB:      "pre-emphasis-result",
		Repository: nil,
		Plot:       nil,
	}

	for _, option := range options {
		option(&args)
	}

	return &PreEmphasis{
		Alfa:       args.Alfa,
		KeyDB:      args.KeyDB,
		Repository: args.Repository,
		Plot:       args.Plot,
	}
}

//PreProcess is a function do do preprocessing steps
//to signal vector
func (pr *PreEmphasis) PreProcess(samples []float32) (*preprocessing.Result, error) {
	n := len(samples)
	res := &preprocessing.Result{}
	for i := 1; i < n; i++ {
		samples[i] = samples[i] - pr.Alfa*samples[i-1]
	}
	res.Samples = samples
	if pr.Repository != nil {
		err := pr.Repository.Save(pr.KeyDB, samples)
		if err != nil {
			return nil, err
		}
	}
	if pr.Plot != nil {
		log.Println("Pfdfdf")
		err := pr.Plot.Draw(samples)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

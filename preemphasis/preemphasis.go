package preemphasis

import (
	"github.com/kennykarnama/experiment-sound-proc/mfcc"
	"github.com/kennykarnama/experiment-sound-proc/preprocessing"
)

//Options for preemphasis
type Options struct {
	Alfa       float32
	KeyDB      string
	Repository mfcc.KeyValueRepository
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
func WithRepository(keydb mfcc.KeyValueRepository) Option {
	return func(o *Options) {
		o.Repository = keydb
	}
}

//Option is complier for optional pattern
type Option func(*Options)

//PreEmphasis represents the class
type PreEmphasis struct {
	Alfa       float32
	KeyDB      string
	Repository mfcc.KeyValueRepository
}

//NewPreEmphasis will construct new preemphasis signal
func NewPreEmphasis(options ...Option) preprocessing.PreProcessing {
	args := Options{
		Alfa:       1.0,
		KeyDB:      "pre-emphasis-result",
		Repository: nil,
	}

	for _, option := range options {
		option(&args)
	}

	return &PreEmphasis{Alfa: args.Alfa}
}

//PreProcess is a function do do preprocessing steps
//to signal vector
func (pr *PreEmphasis) PreProcess(samples []float32) ([]float32, error) {
	n := len(samples)
	for i := 1; i < n; i++ {
		samples[i] = samples[i] - pr.Alfa*samples[i-1]
	}
	if pr.Repository != nil {
		err := pr.Repository.Save(pr.KeyDB, samples)
		if err != nil {
			return nil, err
		}
	}
	return samples, nil
}

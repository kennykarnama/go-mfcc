package framing

import (
	"errors"
	"fmt"

	"github.com/kennykarnama/go-mfcc/mfcc/repository"
)

//Options represents each configurable based on the
//struct props
type Options struct {
	KeyPrefix  string
	Repository repository.KeyValueRepository
}

//Option as abstraction of options
type Option func(*Options)

//SetKeyPrefix will set key used to store data
func SetKeyPrefix(keyprefix string) Option {
	return func(args *Options) {
		args.KeyPrefix = keyprefix
	}
}

//WithRepository sets the key value repository
//will be used
func WithRepository(kv repository.KeyValueRepository) Option {
	return func(args *Options) {
		args.Repository = kv
	}
}

//Framing represents the framing process
type Framing struct {
	KeyPrefix  string
	Repository repository.KeyValueRepository
	M          int32
	N          int32
}

//Result represents framing result
type Result struct {
	Frames    int
	KeyPrefix string
}

//NewFraming constructs new framing
func NewFraming(M int32, N int32, options ...Option) *Framing {
	opts := Options{
		KeyPrefix:  "framing-result",
		Repository: nil,
	}
	for _, option := range options {
		option(&opts)
	}
	return &Framing{
		KeyPrefix:  opts.KeyPrefix,
		Repository: opts.Repository,
		M:          M,
		N:          N,
	}
}

//Run the framing process
//Each frames will be saved to database
func (f *Framing) Run(samples []float32) (*Result, error) {

	if f.Repository == nil {
		err := errors.New("Repository must be specified")
		return nil, err
	}

	res := &Result{}

	N := f.N
	M := f.M

	n := len(samples)
	overlap := N - M
	var startIdx int32
	framesProcessed := 0
	for startIdx < int32(n) {
		endIdx := int32(startIdx) + N
		if endIdx >= int32(n) {
			endIdx = int32(n - 1)
		}
		subvector := samples[startIdx:endIdx]
		key := f.KeyPrefix + "-" + fmt.Sprint((framesProcessed + 1))
		err := f.Repository.Save(key, subvector)
		if err != nil {
			return nil, err
		}
		framesProcessed++
		startIdx = endIdx + overlap + 1
	}
	res.Frames = framesProcessed
	res.KeyPrefix = f.KeyPrefix
	return res, nil
}

package mfcc

import (
	"errors"
	"os"

	"github.com/kennykarnama/go-mfcc/helper"

	"github.com/kennykarnama/go-mfcc/preprocessing"

	"github.com/mjibson/go-dsp/wav"
)

//Options represent configurable options
//for mfcc struct
type Options struct {
	Preprocessing preprocessing.PreProcessing
	Filepath      string
	Repository    KeyValueRepository
}

//Option represents compliant func to configuration
type Option func(*Options)

//WithPreProcessing sets the preemphasis processor
func WithPreProcessing(pr preprocessing.PreProcessing) Option {
	return func(args *Options) {
		args.Preprocessing = pr
	}
}

//WithFilepath sets the wav file path
func WithFilepath(filepath string) Option {
	return func(args *Options) {
		args.Filepath = filepath
	}
}

//WithRepository sets the key value repository
//used to save the result for each processes
func WithRepository(kv KeyValueRepository) Option {
	return func(args *Options) {
		args.Repository = kv
	}
}

//MFCC represents
//Container to do mfcc feature extraction
type MFCC struct {
	Processor     *wav.Wav
	PreProcessing preprocessing.PreProcessing
	Filepath      string
	Repository    KeyValueRepository
}

//NewMFCC creates the object of mfcc
func NewMFCC(options ...Option) *MFCC {
	args := Options{
		Filepath:      "",
		Preprocessing: nil,
		Repository:    nil,
	}
	for _, option := range options {
		option(&args)
	}

	file, err := os.Open(args.Filepath)
	helper.FailOnError(err)
	proc, err := wav.New(file)
	helper.FailOnError(err)
	return &MFCC{
		Processor:     proc,
		PreProcessing: args.Preprocessing,
		Repository:    args.Repository,
	}
}

//Run the process
func (mfcc *MFCC) Run() []float32 {
	n := mfcc.Processor.Samples
	samples, err := mfcc.Processor.ReadSamples(n)
	newSamples, err := conformToArrayFloat32(samples)
	helper.FailOnError(err)
	//PreProcessing step
	if mfcc.PreProcessing != nil {
		samples, err = mfcc.PreProcessing.PreProcess(newSamples)
		helper.FailOnError(err)
		if mfcc.Repository != nil {
			err = mfcc.Repository.Save("pre-processing-result", samples)
			helper.FailOnError(err)
		}

	}
	return newSamples
}

//Conform different bit samples type
//to be []float32 without losing precision
func conformToArrayFloat32(data interface{}) ([]float32, error) {
	switch i := data.(type) {
	case []uint8:
		s := make([]float32, len(i))
		for idx, val := range i {
			s[idx] = float32(val)
		}
		return s, nil
	case []uint16:
		s := make([]float32, len(i))
		for idx, val := range i {
			s[idx] = float32(val)
		}
		return s, nil
	case []float32:
		return i, nil
	default:
		return nil, errors.New("Unknown bit sample format")
	}
}

//GetTime returns the time vector
func (mfcc *MFCC) GetTime() float32 {
	//t := ((1.0) / (1.0 * float32(mfcc.Processor.SampleRate))) * float32((uint32(mfcc.Processor.Samples) * 1.0))
	t := float32(mfcc.Processor.Samples) / (float32(mfcc.Processor.SampleRate))
	return t
}

package mfcc

import (
	"log"
	"os"

	"github.com/kennykarnama/go-mfcc/windowing"

	"github.com/kennykarnama/go-mfcc/framing"
	"github.com/kennykarnama/go-mfcc/helper"
	"github.com/kennykarnama/go-mfcc/mfcc/repository"

	"github.com/kennykarnama/go-mfcc/preprocessing"

	"github.com/mjibson/go-dsp/wav"
)

//Options represent configurable options
//for mfcc struct
type Options struct {
	Preprocessing preprocessing.PreProcessing
	Framing       *framing.Framing
	Filepath      string
	Repository    repository.KeyValueRepository
}

//Option represents compliant func to configuration
type Option func(*Options)

//WithPreProcessing sets the preemphasis processor
func WithPreProcessing(pr preprocessing.PreProcessing) Option {
	return func(args *Options) {
		args.Preprocessing = pr
	}
}

//WithFraming sets the framing processor to perform frame blocking
//May be improved using interface
func WithFraming(fr *framing.Framing) Option {
	return func(args *Options) {
		args.Framing = fr
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
func WithRepository(kv repository.KeyValueRepository) Option {
	return func(args *Options) {
		args.Repository = kv
	}
}

//MFCC represents
//Container to do mfcc feature extraction
type MFCC struct {
	Processor     *wav.Wav
	PreProcessing preprocessing.PreProcessing
	Framing       *framing.Framing
	Windowing     windowing.Windowing
	Filepath      string
	Repository    repository.KeyValueRepository
}

//NewMFCC creates the object of mfcc
func NewMFCC(fr *framing.Framing, wndw windowing.Windowing, options ...Option) *MFCC {
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
		Framing:       fr,
		Windowing:     wndw,
	}
}

//Run the process
func (mfcc *MFCC) Run() []float32 {
	n := mfcc.Processor.Samples
	samples, err := mfcc.Processor.ReadSamples(n)
	newSamples, err := helper.ConformToArrayFloat32(samples)
	helper.FailOnError(err)
	//PreProcessing step
	if mfcc.PreProcessing != nil {
		res, err := mfcc.PreProcessing.PreProcess(newSamples)
		helper.FailOnError(err)
		newSamples = res.Samples
	}

	//Frame blocking (mandatory)
	frames, err := mfcc.Framing.Run(newSamples)
	helper.FailOnError(err)
	log.Println(frames.Frames)
	return newSamples
}

//GetTime returns the time vector
func (mfcc *MFCC) GetTime() float32 {
	//t := ((1.0) / (1.0 * float32(mfcc.Processor.SampleRate))) * float32((uint32(mfcc.Processor.Samples) * 1.0))
	t := float32(mfcc.Processor.Samples) / (float32(mfcc.Processor.SampleRate))
	return t
}

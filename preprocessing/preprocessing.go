package preprocessing

type Result struct {
	Samples []float32
}

//PreProcessing provide signal pre-processing steps
type PreProcessing interface {
	PreProcess(samples []float32) (*Result, error)
}

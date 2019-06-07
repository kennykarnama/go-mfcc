package preprocessing

//PreProcessing provide signal pre-processing steps
type PreProcessing interface {
	PreProcess(samples []float32) ([]float32, error)
}

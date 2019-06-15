package windowing

import (
	"errors"

	"github.com/mjibson/go-dsp/window"
)

const (
	//SymetricFuncErr is an error that happends when user define function
	//to generate Symetric points for more than one function
	SymetricFuncErr = "Symetric Function Should Only Be Defined Once"
)

//Windowing interface define common operation
//should be supported when implementing this interface
type Windowing interface {
	Apply(signal []float32, options ...func() int) ([]float64, error)
}

type hamming struct {
	L int32
}

//NewHamming constructs new windowing function
func NewHamming(symetricL int32) Windowing {
	return &hamming{L: symetricL}
}

func (h *hamming) Apply(signal []float32, options ...func() int) ([]float64, error) {

	var symetricL int
	symetricL = int(h.L)
	if len(options) == 1 {
		symetricL = (options[0]())
	} else {
		return nil, errors.New(SymetricFuncErr)
	}

	var s []float64
	for _, val := range signal {
		converted := float64(val)
		s = append(s, converted)
	}
	s = window.Hamming(symetricL)
	return s, nil
}

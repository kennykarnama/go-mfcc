package helper

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

//FailOnError logs to stdout and os.Exit(1)
func FailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//ConformToArrayFloat32 standardize different bit samples type
//to be []float32 without losing precision
func ConformToArrayFloat32(data interface{}) ([]float32, error) {
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

//ConvertToDelimited converts corresponding []float32 to be
//delimited string
func ConvertToDelimited(value interface{}, delim string) (string, error) {
	samples, err := ConformToArrayFloat32(value)
	if err != nil {
		return "", err
	}
	val := strings.Trim(strings.Join(strings.Split(fmt.Sprint(samples), " "), delim), "[]")
	return val, nil
}

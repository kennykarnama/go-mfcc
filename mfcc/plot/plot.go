package plot

import (
	"github.com/kennykarnama/go-mfcc/helper"
	"github.com/kennykarnama/go-mfcc/mfcc/repository"
)

//Plot mimics plot process
type Plot interface {
	Draw(interface{}) error
}

//plot using message queue :_:
type plot struct {
	MessageQueue repository.MessageQueue
}

//NewPlot constructs plot that will use
//MessageQueue to do plotting
func NewPlot(mq repository.MessageQueue) Plot {
	return &plot{MessageQueue: mq}
}

func (p *plot) Draw(data interface{}) error {
	val, err := helper.ConvertToDelimited(data, ",")
	if err != nil {
		return err
	}
	err = p.MessageQueue.Publish(val)
	return err
}

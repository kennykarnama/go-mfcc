package repository

import (
	"log"

	"github.com/streadway/amqp"
)

//MessageQueue is interface to provide
//interaction with messagequeue
type MessageQueue interface {
	Publish(message string) error
}

type rabbitmq struct {
	Channel     *amqp.Channel
	ChannelName string
}

//NewRabbitMQ constructs new message queue that
//uses rabbitmq
func NewRabbitMQ(mq *amqp.Channel, channelName string) MessageQueue {
	return &rabbitmq{Channel: mq, ChannelName: channelName}
}

func (rq *rabbitmq) Publish(message string) error {
	q, err := rq.Channel.QueueDeclare(
		rq.ChannelName, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	body := []byte(message)
	err = rq.Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err == nil {
		log.Println("Sent message")
	}
	return err

}

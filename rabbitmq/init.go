package rabbitmq

import (
	"errors"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
func Connect() error {
	rabbitmqUsername := os.Getenv("rabbitmqusername")
	rabbitmqPassword := os.Getenv("rabbitmqpassword")
	rabbitmqEndpoint := os.Getenv("rabbitmqendpoint")
	endpoint := fmt.Sprintf("amqps://%v:%v@%v", rabbitmqUsername, rabbitmqPassword, rabbitmqEndpoint)
	conn, err := amqp.Dial(endpoint)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("unable to connect to rabbitmq")
	}
	Conn = conn
	return nil
}
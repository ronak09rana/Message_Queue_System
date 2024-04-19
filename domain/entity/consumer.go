package entity

type Consumer struct {
	QueueName    string
	ExchangeName string
	ExchangeType string
	RoutingKey   string
}

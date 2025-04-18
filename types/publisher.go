package types

// PublisherArgument payload argument for sending into queue via topics
type PublisherArgument struct {
	Topic        string
	Key          string
	Queue        string
	ContentType  string
	ExchangeName string
	Header       map[string]interface{}
	Message      []byte
}

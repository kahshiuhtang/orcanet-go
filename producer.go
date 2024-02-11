package client

type Producer struct {
}

func (prod *Producer) SetupServer() bool {
	return false
}

func (prod *Producer) SendFile() bool {
	return false
}

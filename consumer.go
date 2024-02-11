package client

type Consumer struct {
	currentCoins float64
}

func (cons *Consumer) SetupConsumer() bool {
	return false
}

// Use requestFile if the function is internal, otherwise name it RequestFile
func (cons *Consumer) RequestFileFromMarket() bool {
	return false
}

func (cons *Consumer) RequestFileFromProducer() bool {
	return false
}

func (cons *Consumer) SendCurrency() bool {
	return false
}

package types

type TransferRecord struct {
	Id          string  `bson:"_id"`
	Tx          string  `bson:"tx"`
	BlockNumber uint64  `bson:"blockNumer"`
	From        string  `bson:"from"`
	To          string  `bson:"to"`
	Amount      float64 `bson:"amount"`
	Timestamp   string  `bson:"timestamp"`
	TxTime      string  `bson:"txTime"`
	LogTime     string  `bson:"logTime"`
	NotifyTime  string  `bson:"notifyTime"`
}

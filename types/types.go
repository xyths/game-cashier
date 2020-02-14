package types

import "time"

type TransferRecord struct {
	Id          string    `bson:"_id"` // tx id + seq
	Tx          string    `bson:"tx"`  // tx id
	BlockNumber uint64    `bson:"blockNumer"`
	From        string    `bson:"from"`
	To          string    `bson:"to"`
	Amount      float64   `bson:"amount"`
	Memo        string    `bson:"memo"`
	Timestamp   string    `bson:"timestamp"`
	TxTime      time.Time `bson:"txTime"`
	LogTime     time.Time `bson:"logTime"`
	NotifyTime  string    `bson:"notifyTime"`
}

const TransferCollName = "transfer"

type NotifyElement struct {
	Network string
	Memo    string
	Amount  float64
	Tx      string
}

type WithdrawLog struct {
	// auto id
	From   string
	To     string
	Amount string
	Tx     string
	Time   time.Time `bson:"logTime"`
}

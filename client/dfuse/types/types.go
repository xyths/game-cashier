package types

import (
	"time"
)

type EosioDocument struct {
	SearchTransactionsForward struct {
		Cursor               string
		Undo                 bool
		IsIrreversible       bool
		IrreversibleBlockNum uint32
		Block                BlockHeader
		Trace                TransactionTrace
	}
}

type BlockHeader struct {
	Id        string
	Num       uint32
	Timestamp time.Time
	Producer  string
	Confirmed uint32
	Previous  string

	TransactionMRoot string
	ActionMRoot      string
	ScheduleVersion  uint32
	NewProducers     ProducerSchedule
}

type ProducerSchedule struct {
	Version   uint32
	Producers ProducerKey
}

type ProducerKey struct {
	ProducerName    string
	BlockSigningKey string
}

type TransactionStatus string

type TransactionTrace struct {
	Id              string
	Block           BlockHeader
	Status          TransactionStatus
	Receipt         TransactionReceiptHeader
	Elapsed         int64
	NetUsage        uint64
	Scheduled       bool
	ExecutedActions []ActionTrace
	MatchingActions []ActionTrace
	topLevelActions []ActionTrace
	exceptJSON      string // TODO: JSON format
}

type TransactionReceiptHeader struct {
}

type ActionTrace struct {
	JSON map[string]interface{}
}

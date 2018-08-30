package types

import (
	"math/big"
	"time"

	"github.com/go-ShardingBlockchain/common"
)

const (
	SEEK_TASK TaskNumber = iota
	ANNOUNCE_TASK
	CASH_TASK
	COMMENT_TASK
)

type TaskNumber uint32
type Transactions []*Transaction
type BlockNonce [8]byte

type SubBlockHead struct {
	ParentHash common.Hash
	UncleHash  common.Hash
	Coinbase   common.Address
	Root       common.Hash
	TxHash     common.Hash
	TaskNum    TaskNumber
	BlockGroup common.Hash
	Difficulty *big.Int
	Number     *big.Int
	Nonce      BlockNonce
}

type SubBlockBody struct {
	Transactions []*Transaction
	Uncles       []*SubBlockHead
}

//子区块结构
type SubBlock struct {
	subhead      *SubBlockHead
	uncles       []*SubBlockHead
	transactions Transactions

	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

//定义主链区块
type Head struct {
	PreHash       common.Hash
	LeaderAddress common.Address
	TaskList      map[TaskNumber]common.Hash
	BlockRoot     common.Hash
}

//区块结构
type Block struct {
	head  *Head
	block []*SubBlock

	ObtainTime time.Time
	ObtainFrom interface{}
}

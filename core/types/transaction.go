package types

import (
  "math/big"
)

type TransactionType uint64

//add the transaction Type
const (
  SERARCH         TransactionType = iota
  CLAIMINFO
  QOSSERVICE
  CASH
)

//Define the Transaction struct
type Transaction struct {
  data    txdata
}

//Transaction Data
type txdata struct {
  //Transaction Type
  TxType      TransactionType
  Announce    uint64
  From        string
  To          string
  Amount      *big.Int
  info        interface{}

  V           *big.Int
  R           *big.Int
  S           *big.Int
}

func NewTransaction(nonce)

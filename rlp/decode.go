package rlp

import "errors"

var (
	ErrExpectedString = errors.New("rlp: expected String or Byte")
	ErrExpectedList   = errors.New("rlp: expected List")

	ErrCanonSize     = errors.New("rlp: non-canonical size information")
	ErrValueTooLarge = errors.New("rlp: value size exceeds available input length")
)

type Kind int

const (
	Byte Kind = iota
	String
	List
)

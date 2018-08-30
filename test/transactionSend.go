package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/go-ShardingBlockchain/consensus/vrf/vrfcurve"
)

var (
	curve  = elliptic.P256()
	params = curve.Params()
)

func main() {
	fmt.Println("----------------------Transaction Test-----------------------")
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)
	fmt.Println("key = [", key, "]")
	_, proof := vrfcurve.Evaluate(key, []byte("Hello vrf"))
	fmt.Println("proof = [", proof, "]")
	fmt.Println("----------------------Test         End-----------------------")
}

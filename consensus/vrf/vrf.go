package vrf

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"

	"github.com/go-ShardingBlockchain/consensus/vrf/vrfcurve"
)

func VRF(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, []byte, error) {
	_, proof := vrfcurve.Evaluate(privateKey, message)
	if proof == nil {
		return nil, nil, errors.New("Generate Proof failed")
	}
	return proof[len(proof)-32:], proof[:len(proof)-32], nil
}

func VerifyVRF(publicKey *ecdsa.PublicKey, message []byte, value []byte, proof []byte) bool {
	if len(value) != 32 {
		return false
	}
	proof = append(proof, value...)
	if len(proof) != 129 {
		return false
	}
	_, err := vrfcurve.ProofToHash(publicKey, message, proof)
	if err != nil {
		return false
	}
	return true
}

func UniqueID(userID, appID string) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, uint32(len(userID)))
	b.WriteString(userID)
	binary.Write(b, binary.BigEndian, uint32(len(appID)))
	b.WriteString(appID)
	return b.Bytes()
}

package vrf

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/go-ShardingBlockchain/consensus/vrf"
)

var (
	curve  = elliptic.P256()
	params = curve.Params()
)

func TestGenerateKey(t *testing.T) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Error("Failed")
	}
	t.Log("privateKey = [", key, "]")
	t.Log("publicKey = [", key.PublicKey, "]")
}

func TestVRF(t *testing.T) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Error("Failed")
	}
	key1, err2 := ecdsa.GenerateKey(curve, rand.Reader)
	if err2 != nil {
		t.Error("Failed")
	}

	value, proof, err1 := vrf.VRF(key, []byte("Hello VRF!!!"))
	t.Log("len(proof) = [", len(proof), "]")
	t.Log("value = [", value, "]")
	if err1 != nil {
		t.Error("VRF failed; err = [", err1, "]")
	} else {
		t.Log("Pass, the proof = [", proof, "]")
	}
	//验证正确性
	if vrf.VerifyVRF(&key.PublicKey, []byte("Hello VRF!!!"), value, proof) {
		t.Log("Pass VRFVerify")
	} else {
		t.Error("Verify failed")
	}
	//验证消息错误
	if vrf.VerifyVRF(&key.PublicKey, []byte("Hello VRF!!"), value, proof) {
		t.Log("Pass VRFVerify")
	} else {
		t.Error("Verify failed")
	}
	//验证秘钥错误
	if vrf.VerifyVRF(&key1.PublicKey, []byte("Hello VRF"), value, proof) {
		t.Log("Pass VRFVerify")
	} else {
		t.Error("Verify failed")
	}
}

func TestVRF1(t *testing.T) {
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Error("Failed")
	}
	value, proof, err1 := vrf.VRF(key, []byte("Hello VRF!!!"))
	t.Log("len(proof) = [", len(proof), "]")
	t.Log("value = [", value, "]")
	if err1 != nil {
		t.Error("VRF failed; err = [", err1, "]")
	} else {
		t.Log("Pass, the proof = [", proof, "]")
	}
}

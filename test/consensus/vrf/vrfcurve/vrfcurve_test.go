package vrfcurve

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/go-ShardingBlockchain/consensus/vrf/vrfcurve"
)

var (
	//采用的是P256曲线，后续对曲线进行量化
	curve  = elliptic.P256()
	params = curve.Params()
)

func TestEvaluate(t *testing.T) {
	//生成不同的私钥进行验证
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil && key == nil {
		t.Error("PrivateKey generate error")
	}
	key1, err1 := ecdsa.GenerateKey(curve, rand.Reader)
	if err1 != nil {
		t.Error("PrivateKey generate error")
	}
	Index, proof := vrfcurve.Evaluate(key, []byte("hello"))
	t.Log("proof = [", proof[:], "]")
	t.Log("Index = [", Index, "]")
	//对VRF的产生值进行验证
	//1. 正确性验证
	_, err = vrfcurve.ProofToHash(&key.PublicKey, []byte("hello"), proof)
	if err != nil {
		t.Error("Failed")
	} else {
		t.Log("Pass")
	}
	//2.公钥与私钥不批配验证
	_, err = vrfcurve.ProofToHash(&key1.PublicKey, []byte("hello"), proof)
	if err != nil {
		t.Error("Failed:The publickey not match")
	} else {
		t.Log("Pass")
	}

	//3.消息不匹配验证
	_, err = vrfcurve.ProofToHash(&key.PublicKey, []byte("Hello"), proof)
	if err != nil {
		t.Error("Failed: The message not match")
	} else {
		t.Log("Pass")
	}
}

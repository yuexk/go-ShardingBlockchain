package vrfcurve

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"math/big"
)

var (
	ErrPointNotOnCurve = errors.New("the point is not on the p256 curve")
	ErrWrongKeyType    = errors.New("not an ECDSA key")
	ErrNoPEMFound      = errors.New("no PEM block found")
	ErrInvalidVRF      = errors.New("invalid VRF proof")
)

func H1(curve elliptic.Curve, m []byte) (x, y *big.Int) {
	h := sha512.New()
	var i uint32
	byteLen := (curve.Params().BitSize + 7) >> 3
	//x == nil 第一次执行
	//只要第一次执行不失败，就执行一次
	for x == nil && i < 100 {
		h.Reset()
		binary.Write(h, binary.BigEndian, i)
		h.Write(m)
		r := []byte{2}
		r = h.Sum(r)
		x, y = Unmarshal(curve, r[:byteLen+1])
		i++
	}
	return
}

var one = big.NewInt(1)

func H2(curve elliptic.Curve, m []byte) *big.Int {
	byteLen := (curve.Params().BitSize + 7) >> 3
	h := sha512.New()
	for i := uint32(0); ; i++ {
		h.Reset()
		binary.Write(h, binary.BigEndian, i)
		h.Write(m)
		b := h.Sum(nil)
		k := new(big.Int).SetBytes(b[:byteLen])
		if k.Cmp(new(big.Int).Sub(curve.Params().N, one)) == -1 {
			return k.Add(k, one)
		}
	}
}

func Evaluate(pri *ecdsa.PrivateKey, m []byte) (index [32]byte, proof []byte) {
	curve := pri.Curve
	params := curve.Params()
	nilIndex := [32]byte{}

	//生成一个曲线，产生一个私钥
	r, _, _, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nilIndex, nil
	}

	ri := new(big.Int).SetBytes(r)
	//把消息转换为曲线上的点
	Hx, Hy := H1(curve, m)
	//求该点的一个倍点
	sHx, sHy := params.ScalarMult(Hx, Hy, pri.D.Bytes())
	//将点编码成ANSI X9.62指定的格式
	vrf := elliptic.Marshal(curve, sHx, sHy)

	//ScalarBaseMult(r) --->ScalarMult(curve.Gx,curve.Gy,r)
	rGx, rGy := params.ScalarBaseMult(r)
	rHx, rHy := params.ScalarMult(Hx, Hy, r)

	var b bytes.Buffer
	b.Write(elliptic.Marshal(curve, params.Gx, params.Gy))
	b.Write(elliptic.Marshal(curve, Hx, Hy))
	b.Write(elliptic.Marshal(curve, pri.PublicKey.X, pri.PublicKey.Y))
	b.Write(vrf)
	b.Write(elliptic.Marshal(curve, rGx, rGy))
	b.Write(elliptic.Marshal(curve, rHx, rHy))
	s := H2(curve, b.Bytes())

	// t = r−s*k mod N
	t := new(big.Int).Sub(ri, new(big.Int).Mul(s, pri.D))
	t.Mod(t, params.N)

	// Index = H(vrf)
	index = sha256.Sum256(vrf)

	// Write s, t, and vrf to a proof blob. Also write leading zeros before s and t
	// if needed.
	var buf bytes.Buffer
	buf.Write(make([]byte, 32-len(s.Bytes())))
	buf.Write(s.Bytes())
	buf.Write(make([]byte, 32-len(t.Bytes())))
	buf.Write(t.Bytes())
	buf.Write(vrf)

	return index, buf.Bytes()
}

func ProofToHash(pk *ecdsa.PublicKey, m, proof []byte) (index [32]byte, err error) {
	curve := pk.Curve
	params := curve.Params()

	nilIndex := [32]byte{}
	// verifier checks that s == H2(m, [t]G + [s]([k]G), [t]H1(m) + [s]VRF_k(m))
	if got, want := len(proof), 64+65; got != want {
		return nilIndex, ErrInvalidVRF
	}

	// Parse proof into s, t, and vrf.
	s := proof[0:32]
	t := proof[32:64]
	vrf := proof[64 : 64+65]

	uHx, uHy := elliptic.Unmarshal(curve, vrf)
	if uHx == nil {
		return nilIndex, ErrInvalidVRF
	}

	// [t]G + [s]([k]G) = [t+ks]G
	tGx, tGy := params.ScalarBaseMult(t)
	ksGx, ksGy := params.ScalarMult(pk.X, pk.Y, s)
	tksGx, tksGy := params.Add(tGx, tGy, ksGx, ksGy)

	// H = H1(m)
	// [t]H + [s]VRF = [t+ks]H
	Hx, Hy := H1(curve, m)
	tHx, tHy := params.ScalarMult(Hx, Hy, t)
	sHx, sHy := params.ScalarMult(uHx, uHy, s)
	tksHx, tksHy := params.Add(tHx, tHy, sHx, sHy)

	//   H2(G, H, [k]G, VRF, [t]G + [s]([k]G), [t]H + [s]VRF)
	// = H2(G, H, [k]G, VRF, [t+ks]G, [t+ks]H)
	// = H2(G, H, [k]G, VRF, [r]G, [r]H)
	var b bytes.Buffer
	b.Write(elliptic.Marshal(curve, params.Gx, params.Gy))
	b.Write(elliptic.Marshal(curve, Hx, Hy))
	b.Write(elliptic.Marshal(curve, pk.X, pk.Y))
	b.Write(vrf)
	b.Write(elliptic.Marshal(curve, tksGx, tksGy))
	b.Write(elliptic.Marshal(curve, tksHx, tksHy))
	h2 := H2(curve, b.Bytes())

	// Left pad h2 with zeros if needed. This will ensure that h2 is padded
	// the same way s is.
	var buf bytes.Buffer
	buf.Write(make([]byte, 32-len(h2.Bytes())))
	buf.Write(h2.Bytes())

	if !hmac.Equal(s, buf.Bytes()) {
		return nilIndex, ErrInvalidVRF
	}
	return sha256.Sum256(vrf), nil
}

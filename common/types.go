package common

const (
	UINT256_SIZE  = 32
	HashLength    = 32
	AddressLength = 32
)

type Uint256 [UINT256_SIZE]byte
type Hash [HashLength]byte
type Address [AddressLength]byte

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)
	return
}

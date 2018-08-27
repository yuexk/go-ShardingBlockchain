package common

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestCopyBytes(t *testing.T) {
	b := common.CopyBytes([]byte("Hello world"))
	if string(b) != "Hello world" {
		fmt.Println("Test error!!!")
	}
	fmt.Println("Test Rigth!!!!")
}

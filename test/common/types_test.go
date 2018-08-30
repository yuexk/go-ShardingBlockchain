package common

import (
	"testing"

	"github.com/go-ShardingBlockchain/common"
)

func TestCopyBytes(t *testing.T) {
	b := common.CopyBytes([]byte("Hello world"))
	if string(b) != "Hello world" {
		t.Error("Faild")
	}
	t.Log("Pass")
}

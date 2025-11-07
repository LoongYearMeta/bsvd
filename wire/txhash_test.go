// Copyright (c) 2024 The BSV developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"testing"

	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
)

// TestTxHashVersion10 测试版本10交易的三层哈希计算
func TestTxHashVersion10(t *testing.T) {
	// 创建一个版本10的测试交易
	msgTx := NewMsgTx(10)

	// 添加一个测试输入
	prevHash, err := chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("NewHashFromStr: %v", err)
	}

	prevOut := NewOutPoint(prevHash, 0)
	txIn := NewTxIn(prevOut, []byte{0x00})
	msgTx.AddTxIn(txIn)

	// 添加一个测试输出
	txOut := NewTxOut(5000000000, []byte{0x76, 0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0xac})
	msgTx.AddTxOut(txOut)

	// 计算哈希
	hash := msgTx.TxHash()

	// 验证哈希不为空
	if hash.IsEqual(&chainhash.Hash{}) {
		t.Error("TxHash returned empty hash for version 10 transaction")
	}

	t.Logf("Version 10 transaction hash: %s", hash.String())
}

// TestTxHashVersion1 测试版本1交易的标准哈希计算
func TestTxHashVersion1(t *testing.T) {
	// 创建一个版本1的测试交易
	msgTx := NewMsgTx(1)

	// 添加一个测试输入
	prevHash, err := chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("NewHashFromStr: %v", err)
	}

	prevOut := NewOutPoint(prevHash, 0)
	txIn := NewTxIn(prevOut, []byte{0x00})
	msgTx.AddTxIn(txIn)

	// 添加一个测试输出
	txOut := NewTxOut(5000000000, []byte{0x76, 0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0xac})
	msgTx.AddTxOut(txOut)

	// 计算哈希
	hash := msgTx.TxHash()

	// 验证哈希不为空
	if hash.IsEqual(&chainhash.Hash{}) {
		t.Error("TxHash returned empty hash for version 1 transaction")
	}

	t.Logf("Version 1 transaction hash: %s", hash.String())
}

// TestTxHashDifferentVersions 测试不同版本的交易产生不同的哈希
func TestTxHashDifferentVersions(t *testing.T) {
	// 创建两个相同内容但版本不同的交易
	createTx := func(version int32) *MsgTx {
		msgTx := NewMsgTx(version)

		prevHash, _ := chainhash.NewHashFromStr("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		prevOut := NewOutPoint(prevHash, 1)
		txIn := NewTxIn(prevOut, []byte{0x76, 0xa9, 0x14})
		msgTx.AddTxIn(txIn)

		txOut := NewTxOut(1000000, []byte{0x76, 0xa9, 0x14, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x88, 0xac})
		msgTx.AddTxOut(txOut)

		return msgTx
	}

	tx1 := createTx(1)
	tx10 := createTx(10)

	hash1 := tx1.TxHash()
	hash10 := tx10.TxHash()

	// 验证不同版本产生不同的哈希
	if hash1.IsEqual(&hash10) {
		t.Error("Version 1 and version 10 transactions should produce different hashes")
	}

	t.Logf("Version 1 hash:  %s", hash1.String())
	t.Logf("Version 10 hash: %s", hash10.String())
}

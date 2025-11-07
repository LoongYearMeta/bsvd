// Copyright (c) 2024 The BSV developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/bitcoinsv/bsvd/chaincfg/chainhash"
)

// calculateVersion10TxHash 计算版本10交易的三层哈希
func (msg *MsgTx) calculateVersion10TxHash() chainhash.Hash {
	// 1. 准备各部分数据
	var (
		serialization1 []byte // 输入部分
		serialization2 []byte // 脚本部分
		serialization3 []byte // 输出部分
	)

	// 处理输入部分
	for _, input := range msg.TxIn {
		// 序列化: TXID(小端) + VOUT + Sequence
		serialization1 = append(serialization1, input.PreviousOutPoint.Hash[:]...)

		indexBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(indexBytes, input.PreviousOutPoint.Index)
		serialization1 = append(serialization1, indexBytes...)

		sequenceBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(sequenceBytes, input.Sequence)
		serialization1 = append(serialization1, sequenceBytes...)

		// 脚本哈希
		scriptHash := sha256.Sum256(input.SignatureScript)
		serialization2 = append(serialization2, scriptHash[:]...)
	}

	// 处理输出部分
	for _, output := range msg.TxOut {
		valueBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(valueBytes, uint64(output.Value))
		serialization3 = append(serialization3, valueBytes...)

		scriptHash := sha256.Sum256(output.PkScript)
		serialization3 = append(serialization3, scriptHash[:]...)
	}

	// 计算各部分哈希
	hash1 := sha256.Sum256(serialization1)
	hash2 := sha256.Sum256(serialization2)
	hash3 := sha256.Sum256(serialization3)

	// 准备头部数据
	versionBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionBytes, uint32(msg.Version))

	locktimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktimeBytes, msg.LockTime)

	inputCountBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(inputCountBytes, uint32(len(msg.TxIn)))

	outputCountBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(outputCountBytes, uint32(len(msg.TxOut)))

	// 构建最终序列化数据
	finalSerialization := bytes.Join([][]byte{
		versionBytes,
		locktimeBytes,
		inputCountBytes,
		outputCountBytes,
		hash1[:],
		hash2[:],
		hash3[:],
	}, nil)

	// 计算最终TXID (SHA256d)
	firstHash := sha256.Sum256(finalSerialization)
	finalHash := sha256.Sum256(firstHash[:])

	// 将字节数组转换为chainhash.Hash
	var hash chainhash.Hash
	copy(hash[:], finalHash[:])
	return hash
}

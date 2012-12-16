// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This LevelDB Go implementation is based on LevelDB C++ implementation.
// Which contains the following header:
//   Copyright (c) 2011 The LevelDB Authors. All rights reserved.
//   Use of this source code is governed by a BSD-style license that can be
//   found in the LEVELDBCPP_LICENSE file. See the LEVELDBCPP_AUTHORS file
//   for names of contributors.

package block

import (
	"bytes"
	"encoding/binary"
)

type Writer struct {
	restartInterval int

	buf      *bytes.Buffer
	restarts []uint32
	lastKey  []byte
	counter  int
}

func NewWriter(restartInterval int) *Writer {
	return &Writer{
		restartInterval: restartInterval,
		buf:             new(bytes.Buffer),
		restarts:        make([]uint32, 1),
	}
}

func (b *Writer) Add(key, value []byte) {
	// Get key shared prefix
	shared := 0
	if b.counter < b.restartInterval {
		n := len(b.lastKey)
		if n > len(key) {
			n = len(key)
		}
		for shared < n && b.lastKey[shared] == key[shared] {
			shared++
		}
	} else {
		b.counter = 0
		b.restarts = append(b.restarts, uint32(b.buf.Len()))
	}

	// Add "<shared><non_shared><value_size>" to buffer
	var n int
	var tmp = make([]byte, binary.MaxVarintLen32)
	n = binary.PutUvarint(tmp, uint64(shared))
	b.buf.Write(tmp[:n])
	n = binary.PutUvarint(tmp, uint64(len(key)-shared))
	b.buf.Write(tmp[:n])
	n = binary.PutUvarint(tmp, uint64(len(value)))
	b.buf.Write(tmp[:n])

	// Add string delta to buffer followed by value
	b.buf.Write(key[shared:])
	b.buf.Write(value)

	b.lastKey = key
	b.counter++
}

func (b *Writer) Len() int {
	return b.buf.Len() + len(b.restarts)*4 + 4
}

func (b *Writer) Finish() []byte {
	for _, restart := range b.restarts {
		binary.Write(b.buf, binary.LittleEndian, restart)
	}
	binary.Write(b.buf, binary.LittleEndian, uint32(len(b.restarts)))

	return b.buf.Bytes()
}

func (b *Writer) Reset() {
	b.buf.Reset()
	b.restarts = make([]uint32, 1)
	b.lastKey = nil
	b.counter = 0
}
// SPDX-FileCopyrightText: 2017 Avi Deitcher
// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: MIT

package table

// move all bytes to big endian to fix how GPT stores UUIDs
func bytesToUUIDBytes(in []byte) []byte {
	// first 3 sections (4 bytes, 2 bytes, 2 bytes) are little-endian, last 2 section are big-endian
	b := make([]byte, 0, 16)
	b = append(b, in[0:16]...)
	tmpb := b[0:4]
	reverseSlice(tmpb)

	tmpb = b[4:6]
	reverseSlice(tmpb)

	tmpb = b[6:8]
	reverseSlice(tmpb)
	return b
}

package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
)

// Float32ToByte convert float32 to byte
func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

// ByteToFloat32 convert byte to float32
func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

// Float64ToByte convert float64 to byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

// ByteToFloat64 convert byte to float32
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// Round round
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

// RandomInt 随机数 int
func RandomInt(num int) int {
	return rand.Intn(65536) % num
}

// SubString sub string include chinese character
func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// SubStringFromEnd sub string include chinese character
func SubStringFromEnd(str string, begin, end int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// Int2Byte ...
func Int2Byte(data int) []byte {
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	// 网络字节序为大端字节序
	binary.Write(buf, binary.BigEndian, data)
	return buf.Bytes()
}

// Int2ByteLittleEndian ...
func Int2ByteLittleEndian(data int) []byte {
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

// Byte2Int ...
func Byte2Int(data []byte) int {
	buf := bytes.NewBuffer(data)
	// 网络字节序为大端字节序
	var i2 int
	binary.Read(buf, binary.BigEndian, &i2)
	return i2
}

// Byte2IntLittleEndian ...
func Byte2IntLittleEndian(data []byte) int {
	buf := bytes.NewBuffer(data)
	var i2 int
	binary.Read(buf, binary.LittleEndian, &i2)
	return i2
}

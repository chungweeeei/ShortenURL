package helpers

import (
	"errors"
	"strings"
	"sync"
	"time"
)

/*
   Unique ID Requirements:
   - Unique ID
   - Must fit in 64 bit var
   - Numerical ID
   - ID must increase with time
   - Handle large number of requests
   - Availability

   Snow Flake Algorithm
   1. Timestamp Ordering
   2. Uniqueness
   3. DS support

   1 bit	-> unused (符號位，永遠為0)
   41 bits	-> timestamp (毫秒級時間戳)
   10 bits	-> machine id (機器ID，支援1024台機器)
   12 bits	-> sequence number (序列號，每毫秒可生成4096個ID)
*/

const (
	// Epoch start time (2024-01-01 00:00:00 UTC)
	epoch int64 = 1704067200000 // user defined

	// the number of bits allocated for each part
	timestampBits = 41
	machineIDBits = 10
	sequenceBits  = 12

	// maximum values
	maxMachineID = -1 ^ (-1 << machineIDBits) // 1023
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 4095

	// shifts
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits

	// Base62 字符集 (0-9, a-z, A-Z)
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Generator struct {
	mu            sync.Mutex
	lastTimestamp int64
	machineID     int64
	sequence      int64
}

type IDResponse struct {
	ID       int64
	ShortURL string
}

func NewGenerator(machineID int64) (*Generator, error) {

	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID must be between 0 and 1023")
	}

	return &Generator{
		machineID:     machineID,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}

// Generate
func (g *Generator) Generate() (int64, error) {

	g.mu.Lock()
	defer g.mu.Unlock()

	// get current timestamp in milliseconds
	timestamp := g.getCurrentTimestamp()

	// if timestamp is less than last timestamp, clock moved backwards
	if timestamp < g.lastTimestamp {
		return 0, errors.New("clock moved backwards, refusing to generate id")
	}

	// if timestamp is same as last timestamp, increment sequence
	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence

		// sequence exhausted in the same millisecond, wait for next millisecond
		if g.sequence == 0 {
			// wait until next millisecond
			timestamp = g.waitNextMillis(g.lastTimestamp)
		}
	} else {
		// new millisecond, reset sequence
		g.sequence = 0
	}

	// update last timestamp
	g.lastTimestamp = timestamp

	// construct unique ID
	id := ((timestamp - epoch) << timestampShift) |
		(g.machineID << machineIDShift) |
		g.sequence

	return id, nil
}

func (g *Generator) getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

func (g *Generator) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := g.getCurrentTimestamp()
	for timestamp <= lastTimestamp {
		timestamp = g.getCurrentTimestamp()
	}
	return timestamp
}

func (g *Generator) GenerateShortURL() (IDResponse, error) {

	id, err := g.Generate()
	if err != nil {
		return IDResponse{}, err
	}
	return IDResponse{
		ID:       id,
		ShortURL: g.EncodeToBase62(id),
	}, nil
}

func (g *Generator) EncodeToBase62(num int64) string {

	if num == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder
	base := int64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		result.WriteByte(base62Chars[remainder])
		num = num / base
	}

	// 反轉字串（因為我們是從低位到高位計算的）
	encoded := result.String()
	return reverseString(encoded)
}

func (g *Generator) DecodeFromBase62(encoded string) (int64, error) {
	var num int64
	base := int64(len(base62Chars))

	for _, char := range encoded {
		index := strings.IndexRune(base62Chars, char)
		if index == -1 {
			return 0, errors.New("invalid character in encoded string")
		}
		num = num*base + int64(index)
	}

	return num, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package hashtree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	HashSize    = sha256.Size
	HashHexSize = sha256.Size * 2
)

var (
	emptyHash = [HashSize]byte{}
)

func CalculateHash(data []byte) [HashSize]byte {
	return sha256.Sum256(data)
}

func HashFromString(hexHash string) ([HashSize]byte, error) {
	if len(hexHash) != HashHexSize {
		return emptyHash, fmt.Errorf("expecting hash to be length of %d, but provided %d", HashHexSize, len(hexHash))
	}

	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return emptyHash, err
	}
	var fixedSize [HashSize]byte
	copy(fixedSize[:], hash[:])
	return fixedSize, nil
}

func ToHashString(hash [HashSize]byte) string {
	return hex.EncodeToString(hash[:])
}
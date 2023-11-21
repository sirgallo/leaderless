package utils

import "crypto/rand"
import "crypto/sha256"


func GenerateSHA256HashRandom() ([32]byte, error) {
	randomData := make([]byte, 32)
	_, readErr := rand.Read(randomData)
	if readErr != nil { return GetZero[[32]byte](), readErr }

	hasher := sha256.New()
	_, writeErr := hasher.Write(randomData)
	if writeErr != nil { return GetZero[[32]byte](), writeErr }
	
	var hash [32]byte
	copy(hash[:], hasher.Sum(nil))

	return hash, nil
}

func GenerateSHA256HashWithSeed(seed []byte) ([32]byte, error) {
	hasher := sha256.New()
	_, writeErr := hasher.Write(seed)
	if writeErr != nil { return GetZero[[32]byte](), writeErr }
	
	var hash [32]byte
	copy(hash[:], hasher.Sum(nil))

	return hash, nil
}
package utils

import "crypto/rand"
import "crypto/sha256"


func GenerateSHA256HashRandom() ([]byte, error) {
	randomData := make([]byte, 32)
	_, readErr := rand.Read(randomData)
	if readErr != nil { return nil, readErr }

	hasher := sha256.New()
	_, writeErr := hasher.Write(randomData)
	if writeErr != nil { return nil, writeErr }
	
	return hasher.Sum(nil), nil
}

func GenerateSHA256HashWithSeed(seed []byte) ([]byte, error) {
	hasher := sha256.New()
	_, writeErr := hasher.Write(seed)
	if writeErr != nil { return nil, writeErr }
	
	return hasher.Sum(nil), nil
}
package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/blake2b"
)

const (
	handle = "example_user"  // User's chosen handle
	seed   = "my_chain_seed" // Seed for PRNG
)

func main() {
	// Filter handle
	filteredHandle := handle

	// Generate numeric suffix based on PRNG seed
	suffix := generateSuffix(seed)

	// Generate MSA Id using handle and suffix
	msaId := generateMsaId(filteredHandle, suffix)

	// Check if MSA Id is available
	msaAvailable := true

	// If MSA Id is not available, keep generating new suffixes until an available one is found
	for !msaAvailable {
		suffix = generateSuffix(seed)
		msaId = generateMsaId(filteredHandle, suffix)
		msaAvailable = false
	}

	// Create MSA on chain with handle and suffix
	msaCreated := true //createMsaOnChain(filteredHandle, suffix)
	if !msaCreated {
		panic("Failed to create MSA on chain:" + msaId)
	}
	// // If MSA creation failed, retry with new suffix
	// for !msaCreated {
	// 	suffix = generateSuffix(seed)
	// 	msaId = generateMsaId(filteredHandle, suffix)
	// 	msaAvailable = true

	// 	for !msaAvailable {
	// 		suffix = generateSuffix(seed)
	// 		msaId = generateMsaId(filteredHandle, suffix)
	// 		msaAvailable = false
	// 	}

	// 	msaCreated = true //createMsaOnChain(filteredHandle, suffix)
	// }
}

// generateSuffix generates a numeric suffix based on PRNG seed
func generateSuffix(seed string) int {
	// Create a new PRNG with the seed
	source := rand.NewSource(hashString(seed))
	prng := rand.New(source)

	// Generate a random integer between 1 and 1000
	return prng.Intn(1000) + 1
}

// hashString hashes a string using Blake2b and returns it as an int64
func hashString(str string) int64 {
	hasher, _ := blake2b.New256(nil)
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return int64(binary.BigEndian.Uint64(hash))
}

// generateMsaId generates an MSA Id using handle and suffix
func generateMsaId(handle string, suffix int) string {
	handleBytes := []byte(handle)
	suffixBytes := []byte(fmt.Sprintf("%d", suffix))

	// Compute hash of handle + suffix_bytes using Blake2b
	hasher, _ := blake2b.New256(nil)
	hasher.Write(handleBytes)
	hasher.Write(suffixBytes)
	hash := hasher.Sum(nil)

	// Use first 32 bits of hash as MSA Id
	return hex.EncodeToString(hash[:4])
}

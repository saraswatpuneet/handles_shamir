package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

func main() {
	// User ID and suffix
	userID := 12345
	suffix := "handle"

	// Number of handles to generate
	numHandles := 5

	// Generate the Shamir shares for the user ID
	threshold := 3
	shares, err := shamir.Split([]byte(fmt.Sprintf("%d", userID)), numHandles, threshold)
	if err != nil {
		fmt.Println("Error generating shares:", err)
		return
	}

	// Generate the handles from the shares
	handles := make([]string, numHandles)
	for i, share := range shares {
		handle := fmt.Sprintf("%02d-%s-%s", i, suffix, string(share))
		handles[i] = handle
	}

	// Print the handles
	fmt.Println(strings.Join(handles, ", "))

	// Example of how to recover the user ID from 3 handles
	recoveryHandles := []string{handles[0], handles[2], handles[4]}
	recoveryShares := make([][]byte, len(recoveryHandles))
	for i, handle := range recoveryHandles {
		// Parse the handle to get the share index and value
		parts := strings.Split(handle, "-")
		index, _ := strconv.Atoi(parts[0])
		// value is join from index 2 to the end
		value := strings.Join(parts[2:], "")

		// Use the index to recover the share
		recoveryShares[i] = shares[index]

		// Convert the value to a byte slice and check if it matches the recovered share
		valueBytes := make([]byte, len(value))
		for j := 0; j < len(value); j++ {
			valueBytes[j] = value[j]
		}
		if !bytes.Equal(recoveryShares[i], valueBytes) {
			fmt.Println("Failed to recover user ID from handles")
			return
		}
	}

	recoveredBytes, err := shamir.Combine(recoveryShares)
	if err != nil {
		fmt.Println("Failed to recover user ID from handles")
		return
	}
	recoveredUserID, err := strconv.Atoi(string(recoveredBytes))
	if err != nil {
		fmt.Println("Failed to recover user ID from handles")
		return
	}
	fmt.Printf("Recovered user ID: %d\n", recoveredUserID)
}

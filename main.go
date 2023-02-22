package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func main() {
	userID := 12345
	suffix := "@example.com"
	numHandles := 5

	// Split the user ID into shares
	secret := big.NewInt(int64(userID))
	threshold := 3 // would need at least 3 shares to recover the secret
	prime := big.NewInt(2<<10 - 1)
	shares := generateShares(secret, threshold, numHandles, prime)

	// Generate handles from the shares
	handles := make([]string, numHandles)
	for i, share := range shares {
		handleNum := i + 1
		handle := "handle" + strconv.Itoa(handleNum) + suffix
		handle += "-" + share.String()
		handles[i] = handle
	}

	fmt.Println("Generated handles:")
	fmt.Println(strings.Join(handles, ", "))

	// Attempt to recover the secret from all possible combinations of shares
	combinations := getCombinations(shares, threshold)
	var recovered *big.Int
	for _, subset := range combinations {
		recovered = recoverSecret(subset, prime)
		if recovered != nil {
			break
		}
	}

	if recovered == nil {
		fmt.Println("Failed to recover secret.")
		return
	}

	// Generate new shares for the original user ID
	secret = recovered
	shares = generateShares(secret, threshold, numHandles, prime)

	// Generate handles from the new shares
	handles = make([]string, numHandles)
	for i, share := range shares {
		handleNum := i + 1
		handle := "handle" + strconv.Itoa(handleNum) + suffix
		handle += "-" + share.String()
		handles[i] = handle
	}

	fmt.Println("Recovered secret and generated new handles:")
	fmt.Println("Secret:", secret.String())
	fmt.Println("Handles:")
	fmt.Println(strings.Join(handles, ", "))
}

func getCombinations(shares []*big.Int, threshold int) [][]*big.Int {
	if len(shares) < threshold {
		return nil
	}

	var combinations [][]*big.Int

	if threshold == 1 {
		for _, share := range shares {
			combinations = append(combinations, []*big.Int{share})
		}
		return combinations
	}

	for i := 0; i <= len(shares)-threshold; i++ {
		for _, combination := range getCombinations(shares[i+1:], threshold-1) {
			combinations = append(combinations, append(combination, shares[i]))
		}
	}

	return combinations
}

func generateShares(secret *big.Int, threshold int, numShares int, prime *big.Int) []*big.Int {
	// Generate a random set of coefficients for the polynomial
	coefficients := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		coefficients[i] = randomInt(prime)
	}

	// Evaluate the polynomial at each x value to create the shares
	shares := make([]*big.Int, numShares)
	for i := 1; i <= numShares; i++ {
		x := big.NewInt(int64(i))
		shares[i-1] = evaluatePolynomial(secret, coefficients, x, prime)
	}

	return shares
}

func recoverSecret(shares []*big.Int, prime *big.Int) *big.Int {
	threshold := len(shares)/2 + 1
	// Generate the Lagrange basis polynomials for each share
	basisPolynomials := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		xs := make([]*big.Int, threshold)
		ys := make([]*big.Int, threshold)
		for j := 0; j < threshold; j++ {
			idx := i*threshold + j
			xs[j] = big.NewInt(int64(j + 1))
			ys[j] = shares[idx]
		}
		basisPolynomials[i] = interpolatePolynomial(xs, ys, prime)
	}

	// Combine the basis polynomials to recover the secret
	secret := big.NewInt(0)
	for i := 0; i < threshold; i++ {
		term := new(big.Int).Mul(shares[i], basisPolynomials[i])
		secret = new(big.Int).Add(secret, term)
	}
	secret.Mod(secret, prime)

	return secret
}

func randomInt(max *big.Int) *big.Int {
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return n
}

func evaluatePolynomial(secret *big.Int, coefficients []*big.Int, x *big.Int, prime *big.Int) *big.Int {
	y := new(big.Int).Set(secret)
	for i := len(coefficients) - 1; i >= 0; i-- {
		y.Mul(y, x)
		y.Add(y, coefficients[i])
		y.Mod(y, prime)
	}
	return y
}

func interpolatePolynomial(xs []*big.Int, ys []*big.Int, prime *big.Int) *big.Int {
	sum := big.NewInt(0)
	for i := 0; i < len(xs); i++ {
		term := new(big.Int).Set(ys[i])
		for j := 0; j < len(xs); j++ {
			if i == j {
				continue
			}
			divisor := new(big.Int).Sub(xs[i], xs[j])
			term.Mul(term, divisor)
			term.Mod(term, prime)
			divisor.ModInverse(divisor, prime)
			term.Mul(term, divisor)
			term.Mod(term, prime)
		}
		sum.Add(sum, term)
		sum.Mod(sum, prime)
	}
	return sum
}

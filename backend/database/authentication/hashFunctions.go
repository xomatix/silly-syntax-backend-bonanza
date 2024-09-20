package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash verifies a plain text password against a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func swapOddWithLastOdd(input string) string {
	runes := []rune(input)
	length := len(runes)

	// Swapping odd-indexed characters with their corresponding characters from the end
	for i := 1; i < length/2; i += 2 {
		// Compute the corresponding index from the end
		oppositeIndex := length - i - 1

		// Swap the characters
		runes[i], runes[oppositeIndex] = runes[oppositeIndex], runes[i]
	}

	return string(runes)
}

func truncateStringRandomly(input string) string {
	length := len(input)

	if length <= 4 {
		return input
	}

	rand.NewSource(time.Now().UnixNano())
	minLength := 4
	maxLength := length
	truncatedLength := rand.Intn(maxLength-minLength) + minLength

	return input[:truncatedLength]
}

func encodeAndShuffleTimestamp(timestamp int64) string {
	timestampStr := fmt.Sprintf("%d", timestamp)

	shuffledTimestamp := swapOddWithLastOdd(timestampStr)
	firstEncoded := base64.StdEncoding.EncodeToString([]byte(shuffledTimestamp))

	secondShuffle := swapOddWithLastOdd(firstEncoded)
	secondEncoded := base64.StdEncoding.EncodeToString([]byte(secondShuffle))

	return truncateStringRandomly(secondEncoded)
}

// CreateToken creates a Base64-encoded token with shuffled characters
func CreateToken(userID int64, username string) (string, error) {
	// Get the current time and format it
	currentTime := time.Now().Unix() // Use Unix timestamp for simplicity
	shuffledTime := encodeAndShuffleTimestamp(currentTime)
	shuffledTime2 := encodeAndShuffleTimestamp(currentTime)
	// Create a string with userID, username, and timestamp separated by a colon
	data := fmt.Sprintf("%d:%s:%s:%s", userID, shuffledTime, username, shuffledTime2)

	// First encoding and shuffling
	firstShuffle := swapOddWithLastOdd(data)
	firstEncoded := base64.StdEncoding.EncodeToString([]byte(firstShuffle))

	// Second encoding and shuffling
	secondShuffle := swapOddWithLastOdd(firstEncoded)
	secondEncoded := base64.StdEncoding.EncodeToString([]byte(secondShuffle))

	return secondEncoded, nil
}

// DecodeToken decodes a double-encoded Base64 token, extracts userID, and validates it
func DecodeToken(token string) (int64, string, error) {
	// Decode the second Base64 encoding
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, "", err
	}

	secondShuffle := string(decodedBytes)
	secondDecoded := swapOddWithLastOdd(secondShuffle)

	firstDecodedBytes, err := base64.StdEncoding.DecodeString(secondDecoded)
	if err != nil {
		return 0, "", err
	}

	firstShuffle := string(firstDecodedBytes)
	originalData := swapOddWithLastOdd(firstShuffle)

	parts := strings.SplitN(originalData, ":", 4)
	if len(parts) != 4 {
		return 0, "", errors.New("invalid token format")
	}
	var userID int64
	_, err = fmt.Sscanf(parts[0], "%d", &userID)
	if err != nil {
		return 0, "", err
	}

	username := parts[2]

	// The timestamp is in parts[1 and 3], which has been shuffled, truncated, and encoded, but we don't need to decode it here.

	return userID, username, nil
}

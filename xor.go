package kasa

func Scramble(plaintext string) []byte {
	n := len(plaintext)
	payload := make([]byte, n)

	key := byte(0xAB)
	for i := 0; i < n; i++ {
		payload[i] = plaintext[i] ^ key
		key = payload[i]
	}

	return payload
}

// Unscramble turns the response from the Kasa into parsable JSON
// it works in place -- be careful with your buffers
func Unscramble(ciphertext []byte) []byte {
	key := byte(0xAB)
	var nextKey byte

	for i := 0; i < len(ciphertext); i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}
	return ciphertext
}

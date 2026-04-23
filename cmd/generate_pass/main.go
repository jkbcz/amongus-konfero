package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	const count = 200
	const password = ""
	const passLength = 20

	for i := range count {
		hmac := hmac.New(sha256.New, []byte(password))
		rawPlayerId := fmt.Sprintf("%04d", i)
		hmac.Write([]byte(rawPlayerId))
		res := hmac.Sum(nil)
		resHex := hex.EncodeToString(res)[:passLength-4]
		fmt.Println(i, rawPlayerId+resHex)
	}
}

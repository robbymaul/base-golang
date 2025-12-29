package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog/log"
)

const key = "klink-payment-1234" // bebas panjang

// key bebas panjang → di-hash jadi 32 byte AES-256
func getAESKey() []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func EncryptAES(text string) string {
	if text == "" {
		return ""
	}

	block, _ := aes.NewCipher(getAESKey()) // pasti valid karena 32 byte
	iv := getAESKey()[:aes.BlockSize]      // IV fix

	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptAES(cryptoText string) string {
	log.Debug().Msg(fmt.Sprintf("[decrypt aes] crypto text = [%s]", cryptoText))
	if cryptoText == "" {
		return ""
	}

	block, _ := aes.NewCipher(getAESKey())
	iv := getAESKey()[:aes.BlockSize]

	ciphertext, _ := base64.StdEncoding.DecodeString(cryptoText)
	log.Debug().Msg(fmt.Sprintf("[decrypt aes] ciphertext = [%v]", ciphertext))
	plaintext := make([]byte, len(ciphertext))
	log.Debug().Msg(fmt.Sprintf("[decrypt aes] plaintext = [%v]", plaintext))
	cfb := cipher.NewCFBDecrypter(block, iv)
	log.Debug().Msg(fmt.Sprintf("[decrypt aes] cfb = [%v]", cfb))
	cfb.XORKeyStream(plaintext, ciphertext)
	log.Debug().Msg(fmt.Sprintf("[decrypt aes] plaintext = [%v]", plaintext))

	return string(plaintext)
}

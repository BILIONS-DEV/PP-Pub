package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	cmap "source/pkg/crypt/concurrent-map"
	"strings"
)

var decodeChar = map[string]string{"0": "k", "1": "A", "2": "T", "3": "Y", "4": "B", "5": "P", "6": "Z", "7": "U", "8": "K", "9": "M", "A": "0", "T": "1", "Y": "2", "B": "3", "P": "4", "Z": "5", "U": "6", "K": "7", "M": "8", "a": "9", "q": "a", "w": "b", "e": "c", "r": "d", "t": "e", "y": "f", "u": "g", "i": "h", "o": "i", "p": "j", "s": "l", "d": "m", "f": "n", "g": "o", "h": "p", "j": "q", "k": "r", "l": "s", "z": "t", "x": "u", "c": "v", "v": "w", "b": "x", "n": "y", "m": "z", "R": "&", "N": "=", "G": "."}
var encodeChar = map[string]string{"&": "R", ".": "G", "0": "A", "1": "T", "2": "Y", "3": "B", "4": "P", "5": "Z", "6": "U", "7": "K", "8": "M", "9": "a", "=": "N", "A": "1", "B": "4", "K": "8", "M": "9", "P": "5", "T": "2", "U": "7", "Y": "3", "Z": "6", "a": "q", "b": "w", "c": "e", "d": "r", "e": "t", "f": "y", "g": "u", "h": "i", "i": "o", "j": "p", "k": "0", "l": "s", "m": "d", "n": "f", "o": "g", "p": "h", "q": "j", "r": "k", "s": "l", "t": "z", "u": "x", "v": "c", "w": "v", "x": "b", "y": "n", "z": "m"}

var cpcDecodeChar = map[string]string{"K": ".", "A": "0", "X": "1", "u": "2", "v": "3", "J": "4", "D": "5", "P": "6", "g": "7", "e": "8", "z": "9"}
var cpcEncodeChar = map[string]string{".": "K", "0": "A", "1": "X", "2": "u", "3": "v", "4": "J", "5": "D", "6": "P", "7": "g", "8": "e", "9": "z"}

func NewVliCrypt() VliCrypt {
	result := VliCrypt{
		decodeChar: cmap.New(),
		encodeChar: cmap.New(),
	}
	for k, v := range decodeChar {
		result.decodeChar.Set(k, v)
	}
	for k, v := range encodeChar {
		result.encodeChar.Set(k, v)
	}
	return result
}

func NewCPCCrypt() VliCrypt {
	result := VliCrypt{
		decodeChar: cmap.New(),
		encodeChar: cmap.New(),
	}
	for k, v := range cpcDecodeChar {
		result.decodeChar.Set(k, v)
	}
	for k, v := range cpcEncodeChar {
		result.encodeChar.Set(k, v)
	}
	return result
}

type VliCrypt struct {
	decodeChar cmap.ConcurrentMap
	encodeChar cmap.ConcurrentMap
}

func (v *VliCrypt) Decode(encodedStr string) string {

	s := strings.Split(encodedStr, "")
	var result bytes.Buffer

	for _, char := range s {
		if vl, ok := v.decodeChar.Get(char); ok {
			result.WriteString(vl.(string))
		} else {
			result.WriteString(char)
		}
	}

	return result.String()

}

func (v *VliCrypt) Encode(str string) string {

	str = strings.ToLower(str)
	s := strings.Split(str, "")
	var result bytes.Buffer

	for _, char := range s {
		if vl, ok := v.encodeChar.Get(char); ok {
			result.WriteString(vl.(string))
		} else {
			result.WriteString(char)
		}
	}

	return result.String()

}

func CipherEncrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// IV needs to be unique, but doesn't have to be secure.
	// It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	// returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func CipherDecrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	// IV needs to be unique, but doesn't have to be secure.
	// It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

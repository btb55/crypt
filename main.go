// This app encrypts text / decrypts with AES 256.
// Encryption:
// If the key < 32bit, the padding is simply the key
// repeated until it meets such length.
// Decryption:
// If the key < 32bit, the key is 'depadded' until
// a unique, unpadded key is revealed.
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
)

var (
	e, d    bool
	k, text string
)

func init() {
	flag.StringVar(&k, "k", "", "Your encryption key.")
	flag.BoolVar(&e, "e", false, "Flag to encrypt your text.")
	flag.BoolVar(&d, "d", false, "Flag to decrypt your cipher.")
	flag.Parse()

	text = flag.Arg(0)
}

func main() {
	if k == "" || text == "" {
		printHelp()
		return
	}

	if e && !d {
		encrypt()
	} else if d && !e {
		decrypt()
	} else {
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`                                                     
                                             _/      
    _/_/_/  _/  _/_/  _/    _/  _/_/_/    _/_/_/_/   
 _/        _/_/      _/    _/  _/    _/    _/        
_/        _/        _/    _/  _/    _/    _/         
 _/_/_/  _/          _/_/_/  _/_/_/        _/_/      
                        _/  _/                       
		     _/_/  _/ 
                       
A little tool to encrypt / decrypt text.`)
	fmt.Println()
	flag.CommandLine.Usage()
	fmt.Println()
}

func encrypt() {
	k = make32bitkey(k)

	encrypted := encryptAES(k, text)

	var lines string
	for i := 0; i < len(encrypted); i++ {
		lines += "-"
	}
	color.New(color.FgHiRed).Println(lines)
	fmt.Println(encrypted)
	color.New(color.FgHiRed).Println(lines)
}

func decrypt() {
	k = make32bitkey(k)

	decrypted := decryptAES(k, text)

	var lines string
	for i := 0; i < len(decrypted); i++ {
		lines += "-"
	}
	color.New(color.FgGreen).Println(lines)
	fmt.Println(decrypted)
	color.New(color.FgGreen).Println(lines)
}

func make32bitkey(key string) string {
	for len(key) < 32 {
		key += key
	}

	return key[:32]
}

func encryptAES(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("creating cipher: %v", err)
	}

	msg := pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalf("reading iv: %v", err)
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return finalMsg
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func decryptAES(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("creating cipher: %v", err)
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		log.Fatalf("decoding string: %v", err)
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		log.Fatalf("creating cipher: %v", errors.New("blocksize must be multipe of decoded message length"))
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		log.Fatalf("unpadding message: %v", err)
	}

	return string(unpadMsg)
}

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("potentially incorrect key provided")
	}

	return src[:(length - unpadding)], nil
}

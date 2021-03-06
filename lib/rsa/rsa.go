package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/tsaridas/salt-golang/lib/utils"
	"io/ioutil"
	"os"
)

// GeneratePEMKeys and save them
func GeneratePEMKeys(priKeyPath string, pubKeyPath string) {
	if file.Exists(priKeyPath) && file.Exists(pubKeyPath) {
		return
	}
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey
	SavePEMKey(priKeyPath, key)
	SavePublicPEMKey(pubKeyPath, publicKey)
	return

}

// LoadPemKeyFromFile function
func LoadPemKeyFromFile(fileName string) (privateKey *rsa.PrivateKey, err error) {
	privKey, er := ioutil.ReadFile(fileName)
	if er != nil {
		return privateKey, er
	}
	privPem, _ := pem.Decode(privKey)
	privateKey, er = x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if er != nil {
		return privateKey, er
	}
	return

}

// LoadPrivKeyFromString function
func LoadPrivKeyFromString(privKey []byte) (privateKey *rsa.PrivateKey, err error) {
	privPem, _ := pem.Decode(privKey)
	privateKey, er := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if er != nil {
		return privateKey, er
	}
	return
}

// LoadPubKeyFromString function
func LoadPubKeyFromString(pubKey string) (publicKey *rsa.PublicKey, err error) {
	pubPem, _ := pem.Decode([]byte(pubKey))
	parsedKey, er := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if er != nil {
		return publicKey, er
	}
	publicKey, _ = parsedKey.(*rsa.PublicKey)
	if er != nil {
		return publicKey, er
	}
	return
}

// LoadPemPubKeyFromFile function
func LoadPemPubKeyFromFile(fileName string) (publicKey *rsa.PublicKey, err error) {
	pubKey, er := ioutil.ReadFile(fileName)
	if er != nil {
		return publicKey, er
	}
	pubPem, _ := pem.Decode(pubKey)
	parsedKey, er := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if er != nil {
		return publicKey, er
	}
	publicKey, _ = parsedKey.(*rsa.PublicKey)
	if er != nil {
		return publicKey, er
	}
	return

}

// SavePEMKey to file
func SavePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

// SavePublicPEMKey to file
func SavePublicPEMKey(fileName string, pubKey rsa.PublicKey) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(&pubKey)
	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubkeyBytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

// EncryptWithPublicKey function
func EncryptWithPublicKey(ciphertext []byte, pubKey *rsa.PublicKey) []byte {
	hash := sha1.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pubKey, ciphertext, nil)
	if err != nil {
		fmt.Println("Could not encrypt with pub key")
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha1.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		fmt.Printf("Could not decrypt with priv key with error %s\n", err)
	}
	return plaintext
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

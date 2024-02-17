package helpers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string)(string , error){
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil{
		return "", err
	}

	hashedPassword := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	encodedPassword := fmt.Sprintf("%s,%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashedPassword))

	return encodedPassword,  nil
}

func ComparePassword(encodedPassword, password string)bool{
	enodedSaltAndPassword := password
	parts := strings.Split(enodedSaltAndPassword, ".")
	decodedHashedPassword, err := base64.RawStdEncoding.DecodeString(parts[1])

	if err != nil{
		return false
	}

	decodedSalt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil{
		return false
	}

	hashedPassword := argon2.IDKey([]byte(encodedPassword), decodedSalt, 1, 64*1024,4,32)

	return bytes.Equal(hashedPassword, decodedHashedPassword)
}
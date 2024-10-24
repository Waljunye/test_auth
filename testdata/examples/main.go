package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
)

func main() {
	privateKeyBytes, err := ioutil.ReadFile("/Users/matysha/PetProjects/auth/private_key.pem")

	obj := struct {
		foo string
		Bar int
	}{
		foo: "Bar",
		Bar: 100,
	}
	bts, err := json.Marshal(&obj)
	if err != nil {
		panic(err)
	}
	tokenNew := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Subject: string(bts),
	})
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	tokenString, err := tokenNew.SignedString(privateKey)

	fmt.Println(tokenString, err)

	publicKeyBytes, err := ioutil.ReadFile("/Users/matysha/PetProjects/auth/public_key.pem")

	block2, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block2.Bytes)
	if err != nil {
		panic(err)
	}

	tokenParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS512.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		panic(err)
	}
	if claims, ok := tokenParsed.Claims.(jwt.MapClaims); ok {
		subj := claims["sub"].(string)
		var obj struct {
			foo string
			Bar int
		}
		json.Unmarshal([]byte(subj), &obj)

		fmt.Println(obj)
	} else {
		fmt.Println(err)
	}
}

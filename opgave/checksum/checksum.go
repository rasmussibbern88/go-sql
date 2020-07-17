package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	//"golang.org/x/crypto/blake2b"
)

// User is the user

func main() {

	file, err := os.Open("text.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	digest := md5.New()

	if _, err := io.Copy(digest, file); err != nil {
		panic(err)
	}

	hash := digest.Sum(nil)
	//fmt.Println()
	fmt.Printf("%x\n", hash)
	//println(ID, FirstName, LastName, ElectronicMail)
}

//db, err := sql.Open("postgres", "user:password@/dbname")

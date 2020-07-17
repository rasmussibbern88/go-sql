package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func check(e error) { // Convenience function
	if e != nil {
		panic(e)
	}
}

// User is the user, format for the databse
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

const (
	// Username : for database
	Username = "rasmus"
	// Password for database
	Password = "1243"
	// Db name of database default is postgres
	Db = "postgres"
)

func filehash(filename string) []byte {
	//Hash a file based on the filename parameter the hash returns a different value for a single byte difference, which in this case is used to verify that the contents of the file are valid
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	digest := md5.New()

	//io.Copy works in chunks, for working with large files
	if _, err := io.Copy(digest, file); err != nil {
		panic(err)
	}
	hash := digest.Sum(nil)
	fmt.Printf("%x\n", hash)
	return hash
}

func createfile() string {
	//get unused filename in the form of data1.csv, data2.csv etc.
	filename := "data"
	suffix := ".csv"
	count := 1
	for true {
		if _, err := os.Stat(filename + strconv.Itoa(count) + suffix); err != nil {
			filename = filename + strconv.Itoa(count) + suffix
			break
		} else {
			count++
		}
	}
	return filename
}

func main() {
	//create file with 0644 read write permissions, and open for write only, and create file if it does not exist. append will never apply as the file is unique
	filename := createfile()
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		Username, Password, Db)

	db, err := sql.Open("postgres", dbinfo) // Driver and formatted information
	defer db.Close()
	check(err)
	// TODO: User input
	id := "0D00A443-93D3-8573-135D-8946397866A1"
	sqlStatement := `SELECT * FROM users where id=$1`
	rows, err := db.Query(sqlStatement, id) //query row, men er åben til flere linjer
	check(err)
	defer rows.Close()

	w := bufio.NewWriter(f)

	digest := md5.New()

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		check(err)
		//fmt.Println(user)
		towrite := user.ID + "," + user.FirstName + "," + user.LastName + "," + user.Email + "\n"

		io.WriteString(digest, towrite)
		n, err := w.WriteString(towrite)
		check(err)
		fmt.Printf("Wrote %d bytes\n", n)
		//// TODO: Total
	}
	//Flush to ensure everything in the buffer has been written
	w.Flush()

	f.Close()

	writehash := digest.Sum(nil)

	readhash := filehash(filename)

	if bytes.Equal(writehash, readhash) {
		fmt.Printf("equal")
		result, err := db.Exec(`DELETE FROM users where id=$1`, id)
		check(err)
		fmt.Print(result)
	}
}

//db, err := sql.Open("postgres", "user:password@/dbname")

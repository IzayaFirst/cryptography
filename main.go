package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/IzayaFirst/cryptography/cipher"
)

func main() {
	filerc, err := os.Open("data/fng.1000.csv.rot128")
	if err != nil {
		log.Fatal(err)
	}
	defer filerc.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(filerc)
	decryption, err2 := cipher.NewRot128Reader(buf)
	if err2 != nil {
		log.Fatal(err2)
	}
	bufReader := make([]byte, 1, 1)
	stringOfCSV := new(bytes.Buffer)
	for {
		n, err3 := decryption.Read(bufReader)
		if err3 == io.EOF {
			break
		} else {
			stringOfCSV.Write(bufReader)
		}
		log.Println("n", n)
	}
	log.Println(stringOfCSV.String())
}

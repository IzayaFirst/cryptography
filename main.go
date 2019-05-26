package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IzayaFirst/cryptography/cipher"
)

// Charging struct for collect transaction data
type Charging struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	CCNumber string  `json:"cc_number"`
	CVV      string  `json:"cvv"`
	ExpMonth int     `json:"expire_month"`
	ExpYear  int     `json:"expire_year"`
}

func main() {
	cipherText, err := os.Open("data/fng.1000.csv.rot128")
	if err != nil {
		log.Fatal(err)
	}
	defer cipherText.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(cipherText)
	decryption, errDecryption := cipher.NewRot128Reader(buf)
	if errDecryption != nil {
		log.Fatal(errDecryption)
	}
	bufReader := make([]byte, 1)
	planTextBuff := new(bytes.Buffer)
	for {
		_, errReadingDecryption := decryption.Read(bufReader)
		if errReadingDecryption == io.EOF {
			break
		} else {
			planTextBuff.Write(bufReader)
		}
	}
	plainText := planTextBuff.String()
	scanner := bufio.NewScanner(strings.NewReader(plainText))
	rowNum := 0
	charges := make([]Charging, 0)
	for scanner.Scan() {
		rowNum++
		if rowNum == 1 {
			continue
		} else {
			transaction := scanner.Text()
			log.Println(transaction)
			transactionAttribute := strings.Split(transaction, ",")
			customerName := transactionAttribute[0]
			amount, _ := strconv.ParseFloat(transactionAttribute[1], 32)
			log.Println(amount)
			ccNumber := transactionAttribute[2]
			ccv := transactionAttribute[3]
			expMonth, _ := strconv.Atoi(transactionAttribute[4])
			expYear, _ := strconv.Atoi(transactionAttribute[5])
			charge := &Charging{
				Name:     customerName,
				Amount:   amount / 100,
				CCNumber: ccNumber,
				CVV:      ccv,
				ExpMonth: expMonth,
				ExpYear:  expYear,
			}
			charges = append(charges, *charge)
		}
	}
}

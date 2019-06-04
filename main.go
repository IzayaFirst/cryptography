package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

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

func convertToArray(scanner *bufio.Scanner, charges chan []Charging, wg *sync.WaitGroup) {
	rowNum := 0
	charger := make([]Charging, 0)
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
			charger = append(charger, *charge)
		}
	}
	charges <- charger
	close(charges)
	wg.Done()
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
	ch := make(chan string)
	go readBufferFile(decryption, ch)
	var wg sync.WaitGroup
	charges := make(chan []Charging)
	for c := range ch {
		wg.Add(1)
		scanner := bufio.NewScanner(strings.NewReader(c))
		go convertToArray(scanner, charges, &wg)
	}
	for c := range charges {
		fmt.Print(c)
		fmt.Print(len(c))
	}
	wg.Wait()
}

func readBufferFile(cipherTextBuffer *cipher.Rot128Reader, c chan string) {
	bufReader := make([]byte, 1)
	planTextBuff := new(bytes.Buffer)
	for {
		_, errReadingDecryption := cipherTextBuffer.Read(bufReader)
		if errReadingDecryption == io.EOF {
			break
		} else {
			planTextBuff.Write(bufReader)
		}
	}
	plainText := planTextBuff.String()
	c <- plainText
	close(c)
}

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"github.com/chanonchanpiwat/challenge.git/cipher"
	"github.com/chanonchanpiwat/challenge.git/constant"
	"github.com/chanonchanpiwat/challenge.git/utils"
	"github.com/chanonchanpiwat/challenge.git/logger"
	"github.com/chanonchanpiwat/challenge.git/songphapadonationconsumer"
	"github.com/chanonchanpiwat/challenge.git/songphapadonator"
	"github.com/chanonchanpiwat/challenge.git/songphapagenerator"
	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

func main() {

	godotenv.Load()
	filepath := os.Args[1]
	publicKey := os.Getenv("OMISE_PUBLIC_KEY")
	secretKey := os.Getenv("OMISE_SECRET_KEY")
	if len(publicKey) == 0 || len(secretKey) == 0 {
		panic(errors.New("OMISE_PUBLIC_KEY or OMISE_SECRET_KEY not found in process"))
	}

	data, err := os.Open(filepath)
	logger.LogAndExit(err)

	rot128Reader, err := cipher.NewRot128Reader(data)
	logger.LogAndExit(err)

	client, err := omise.NewClient(publicKey, secretKey)
	logger.LogAndExit(err)

	done := make(chan interface{})
	defer close(done)

	fmt.Println("performing donations...")

	// Set throttle sleep 200 ms
	songPhaPaChannel := songphapagenerator.GenerateSongPhaPaChannel(done, csv.NewReader(rot128Reader), 200)

	// TO DO: request is an IO bound optimal goroutine have to be tested for multiple factor
	multiplePipe := make([]<-chan *songphapadonator.ChargeResponse, 20)
	for i := 0; i < 20; i++ {
		multiplePipe[i] = songphapadonator.ChargeAPICall(done, songPhaPaChannel, client)
	}

	donationChargeChannel := utils.FanIn(done, multiplePipe...)

	donationStat := []*songphapadonator.ChargeResponse{}
	for item := range donationChargeChannel {
		donationStat = append(donationStat, item)
	}

	donationResult := songphapadonationconsumer.Consumer(donationStat, 3)
	totalDonation := donationResult.SuccessAmount + donationResult.FailedAmount
	averageDonation := float64(totalDonation) / float64(donationResult.DonationCount)
	currency := constant.ThaiCurrency

	fmt.Printf("%25s %s %14d.00\n", "total received:", currency, totalDonation)
	fmt.Printf("%25s %s %14d.00\n", "successfully donated:", currency, donationResult.SuccessAmount)
	fmt.Printf("%25s %s %14d.00\n", "faulty donation:", currency, donationResult.FailedAmount)
	fmt.Printf("%25s %s %17.2f\n", "average per person:", currency, averageDonation)
	fmt.Printf("%25s \n", "top donors:")
	for _, topDonation := range donationResult.TopNDonator {
		fmt.Printf("%25s %s\n", "", topDonation.Name)
	}
}

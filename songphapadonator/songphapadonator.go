package songphapadonator

import (
	"github.com/chanonchanpiwat/challenge.git/constant"
	"github.com/chanonchanpiwat/challenge.git/songphapagenerator"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type ChargeResponse struct {
	Charge *omise.Charge
	Error  error
}

func Charge(client *omise.Client, songpahPa *songphapagenerator.SongPahPa) *ChargeResponse {
	token, createToken := &omise.Token{}, &operations.CreateToken{
		ExpirationMonth: songpahPa.ExpMonth,
		ExpirationYear:  songpahPa.ExpYear,
		Name:            songpahPa.Name,
		Number:          songpahPa.CCNumber,
		SecurityCode:    songpahPa.CVV,
	}

	// TO DO: depending on type of error API limit command may need send back to channel for retry
	e := client.Do(token, createToken)
	if e != nil {
		return &ChargeResponse{
			Charge: nil,
			Error:  e,
		}
	}

	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   songpahPa.AmountSubunits,
		Currency: constant.ThaiCurrency,
		Card:     token.ID,
	}

	// TO DO: depending on type of error API limit command may need send back to channel for retry
	e = client.Do(charge, createCharge)
	if e != nil {
		return &ChargeResponse{
			Charge: nil,
			Error:  e,
		}
	}

	return &ChargeResponse{
		Charge: charge,
		Error:  nil,
	}
}


func ChargeAPICall(done <- chan interface{}, songPhaPaChannel <- chan *songphapagenerator.SongPahPa, client *omise.Client) <- chan *ChargeResponse {
	chargeResponseChannel := make(chan *ChargeResponse)
	go func ()  {
		defer close(chargeResponseChannel)
		for songpahPa := range songPhaPaChannel {
			select {
			case <- done:
				return
			case chargeResponseChannel <- Charge(client, songpahPa):
			}
		}
	}()

	return chargeResponseChannel

}


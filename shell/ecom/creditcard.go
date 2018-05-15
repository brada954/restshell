package ecom

import (
	"fmt"
	"time"
)

//
// Test Card Brand	Number
// American Express 	370000000000002
// Discover	6011000000000012
// JCB	3088000000000017
// Diners Club/ Carte Blanche	38000000000006
// Visa	4007000000027
//  	4012888818888
//  	4111111111111111
// MasterCard	5424000000000015
//  	2223000010309703
//  	2223000010309711
//
// Look at the following for more testing guidance
// https://developer.authorize.net/hello_world/testing_guide/
//

type CreditCard struct {
	Number     string
	Csc        string
	Expiration string
}

var validVisaCard = CreditCard{
	Number:     "4111111111111111",
	Csc:        "123",
	Expiration: "09/22",
}

var invalidCard = CreditCard{
	Number:     "4522222222222222",
	Csc:        "123",
	Expiration: "09/22",
}

var expiredCard = CreditCard{
	Number:     "4111111111111111",
	Csc:        "123",
	Expiration: "09/16",
}

type CardType int

const (
	ValidVisa CardType = 1 + iota
	ValidMaster
	ValidDiscover
	Invalid
	Expired
)

func GetCreditCard(t CardType) CreditCard {
	switch t {
	case ValidVisa:
		return validVisaCard
	case Invalid:
		return invalidCard
	case Expired:
		return expiredCard
	default:
		return invalidCard
	}
}

var startTime = time.Now()
var uniqueExpiration = startTime

func (c CreditCard) MakeCardExpired() CreditCard {
	c.Expiration = getExpiredString()
	return c
}

func (c CreditCard) MakeCardExpirationUnique() CreditCard {
	c.Expiration = getUniqueExpiredString()
	return c
}

func getUniqueExpiredString() string {
	// Create incrementing expiration dates up to 15 years in future
	// then reset
	uniqueExpiration = uniqueExpiration.Add(time.Hour * 24 * 31)
	if uniqueExpiration.Sub(startTime) > (time.Hour * 24 * 31 * 12 * 15) {
		startTime = time.Now()
		uniqueExpiration = startTime
	}
	return fmt.Sprintf("%d-%02d", uniqueExpiration.Year(), int(uniqueExpiration.Month()))
}

func getExpiredString() string {
	exp := time.Now().Add(time.Hour * 24 * 90 * -1)
	return fmt.Sprintf("%4d-%2d", exp.Year(), exp.Month())
}

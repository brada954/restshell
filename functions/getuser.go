package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/brada954/restshell/shell"
)

type NameDetails struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type LocationDetails struct {
	Street string      `json:"street"`
	City   string      `json:"city"`
	State  string      `json:"state"`
	Zip    json.Number `json:"postcode"`
}

type Consumer struct {
	Gender   string          `json:"gender"`
	Name     NameDetails     `json:"name"`
	Location LocationDetails `json:"location"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Cell     string          `json:"cell"`
}

type Result struct {
	Consumers []Consumer `json:"results"`
}

var client = &http.Client{Timeout: time.Duration(30 * time.Second)}

func init() {
	// Register substitutes
	shell.RegisterSubstitutionHandler("randomuser", "getfirstname", GetFirstNameSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getlastname", GetLastNameSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getphone", GetPhoneSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getemail", GetEmailSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getstreet", GetStreetSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getcity", GetCitySubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getstate", GetStateSubstitute)
	shell.RegisterSubstitutionHandler("randomuser", "getzip", GetZipSubstitute)
}

// GetFirstNameSubstitute -- Get a random first name
func GetFirstNameSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Name.First, c
}

// GetLastNameSubstitute -- Get a random last name
func GetLastNameSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Name.Last, c
}

// GetPhoneSubstitute -- Get a random first name
func GetPhoneSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	phone := c.Phone
	if strings.ToLower(option) == "cell" {
		phone = c.Cell
	}

	seperator := "-"
	switch format {
	case "raw":
		return getRawPhone(phone), c
	case "dots":
		seperator = "."
		fallthrough
	case "dashes":
		phone = getRawPhone(phone)
		switch len(phone) {
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 7:
			return phone[:3] + seperator + phone[3:], c
		case 9:
			fallthrough
		case 10:
			return phone[:3] + seperator + phone[3:6] + seperator + phone[6:], c
		case 11:
			return "+" + phone[:1] + " " + phone[1:4] + seperator + phone[4:7] + seperator + phone[7:], c
		case 12:
			return "+" + phone[:2] + " " + phone[2:5] + seperator + phone[5:8] + seperator + phone[8:], c
		case 13:
			return "+" + phone[:3] + " " + phone[3:6] + seperator + phone[6:9] + seperator + phone[9:], c
		default:
			return phone, c
		}
	default:
		return phone, c
	}
}

// GetEmailSubstitute -- Get a random last name
func GetEmailSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Email, c
}

// GetStreetSubstitute -- Get a random last name
func GetStreetSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Location.Street, c
}

// GetCitySubstitute -- Get a random city
func GetCitySubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Location.City, c
}

// GetStateSubstitute -- Get a random state
func GetStateSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return c.Location.State, c
}

// GetZipSubstitute -- Get a random zip code
func GetZipSubstitute(cache interface{}, subname string, format string, option string) (value string, data interface{}) {
	c := getConsumer(cache)
	return string(c.Location.Zip), c
}

func getConsumer(cache interface{}) *Consumer {
	if d, ok := cache.(*Consumer); !ok {
		c, err := getRandomUserData()
		if err != nil {
			panic("User data not available")
		}
		return &c
	} else {
		return d
	}
}

func getRandomUserData() (Consumer, error) {
	result := Result{}
	consumer := Consumer{}

	url := "https://randomuser.me/api?nat=us"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return consumer, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return consumer, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return consumer, err
	}

	err = json.Unmarshal(body, &result)
	if len(result.Consumers) == 0 {
		return consumer, errors.New("No users returned")
	}
	fmt.Printf("%v\n", result.Consumers[0])
	return result.Consumers[0], err
}

func getRawPhone(phone string) string {

	// Make a Regex to say we only want
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(phone, "")

}

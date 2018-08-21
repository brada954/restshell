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

var GetFirstNameDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getfirstname",
	FunctionHelp:      "Get a random first name of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetFirstNameSubstitute,
}

var GetLastNameDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getlastname",
	FunctionHelp:      "Get a random last name of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetLastNameSubstitute,
}

var GetPhoneDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getphone",
	FunctionHelp:      "Get a random phone number of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetPhoneSubstitute,
}

var GetEmailDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getemail",
	FunctionHelp:      "Get a random email address of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetEmailSubstitute,
}

var GetStreetDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getstreet",
	FunctionHelp:      "Get a random street of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetStreetSubstitute,
}

var GetCityDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getcity",
	FunctionHelp:      "Get a random city of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetCitySubstitute,
}

var GetStateDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getstate",
	FunctionHelp:      "Get a random state of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetStateSubstitute,
}

var GetZipDefinition = shell.SubstitutionFunction{
	Group:             "randomuser",
	Name:              "getzip",
	FunctionHelp:      "Get a random zip code of a consumer",
	Formats:           nil,
	OptionDescription: "",
	Options:           nil,
	Function:          GetZipSubstitute,
}

func init() {
	// Register substitutes
	shell.RegisterSubstitutionHandler(GetFirstNameDefinition)
	shell.RegisterSubstitutionHandler(GetLastNameDefinition)
	shell.RegisterSubstitutionHandler(GetPhoneDefinition)
	shell.RegisterSubstitutionHandler(GetEmailDefinition)
	shell.RegisterSubstitutionHandler(GetStreetDefinition)
	shell.RegisterSubstitutionHandler(GetCityDefinition)
	shell.RegisterSubstitutionHandler(GetStateDefinition)
	shell.RegisterSubstitutionHandler(GetZipDefinition)
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
			panic("User data not available: " + err.Error())
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

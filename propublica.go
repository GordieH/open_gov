package main

// propublica api - https://api.propublica.org/congress/v1/
// google civic information api - https://www.googleapis.com/civicinfo/v2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tkanos/gonfig"
	//"reflect"
)

var civicKey = "AIzaSyCGq0GWsj2iDVjr-D201iBOLlk5iRNNqlw"

// NormalizedInput ... office of a representative
type NormalizedInput struct {
	Line1 string `json:"line1"`
	City  string `json:"city"`
	State string `json:"state"`
	Zip   string `json:"zip"`
}

// Office ... office of a representative
type Office struct {
	Name            string   `json:"name"`
	DivisionID      string   `json:"divisionId"`
	Levels          []string `json:"levels"`
	Roles           []string `json:"roles"`
	OfficialIndices []int    `json:"officialIndices"`
}

// Official ... official in the matching position
type Official struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Party    string `json:"party"`
	Phone    string `json:"phone"`
	Urls     string `json:"urls"`
	PhotoURL string `json:"photoUrl"`
	Channels string `json:"channels"`
}

type civicResponse struct {
	NormalizedInput NormalizedInput        `json:"normalizedInput"`
	Kind            string                 `json:"kind"`
	Divisions       map[string]interface{} `json:"divisions"`
	Offices         []Office               `json:"offices"`
	Officials       []Official             `json:"officials"`
}

// Configuration ... configuration data
type Configuration struct {
	KeyName  string
	KeyValue string
}

var googleCivic civicResponse

func main() {

	configuration := Configuration{}
	err := gonfig.GetConf("./data/config.json", &configuration)

	// read user input of address
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input Your Address: ")
	address, _ := reader.ReadString('\n')
	address = strings.Replace(address, " ", "%20", -1)
	address = strings.Replace(address, "\n", "", -1)

	resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=" + address + "&includeOffices=true&key=" + configuration.KeyValue)

	// resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=37%20ibis%20dr%20akron%20ohio&includeOffices=true&key=" + civicKey)
	// resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=80204&includeOffices=true&key=" + configuration.KeyValue)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)

	//fmt.Print(string(body))
	// err = ioutil.WriteFile("denver.json", body, 0644)

	// jsonInputFile := "denver.json"

	// Open our jsonFile
	//jsonFile, err := os.Open(jsonInputFile)
	// if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Successfully Opened ibis.json")
	// defer the closing of our jsonFile so that we can parse it later on
	// defer jsonFile.Close()

	// byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &googleCivic)

	officeMap := make(map[int]string) // map for representatives
	officeDivisionMap := make(map[string]string)
	divisionMap := make(map[string]string)

	// map office names to office indices
	for i := range googleCivic.Offices {
		// fmt.Println(googleCivic.Offices[i].Name)
		for j := range googleCivic.Offices[i].OfficialIndices {
			// fmt.Println(strconv.Itoa(googleCivic.Offices[i].OfficialIndices[j]))
			officeMap[googleCivic.Offices[i].OfficialIndices[j]] = googleCivic.Offices[i].Name
		}
		officeDivisionMap[googleCivic.Offices[i].Name] = googleCivic.Offices[i].DivisionID
	}

	for key, value := range googleCivic.Divisions {
		fmt.Println("key: ", key)
		fmt.Println("RAW: ", value)
		for key1, value1 := range value.(map[string]interface{}) {
			fmt.Println("key1: ", key1)
			fmt.Println("value1: ", value1)
			if key1 == "name" {
				divisionName := fmt.Sprintf("%v", value1)
				divisionMap[key] = divisionName
			}
		}
	}

	fmt.Println("")
	fmt.Println("")
	for key, value := range divisionMap {
		// fmt.Println(i)
		fmt.Println("Key: " + key + " value: " + value)
	}

	fmt.Println("")
	fmt.Println("")
	for i := 0; i < len(googleCivic.Officials); i++ {
		// fmt.Println(i)
		fmt.Println(strconv.Itoa(i) + " - " + officeMap[i] + " - " + googleCivic.Officials[i].Name + " - " + divisionMap[officeDivisionMap[officeMap[i]]])
	}

}

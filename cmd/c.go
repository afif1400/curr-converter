package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// cCmd represents the c command
var cCmd = &cobra.Command{
	Use:   "convert",
	Short: "use this command to convert",
	Long:  `This command is used to convert INR to any other currency in the world`,
	Run: func(cmd *cobra.Command, args []string) {
		// to run convert function
		inrToOther(args)
	},
}

type Converted map[string]float64

func init() {
	rootCmd.AddCommand(cCmd)
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}

func inrToOther(args []string) {
	val, err := strconv.ParseFloat(args[0], 64)

	if err != nil {
		log.Printf("could not parse the float value - %v", err)
		return
	}

	to := args[1]

	if &to == nil {
		log.Printf("please specify the currency")
		return
	}

	responseBytes := conversion(to)
	converted := Converted{}
	err = json.Unmarshal(responseBytes, &converted)
	if len(converted) == 0 {
		log.Printf("Enter a valid currency")
		return
	}
	if err != nil {
		log.Printf("could not marshal response - %v", err)
	}

	fmt.Println(converted["INR_"+to] * val)
}

func conversion(to string) []byte {
	url := "https://free.currconv.com/api/v7/convert?q=INR_"
	APIKEY := os.Getenv("CURR_APIKEY")
	urlParams := "&compact=ultra&apiKey=" + APIKEY
	requestUrl := url + to + urlParams

	request, err := http.NewRequest(
		http.MethodGet, requestUrl, nil)

	if err != nil {
		log.Printf("could not request a conversion rate - %v", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("could not make a request to server - %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf("could not read the response body - %v", err)
	}

	return responseBytes
}

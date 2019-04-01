package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func executeQuery(medallions string, date string, ignoreCache string) {
	ignoreCacheParam := "false"
	if strings.Compare("Y", ignoreCache) == 0 {
		ignoreCacheParam = "true"
	}

	var buf bytes.Buffer
	buf.WriteString("http://152.28.1.2:8080/insights/v1/trips?medallionIDs=")
	buf.WriteString(medallions)
	buf.WriteString("&tripDate=")
	buf.WriteString(date)
	buf.WriteString("&ignoreCache=")
	buf.WriteString(ignoreCacheParam)
	url := buf.String()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func clearCache(reader *bufio.Reader) {
	resp, err := http.Post("http://152.28.1.2:8080/insights/v1/trips/clearCache", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	fmt.Println("To start again enter 1 for Query")
}

func handleQueryInputs(reader *bufio.Reader) {
	fmt.Println("Enter Trip Date in dd/mm/yyyy format")
	date, _ := reader.ReadString('\n')
	date = strings.Replace(date, "\n", "", -1)
	fmt.Println("Enter 1 or more Medallions in a comma separated format, without spaces")
	medallions, _ := reader.ReadString('\n')
	medallions = strings.Replace(medallions, "\n", "", -1)
	fmt.Println("Ignore cached data? Y to ignore or N to use cache")
	ignoreCache, _ := reader.ReadString('\n')
	ignoreCache = strings.Replace(ignoreCache, "\n", "", -1)
	executeQuery(medallions, date, ignoreCache)
	fmt.Println("To start again enter 1 for Query or 2 for Clear Cache")
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Get NY Cab Insights")
	fmt.Println("---------------------")
	fmt.Println("1 Query")
	fmt.Println("2 Clear Cache")
	fmt.Println("Please enter 1 or 2 to proceed")

	for {
		fmt.Print("-> ")
		option, _ := reader.ReadString('\n')
		// convert CRLF to LF
		option = strings.Replace(option, "\n", "", -1)

		if strings.Compare("1", option) == 0 {
			handleQueryInputs(reader)
		}

		if strings.Compare("2", option) == 0 {
			clearCache(reader)
		}

	}

}

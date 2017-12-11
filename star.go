package main

import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/jaytaylor/html2text"
	"github.com/kyokomi/emoji"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var githubURL = "https://github.com/#/#"

// command params
var (
	username, reponame string
	names              bool
)

// results
var (
	stars int
)

var (
	boldred   = color.New(color.FgRed, color.Bold)
	boldgreen = color.New(color.FgHiGreen, color.Bold)
)

// bind flags to params
func bindFlags() {
	flag.StringVar(&username, "u", "", "Give a username/organisation ex: dragonzurfer")
	flag.StringVar(&reponame, "r", "", "Give a repository name")
	flag.Parse()
}

// Get HTML page as string
func getContent(URL string) (string, bool) {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error fetching page")
		return "", true
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ret := string(body)

	if err != nil {
		fmt.Println("Error :( Try Again")
		return "", true
	}

	return ret, false
}

func loader() {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner
	time.Sleep(2 * time.Second)                                 // Run for some time to simulate work
	s.Stop()
}

// Print to terminal
func PrintParams() {
	// print the number of stars recieved

	boldgreen.Printf("  %v's got %d", reponame, stars)
	emoji.Printf(" :star2:'s \n")

	if stars >= 10 && stars < 100 {

		emoji.Println("  AWSM! :clap: ")

	} else if stars < 1000 {

		emoji.Printf("  %v Deserves a :beer:\n", username)

	} else {

		emoji.Printf("  WOW! %v is popular :lollipop:\n", reponame)

	}

}

func main() {
	var url, resp string
	var err bool
	var errconv error

	bindFlags()

	if username == "" || reponame == "" {
		boldred.Println("-u or -r parameter is empty!")
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	//get url content
	go func() {
		defer wg.Done()
		url = strings.Replace(githubURL, "#", username, 1)
		url = strings.Replace(url, "#", reponame, 1)
		resp, err = getContent(url)

		if err {
			boldred.Printf("Error fetching ")
			boldgreen.Printf("stars\n")
			return
		}
	}()

	loader()
	wg.Wait()

	//extract stars
	text, _ := html2text.FromString(resp)
	regex := `\*\s+Star[^)]*\)\s([\d,]+)`
	re := regexp.MustCompile(regex)
	result := re.FindAllStringSubmatch(text, -1)
	stars, errconv = strconv.Atoi(strings.Replace(result[0][1], ",", "", -1))

	if errconv != nil {
		boldred.Println("Error")
		return
	}

	//print to screen
	PrintParams()
}

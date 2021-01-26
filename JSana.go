package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var maxTimeout int
var urlArray []string
var innerHTMlfound []string
var newInnerHTMLfound []string
var htmlFound []string
var newHTMLFound []string
var evalFound []string
var newEvalFound []string
var dangerousFound []string
var newDangeorusFound []string
var s3BucketsFound []string
var newS3BucketsFound []string
var newKeyFound []string
var keyFound []string
var dirName string
var defaultURLsfound string
var jsURLtoFind string
var URL string

var defaultS3Name string
var defaultInnerHTMLName string
var defaultHTMLName string
var defaultEvalName string
var defaultDangerouslySetInnerHTML string
var defaultKeyName string

var reset string = "\033[0m"
var red string = "\033[31m"
var green string = "\033[32m"
var yellow string = "\033[33m"
var blue string = "\033[34m"
var purple string = "\033[35m"
var cyan string = "\033[36m"
var gray string = "\033[37m"
var white string = "\033[97m"

func help() {
	fmt.Println("Displaying the help page")
	fmt.Println("Usage: JSana -u <FILE>")
	fmt.Println("-u <FILE> : A file of js url's one on each line, (eg output from my tool 'wriggle')")
	fmt.Println("-v : verbose mode, not advisiable unless you love spam")
	fmt.Println("-j <URL STRING> : got A link but dont know where it came from? pass the url here")
	fmt.Println("-h : Display this help page")
	os.Exit(3)
}

func inArray(arr []string, toFind string) bool {
	var answer bool = false
	for i := 0; i < len(arr); i++ {
		if arr[i] == toFind {
			answer = true
		}
	}
	return answer
}

func writeToFile(fileName string, listOfStrings []string) {
	if len(listOfStrings) > 0 {
		fileName = dirName + "/" + fileName
		f, err := os.OpenFile(fileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		for i := 0; i < len(listOfStrings); i++ {
			if _, err := f.WriteString(listOfStrings[i] + "\n"); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func processURL(url string) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Println("GetXYZ")
		if strings.Contains(err.Error(), "Client.Timeout") {
			fmt.Println(red + "[Warning]" + reset + " " + white + "The get request has timed out, either increase max timeout or check if the site is up : " + url + reset)
			return
		}
		if strings.Contains(err.Error(), "connection reset by") {
			fmt.Println(red + "[Warning]" + reset + " " + white + "The connection was reset by peer : " + url + reset)
			return
		}
		fmt.Println(red + "[Warning]" + reset + " " + white + err.Error() + " : " + url + reset)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create safe file name
	base64Name := base64.StdEncoding.EncodeToString([]byte(url))
	fullPath := dirName + "/" + base64Name
	putIntoFile(fullPath, resp.Body)

	fileString := string(body)

	extractInterestingStrings(fileString, url)
	clean(fullPath)

}

func clean(nameOfFile string) {
	err := os.Remove(nameOfFile)
	if err != nil {
		fmt.Println(err)
	}
}

func extractInterestingStrings(fileString string, url string) {
	fileString = strings.ToLower(fileString)
	//InnerHTML finder
	if strings.Contains(fileString, ".innerhtml=") || strings.Contains(fileString, ".innerhtml =") {
		tmpString := "innerHTML : " + url
		if !inArray(newInnerHTMLfound, tmpString) && !inArray(innerHTMlfound, tmpString) {
			newInnerHTMLfound = append(newInnerHTMLfound, tmpString)
		}
	}

	//Jquery Html() finder
	myRegex, _ := regexp.Compile(`\.html\(.+\"*\".+\)`)
	found := myRegex.FindAllString(fileString, -1)
	if len(found) > 0 {
		tmpString := strings.Join(found, " ") + " : " + url
		if !inArray(newHTMLFound, tmpString) && !inArray(htmlFound, tmpString) {
			newHTMLFound = append(newHTMLFound, tmpString)
			fmt.Println("boo")
		}
	}

	//Eval finder
	if strings.Contains(fileString, "eval(") {

		tmpString := "eval( : " + url
		if !inArray(newEvalFound, tmpString) && !inArray(evalFound, tmpString) {
			newEvalFound = append(newEvalFound, tmpString)
		}
	}

	//dangerouslySetInnerHTML finder
	if strings.Contains(fileString, ".dangerouslysetinnerhtml=") || strings.Contains(fileString, ".dangerouslysetinnerhtml =") {

		tmpString := ".dangerouslySetInnerHTML : " + url
		if !inArray(newDangeorusFound, tmpString) && !inArray(dangerousFound, tmpString) {
			newDangeorusFound = append(newDangeorusFound, tmpString)
		}
	}

	//s3 buckets finder
	if strings.Contains(fileString, "s3.") || strings.Contains(fileString, "s3-") {
		tmpString := "s3. / s3- : " + url
		if !inArray(newS3BucketsFound, tmpString) && !inArray(s3BucketsFound, tmpString) {
			newS3BucketsFound = append(newS3BucketsFound, tmpString)
		}
	}

	//key finder
	if strings.Contains(fileString, "apikey") || strings.Contains(fileString, "api_key") || strings.Contains(fileString, "api key") {
		tmpString := "apikey / api_key / api key : " + url
		if !inArray(newKeyFound, tmpString) && !inArray(keyFound, tmpString) {
			newKeyFound = append(newKeyFound, tmpString)
		}
	}

}

func putIntoFile(fullPath string, respBody io.ReadCloser) {
	// create the file
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, respBody)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	startTimetime := time.Now()
	startTime := startTimetime.String()
	startTime = startTime[:19]
	startTime = strings.Replace(startTime, " ", "_", 1)
	fmt.Println(startTime)

	defaultS3Name = "S3of" + startTime
	defaultInnerHTMLName = "InnerHTMLof" + startTime
	defaultHTMLName = "htmlOf" + startTime
	defaultEvalName = "EvalOf" + startTime
	defaultDangerouslySetInnerHTML = "DangerouslySetInnerHTMLof" + startTime
	defaultKeyName = "KeyOf" + startTime
	dirName = "_resultsOf" + startTime
	defaultURLsfound = "foundURLcontainingJS" + startTime
	os.Mkdir(dirName, 0777)

	wantHelp := flag.Bool("h", false, "display help page")
	maxTimeoutOption := flag.String("t", "20", "max timeout for connection timeouts")
	jsFileName := flag.String("u", "", "js input")
	jsURLtoFind := flag.String("j", "", "js url to find")
	flag.Parse()

	maxTimeout, _ = strconv.Atoi(*maxTimeoutOption)

	if *wantHelp {
		help()
	}

	if len(*jsURLtoFind) > 0 {
		URL = *jsURLtoFind
	}

	if *jsFileName == "" {
		fmt.Println(red + "[ERROR]" + reset + " " + white + ": No URL file selected")
		os.Exit(3)
	}

	allOfURLfile, err := ioutil.ReadFile(*jsFileName)
	if err != nil {
		fmt.Println(err)
	}
	urlArray = strings.Split(string(allOfURLfile), "\n")
	urlArray = urlArray[:len(urlArray)-1]

	for i := 0; i < len(urlArray); i++ {

		processURL(urlArray[i])
		if err != nil {
			fmt.Println(err)
		}

		if i%50 == 0 {
			fmt.Println(green + "[Progress]" + reset + " " + white + strconv.Itoa(i) + "/" + strconv.Itoa(len(urlArray)))
		}

		a := append(innerHTMlfound, newInnerHTMLfound...)
		innerHTMlfound = a
		b := append(htmlFound, newHTMLFound...)
		htmlFound = b
		c := append(evalFound, newEvalFound...)
		evalFound = c
		d := append(dangerousFound, newDangeorusFound...)
		dangerousFound = d
		e := append(s3BucketsFound, newS3BucketsFound...)
		s3BucketsFound = e
		f := append(keyFound, newKeyFound...)
		keyFound = f

		writeToFile(defaultInnerHTMLName, newInnerHTMLfound)
		writeToFile(defaultHTMLName, newHTMLFound)
		writeToFile(defaultEvalName, evalFound)
		writeToFile(defaultDangerouslySetInnerHTML, newDangeorusFound)
		writeToFile(defaultS3Name, newS3BucketsFound)
		writeToFile(defaultKeyName, newKeyFound)

		newInnerHTMLfound = nil
		newHTMLFound = nil
		newEvalFound = nil
		newDangeorusFound = nil
		newS3BucketsFound = nil
		newKeyFound = nil

	}

	elapsed := time.Since(startTimetime).String()
	fmt.Println(cyan + "[Info]" + reset + " " + white + "Scan took " + elapsed + " seconds")
}

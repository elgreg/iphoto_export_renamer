package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
  "github.com/gosexy/exif"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printType(a interface{}) {
	fmt.Println(reflect.TypeOf(a))
}

func trimLeadingDashSpace(s string) string {
	return strings.TrimLeft(s, " -_")
}

func addLeadingDashSpace(s string) string {
	return " - " + s
}

func convYYYYMMDD(name string) []string {
	YYYYMMDD := regexp.MustCompile(`^(\d{8})(\D.*)`)
	// No great reason for MustCompile, just going with docs
	match := YYYYMMDD.FindStringSubmatch(name)
	result := make([]string, 2)
	if len(match) > 0 {
		result = []string{match[1], match[2]}
	}
	return result
}

func convDayMonYr(name string) []string {
	dMonYYYY := regexp.MustCompile(`^(\d+)\s+([A-Z][a-z][a-z])\s+(\d{4})(.*)`)
	monMap := map[string]int{
		"Jan": 1,
		"Feb": 2,
		"Mar": 3,
		"Apr": 4,
		"May": 5,
		"Jun": 6,
		"Jul": 7,
		"Aug": 8,
		"Sep": 9,
		"Oct": 10,
		"Nov": 11,
		"Dec": 12,
	}
	match := dMonYYYY.FindStringSubmatch(name)
	result := make([]string, 2)
	if len(match) > 0 {
		year := match[3]
		monStr := match[2]
		monInt := monMap[monStr]
		monStr = zeroPad(monInt, 2)
		dayStr := match[1]
		dayInt, err := strconv.Atoi(dayStr)
		check(err)
		dayStr = zeroPad(dayInt, 2)
		date := year + monStr + dayStr
		result = []string{date, match[4]}
	}
	return result
}

func convYYYYUnderscore(name string) []string {
	YYYY_MM_DD := regexp.MustCompile(`^(\d{4})_(\d{2})_(\d{2})(.*)`)
	match := YYYY_MM_DD.FindStringSubmatch(name)
	result := make([]string, 2)
	if len(match) > 0 {
		result = []string{match[1] + match[2] + match[3], match[4]}
	}
	return result
}

func isFileJpg

func goIntoDirectory(file os.FileInfo) string {
  dirList, err := ioutil.ReadDir(file.Name())
  for _, file := range dirList {
    isJpg := isFileJpg(file)
    if isJpg {
      jpg
    }
  }
}

func getDateBasedFilename(file os.FileInfo) string {
	fmt.Println("Processing " + file.Name())
	dateFilename := convYYYYMMDD(file.Name())
	if len(dateFilename[0]) == 0 {
		dateFilename = convDayMonYr(file.Name())
		if len(dateFilename[0]) == 0 {
			dateFilename = convYYYYUnderscore(file.Name())
			if len(dateFilename[0]) == 0 {
        dateFilename = goIntoDirectory(file)
				fmt.Println("No conversion for: " + file.Name())
			}
		}
	}
	fileDate := dateFilename[0]
	fileName := dateFilename[1]
	if len(fileName) > 0 {
		fileName = trimLeadingDashSpace(fileName)
		fileName = addLeadingDashSpace(fileName)
	}
	newFileName := fileDate + fileName
	return newFileName
}

// Check map to see if key exists and increment if it already does
func getUniqueMapKey(key string, mapToProcess map[string]string) string {
	_, prs := mapToProcess[key]
	if prs {
		newKey := incrementKey(key)
		return getUniqueMapKey(newKey, mapToProcess)
	} else {
		return key
	}
}

// Look for _DD at the end. If it's already there,
// increment it by 1. Use 0 padding
func incrementKey(key string) string {
	split := strings.Split(key, "_")
	last := split[len(split)-1]
	i, err := strconv.Atoi(last)
	if err != nil {
		return key + "_00"
	} else {
		last = zeroPad(i+1, 2)
		return key + "_" + last
	}
}

func zeroPad(i int, n int) string {
	iAsString := strconv.Itoa(i)
	return strings.Repeat("0", n-len(iAsString)) + iAsString
}

func main() {
	app := cli.NewApp()
	app.Name = "iphoto_export_renamer"
	app.Usage = "Rename a folder from iPhoto events export by date"
	app.Action = func(c *cli.Context) {
		directory := "."
		if len(c.Args()) > 0 {
			directory = c.Args()[0]
		}

		fmt.Println("Processing Directory " + directory + "...\n")

		dirList, err := ioutil.ReadDir(directory)
		check(err)
		var newDirList = make(map[string]string)

		// Some vars we'll use in our loop
		var fileName string
		for _, file := range dirList {

			fileName = getDateBasedFilename(file)
			fmt.Println(file.Name() + " ---> " + fileName)
			fileName = getUniqueMapKey(fileName, newDirList)
			newDirList[fileName] = file.Name()

			// // Go into directory, list files and repeat this loop to name

			// // If no files are datelike, then try their exif data

		}

		fmt.Println(newDirList)
	}

	app.Run(os.Args)
}

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printType(a interface{}) {
	fmt.Println(reflect.TypeOf(a))
}

func getDateBasedFilename(file os.FileInfo) string {
	// Check if filename.Name() already starts with YYYMMDD and bail if so
	/* isDate := blah(dir.Name())
	   if someCheck {
	     go equiv of next()
	  }*/

	// Check if dir.Name() converts easily to YYYMMDD in format
	// 1 Apr YYYY or YYYY_MM_DD

	/* getDateFromName := blah(dir.Name())
	   if getDateFromName which returns a string then pass that to move
	*/
	fmt.Println("Processing " + file.Name())

	return file.Name()
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
			fileName = getUniqueMapKey(fileName, newDirList)
			newDirList[fileName] = file.Name()

			// // Go into directory, list files and repeat this loop to name

			// // If no files are datelike, then try their exif data

		}

		fmt.Println(newDirList)
	}

	app.Run(os.Args)
}

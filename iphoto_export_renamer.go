package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gosexy/exif"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkNoPanic(e error) {
	if e != nil {
		fmt.Println("Yikes! " + e.Error())
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

func convYYYYSep(name string, sep string) []string {
	YYYY_MM_DD := regexp.MustCompile(`^(\d{4})` + sep + `(\d{2})` + sep + `(\d{2})(.*)`)
	match := YYYY_MM_DD.FindStringSubmatch(name)
	result := make([]string, 2)
	if len(match) > 0 {
		result = []string{match[1] + match[2] + match[3], match[4]}
	}
	return result
}

func isFileJpg(file os.FileInfo) bool {
	var isJpg bool = false
	fileType := mime.TypeByExtension(path.Ext(file.Name()))
	fmt.Println(" Type for " + file.Name() + " is " + fileType)
	if fileType == "image/jpeg" {
		isJpg = true
	}
	return isJpg
}

func getImgExifDate(filePath string) string {
	dateExif := "YYYYMMDD"
	exifReader := exif.New()
	imgData, errOpen := os.Open(filePath)
	checkNoPanic(errOpen)
	if errOpen != nil {
		return dateExif
	}

	defer imgData.Close()

	_, errCopy := io.Copy(exifReader, imgData)
	if errCopy != nil && errCopy != exif.FoundExifInData {
		fmt.Println("YIKES copying exif: " + errCopy.Error())
		return dateExif
	}

	parseErr := exifReader.Parse()
	checkNoPanic(parseErr)
	if parseErr != nil {
		return dateExif
	}

	// print out key,val pair for each tag
	for key, val := range exifReader.Tags {
		if key == "Date and Time" {
			dateExifSlice := convYYYYSep(val, ":")
			dateExif = dateExifSlice[0]
			break
		}
	}

	return dateExif
}

func getDateFromFilesInDir(dir os.FileInfo, basePath string) []string {
	dirList, err := ioutil.ReadDir(dir.Name())
	check(err)
	result := make([]string, 2)

	dateExif := "YYYYMMDD"
	var isJpg bool = false
	for _, imgFile := range dirList {
		isJpg = isFileJpg(imgFile)
		filePath := basePath + "/" + imgFile.Name()
		if isJpg {
			dateExif = getImgExifDate(filePath)

			if dateExif != "YYYYMMDD" {
				break
			}
		}
	}
	result = []string{dateExif, dir.Name()}
	return result
}

func getDateBasedFilename(dir os.FileInfo, absPath string) string {
	dirName := dir.Name()
	dateFilename := convYYYYMMDD(dirName)
	if len(dateFilename[0]) == 0 {
		dateFilename = convDayMonYr(dirName)
		if len(dateFilename[0]) == 0 {
			dateFilename = convYYYYSep(dirName, "_")
			if len(dateFilename[0]) == 0 {
				dateFilename = getDateFromFilesInDir(dir, absPath+"/"+dirName)
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "outfile, o",
			Value: "output.sh",
			Usage: "Name of output file",
		},
	}

	app.Action = func(c *cli.Context) {
		directory := "."
		if len(c.Args()) > 0 {
			directory = c.Args()[0]
		}
		absPath, pathErr := filepath.Abs(directory)
		check(pathErr)

		outfilePath := absPath + "/" + c.String("outfile")

		fmt.Println("Processing Directory " + absPath + "...\n")

		dirList, err := ioutil.ReadDir(absPath)
		check(err)
		var newDirList = make(map[string]string)

		// Some vars we'll use in our loop
		var fileName string
		for _, dir := range dirList {
			if dir.IsDir() {
				fileName = getDateBasedFilename(dir, absPath)
				fmt.Println(dir.Name() + " ---> " + fileName)
				fileName = getUniqueMapKey(fileName, newDirList)
				newDirList[fileName] = dir.Name()
			}
		}

		// Write to a file disignated by the -o flag
		f, err := os.Create(outfilePath)
		check(err)
		defer f.Close()
		for key, value := range newDirList {
			f.WriteString("mv \"" + value + "\" \"" + key + "\"\n")
		}
		f.Sync()

	}

	app.Run(os.Args)
}

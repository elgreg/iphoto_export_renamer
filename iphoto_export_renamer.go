package main

import (
  "os"
  "io/ioutil"
  "github.com/codegangsta/cli"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
  app := cli.NewApp()
  app.Name = "iphoto_export_renamer"
  app.Usage = "Rename a folder from iPhoto events export by date"
  app.Action = func(c *cli.Context){
    directory := "."
    if len(c.Args()) > 0 {
      directory = c.Args()[0]
    }

    println(directory)
    dirList, err := ioutil.ReadDir(directory)
    check(err)
    for _, dir := range dirList{
      // Check if dir.Name() already starts with YYYMMDD and bail if so
      /* isDate := blah(dir.Name())
       if someCheck {
         go equiv of next()
      }*/

      // Check if dir.Name() converts easily to YYYMMDD in format
      // 1 Apr YYYY or YYYY_MM_DD

      /* getDateFromName := blah(dir.Name())
      if getDateFromName which returns a string then pass that to move
      */

      // Go into directory, list files and repeat this loop to name

      // If no files are datelike, then try their exif data

      println(dir.Name())
    }

  }

  app.Run(os.Args)
}


func printHi() {
  println("hi")
}
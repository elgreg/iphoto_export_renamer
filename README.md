iPhoto Export Renamer
===========================

## In short
I wanted to get off of iPhoto and switch back to just using folders and files. Easier said than done with the format that I had existing pictures in and the way that iPhoto exports them. When I exported the entire library, I ended up with lots of different types of folder names:

* YYYY_MM_DD Name of Event
* Apr 3 2016 Name of Event
* etc.

So the purpose of this script is:

1. For me to learn a bit of go
2. To get all the folders named appropriately.

Your mileage may vary given that I don't know why some of my folders
were named one way and some another. Oh well.

VIDEOS didn't work at all.

## Long

Ideal output: folder name with format YYYYMMDD - Event Name

Currently there are some folder with names like this.
Others contained scanned images that have names starting with YYYYMMDD
The folders are at least organized by event - which means generally pictures
in one folder should stay in that folder.

1. Check if the folder already has a name like YYYYMMDD.* If so, STOP.
2. Check if the folder name already looks like a date - for example 1 Apr 1973. If so, set YYYYMMDD and GOTO moveFolder function with old name and new name:  D M YYYY,   YYYY_MM_DD_(Title)
3. Check if any files in the folder have names like YYYYMMDD.* If so, call moveFolder with that as the date.
4. Run exif on the first file in the folder and extract it. Use this as the date for the folder


Rather than actually write a function that does the file moving, this just spits out a bash script. What better!

Usage
======

```bash
go install github.com/elgreg/iphoto_export_renamer
# cd into folder you
iphoto_export_renamer -o output-filename.sh
./output-filename.sh
```

This will create a directory called "renamed_folders" with all the new folders



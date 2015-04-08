# iPhoto Export Renamer
# ===========================

# Ideal output: folder name with format YYYYMMDD-Event-Name


# Currently there are some folder with names like this.
# Others contained scanned images that have names starting with YYYYMMDD
# The folders are at least organized by event - which means generally pictures
# in one folder should stay in that folder.

# 1. Check if the folder already has a name like YYYYMMDD.* If so, STOP.
# 2. Check if the folder name already looks like a date - for example 1 Apr 1973. If so, set YYYYMMDD and GOTO moveFolder function with old name and new name
#   D M YYYY
#   YYYY_MM_DD_(Title)
# 3. Check if any files in the folder have names like YYYYMMDD.* If so, call moveFolder with that as the date.
# 4. Run exif on the first file in the folder and extract it. Use this as the date for the folder


# moveFolder function needs to take in old name and new name, check to make sure new name doesn't
# already exist. If it does, re-call yourself with old name and new name\ +1

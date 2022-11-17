package tools

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func CheckFilterPath(filterFilePath string) bool {
	/*
		Will return true if filter path can get written to.
		If the .cache dir does not exist it will try and create the directory.
		If the input is a custom path, this function will only validate that path it is writeable.
	*/
	dirPath := filepath.Dir(filterFilePath)
	if _, err := os.Stat(filterFilePath); err == nil { // File exists
		if IsWriteable(filterFilePath) {
			return true // Can write to file
		} else {
			return false // Can not write to file
		}
	} else { // File does not exist
		if _, err := os.Stat(dirPath); err == nil { // Directory exists
			if IsWriteable(dirPath) {
				return true // Directory path is writeable
			} else {
				return false // Directory is not writeable
			}
		} else { // Path does not exist, if it ends with .cache then we assume it's the default setting and that dir should get created.
			if strings.HasSuffix(dirPath, ".cache") {
				err := os.Mkdir(dirPath, 0755)
				if err != nil {
					return false // Can not create dir
				}
				return true // Could create dir
			}
		}

	}
	return false // Catch all false
}

func IsWriteable(filePath string) bool {
	// This function will check wether file is owned by user and user has write permissions.
	// It will work on directories as well.
	// Todo: check group as well.
	uid := os.Geteuid()          // get user id
	groups, _ := os.Getgroups()  // get groups  of user
	info, _ := os.Stat(filePath) // get metadata from file
	mode := info.Mode()          // returns -rw-r--r-- or similar

	fileStat, _ := info.Sys().(*syscall.Stat_t)

	// This works, but is it best practice ??
	// Check if user owns file
	if uid == int(fileStat.Uid) {
		for i := 1; i < 4; i++ { // Will search for "w" in user range
			perm := mode.String()[i]
			if perm == 119 { // if w, then writeable
				return true
			}
		}
	}
	for _, g := range groups { // Loop over all groups user is in .
		if g == int(fileStat.Gid) { // If group is the same as file
			for i := 4; i < 7; i++ { // Will search for "w" in group range
				perm := mode.String()[i]
				if perm == 119 { // if w, then writeable
					return true
				}
			}
		}
	}
	// The final any catch all
	for i := 7; i < 10; i++ { // Will search for "w" in any range
		perm := mode.String()[i]
		if perm == 119 { // if w, then writeable
			return true
		}
	}
	return false // Not writeable
}

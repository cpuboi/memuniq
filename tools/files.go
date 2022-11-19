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
		if mode.String()[2] == 119 { // -r[2]-r--r-- if 2 = w, then writeable
			return true
		} else {
			return false // Linux is interesting, if permissions are 577 and you are the owner of the file, you still cant write.
		}
	}

	for _, g := range groups { // Loop over all groups user is in .
		if g == int(fileStat.Gid) { // If group is the same as file
			if mode.String()[5] == 119 { // -rw-r[5]-r-- if 2 = w, then writeable
				return true
			} else {
				return false // You are in group but not allowed to write
			}
		}
	}
	// The final any catch all
	if mode.String()[8] == 119 { // -r--r--r[8]- if 2 = w, then writeable
		return true
	} else {
		return false // You are in group but not allowed to write
	}
}

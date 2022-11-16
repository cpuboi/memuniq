package tools

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"syscall"
)

// TODO: Create directory with correct permissions
// TODO: Check if last directory == memuniq, if not, create memuniq..

func checkPermissions(filePath string, createDir bool) bool {
	// Check if file exists and is writeable
	if _, err := os.Stat(filePath); err == nil {
		if isWriteable(filePath) {
			return true
		}
	}
	dirPath := string(path.Dir(filePath))
	// Check if dir exists and is writeable so that a file can get created.
	if _, err := os.Stat(dirPath); err == nil {
		if isWriteable(dirPath) {
			return true
		}
	}
	return false
}

func checkCacheDir(filterFilePath string) string {
	/* Checks if the default cache directory exists
	Will try alternatives
	If user uses a custom bloomfilter path then this should not get used.
	First it checks ~./.cache/bloomfilter.bin
	If that file is writeable, return filename.
	*/
	if checkPermissions(filterFilePath, false) { // First check file exists,
		fmt.Println("Can write to path", filterFilePath)
		return filterFilePath
	} else {
		fmt.Println("Can NOT write to path", filterFilePath)
		return filterFilePath
	}
	var cacheDirs [4]string
	//cacheDirs[0] = path.Dir(filterFilePath)
	cacheDirs[1] = "/var/cache/memuniq"
	cacheDirs[2] = "/var/tmp/memuniq"
	cacheDirs[3] = "/tmp/memuniq"

	//filename := path.Base(filterFilePath)

	// Create a user unique filter name
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	var filename string = string(user.Username) + "_" + path.Base(filterFilePath)

	// Loop over all directories, if it is possible to write the filter file there, return path.
	for _, element := range cacheDirs {
		checkCacheDir(element)
		fmt.Println(path.Join(element, filename))
		// Try and create mem
	}
	return "Done"
}

func isWriteable(filePath string) bool {
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

//func main() {
//	//checkCacheDir("/home/user/.cache/bloomfilter.bin")
//	x := checkCacheDir("/tmp/memuniq/testfile.bin")
//	fmt.Println(x)
//}

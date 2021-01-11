package files

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"nenvoy.com/pkg/utils/osystem"
)

//Copy - copies a file to another location
func Copy(srcPath string, dstPath string) (err error) {

	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "could not open source file:")
	}

	//Defer closing the file until the func returns
	defer srcFile.Close()

	// Creates the new file if it doesnt exist
	destFile, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrap(err, "could not create destination file:")
	}

	//Defer closing the file until the func returns
	defer destFile.Close()

	//Copy the file
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "failed while copying:")
	}

	// Sync and save the file
	err = destFile.Sync()
	if err != nil {
		return errors.Wrap(err, "could not sync destination file:")
	}

	return nil
}

//CreateDirectories - Creates any directories needed for the setup
func CreateDirectories(dirs []string) (err error) {
	// Loop throigh the dirs and create them
	for _, dir := range dirs {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

//RemoveDirectories - Removes a list of directories
func RemoveDirectories(dirs []string) (err error) {
	// Loop throigh the dirs and remove them
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

//RemoveFiles - removes a list of files
func RemoveFiles(files []string) (err error) {
	for _, file := range files {
		// Check if the file exists
		exists, err := osystem.PathExists(file)
		if exists {
			err := os.Remove(file)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

//Distribute - Distributes the deployment files to their destinations
func Distribute(files map[string][]string, permissions os.FileMode) (err error) {
	//Loop through the map and list and copy the files
	for file, destinations := range files {
		for _, destination := range destinations {
			//Copy the file
			err := Copy(file, destination)
			if err != nil {
				return err
			}
			//Set the permissions
			err = os.Chmod(destination, permissions)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

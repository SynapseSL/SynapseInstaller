package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
)

func ShouldIPanic(err error, message string) {
	if err != nil {
		logger.Critical(message)
	}
}

func IsErrorHarmfulSteamcmd(err error) bool {
	if exitError, ok := err.(*exec.ExitError); ok {
		return !(exitError.ExitCode() == 0 || exitError.ExitCode() == 7)
	}
	return true
}

func RecursiveCopyAndDelete(folder string, target string) {
	logger.Info("Recursively copying " + folder + " ...")
	items, err := ioutil.ReadDir(folder + "/")
	ShouldIPanic(err, "Failed recursively copying.")
	for i, item := range items {
		if item.IsDir() {
			logger.Debug(fmt.Sprintf("%d - Found directory %s", i, item.Name()))
			CreateFolderIfNotExist(target + "/" + item.Name())
			logger.Debug(fmt.Sprintf("%d - Created directory %s", i, target+item.Name()))
			RecursiveCopyAndDelete(folder+"/"+item.Name(), target+"/"+item.Name())
			err = os.Remove(folder + "/" + item.Name())
			ShouldIPanic(err, fmt.Sprintf("Failed removing directory: %s", err))
		} else {
			logger.Debug(fmt.Sprintf("%d - Found file %s", i, item.Name()))
			err = os.Rename(folder+"/"+item.Name(), target+"/"+item.Name())
			ShouldIPanic(err, fmt.Sprintf("Failed to move file: %s", err))
			logger.Debug(fmt.Sprintf("%d - Moved file %s to %s", i, folder+"/"+item.Name(), target+"/"+item.Name()))
		}
	}
}

func CreateFolderIfNotExist(path string) {
	logger.Debug(fmt.Sprintf("Looking for %s", path))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		logger.Debug(fmt.Sprintf("%s was not found. Creating...", path))
		err := os.Mkdir(path, os.ModePerm)
		ShouldIPanic(err, "Failed to create Synapse folders")
	}
}

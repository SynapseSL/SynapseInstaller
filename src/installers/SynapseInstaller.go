package installers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/utils"
)

func Download_Synapse() string {
	logger.Info("Downloading Synapse...")
	resp, err := http.Get("https://cdn.culabs.eu/synapseinstaller/Synapse.zip")
	utils.ShouldIPanic(err, "Failed to download Synapse.zip from provider cdn.culabs.eu")

	file, err := os.Create("Synapse.zip")
	utils.ShouldIPanic(err, "Failed to create file for Synapse.zip")

	_, err = io.Copy(file, resp.Body)
	utils.ShouldIPanic(err, "Failed to copy stream")
	logger.Debug("Saved as " + file.Name())

	resp.Body.Close()
	utils.ShouldIPanic(file.Close(), "Failed to close Synapse.zip properly.")
	return file.Name()
}

func Install_Synapse_To(binaries string, files string, path string, unzip_cmd string, unzip_args string) {
	unzipSynapse(path, unzip_cmd, unzip_args)

	// Copy Assembly-CSharp.dll to where it should be.

	var assemblycs string = binaries + "SCPSL_Data/Managed/Assembly-CSharp.dll"
	logger.Info("Moving assemblies...")
	err := os.Rename(assemblycs, assemblycs+".bak")
	utils.ShouldIPanic(err, "Failed to rename file Assembly-CSharp.dll - your installation is likely corrupt")
	err = os.Rename("Assembly-CSharp.dll", assemblycs)
	utils.ShouldIPanic(err, "Failed to move Assembly-CSharp.dll to game directory")
	logger.Info("Installed SynapseLoader.")

	// Copy Synapse files.
	utils.CreateFolderIfNotExist(files + "/Synapse")
	utils.ShouldIPanic(err, "Could not create Synapse directory")

	// Create Synapse folders if they do not exist.
	utils.RecursiveCopyAndDelete("Synapse", files+"/Synapse")

	logger.Info("Synapse is now installed.")
}

func unzipSynapse(path string, unzip_cmd string, unzip_args string) {
	var cmd string
	var args string = ""
	if unzip_cmd == "" {
		logger.Debug("Falling back to default zip command")
		if runtime.GOOS == "windows" {
			logger.Info("Detected OS: Windows. Using 7za for unzip.")
			cmd = "7za"
			args = "x -y"
		} else if runtime.GOOS == "linux" {
			logger.Info("Detected OS: Linux. Using unzip for unzip.")
			cmd = "unzip"
			args = "-o"
		} else {
			logger.Warn("Your OS seems to be " + runtime.GOOS + ", but this is not natively supported by SynapseInstaller. Falling back to 7z, if that doesn't work please specify a custom unzip command.")
			logger.Info("Detected OS: " + runtime.GOOS + ". Using 7za for unzip.")
			cmd = "7za"
			args = "x -y"
		}
	} else {
		logger.Debug("Custom Unzip: " + unzip_cmd)
		cmd = unzip_cmd

		if unzip_args != "" {
			logger.Debug("Custom args: " + unzip_args)
			args = unzip_args
		}
	}

	pwd, err := os.Getwd()
	utils.ShouldIPanic(err, "Failed to determine working directory - this should not happen!")

	if args == "" {
		args = path
	} else {
		args += " " + path
	}
	largs := strings.Split(args, " ")
	logger.Debug(fmt.Sprintf("Running %s %s ...", cmd, args))
	e_unzip_cmd := exec.Command(cmd, largs...)
	e_unzip_out, err := e_unzip_cmd.Output()
	if err != nil && runtime.GOOS == "windows" && unzip_cmd == "" {
		logger.Warn("Failed to find 7za in your PATH. Falling back to bundled 7za.")
		cmd = pwd + "/bundled/7za.exe"
		logger.Debug(fmt.Sprintf("Running %s %s ...", cmd, args))
		e_unzip_cmd := exec.Command(cmd, largs...)
		e_unzip_out, err = e_unzip_cmd.Output()
		utils.ShouldIPanic(err, fmt.Sprintf("Failed fallback: %s\n%s", err, e_unzip_out))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Failed to unzip: %s\n%s", err, e_unzip_out))
	}

	logger.Output(string(e_unzip_out))
}

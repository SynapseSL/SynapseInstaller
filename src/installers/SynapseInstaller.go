package installers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SynapseSL/SynapseInstaller/src/logger"
	"github.com/SynapseSL/SynapseInstaller/src/utils"
)

func DownloadSynapse() string {
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

func InstallSynapseTo(binaries string, files string, path string, unzipCmd string, unzipArgs string) {
	unzipSynapse(path, unzipCmd, unzipArgs)

	// Copy Assembly-CSharp.dll to where it should be.

	assemblycs := binaries + "SCPSL_Data/Managed/Assembly-CSharp.dll"

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

func unzipSynapse(path string, unzipCmd string, unzipArgs string) {
	var (
		cmd  string
		args string
	)

	if unzipCmd == "" {
		logger.Debug("Falling back to default zip command")

		switch runtime.GOOS {
		case win:
			logger.Info("Detected OS: Windows. Using 7za for unzip.")

			cmd = "7za"
			args = "x -y"
		case lnx:
			logger.Info("Detected OS: Linux. Using unzip for unzip.")

			cmd = "unzip"
			args = "-o"
		default:
			logger.Warn("Your OS seems to be " + runtime.GOOS + ", but this is not natively supported by SynapseInstaller. Falling back to 7z, if that doesn't work please specify a custom unzip command.")
			logger.Info("Detected OS: " + runtime.GOOS + ". Using 7za for unzip.")

			cmd = "7za"
			args = "x -y"
		}
	} else {
		logger.Debug("Custom Unzip: " + unzipCmd)
		cmd = unzipCmd

		if unzipArgs != "" {
			logger.Debug("Custom args: " + unzipArgs)
			args = unzipArgs
		}
	}

	pwd, err := os.Getwd()
	utils.ShouldIPanic(err, "Failed to determine working directory - this should not happen!")

	if args == "" {
		args = path
	} else {
		args += " " + path
	}

	logger.Debug(fmt.Sprintf("Running %s %s ...", cmd, args))

	largs := strings.Split(args, " ")
	eUnzipCmd := exec.Command(cmd, largs...)
	eUnzipOut, err := eUnzipCmd.Output()

	if err != nil && runtime.GOOS == "windows" && unzipCmd == "" {
		logger.Warn("Failed to find 7za in your PATH. Falling back to bundled 7za.")

		cmd = pwd + "/bundled/7za.exe"
		logger.Debug(fmt.Sprintf("Running %s %s ...", cmd, args))

		eUnzipCmd := exec.Command(cmd, largs...)
		eUnzipOut, err = eUnzipCmd.Output()
		utils.ShouldIPanic(err, fmt.Sprintf("Failed fallback: %s\n%s", err, eUnzipOut))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Failed to unzip: %s\n%s", err, eUnzipOut))
	}

	logger.Output(string(eUnzipOut))
}

package main

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/installers"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/utils"
)

//nolint:gochecknoglobals,lll
var (
	flagScripted             = flag.Bool("scripted", false, "If enabled, no interaction will be required to install.")
	flagInstallSynapse       = flag.Bool("install-synapse", false, "Install/Update Synapse?")
	flagInstallGame          = flag.Bool("install-game", false, "Install/Update SCP: Secret Laboratory?")
	flagNoAnsi               = flag.Bool("noansi", false, "Should ansi output be disabled? Recommended for Cmd/PowerShell and other shells without ansi support.")
	flagVerbosity            = flag.Int("verbosity", logger.LogLevelWARN, "How verbose should output be? Lower number = more verbose.")
	flagGameBinariesLocation = flag.String("binaries", "./scpsl_dedicatedserver/", "Where are the game binaries located?")
	flagGameFilesLocation    = flag.String("files", "~/.config/", "Where are your config files located?")
	flagSynapseLocation      = flag.String("synapsezip", "", "If you have already downloaded the Synapse.zip, where is it?")
	flagCustomUnzipCmd       = flag.String("unzip-cmd", "", "Want to use a custom unzip command? Input its command here.")
	flagCustomUnzipArgs      = flag.String("unzip-args", "", "Custom Unzip Args")
)

func main() {
	// Logger setup - a one time step we need to do.
	flag.Parse()
	logger.UsedLogger.SetLogLevel(*flagVerbosity)
	logger.UsedLogger.SetAnsi(!(*flagNoAnsi))

	pwd, err := os.Getwd()
	utils.ShouldIPanic(err, "Could not figure out current working directory - this shouldn't happen! (GOOS_GETWD_ERR)")

	updateFlagPaths(pwd)

	if *flagScripted {
		scripted()
	} else {
		interactive()
	}
}

func scripted() {
	logger.Info("Running in scripted mode.")

	if *flagInstallGame {
		installers.InstallGame(*flagGameBinariesLocation)
	} else {
		logger.Info("Skipped - Game was not installed.")
	}

	if *flagInstallSynapse {
		installSynapse()
	} else {
		logger.Info("Skipped - Synapse was not installed.")
	}
}

func installSynapse() {
	logger.Info("Attempting to install Synapse...")

	var path string
	if *flagSynapseLocation == "" {
		path = installers.DownloadSynapse()
	} else {
		path = *flagSynapseLocation
	}

	logger.Debug(path)

	installers.InstallSynapseTo(
		*flagGameBinariesLocation,
		*flagGameFilesLocation,
		path,
		*flagCustomUnzipCmd,
		*flagCustomUnzipArgs)
	logger.Ok("Installed Synapse.")
}

func interactive() {
	panic("Interactive not implemented o.o\nPlease use scripted mode for now!")
}

func updateFlagPaths(pwd string) {
	if strings.HasPrefix(*flagGameBinariesLocation, "./") {
		*flagGameBinariesLocation = strings.Replace(*flagGameBinariesLocation, "./", pwd+"/", 1)
	}

	if strings.HasPrefix(*flagGameFilesLocation, "./") {
		*flagGameFilesLocation = strings.Replace(*flagGameFilesLocation, "./", pwd+"/", 1)
	}

	if strings.HasPrefix(*flagGameFilesLocation, "~/") {
		switch runtime.GOOS {
		case "windows":
			logger.Info("Detected OS: Windows. Adjusting path.")

			if *flagGameFilesLocation == "~/.config/" {
				logger.Info("Detected Linux Default Install Directory! Fixing path for Windows.")

				*flagGameFilesLocation = os.Getenv("appdata") + "/"
			} else {
				logger.Warn(
					"Detected UNIX Home directive, but OS is Windows. " +
						"We're fixing this up, but this might cause problems later on.")

				*flagGameFilesLocation = strings.Replace(*flagGameFilesLocation, "~/", os.Getenv("UserProfile")+"/", 1)
			}
		case "linux":
			logger.Info("Detected OS: Linux. Adjusting path.")

			*flagGameFilesLocation = strings.Replace(*flagGameFilesLocation, "~/", os.Getenv("HOME")+"/", 1)
		default:
			logger.Warn(
				"Detected OS to be " + runtime.GOOS +
					", but this is not natively supported by SynapseInstaller (Expected: Windows, Linux). " +
					"Issues may arise.")
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/installers"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/utils"
)

var (
	flag_scripted               = flag.Bool("scripted", false, "If enabled, no interaction will be required to install.")
	flag_game_binaries_location = flag.String("binaries", "./scpsl_dedicatedserver/", "Where are the game binaries located?")
	flag_game_files_location    = flag.String("files", "~/.config/", "Where are your config files located?")
	flag_install_game           = flag.Bool("install-game", false, "Install/Update SCP: Secret Laboratory?")
	flag_install_synapse        = flag.Bool("install-synapse", false, "Install/Update Synapse?")
	flag_verbosity              = flag.Int("verbosity", 2, "How verbose should output be? Lower number = more verbose.")
	flag_synapse_location       = flag.String("synapsezip", "", "If you have already downloaded the Synapse.zip, where is it?")
	flag_custom_unzip_cmd       = flag.String("unzip-cmd", "", "Want to use a custom unzip command? Input its command here.")
	flag_custom_unzip_args      = flag.String("unzip-args", "", "Custom Unzip Args")
	flag_noansi                 = flag.Bool("noansi", false, "Should ansi output be disabled? Recommended for Cmd/PowerShell and other shells without ansi suport.")
)

func main() {
	// Logger setup - a one time step we need to do.
	flag.Parse()
	logger.UsedLogger.SetLogLevel(*flag_verbosity)
	logger.UsedLogger.SetAnsi(!(*flag_noansi))

	pwd, err := os.Getwd()
	utils.ShouldIPanic(err, "Could not figure out current working directory - this shouldn't happen! (GOOS_GETWD_ERR)")

	updateFlagPaths(pwd)

	if *flag_scripted {
		scripted()
	} else {
		interactive()
	}
}

func scripted() {
	logger.Info("Running in scripted mode.")
	if *flag_install_game {
		installers.Install_Game(*flag_game_binaries_location)
	} else {
		logger.Info("Skipped - Game was not installed.")
	}

	if *flag_install_synapse {
		install_synapse()
	} else {
		logger.Info("Skipped - Synapse was not installed.")
	}
}

func install_synapse() {
	logger.Info("Attempting to install Synapse...")

	var path string
	if *flag_synapse_location == "" {
		path = installers.Download_Synapse()
	} else {
		path = *flag_synapse_location
	}

	logger.Debug(path)

	installers.Install_Synapse_To(*flag_game_binaries_location, *flag_game_files_location, path, *flag_custom_unzip_cmd, *flag_custom_unzip_args)
	logger.Ok("Installed Synapse.")
}

func interactive() {
	fmt.Println("SynapseInstaller 1.0.0b")
	fmt.Println("============================")
	fmt.Println("Running in interactive mode.")
	fmt.Println("\nWelcome to the interactive SynapseSL installer.")
	panic("Interactive not implemented o.o\nPlease use scripted mode for now!")
}

func updateFlagPaths(pwd string) {
	if strings.HasPrefix(*flag_game_binaries_location, "./") {
		*flag_game_binaries_location = strings.Replace(*flag_game_binaries_location, "./", pwd+"/", 1)
	}

	if strings.HasPrefix(*flag_game_files_location, "./") {
		*flag_game_files_location = strings.Replace(*flag_game_files_location, "./", pwd+"/", 1)
	}

	if strings.HasPrefix(*flag_game_files_location, "~/") {
		if runtime.GOOS == "windows" {
			logger.Info("Detected OS: Windows. Adjusting path.")
			if *flag_game_files_location == "~/.config/" {
				logger.Info("Detected Linux Default Install Directory! Fixing path for Windows.")
				*flag_game_files_location = os.Getenv("appdata") + "/"
			} else {
				logger.Warn("Detected UNIX Home directive, but OS is Windows. We're fixing this up, but this might cause problems later on.")
				*flag_game_files_location = strings.Replace(*flag_game_files_location, "~/", os.Getenv("UserProfile")+"/", 1)
			}
		} else if runtime.GOOS == "linux" {
			logger.Info("Detected OS: Linux. Adjusting path.")
			*flag_game_files_location = strings.Replace(*flag_game_files_location, "~/", os.Getenv("HOME")+"/", 1)
		} else {
			logger.Warn("Detected OS to be " + runtime.GOOS + ", but this is not natively supported by SynapseInstaller (Expected: Windows, Linux). Issues may arise.")
		}
	}
}

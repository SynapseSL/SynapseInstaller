package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
)

var flag_scripted = flag.Bool("scripted", false, "If enabled, no interaction will be required to install.")
var flag_game_binaries_location = flag.String("binaries", "./SCPSL_DEDICATEDSERVER/", "Where are the game binaries located?")
var flag_game_files_location = flag.String("files", "~/.config/", "Where are your config files located?")
var flag_install_game = flag.Bool("install-game", false, "Install/Update SCP: Secret Laboratory?")
var flag_install_synapse = flag.Bool("install-synapse", false, "Install/Update Synapse?")
var flag_verbosity = flag.Int("verbosity", 2, "How verbose should output be? Lower number = more verbose.")
var flag_synapse_location = flag.String("synapsezip", "", "If you have already downloaded the Synapse.zip, where is it?")
var flag_custom_unzip_cmd = flag.String("unzipcmd", "", "Want to use a custom unzip command? Input its command here.")
var flag_custom_unzip_args = flag.String("unzipargs", "", "Custom Unzip Args")

func main() {
	// Logger setup - a one time step we need to do.
	flag.Parse()
	logger.UsedLogger.SetLogLevel(*flag_verbosity)

	pwd, err := os.Getwd()
	shouldIPanic(err, "Could not figure out current working directory - this shouldn't happen! (GOOS_GETWD_ERR)")

	updateFlagPaths(pwd)

	if *flag_scripted {
		scripted()
	} else {
		interactive()
	}
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
				*flag_game_files_location = os.Getenv("appdata") + "/Synapse/"
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

func scripted() {
	logger.Info("Running in scripted mode.")
	if *flag_install_game {
		install_game()
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
		path = download_synapse()
	} else {
		path = *flag_synapse_location
	}

	install_synapse_to(*flag_game_binaries_location, *flag_game_files_location, path, *flag_custom_unzip_cmd, *flag_custom_unzip_args)
	logger.Ok("Installed Synapse.")
}

func install_game() {
	logger.Info("Attempting to install game...")
	test_steamcmd()
	logger.Debug("SteamCMD seems to work! Hooray!")
	logger.Info("Installing SCPSL_DEDICATED to " + *flag_game_binaries_location + " ...")
	steamcmd_install_to(*flag_game_binaries_location)
	logger.Ok("Installed SCP:SL.")
}

func download_synapse() string {
	logger.Info("Downloading Synapse...")
	resp, err := http.Get("https://cdn.culabs.eu/synapseinstaller/Synapse.zip")
	shouldIPanic(err, "Failed to download Synapse.zip from provider cdn.culabs.eu")
	defer resp.Body.Close()

	file, err := os.Create("Synapse.zip")
	shouldIPanic(err, "Failed to create file for Synapse.zip")

	_, err = io.Copy(file, resp.Body)
	shouldIPanic(err, "Failed to copy stream")
	logger.Debug("Saved as " + file.Name())

	return file.Name()
}

func install_synapse_to(binaries string, files string, path string, unzip_cmd string, unzip_args string) {
	unzipSynapse(unzip_cmd, unzip_args)

	// Copy Assembly-CSharp.dll to where it should be.

	var assemblycs string = binaries + "SCPSL_Data/Managed/Assembly-CSharp.dll"
	logger.Info("Moving assemblies...")
	err := os.Rename(assemblycs, assemblycs+".bak")
	shouldIPanic(err, "Failed to rename file Assembly-CSharp.dll - your installation is likely corrupt")
	err = os.Rename("Assembly-CSharp.dll", assemblycs)
	shouldIPanic(err, "Failed to move Assembly-CSharp.dll to game directory")
	logger.Info("Installed SynapseLoader.")

	// Copy Synapse files.
	createFolderIfNotExist(files + "/Synapse")
	shouldIPanic(err, "Could not create Synapse directory")

	// Create Synapse folders if they do not exist.
	recursiveCopy("Synapse/", files)

	logger.Info("Synapse is now installed.")
}

func unzipSynapse(unzip_cmd string, unzip_args string) error {
	var cmd string = ""
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

	if args == "" {
		args = "Synapse.zip"
	} else {
		args += " Synapse.zip"
	}
	largs := strings.Split(args, " ")
	logger.Debug(fmt.Sprintf("Running %s %s ...", cmd, args))
	e_unzip_cmd := exec.Command(cmd, largs...)
	e_unzip_out, err := e_unzip_cmd.Output()
	shouldIPanic(err, "Failed to unzip - your unzipper was presumably called with invalid arguments.")

	logger.Output(string(e_unzip_out))
	return err
}

func recursiveCopy(folder string, target string) {
	logger.Info("Recursively copying " + folder + " ...")
	items, err := ioutil.ReadDir(folder)
	shouldIPanic(err, "Failed recursively copying.")
	for i, item := range items {
		if item.IsDir() {
			logger.Debug(fmt.Sprintf("%d - Found directory %s", i, item.Name()))
			createFolderIfNotExist(target + "/" + item.Name())
			logger.Debug(fmt.Sprintf("%d - Created directory %s", i, target+item.Name()))
			recursiveCopy(folder+item.Name(), target+item.Name())
		} else {
			logger.Debug(fmt.Sprintf("%d - Found file %s", i, item.Name()))
			os.Rename(folder+item.Name(), target+"/"+item.Name())
			logger.Debug(fmt.Sprintf("%d - Moved file %s to %s", i, folder+item.Name(), target+"/"+item.Name()))
		}
	}
}

func createFolderIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		shouldIPanic(err, "Failed to create Synapse folders")
	}
}

func steamcmd_install_to(s string) {
	logger.Debug("Calling SteamCMD with:")
	logger.Debug("steamcmd +force_install_dir " + strings.Replace(s, " ", "\\ ", -1) + " +login anonymous +app_update 996560 validate +quit")
	steamcmd_test_cmd := exec.Command("steamcmd", "+force_install_dir", strings.Replace(s, " ", "\\ ", -1), "+login", "anonymous", "+app_update", "996560", "validate", "+quit")
	steamcmd_test_out, err := steamcmd_test_cmd.Output()
	shouldIPanic(err, "Something went wrong calling SteamCMD (err: SteamcmdNotZero)!")
	logger.Output(string(steamcmd_test_out))
}

func test_steamcmd() {
	logger.Debug("Testing steamcmd availability...")
	steamcmd_test_cmd := exec.Command("steamcmd", "+quit")
	steamcmd_test_out, err := steamcmd_test_cmd.Output()
	shouldIPanic(err, "Something went wrong calling SteamCMD (err: SteamcmdNotZero)!")
	logger.Output(string(steamcmd_test_out))
}

func shouldIPanic(err error, message string) {
	if err != nil {
		logger.Critical(message)
	}
}

func interactive() {
	fmt.Println("SynapseInstaller 1.0.0b")
	fmt.Println("============================")
	fmt.Println("Running in interactive mode.")
	fmt.Println("\nWelcome to the interactive SynapseSL installer.")
}
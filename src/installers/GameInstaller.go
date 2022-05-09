package installers

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"git.culabs.eu/cubuzz/SynapseInstaller/src/logger"
	"git.culabs.eu/cubuzz/SynapseInstaller/src/utils"
)

func Install_Game(game_binaries_location string) {
	logger.Info("Attempting to install game...")
	test_steamcmd()
	logger.Debug("SteamCMD seems to work! Hooray!")
	logger.Info("Installing SCPSL_DEDICATED to " + game_binaries_location + " ...")
	steamcmd_install_to(game_binaries_location)
	logger.Ok("Installed SCP:SL.")
}

func test_steamcmd() {
	logger.Debug("Testing steamcmd availability...")
	steamcmd_test_cmd := exec.Command("steamcmd", "+quit")
	steamcmd_test_out, err := steamcmd_test_cmd.Output()
	if err != nil && runtime.GOOS == "windows" && utils.IsErrorHarmfulSteamcmd(err) {
		logger.Warn("Failed calling steamcmd, falling back to bundled SteamCMD.")
		steamcmd_test_cmd := exec.Command("./bundled/steamcmd.exe", "+quit")
		steamcmd_test_out, err := steamcmd_test_cmd.Output()
		if utils.IsErrorHarmfulSteamcmd(err) {
			utils.ShouldIPanic(err, fmt.Sprintf("Failed calling bundled SteamCMD: %s\n%s", err, steamcmd_test_out))
		}
		logger.Output(string(steamcmd_test_out))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Something went wrong calling SteamCMD: %s\n%s", err, steamcmd_test_out))
		logger.Output(string(steamcmd_test_out))
	}
}

func steamcmd_install_to(s string) {
	logger.Debug("Calling SteamCMD with:")
	logger.Debug("steamcmd +force_install_dir " + strings.Replace(s, " ", "\\ ", -1) + " +login anonymous +app_update 996560 validate +quit")
	steamcmd_cmd := exec.Command("steamcmd", "+force_install_dir", strings.Replace(s, " ", "\\ ", -1), "+login", "anonymous", "+app_update", "996560", "validate", "+quit")
	steamcmd_out, err := steamcmd_cmd.Output()
	if err != nil && runtime.GOOS == "windows" {
		logger.Warn("Something went wrong calling SteamCMD. Falling back to bundled SteamCMD.")
		steamcmd_cmd := exec.Command("./bundled/steamcmd.exe", "+force_install_dir", strings.Replace(s, " ", "\\ ", -1), "+login", "anonymous", "+app_update", "996560", "validate", "+quit")
		steamcmd_out, err := steamcmd_cmd.Output()
		utils.ShouldIPanic(err, fmt.Sprintf("Failed calling bundled SteamCMD: %s\n%s", err, steamcmd_out))
		logger.Output(string(steamcmd_out))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Something went wrong calling SteamCMD: %s\n%s", err, steamcmd_out))
		logger.Output(string(steamcmd_out))
	}
}

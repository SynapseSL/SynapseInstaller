package installers

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SynapseSL/SynapseInstaller/src/logger"
	"github.com/SynapseSL/SynapseInstaller/src/utils"
)

func InstallGame(gameBinariesLocation string) {
	logger.Info("Attempting to install game...")
	testSteamcmd()
	logger.Debug("SteamCMD seems to work! Hooray!")
	logger.Info("Installing SCPSL_DEDICATED to " + gameBinariesLocation + " ...")
	steamcmdInstallTo(gameBinariesLocation)
	logger.Ok("Installed SCP:SL.")
}

func testSteamcmd() {
	logger.Debug("Testing steamcmd availability...")

	steamcmdTestCmd := exec.Command("steamcmd", "+quit")
	steamcmdTestOut, err := steamcmdTestCmd.Output()

	if err != nil && runtime.GOOS == win && utils.IsErrorHarmfulSteamcmd(err) {
		logger.Warn("Failed calling steamcmd, falling back to bundled SteamCMD.")

		steamcmdTestCmd := exec.Command("./bundled/steamcmd.exe", "+quit")
		steamcmdTestOut, err := steamcmdTestCmd.Output()

		if utils.IsErrorHarmfulSteamcmd(err) {
			utils.ShouldIPanic(err, fmt.Sprintf("Failed calling bundled SteamCMD: %s\n%s", err, steamcmdTestOut))
		}

		logger.Output(string(steamcmdTestOut))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Something went wrong calling SteamCMD: %s\n%s", err, steamcmdTestOut))
		logger.Output(string(steamcmdTestOut))
	}
}

func steamcmdInstallTo(installDir string) {
	installDir = strings.ReplaceAll(installDir, " ", "\\ ")

	logger.Debug("Calling SteamCMD with:")
	logger.Debug("steamcmd +force_install_dir " + installDir + " +login anonymous +app_update 996560 validate +quit")

	steamcmdCmd := exec.Command("steamcmd", "+force_install_dir", installDir, "+login", "anonymous", "+app_update", "996560", "validate", "+quit")
	steamcmdOut, err := steamcmdCmd.Output()

	if err != nil && runtime.GOOS == win {
		logger.Warn("Something went wrong calling SteamCMD. Falling back to bundled SteamCMD.")
		steamcmdCmd := exec.Command("./bundled/steamcmd.exe", "+force_install_dir", installDir, "+login", "anonymous", "+app_update", "996560", "validate", "+quit")
		steamcmdOut, err := steamcmdCmd.Output()
		utils.ShouldIPanic(err, fmt.Sprintf("Failed calling bundled SteamCMD: %s\n%s", err, steamcmdOut))
		logger.Output(string(steamcmdOut))
	} else {
		utils.ShouldIPanic(err, fmt.Sprintf("Something went wrong calling SteamCMD: %s\n%s", err, steamcmdOut))
		logger.Output(string(steamcmdOut))
	}
}

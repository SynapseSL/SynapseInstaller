# SynapseInstaller
A simple cli tool to install SCP:SL, the Synapse modloader, plugins and dependencies all from the same location.

## Installation
Grab yourself a binary from the releases tab.
Alternatively, you can build and install this yourself by running:
`go install github.com/SynapseSL/SynapseInstaller@latest`
Note: Building yourself requires Go installed on your machine.

### Dependencies
In order to download SCP:SL using `-install-game`, `steamcmd` is required to be in your path.

To extract Synapse and its plugins, a zip tool with cli support must be installed. On linux we use to `unzip`, on Windows `7za`.
If you want to use a different cli tool to unpack the `.zip` files, you may use the `-unzip-cmd` and `-unzip-args` switches.

## Usage:
To run SynapseInstaller interactively, just launch the executable. If you want to use SynapseInstaller from a script, you may do so by passing flags.
The following flags are recognized:

```
    -scripted
        Runs the installer in scripted mode, using the parameters specified by these flags.
        In interactive mode these flags will merely determine the recommended action to the user.
    -binaries <string>
        This is where SynapseInstaller will look for the game files. Required to install and/or update Synapse.
        If you install the game using "-install-game", this is where it'll be installed to.
        Defaults to "./SCPSL_DEDICATEDSERVER/"
    -files <string>
        This is where your local Synapse configs are stored - please point this to the PARENT directory.
        Example: On Linux, Synapse gets loaded from "~/.config/Synapse". Please set this to "~/.config/".
        Defaults to its respective folders on Windows and Linux.
    -install-game
        When passed, this will install/update SCP:SL in the location specified by "-binaries"
        Requires SteamCMD to be in your PATH.
    -install-synapse
        When passed, this will install/update Synapse.
    -scripted
        When passed, this will set the installer to scripted mode.
    -synapsezip <string>
        When passed, instead of downloading Synapse.zip from our CDN we will use the .zip provided.
    -unzip-args
        When passed, this will be passed to the specified unzip-cmd.
        This will only be passed if a custom unzip command was passed at all.
        NOTE: We auto-append the file name in question as its last argument.
    -unzip-cmd <string>
        When passed, this is the command that will be run to unzip the downloaded .zip files
    -verbosity <int>
        Determines how verbose logging should be. The higher the int, the less logs will be generated.
        Log Level 0:  Debug
        Log Level 1:  Info, Output
        Log Level 2:  Warn, OK
        Log Level 3:  Error
        Log Level >3: Only fatal errors.
```

## Platform specific notes:
### Windows
On Windows, you need to run SynapseInstaller from the disk where both game files and config files are stored. Usually this should be your `C:\` drive. If this is not possible, you'll probably need to install Synapse by hand. Don't worry, it's not hard either! Check out the [installation guide](https://github.com/SynapseSL/Synapse/blob/master/README.md#installation and you'll be good to go in no time!

## Honorable mentions
[AlmightyLks](https://github.com/AlmightyLks) has approved this project. Partially because he's affiliated with it.
Developed with <3 by [cubuzz](https://github.com/cubuzz).
Part of the [Synapse Modloader Project](https://github.com/SynapseSL/Synapse).
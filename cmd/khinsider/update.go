package khinsider

import (
	"github.com/marcus-crane/khinsider/v2/internal/updater"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"os"

	"github.com/marcus-crane/khinsider/v2/pkg/update"
)

func isUpdaterDisabled() bool {
	return os.Getenv("CI") == "true" || os.Getenv("KHINSIDER_NO_UPDATE") == "true"
}

func CheckForUpdates(c *cli.Context, currentVersion string, prerelease bool) (bool, string) {
	// TODO: Fix for go install which doesn't honour ldflags it seems?
	if isUpdaterDisabled() && c.Command.Name != "update" || currentVersion == "dev" {
		pterm.Debug.Println("Updater is disabled. Skipping update check.")
		return false, ""
	}
	remoteVersion := ""
	if !prerelease {
		remoteVersion = update.GetRemoteAppVersion()
		pterm.Debug.Printfln("Found remote version: %s", remoteVersion)
	} else {
		pterm.Debug.Printfln("Found remote prerelease version: %s", remoteVersion)
		remoteVersion = update.GetRemoteAppPrerelease()
	}
	isUpdateAvailable := update.IsRemoteVersionNewer(currentVersion, remoteVersion)
	pterm.Debug.Printfln("Current is %s while remote is %s. Update is available: %t", currentVersion, remoteVersion, isUpdateAvailable)
	if isUpdateAvailable {
		return true, remoteVersion
	}
	return false, remoteVersion
}

func UpdateAction(c *cli.Context, currentVersion string, prerelease bool) error {
	releaseAvailable, remoteVersion := CheckForUpdates(c, currentVersion, prerelease)
	pterm.Debug.Printfln("Release is available: %t. Remote version is %s", releaseAvailable, remoteVersion)
	if !releaseAvailable && !isUpdaterDisabled() {
		pterm.Info.Printfln("Sorry, no updates are available. The latest version is %s and you're on %s", remoteVersion, currentVersion)
	}
	return updater.UpgradeInPlace(c.App.Writer, c.App.ErrWriter, remoteVersion)
}

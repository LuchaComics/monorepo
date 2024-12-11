package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	osRuntime "runtime"
	"syscall"
)

func (a *App) GetPreferences() *Preferences {
	preferences := PreferencesInstance()
	return preferences
}

func (a *App) SavePreferences(pref *Preferences) error {
	//
	// STEP 1: Validation.
	//

	// Defensive code
	if pref == nil {
		return errors.New("Nothing inputted")
	}

	a.logger.Debug("Updating preferences...",
		slog.Any("DataDirectory", pref.DataDirectory),
		slog.Any("DefaultWalletAddress", pref.DefaultWalletAddress),
		slog.Any("NFTStorageAddress", pref.NFTStorageAddress),
		slog.Any("ChainID", pref.ChainID),
		slog.Any("AuthorityAddress", pref.AuthorityAddress),
	)

	//
	// STEP 2: Update preferences
	//

	if err := a.SetNFTStorageAddress(pref.NFTStorageAddress); err != nil {
		return err
	}
	if err := a.SetAuthorityAddress(pref.AuthorityAddress); err != nil {
		return err
	}

	//
	// STEP 3: Restart the server.
	//

	return restartSelf()
}

func restartSelf() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()
	// Windows does not support exec syscall.
	if osRuntime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Start()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, args, env)
}

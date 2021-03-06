// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fouwels/os-next/init/config"
	"github.com/fouwels/os-next/init/journal"
	"github.com/fouwels/os-next/init/shell"
)

func Login(auth config.Authenticators) error {

	success := false
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	for !success {
		fmt.Printf("\nenter authenticator for shell\n> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		santext := strings.TrimSuffix(text, "\n")

		switch auth.Root.Mode {
		case config.AuthenticatorsModeHash:
			result := checkAuthenticator(auth.Root.Value, santext)
			if result {
				success = true
			}
		case config.AuthenticatorsModeTOTP:
			result, err := checkTotp(auth.Root.Value, santext)
			if err == nil && result {
				success = true
			}
		}

		if !success {
			journal.Logfln("user failed to authenticate")
			time.Sleep(2 * time.Second)
		} else {
			journal.Logfln("user succeeded to authenticate")
		}
	}

	return nil
}

func Shell() error {

	fmt.Printf("\n") // Add final newline before dropping to shell

	commands := []shell.Command{
		{Executable: shell.Ash, Arguments: []string{}},
	}

	err := shell.Executor.ExecuteInteractive(commands)
	if err != nil {
		return err
	}

	return nil
}

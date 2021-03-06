// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"fmt"
	"os"
	"path"
)

func WriteSync(filename string, content []byte) error {

	f, err := os.OpenFile(path.Clean(filename), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to create: %v", err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	err = f.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}

	return nil
}

// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 K. Fouwels <k@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"init/config"
	"testing"
)

func TestLoads(t *testing.T) {

	c := config.Config{}
	configPrimary := config.PrimaryFile{}
	err := config.LoadConfig("primary_example.yml", &configPrimary)
	if err != nil {
		t.Fatalf("%v", err)
	}

	c.Primary = configPrimary.Primary

	configSecondary := config.SecondaryFile{}
	err = config.LoadConfig("secondary_example.yml", &configSecondary)
	if err != nil {
		t.Fatalf("%v", err)
	}
	c.Secondary = configSecondary.Secondary
}

package stages

import (
	"fmt"
	"init/config"
	"io/ioutil"
)

//Microcode implementes IStage
type Microcode struct {
	finals []string
}

//String ..
func (m *Microcode) String() string {
	return "microcode"
}

//Finalise ..
func (m *Microcode) Finalise() []string {
	return m.finals
}

//Run ..
func (m *Microcode) Run(c config.Config) error {

	err := ioutil.WriteFile("/sys/devices/system/cpu/microcode/reload", []byte("1"), 0644)
	if err != nil {
		return fmt.Errorf("failed to trigger microcode load: %w", err)
	}

	return nil
}

package stages

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"uinit-custom/config"
)

//Docker implementes IStage
type Docker struct {
	finals []string
}

//String ..
func (d Docker) String() string {
	return "Docker"
}

//Finalise ..
func (d Docker) Finalise() []string {
	return d.finals
}

//Run ..
func (d Docker) Run(c config.Config) error {

	// Start Docker

	cmd := exec.Command("dockerd")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_RAMDISK=true")
	cmd.Start()

	response := ""
	for i := 0; i < 5; i++ {

		resp, err := executeOne("docker version", "")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		response = resp
	}

	if response == "" {
		return fmt.Errorf("Failed to get docker version, docker did not start correctly")
	}

	d.finals = append(d.finals, fmt.Sprintf("Docker version: %v", response))

	return nil

}

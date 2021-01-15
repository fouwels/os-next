package stages

import (
	"fmt"
	"init-custom/config"
	"os"
)

//Networking implements IStage
type Networking struct {
	finals []string
}

//String ..
func (n *Networking) String() string {
	return "Networking"
}

//Finalise ..
func (n *Networking) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Networking) Run(c config.Config) (e error) {

	commands := []command{}
	commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", "lo", "up"}})
	for _, nd := range c.Secondary.Networking.Networks {

		if nd.Type != "" {
			// If type not default, create as specified
			commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "add", "dev", nd.Device, "type", nd.Type}})
		}

		if nd.DHCP {
			if nd.IPV6 {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", nd.Device, "up"}})
				commands = append(commands, command{command: "/sbin/udhcpc", arguments: []string{"-b", "-i", nd.Device, "-p", "/var/run/udhcpc.pid"}})
			} else {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", nd.Device, "up"}})
				commands = append(commands, command{command: "/sbin/udhcpc", arguments: []string{"-b", "-i", nd.Device, "-p", "/var/run/udhcpc.pid"}})
			}
		} else {

			commands = append(commands, command{command: "/sbin/ip", arguments: []string{"link", "set", "dev", nd.Device, "up"}})

			for _, v := range nd.Addresses {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"addr", "add", v, "dev", nd.Device}})
			}

			if nd.DefaultGateway != "" {
				commands = append(commands, command{command: "/sbin/ip", arguments: []string{"route", "add", "default", "via", nd.DefaultGateway, "dev", nd.Device}})
			}
		}
	}

	err := execute(commands)
	if err != nil {
		return err
	}

	commands = []command{}
	for _, rt := range c.Secondary.Networking.Routes {
		commands = append(commands, command{command: "/sbin/ip", arguments: []string{"route", "add", rt.Address, "dev", rt.Device}})
	}

	err = execute(commands)
	if err != nil {
		return err
	}

	if len(c.Secondary.Networking.Nameservers) != 0 {
		// #nosec G302 (CWE-276). 644 is intentional.
		f, err := os.OpenFile("/etc/resolv.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 644)
		if err != nil {
			return fmt.Errorf("Failed to open file to write nameservers: %v", err)
		}
		// #nosec G307. Double defer is safe for file.Writer
		defer f.Close()

		for _, ns := range c.Secondary.Networking.Nameservers {
			_, err = fmt.Fprintf(f, "nameserver %v\n", ns)
			if err != nil {
				return fmt.Errorf("Failed to write nameserver: %v", err)
			}
		}

		err = f.Sync()
		if err != nil {
			return fmt.Errorf("Failed to sync on %v: %v", f.Name(), err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("Failed to close on %v: %v", f.Name(), err)
		}

		n.finals = append(n.finals, fmt.Sprintf("nameservers configured to %v", c.Secondary.Networking.Nameservers))
	}
	return nil
}
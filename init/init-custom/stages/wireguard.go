package stages

import (
	"bytes"
	"fmt"
	"init-custom/config"
	"init-custom/util"
	"io/ioutil"
	"net"
	"path/filepath"
	"time"

	"github.com/mdp/qrterminal"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

//Wireguard implements IStage
type Wireguard struct {
	finals []string
}

//String ..
func (n *Wireguard) String() string {
	return "wireguard"
}

//Finalise ..
func (n *Wireguard) Finalise() []string {
	return n.finals
}

//Run ..
func (n *Wireguard) Run(c config.Config) error {

	const _keyroot = "/var/lib/docker"

	for _, wloop := range c.Secondary.Wireguard {
		wg := wloop //Prevent loop ref capture

		keypath := fmt.Sprintf("%v/%v.private", _keyroot, wg.Device)

		wgc, err := wgctrl.New()
		if err != nil {
			return err
		}

		wgkey := wgtypes.Key{}

		skey, err := ioutil.ReadFile(filepath.Clean(keypath))
		if err != nil {

			n.finals = append(n.finals, fmt.Sprintf("Private key generated for %v", wg.Device))

			wgkey, err = wgtypes.GeneratePrivateKey()
			if err != nil {
				return fmt.Errorf("Failed to generate private key: %v", err)
			}

			err = util.File.SetFile(keypath, wgkey.String(), 0600)
			if err != nil {
				return fmt.Errorf("Failed to save wg key: %v", err)
			}

		} else {
			wgkey, err = wgtypes.ParseKey(string(skey))
			if err != nil {
				return fmt.Errorf("Failed to parse loaded private key: %w", err)
			}

			n.finals = append(n.finals, "Private Key Loaded")
		}

		n.finals = append(n.finals, fmt.Sprintf("Public Key for %v: %v", wg.Device, wgkey.PublicKey().String()))
		n.finals = append(n.finals, fmt.Sprintf("\n%v", n.writeQR(wgkey.PublicKey())))

		wgpeers := []wgtypes.PeerConfig{}

		for _, v := range wg.Peers {

			vkey, err := wgtypes.ParseKey(v.PublicKey)
			if err != nil {
				return fmt.Errorf("Failed to parse key for %v: %w", v.Endpoint, err)
			}

			vudp, err := net.ResolveUDPAddr("udp", v.Endpoint)
			if err != nil {
				return fmt.Errorf("Failed to resolve endpoint for %v: %w", v.Endpoint, err)
			}

			keepalive := time.Duration(v.PersistentKeepalive) * time.Second

			ipnets := []net.IPNet{}

			for _, ap := range v.AllowedIps {
				_, ipnet, err := net.ParseCIDR(ap)
				if err != nil {
					return fmt.Errorf("Failed to parse allowedIP for %v: %w", ap, err)
				}

				ipnets = append(ipnets, *ipnet)
			}

			wgp := wgtypes.PeerConfig{
				PublicKey:                   vkey,
				Endpoint:                    vudp,
				PersistentKeepaliveInterval: &keepalive,
				AllowedIPs:                  ipnets,
			}

			wgpeers = append(wgpeers, wgp)
		}

		cfg := wgtypes.Config{
			PrivateKey: &wgkey,
			ListenPort: &wg.ListenPort,
			Peers:      wgpeers,
		}

		err = wgc.ConfigureDevice(wg.Device, cfg)
		if err != nil {
			return fmt.Errorf("Failed to configure wireguard device: %v", err)
		}
	}

	return nil
}
func (n Wireguard) writeQR(publicKey wgtypes.Key) string {

	var buf bytes.Buffer

	config := qrterminal.Config{
		Level:      qrterminal.L,
		HalfBlocks: false,
		Writer:     &buf,
		BlackChar:  qrterminal.BLACK,
		WhiteChar:  qrterminal.WHITE,
		QuietZone:  1,
	}

	qrterminal.GenerateWithConfig(publicKey.String(), config)

	return string(buf.Bytes())
}

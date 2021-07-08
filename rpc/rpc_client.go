package rpc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yankiwi/aruba_exporter/connector"
	
	"github.com/prometheus/common/log"
)

const (
	IOSXE string = "IOSXE"
	NXOS  string = "NXOS"
	IOS   string = "IOS"
)

// Client sends commands to a Cisco device
type Client struct {
	conn   *connector.SSHConnection
	Debug  bool
	OSType string
}

// NewClient creates a new client connection
func NewClient(ssh *connector.SSHConnection, debug bool) *Client {
	rpc := &Client{conn: ssh, Debug: debug}

	return rpc
}

// Identify tries to identify the OS running on a Aruba device
func (c *Client) Identify() error {
	output, err := c.RunCommand("show version")
	if err != nil {
		return err
	}
	
	log.Debugf("show version output: %s\n", output)
	
	switch {
	case strings.Contains(output, "IOS XE"):
		c.OSType = IOSXE
	case strings.Contains(output, "NX-OS"):
		c.OSType = NXOS
	case strings.Contains(output, "IOS Software"):
		c.OSType = IOS
	default:
		return errors.New("Unknown OS")
	}

	log.Debugf("Host %s identified as: %s\n", c.conn.Host, c.OSType)

	return nil
}

// RunCommand runs a command on a Aruba device
func (c *Client) RunCommand(cmd string) (string, error) {

	log.Debugf("Running command on %s: %s\n", c.conn.Host, cmd)

	output, err := c.conn.RunCommand(fmt.Sprintf("%s", cmd))
	if err != nil {
		println(err.Error())
		return "", err
	}

	return output, nil
}

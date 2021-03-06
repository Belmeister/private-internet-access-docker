package pia

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetPortForward() (port uint16, err error) {
	c.logger.Info("Obtaining port to be forwarded")
	b, err := c.random.GenerateRandomBytes(32)
	if err != nil {
		return 0, err
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	content, status, err := c.client.GetContent(url)
	switch {
	case err != nil:
		return 0, err
	case status != http.StatusOK:
		return 0, fmt.Errorf("status is %d for %s; does your PIA server support port forwarding?", status, url)
	case len(content) == 0:
		return 0, fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding")
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(content, &body); err != nil {
		return 0, fmt.Errorf("port forwarding response: %w", err)
	}
	c.logger.Info("Port forwarded is %d", body.Port)
	return body.Port, nil
}

func (c *configurator) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	c.logger.Info("Writing forwarded port to %s", filepath)
	return c.fileManager.WriteLinesToFile(
		string(filepath),
		[]string{fmt.Sprintf("%d", port)},
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func (c *configurator) AllowPortForwardFirewall(device models.VPNDevice, port uint16) (err error) {
	c.logger.Info("Allowing forwarded port %d through firewall", port)
	return c.firewall.AllowInputTrafficOnPort(device, port)
}

func (c *configurator) ClearPortForward(filepath models.Filepath, uid, gid int) (err error) {
	c.logger.Info("Clearing forwarded port status file %s", filepath)
	return c.fileManager.WriteToFile(string(filepath), nil, files.Ownership(uid, gid), files.Permissions(0400))
}

package pia

import (
	"net"

	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	GetOpenVPNConnections(region models.PIARegion, protocol models.NetworkProtocol,
		encryption models.PIAEncryption, targetIP net.IP) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, encryption models.PIAEncryption, verbosity, uid, gid int, root bool, cipher, auth string) (err error)
	GetPortForward() (port uint16, err error)
	WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error)
	ClearPortForward(filepath models.Filepath, uid, gid int) (err error)
	AllowPortForwardFirewall(device models.VPNDevice, port uint16) (err error)
}

type configurator struct {
	client      network.Client
	fileManager files.FileManager
	firewall    firewall.Configurator
	logger      logging.Logger
	random      random.Random
	verifyPort  func(port string) error
	lookupIP    func(host string) ([]net.IP, error)
}

// NewConfigurator returns a new Configurator object
func NewConfigurator(client network.Client, fileManager files.FileManager, firewall firewall.Configurator, logger logging.Logger) Configurator {
	return &configurator{
		client:      client,
		fileManager: fileManager,
		firewall:    firewall,
		logger:      logger.WithPrefix("PIA configurator: "),
		random:      random.NewRandom(),
		verifyPort:  verification.NewVerifier().VerifyPort,
		lookupIP:    net.LookupIP}
}

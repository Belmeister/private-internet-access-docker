package tinyproxy

import (
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type Configurator interface {
	Version() (string, error)
	MakeConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string, uid, gid int) error
	Start() (stdout io.ReadCloser, waitFn func() error, err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
}

func NewConfigurator(fileManager files.FileManager, logger logging.Logger) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger.WithPrefix("tinyproxy configurator: "),
		commander:   command.NewCommander()}
}

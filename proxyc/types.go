package proxyc

import (
	"errors"

	"github.com/hasura/gotel"
)

var tracer = gotel.NewTracer("relixy")

var (
	errServerURLRequired = errors.New("server url is required")
	errInvalidServerURL  = errors.New("invalid server URL")
	ErrNoAvailableServer = errors.New(
		"failed to initialize servers. Require at least 1 server has URL",
	)
)

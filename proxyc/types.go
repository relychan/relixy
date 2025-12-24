package proxyc

import (
	"errors"

	"github.com/hasura/gotel"
)

var tracer = gotel.NewTracer("relixy")

var (
	errServerURLRequired = errors.New("server url is required")
	errInvalidHTTPURL    = errors.New("invalid http url")
	errForbiddenHost     = errors.New("forbidden host")
	ErrNoAvailableServer = errors.New(
		"failed to initialize servers. Require at least 1 server has URL",
	)
)

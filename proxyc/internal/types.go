package internal

import "github.com/hasura/goenvconf"

// InsertRouteOptions represents options for inserting routes.
type InsertRouteOptions struct {
	GetEnv goenvconf.GetEnvFunc
}

package main

import (
	"github.com/conduitio/conduit/cmd/conduit/cli"
	"github.com/conduitio/conduit/pkg/conduit"

	chaos "github.com/conduitio-labs/conduit-connector-chaos"
	snowflake "github.com/conduitio-labs/conduit-connector-snowflake"
)

func main() {
	// Get the default configuration, including all built-in connectors
	cfg := conduit.DefaultConfig()

	// Add the Snowflake connector to list of built-in connectors
	cfg.ConnectorPlugins["snowflake"] = snowflake.Connector
	cfg.ConnectorPlugins["chaos"] = chaos.Connector

	cli.Run(cfg)
}

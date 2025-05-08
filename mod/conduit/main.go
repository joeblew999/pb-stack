package main

import (
	snowflake "github.com/conduitio-labs/conduit-connector-snowflake"
	"github.com/conduitio/conduit/cmd/conduit/cli"
	"github.com/conduitio/conduit/pkg/conduit"
)

func main() {
	// Get the default configuration, including all built-in connectors
	cfg := conduit.DefaultConfig()

	// Add the Snowflake connector to list of built-in connectors
	cfg.ConnectorPlugins["snowflake"] = snowflake.Connector

	cli.Run(cfg)
}

package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-arsys-baremetal/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version is set by goreleaser at build time.
var version string = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// Address published on the Terraform Registry (namespace/type).
		// Must match the `source` used by consumers.
		Address: "registry.terraform.io/arsys-internet/arsys-baremetal",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

# Terraform Arsys Baremetal Provider

A Terraform provider to manage Arsys Baremetal resources.

## Status

Alpha Status: This provider is under active development and is subject to change, and breaking changes may occur.
Not recommended for production use without proper testing.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.9.7
- [Go](https://go.dev/doc/install) >= 1.23

## Supported Services

    - Infrastructure
      - Servers
      - Server appliances
      - Datacenters
    - Network
      - Public IPs
      - Subnets
      - Firewall Policies
      - Private Networks
      - Public Networks
    - Security
      - SSH Keys

## Installation

### Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

To build the binary, run:
Linux:

```shell
set GOOS=linux && set GOARCH=amd64 && go build -o terraform-provider-arsys-baremetal
```

Windows:

```shell
set GOOS=windows&&set GOARCH=amd64&&go build -o terraform-provider-arsys-baremetal.exe
```

### Installing the provider locally

Since this provider is not yet published to the Terraform Registry, you need to install it locally. The installation
path depends on your operating system and CPU architecture.

Create the appropriate directory and copy the provider binary:
Linux:
For AMD64

```shell
mkdir -p ~/.terraform.d/plugins/local/arsys-baremetal/{provider-version}/linux_amd64/
cp $GOPATH/bin/terraform-provider-arsys-baremetal ~/.terraform.d/plugins/local/arsys-baremetal/{provider-version}/linux_amd64/terraform-provider-arsys-baremetal_v{provider-version}
```

Windows:
For AMD64

```shell
mkdir -p %APPDATA%\terraform.d\plugins\registry.terraform.io\local\arsys-baremetal\{provider-version}\windows_amd64
copy %GOPATH%\bin\terraform-provider-arsys-baremetal.exe %APPDATA%\terraform.d\plugins\registry.terraform.io\local\arsys-baremetal\{provider-version}\windows_amd64\terraform-provider-arsys-baremetal_v{provider-version}.exe
```

## Configuration

Configuring the provider to use it locally.

Generate .terraformrc file in $HOME directory:

```shell
cat > ~/.terraformrc << 'EOF'
provider_installation {
dev_overrides {
"registry.terraform.io/local/arsys-baremetal" = "/$HOME/.terraform.d/plugins/local/arsys-baremetal/0.1/linux_amd64"
}
direct {}
}
EOF
```

**IMPORTANT!**
You need to add the machine ip to your user to allow access via api in your Baremetal panel

Export the API token:

```shell
export BAREMETAL_API_TOKEN="{your-api-token}"
```

Add the provider block to your Terraform configuration file.

```hcl
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "local/arsys-baremetal"
      version = "{provider-version}"
    }
  }
}
```

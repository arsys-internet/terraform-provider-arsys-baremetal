# Terraform Arsys Baremetal Provider

A Terraform provider to manage Arsys Baremetal resources.

## Status

Alpha Status: This provider is under active development and is subject to change, and breaking changes may occur.
Not recommended for production use without proper testing.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.9.7
- [Go](https://go.dev/doc/install) >= 1.23 (only required to build the provider from source)

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

## Usage

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/arsys-internet/arsys-baremetal/latest).
To use it, declare it in the `required_providers` block and run `terraform init`; Terraform
will download it automatically.

```hcl
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.1"
    }
  }
}

provider "arsys-baremetal" {
  # token = "your-api-token"
}
```

**IMPORTANT!**
You need to add the machine IP to your user to allow API access from your Baremetal panel.

Export the API token before running Terraform:

```shell
export BAREMETAL_API_TOKEN="{your-api-token}"
```

See the [documentation](https://registry.terraform.io/providers/arsys-internet/arsys-baremetal/latest/docs)
and the [`examples/`](./examples) directory for resource and data source usage.

## Local development

These steps are only needed to build the provider from source and test it locally
(for example, when contributing).

### Building the provider

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

To test a locally built binary, install it under the Terraform plugins directory. The
installation path depends on your operating system and CPU architecture.

Linux (AMD64):

```shell
mkdir -p ~/.terraform.d/plugins/local/arsys-baremetal/{provider-version}/linux_amd64/
cp $GOPATH/bin/terraform-provider-arsys-baremetal ~/.terraform.d/plugins/local/arsys-baremetal/{provider-version}/linux_amd64/terraform-provider-arsys-baremetal_v{provider-version}
```

Windows (AMD64):

```shell
mkdir -p %APPDATA%\terraform.d\plugins\registry.terraform.io\local\arsys-baremetal\{provider-version}\windows_amd64
copy %GOPATH%\bin\terraform-provider-arsys-baremetal.exe %APPDATA%\terraform.d\plugins\registry.terraform.io\local\arsys-baremetal\{provider-version}\windows_amd64\terraform-provider-arsys-baremetal_v{provider-version}.exe
```

### Using the local build with dev_overrides

Generate a `.terraformrc` file in your `$HOME` directory pointing to the local build:

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

With `dev_overrides` active, point the `source` to the local namespace:

```hcl
terraform {
  required_providers {
    arsys-baremetal = {
      source = "local/arsys-baremetal"
    }
  }
}
```
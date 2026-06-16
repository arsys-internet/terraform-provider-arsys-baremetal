# Terraform Arsys Baremetal Provider

A Terraform provider to manage Arsys Baremetal resources through the
Arsys/CloudBuilder API.

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/arsys-internet/arsys-baremetal/latest).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.9.7
- [Go](https://go.dev/doc/install) >= 1.25 (only required to build the provider from source)

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

Add the provider to your Terraform configuration. Terraform will download it
from the registry automatically on `terraform init`:

```hcl
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.7"
    }
  }
}

provider "arsys-baremetal" {
  # token = "your-api-token"   # prefer the BAREMETAL_API_TOKEN environment variable
}
```

### Authentication

**IMPORTANT!** You need to add the machine IP to your user to allow API access
in your Baremetal panel.

Provide your API token via environment variable (recommended):

```shell
export BAREMETAL_API_TOKEN="{your-api-token}"
```

Optionally override the API base URL (defaults to `https://api.cloudbuilder.es/v1`):

```shell
export BAREMETAL_HOST="https://api.cloudbuilder.es/v1"
```

Then run:

```shell
terraform init
terraform plan
```

See the [`examples/`](./examples) directory and the
[registry documentation](https://registry.terraform.io/providers/arsys-internet/arsys-baremetal/latest/docs)
for per-resource usage.

## Development

These steps are only needed if you want to build and test the provider locally
instead of using the published registry version.

### Building the provider

Clone the repository, enter the directory and build the binary.

Linux:

```shell
GOOS=linux GOARCH=amd64 go build -o terraform-provider-arsys-baremetal
```

Windows:

```shell
set GOOS=windows&&set GOARCH=amd64&&go build -o terraform-provider-arsys-baremetal.exe
```

### Using a local build (dev overrides)

To test a local binary without publishing it, configure `dev_overrides` in your
Terraform CLI config file (`~/.terraformrc`), pointing to the directory that
contains the compiled binary:

```shell
cat > ~/.terraformrc << 'EOF'
provider_installation {
  dev_overrides {
    "arsys-internet/arsys-baremetal" = "/absolute/path/to/your/built/provider/dir"
  }
  direct {}
}
EOF
```

With `dev_overrides` active you do **not** run `terraform init`; just run
`terraform plan` / `terraform apply` directly.

### Running tests

```shell
# Unit tests (no API calls)
TF_ACC=0 go test -v ./...

# Acceptance tests (run against the real API)
TF_ACC=1 BAREMETAL_API_TOKEN=xxx go test -v -timeout=120m ./...
```

---
subcategory: "Infrastructure"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Server Data Sources"
sidebar_current: "docs-datasource-server"
description: |-
  Get information about Servers in Arsys Baremetal.
---

# arsys-baremetal\_server

The **Servers data source** can be used to search and return all existing servers.
Also, it can be used to search for and return an existing server.
You can provide a string for the id parameter which will be compared with the created server.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Important Note

**The attributes returned by `arsys-baremetal_servers` (list) and `arsys-baremetal_server` (individual) datasource are
different:**

### Key Differences:

**`arsys-baremetal_servers` (Server List):**

- `status` object contains only: `state`
- **Does NOT include** recovery-related fields

**`arsys-baremetal_server` (Server Detail):**

- `status` object contains: `state` + `percent`
- **Includes additional fields:**
    - `recovery_mode`
    - `recovery_image_os`
    - `recovery_user`
    - `recovery_password`

Use the individual server datasource when you need complete server information, including recovery details and operation
progress.

## Example Usage

### Get all Servers

```hcl
data "arsys-baremetal_servers" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `servers` - List of servers
* ## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier for the servers' data source
* `servers` - List of servers with the following attributes:
    * `id` - Identifier of the server
    * `name` - The name of the server
    * `description` - The description of the server (nullable)
    * `server_type` - The server type (e.g., "baremetal")
    * `creation_date` - The date when the server was created (ISO 8601 format)
    * `first_password` - The initial password for the server (nullable)
    * `managed` - Whether the server is managed
    * `ssh_password` - Whether SSH password authentication is enabled
    * `rsa_key` - Whether RSA key authentication is enabled
    * `hostname` - The hostname of the server
    * `cloudpanel_id` - The CloudPanel identifier (nullable)
    * `datacenter` - Datacenter information
        * `id` - Datacenter identifier
        * `country_code` - Country code (e.g., "ES")
        * `location` - Datacenter location (e.g., "Spain")
    * `image` - Operating system image
        * `id` - Image identifier
        * `name` - Image name (e.g., "Rocky Linux 9")
    * `hardware` - Hardware configuration of the server
        * `fixed_instance_size_id` - Fixed instance size identifier (nullable)
        * `baremetal_model_id` - Baremetal model identifier (nullable)
        * `vcore` - Number of virtual cores
        * `cores_per_processor` - Number of cores per processor
        * `ram` - RAM size in GB
        * `hdds` - List of hard disk drives
            * `id` - HDD identifier
            * `size` - HDD size in GB
            * `is_main` - Whether this is the main drive
            * `disk_type` - Type of disk (e.g., "HDD", "SSD")
            * `disk_raid` - RAID configuration (e.g., "Software RAID 1")
            * `disk_raid_count` - Number of disks in RAID
    * `status` - Server status information
        * `state` - Current state (e.g., "POWERED_OFF", "POWERED_ON")
    * `ips` - List of assigned IP addresses
        * `id` - IP identifier
        * `ip` - IP address
        * `type` - IP type (e.g., "IPV4", "IPV6")
        * `reverse_dns` - Reverse DNS configuration (nullable)
        * `main` - Whether this is the main IP
        * `firewall_policy` - Firewall policy configuration (nullable)
            * `id` - Firewall policy identifier
            * `name` - Firewall policy name
        * `load_balancers` - List of associated load balancers
            * `id` - Load balancer identifier
            * `name` - Load balancer name
    * `private_networks` - List of private networks
        * `id` - Private network identifier
        * `name` - Private network name
        * `server_ip` - Server IP in the private network
    * `connection_speed` - Network connection speeds (nullable)
        * `private` - Private network speed configuration
            * `available` - Available speeds (array)
            * `current` - Current speed
        * `public` - Public network speed configuration
            * `available` - Available speeds (array)
            * `current` - Current speed
    * `redundancy` - Redundancy configuration (nullable)
        * `available` - Whether redundancy is available
        * `enabled` - Whether redundancy is enabled
    * `dvd` - DVD configuration (nullable)
        * `id` - DVD identifier
        * `name` - DVD name
    * `alerts` - Alert configuration (nullable)
    * `monitoring_policy` - Monitoring policy (nullable)
        * `id` - Monitoring policy identifier
        * `name` - Monitoring policy name
    * `snapshot` - Snapshot information (nullable)

### Get By Id

```hcl
data "arsys-baremetal_server" "example" {
  id = server_id
}
```

## Argument Reference

* `id` - (Required)ID of an existing server that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the server
* `name` - The name of the server
* `description` - The description of the server (nullable)
* `server_type` - The server type (e.g., "baremetal")
* `creation_date` - The date when the server was created (ISO 8601 format)
* `first_password` - The initial password for the server (nullable)
* `managed` - Whether the server is managed
* `ssh_password` - Whether SSH password authentication is enabled
* `rsa_key` - Whether RSA key authentication is enabled
* `hostname` - The hostname of the server
* `cloudpanel_id` - The CloudPanel identifier (nullable)
* `datacenter` - Datacenter information
    * `id` - Datacenter identifier
    * `country_code` - Country code (e.g., "ES")
    * `location` - Datacenter location (e.g., "Spain")
* `image` - Operating system image
    * `id` - Image identifier
    * `name` - Image name (e.g., "Rocky Linux 9")
* `hardware` - Hardware configuration of the server
    * `fixed_instance_size_id` - Fixed instance size identifier (nullable)
    * `baremetal_model_id` - Baremetal model identifier (nullable)
    * `vcore` - Number of virtual cores
    * `cores_per_processor` - Number of cores per processor
    * `ram` - RAM size in GB
    * `hdds` - List of hard disk drives
        * `id` - HDD identifier
        * `size` - HDD size in GB
        * `is_main` - Whether this is the main drive
        * `disk_type` - Type of disk (e.g., "HDD", "SSD")
        * `disk_raid` - RAID configuration (e.g., "Software RAID 1")
        * `disk_raid_count` - Number of disks in RAID
* `status` - Server status information
    * `state` - Current state (e.g., "POWERED_OFF", "POWERED_ON")
    * `percent` - Progress percentage for operations (nullable)
* `ips` - List of assigned IP addresses
    * `id` - IP identifier
    * `ip` - IP address
    * `type` - IP type (e.g., "IPV4", "IPV6")
    * `reverse_dns` - Reverse DNS configuration (nullable)
    * `main` - Whether this is the main IP
    * `firewall_policy` - Firewall policy configuration (nullable)
        * `id` - Firewall policy identifier
        * `name` - Firewall policy name
    * `load_balancers` - List of associated load balancers
        * `id` - Load balancer identifier
        * `name` - Load balancer name
* `private_networks` - List of private networks
    * `id` - Private network identifier
    * `name` - Private network name
    * `server_ip` - Server IP in the private network
    * `vlan_id` - VLAN identifier (nullable)
* `connection_speed` - Network connection speeds (nullable)
    * `private` - Private network speed configuration
        * `available` - Available speeds (array)
        * `current` - Current speed
    * `public` - Public network speed configuration
        * `available` - Available speeds (array)
        * `current` - Current speed
* `redundancy` - Redundancy configuration (nullable)
    * `available` - Whether redundancy is available
    * `enabled` - Whether redundancy is enabled
* `dvd` - DVD configuration (nullable)
    * `id` - DVD identifier
    * `name` - DVD name
* `alerts` - Alert configuration (nullable)
* `monitoring_policy` - Monitoring policy (nullable)
    * `id` - Monitoring policy identifier
    * `name` - Monitoring policy name
* `snapshot` - Snapshot information (nullable)
* `recovery_mode` - Whether server is in recovery mode
* `recovery_image_os` - Recovery image OS (nullable)
* `recovery_user` - Recovery user (nullable)
* `recovery_password` - Recovery password (nullable)
# Changelog

## [0.1.0] - 2025-08-11

- Added Datacenter management
- Added Public IPs management
- Added Firewall Policies management
- Added Private Network management
- Added Public Network management
- Added Server appliances management

## [0.2.0] - 2025-08-21

- Added Baremetal server management
- Added Firewall Policies server IPs management (list, assign, get details)
- Added Firewall Policies rules management (list, create, get details, delete)

## [0.3.0] - 2025-08-28

- Added Subnets management
- Improved Public IPs filtering

## [0.4.0] - 2025-09-05

- Improved Baremetal models filtering
- Added SSH keys management (list, create, get details, delete)

## [0.4.1] - 2025-09-08

- Fixed name regex validation
- Fixed regex validation for id fields
- Fixed error for non-updatable fields in Ips
- Fixed server description to handle null values

## [0.5.0] - 2025-09-11

- Added Public network IPs assign/unassign
- Added Private network Servers assign/unassign

## [0.6.0] - 2025-09-18

- Added Baremetal Models documentation and examples
- Fixed Schema and error messages
- Fixed Public network servers where public network was not returned after servers were assigned

## [0.6.1] - 2025-09-18

- Fix server appliance type schema validation
# Provider Configuration

This directory contains the base provider configuration for the Arsys Baremetal provider used across all examples.

## What This Contains

- Base Terraform configuration with provider requirements
- Provider configuration with authentication setup instructions
- Common settings that other examples can reference

## Authentication

The provider supports authentication via environment variables (recommended) or direct configuration.

### Environment Variables (Recommended)
```bash
export ARSYS_BAREMETAL_TOKEN="your-api-token-here"
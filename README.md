# terraform-plugin-framework-docs

A library to generate rich documents for terraform providers based on [`terraform-plugin-framework`](https://github.com/hashicorp/terraform-plugin-docs).  

## Features

- Full coverage of provider resources, including:
    - Provider
    - Managed Resource
    - Data Source
    - Ephemeral Resource
    - List Resource
    - Action
    - Function
- Rich document of the provider, including descriptions from `CustomType`, `PlanModifiers` and `Validators`.
- Flexible to add examples for all provider resources, the API simply takes a string, which can be derived from the acceptance test's config, to make sure the examples are up-to-date.
- Flexible to add sub-category for the generated documents (This solves the issue: https://github.com/hashicorp/terraform-plugin-docs/issues/156).
- `fwdtypes` package provides a set of types to be used in place of the `terraform-plugin-framework` `basetypes`, but with a description. This is useful to be used in place of a plain `ObjectType`, to add description to the members (This solves the issue: https://github.com/hashicorp/terraform-plugin-docs/issues/333).

## Why not `tfplugindocs`

The Terraform official solution for provider document generation is https://github.com/hashicorp/terraform-plugin-docs, which supports both `terraform-plugin-framework` and `terraform-plugin-sdk` based providers. The current `terraform-plugin-docs` tool is dependent on using Terraform CLI's `terraform providers schema -json` [output](https://developer.hashicorp.com/terraform/cli/commands/providers/schema) which restricts it in what can be generated to the document. See https://github.com/hashicorp/terraform-plugin-framework/issues/625#issuecomment-1424690927 (and any issues referenced by this).

This library works in another way, it read the schema from the provider code base. The user is supposed to create a separate Go package along side the provider's `internal` package to use this library to generate the documents. An example can be found at https://github.com/magodo/terraform-plugin-framework-docs/blob/main/tffwdocs_test.go.

## Example

- https://registry.terraform.io/providers/magodo/restful/latest/docs

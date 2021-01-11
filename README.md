# Terraform Provider IT Glue

This provider is in extremely early stages. Use at your own risk.

It uses our incomplete [Go IT Glue API wrapper](https://github.com/Private-Universe/itglue) for the backend.

## Table of contents

- [Terraform Provider IT Glue](#terraform-provider-it-glue)
  - [Table of contents](#table-of-contents)
  - [Installation](#installation)
  - [Upgrading](#upgrading)
  - [Setup](#setup)
    - [Example providing API key using AWS Parameter Store](#example-providing-api-key-using-aws-parameter-store)
  - [Example Usage](#example-usage)
    - [Flexible assets](#flexible-assets)
      - [Limitations](#limitations)
    - [Passwords](#passwords)

## Installation

Run the following command to build the provider or download a GitHub release.

```shell
go build -o terraform-provider-itglue.exe
```

Move the executable to the below directory (replace x.x.x with the appropriate version number)
```
%AppData%\terraform.d\plugins\privateuniverse.com.au\pu\itglue\x.x.x\windows_amd64\
```

## Upgrading

In the `%AppData%\terraform.d\plugins\privateuniverse.com.au\pu\itglue\` directory, add a new folder structure for the new version of the provider.

In your Terraform directory, bump up the version number of the provider and run the following command:

```shell
terraform init -upgrade
```

## Setup

In your Terraform configuration, add the provider (note: AWS is not required)
```terraform
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.8.0"
    }
    itglue = {
      source  = "privateuniverse.com.au/pu/itglue"
      version = "0.1.18"
    }
  }
}
```

Add the IT Glue API key to the provider
```terraform
provider "itglue" {
  api_key = local.itglue_api_key.key
}
```

**It is recommended to provide the API key from something like AWS Parameter Store**

### Example providing API key using AWS Parameter Store

In AWS Parameter Store as a secure string with the name `/company/itglue/apikey`:
```json
{
    "key":"ITG.exampleapikey.longstringofrandomcharacters"
}
```

In your locals:
```terraform
data "aws_ssm_parameter" "itglueapikey" {
  name = "/company/itglue/apikey"
}

locals {
  itglue_api_key = jsondecode(
    data.aws_ssm_parameter.itglueapikey.value
  )
}
```

Use `local.itglue_api_key.key` where needed.

## Example Usage

Currently only the below are supported.

Refer to the [official IT Glue API documentation](https://api.itglue.com/developer/) for valid fields. Please submit an issue or pull request if something is not supported.

### Flexible assets

For flexible assets, you need to provide traits, organization_id and flexible_asset_type_id.

The traits are based upon how your flexible asset is set up and can be any string.

You also need to provide the flexible asset type ID which has the traits and the organization ID that you want the flexible asset to be listed under.

```terraform
resource "itglue_flexible_asset" "example_server" {
  traits = {
        company-name         = var.tag_company_name
        admin-username       = "test"
        url                  = "https://example.com"
        internal-ip-address  = "1.1.1.1"
        external-ip-address  = "1.1.1.2"
        link-to-organisation = 123456
        license-key          = itglue_flexible_asset.server_license.id
        notes                = "This document is managed by Terraform. Changes will be overridden."
  }
  organization_id        = 123457
  flexible_asset_type_id = 45678
}

resource "itglue_flexible_asset" "server_license" {
  traits = {
        license-key  = "AAAA-BBBB-CCCC-DDDD"
        renewal-date = "2021-09-09"
        renewal-type = "Monthly"
        notes        = "This document is managed by Terraform. Changes will be overridden."
  }
  organization_id        = 123457
  flexible_asset_type_id = 45679
}
```

#### Limitations

*Currently password traits are not supported but all other trait types should work when passed a string, integer or boolean.*

*Currently tag traits can only have one ID specified.*

### Passwords

Passwords are recommended to be passed from something like AWS Parameter Store (refer to [Example providing API key using AWS Parameter Store](#example-providing-api-key-using-aws-parameter-store)).

Example with an embedded password, embedded into the above server flexible asset

```terraform
resource "itglue_password" "server_password" {
  name            = "test server password"
  username        = "testusername"
  password        = "testpassword"
  resource_id     = itglue_flexible_asset.example_server.id
  resource_type   = "Flexible Asset"
  organization_id = 123457
  notes           = "This document is managed by Terraform. Changes will be overridden."
}
```

Example with a standalone password

```terraform
resource "itglue_password" "server_password" {
  name            = "test server password"
  username        = "testusername"
  password        = "testpassword"
  organization_id = 123457
  notes           = "This document is managed by Terraform. Changes will be overridden."
}
```
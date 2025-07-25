---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "AWS: aws_vpc"
description: |-
  Provides a VPC resource.
---

# Resource: aws_vpc

Provides a VPC resource.

## Example Usage

Basic usage:

```terraform
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}
```

Basic usage with tags:

```terraform
resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"

  tags = {
    Name = "main"
  }
}
```

VPC with CIDR from AWS IPAM:

```terraform
data "aws_region" "current" {}

resource "aws_vpc_ipam" "test" {
  operating_regions {
    region_name = data.aws_region.current.name
  }
}

resource "aws_vpc_ipam_pool" "test" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.test.private_default_scope_id
  locale         = data.aws_region.current.name
}

resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "172.20.0.0/16"
}

resource "aws_vpc" "test" {
  ipv4_ipam_pool_id   = aws_vpc_ipam_pool.test.id
  ipv4_netmask_length = 28
  depends_on = [
    aws_vpc_ipam_pool_cidr.test
  ]
}
```

## Argument Reference

This resource supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `cidr_block` - (Optional) The IPv4 CIDR block for the VPC. CIDR can be explicitly set or it can be derived from IPAM using `ipv4_netmask_length`.
* `instance_tenancy` - (Optional) A tenancy option for instances launched into the VPC. Default is `default`, which ensures that EC2 instances launched in this VPC use the EC2 instance tenancy attribute specified when the EC2 instance is launched. The only other option is `dedicated`, which ensures that EC2 instances launched in this VPC are run on dedicated tenancy instances regardless of the tenancy attribute specified at launch. This has a dedicated per region fee of $2 per hour, plus an hourly per instance usage fee.
* `ipv4_ipam_pool_id` - (Optional) The ID of an IPv4 IPAM pool you want to use for allocating this VPC's CIDR. IPAM is a VPC feature that you can use to automate your IP address management workflows including assigning, tracking, troubleshooting, and auditing IP addresses across AWS Regions and accounts. Using IPAM you can monitor IP address usage throughout your AWS Organization.
* `ipv4_netmask_length` - (Optional) The netmask length of the IPv4 CIDR you want to allocate to this VPC. Requires specifying a `ipv4_ipam_pool_id`.
* `ipv6_cidr_block` - (Optional) IPv6 CIDR block to request from an IPAM Pool. Can be set explicitly or derived from IPAM using `ipv6_netmask_length`.
* `ipv6_ipam_pool_id` - (Optional) IPAM Pool ID for a IPv6 pool. Conflicts with `assign_generated_ipv6_cidr_block`.
* `ipv6_netmask_length` - (Optional) Netmask length to request from IPAM Pool. Conflicts with `ipv6_cidr_block`. This can be omitted if IPAM pool as a `allocation_default_netmask_length` set. Valid values are from `44` to `60` in increments of 4.
* `ipv6_cidr_block_network_border_group` - (Optional) By default when an IPv6 CIDR is assigned to a VPC a default ipv6_cidr_block_network_border_group will be set to the region of the VPC. This can be changed to restrict advertisement of public addresses to specific Network Border Groups such as LocalZones.
* `enable_dns_support` - (Optional) A boolean flag to enable/disable DNS support in the VPC. Defaults to true.
* `enable_network_address_usage_metrics` - (Optional) Indicates whether Network Address Usage metrics are enabled for your VPC. Defaults to false.
* `enable_dns_hostnames` - (Optional) A boolean flag to enable/disable DNS hostnames in the VPC. Defaults false.
* `assign_generated_ipv6_cidr_block` - (Optional) Requests an Amazon-provided IPv6 CIDR block with a /56 prefix length for the VPC. You cannot specify the range of IP addresses, or the size of the CIDR block. Default is `false`. Conflicts with `ipv6_ipam_pool_id`
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name (ARN) of VPC
* `id` - The ID of the VPC
* `instance_tenancy` - Tenancy of instances spin up within VPC
* `dhcp_options_id` - DHCP options id of the desired VPC.
* `enable_dns_support` - Whether or not the VPC has DNS support
* `enable_network_address_usage_metrics` - Whether Network Address Usage metrics are enabled for the VPC
* `enable_dns_hostnames` - Whether or not the VPC has DNS hostname support
* `main_route_table_id` - The ID of the main route table associated with
     this VPC. Note that you can change a VPC's main route table by using an
     [`aws_main_route_table_association`](/docs/providers/aws/r/main_route_table_association.html).
* `default_network_acl_id` - The ID of the network ACL created by default on VPC creation
* `default_security_group_id` - The ID of the security group created by default on VPC creation
* `default_route_table_id` - The ID of the route table created by default on VPC creation
* `ipv6_association_id` - The association ID for the IPv6 CIDR block.
* `ipv6_cidr_block_network_border_group` - The Network Border Group Zone name
* `owner_id` - The ID of the AWS account that owns the VPC.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import VPCs using the VPC `id`. For example:

```terraform
import {
  to = aws_vpc.test_vpc
  id = "vpc-a01106c2"
}
```

Using `terraform import`, import VPCs using the VPC `id`. For example:

```console
% terraform import aws_vpc.test_vpc vpc-a01106c2
```

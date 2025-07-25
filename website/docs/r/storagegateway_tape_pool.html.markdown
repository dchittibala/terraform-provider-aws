---
subcategory: "Storage Gateway"
layout: "aws"
page_title: "AWS: aws_storagegateway_tape_pool"
description: |-
  Manages an AWS Storage Gateway Tape Pool
---

# Resource: aws_storagegateway_tape_pool

Manages an AWS Storage Gateway Tape Pool.

## Example Usage

```terraform
resource "aws_storagegateway_tape_pool" "example" {
  pool_name     = "example"
  storage_class = "GLACIER"
}
```

## Argument Reference

This resource supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `pool_name` - (Required) The name of the new custom tape pool.
* `storage_class` - (Required) The storage class that is associated with the new custom pool. When you use your backup application to eject the tape, the tape is archived directly into the storage class that corresponds to the pool. Possible values are `DEEP_ARCHIVE` or `GLACIER`.
* `retention_lock_type` - (Required) Tape retention lock can be configured in two modes. When configured in governance mode, AWS accounts with specific IAM permissions are authorized to remove the tape retention lock from archived virtual tapes. When configured in compliance mode, the tape retention lock cannot be removed by any user, including the root AWS account. Possible values are `COMPLIANCE`, `GOVERNANCE`, and `NONE`. Default value is `NONE`.
* `retention_lock_time_in_days` - (Optional) Tape retention lock time is set in days. Tape retention lock can be enabled for up to 100 years (36,500 days). Default value is 0.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - Volume Amazon Resource Name (ARN), e.g., `aws_storagegateway_tape_pool.example arn:aws:storagegateway:us-east-1:123456789012:tapepool/pool-12345678`.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import `aws_storagegateway_tape_pool` using the volume Amazon Resource Name (ARN). For example:

```terraform
import {
  to = aws_storagegateway_tape_pool.example
  id = "arn:aws:storagegateway:us-east-1:123456789012:tapepool/pool-12345678"
}
```

Using `terraform import`, import `aws_storagegateway_tape_pool` using the volume Amazon Resource Name (ARN). For example:

```console
% terraform import aws_storagegateway_tape_pool.example arn:aws:storagegateway:us-east-1:123456789012:tapepool/pool-12345678
```

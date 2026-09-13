## Requirements

| Name | Version |
| ---- | ------- |
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.12.0 |
| <a name="requirement_helm"></a> [helm](#requirement\_helm) | 3.2.0 |

## Providers

| Name | Version |
| ---- | ------- |
| <a name="provider_helm"></a> [helm](#provider\_helm) | 3.2.0 |

## Modules

No modules.

## Resources

| Name | Type |
| ---- | ---- |
| [helm_release.this](https://registry.terraform.io/providers/hashicorp/helm/3.2.0/docs/resources/release) | resource |

## Inputs

| Name | Description | Type | Default | Required |
| ---- | ----------- | ---- | ------- | :------: |
| <a name="input_name"></a> [name](#input\_name) | Release name | `string` | `"traefik"` | no |
| <a name="input_namespace"></a> [namespace](#input\_namespace) | Namespace to install the release into | `string` | `"traefik"` | no |
| <a name="input_values"></a> [values](#input\_values) | List of values in raw yaml format to pass to helm | `list(string)` | `[]` | no |

## Outputs

No outputs.

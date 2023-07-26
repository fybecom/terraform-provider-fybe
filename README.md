# Fybe terraform provider

Manage your Fybe infrastructure not only via our [cli](https://github.com/fybecom/fybe) or [APIs](https://api.fybe.com/#/), but also via [Terraform](https://www.terraform.io/)!  

## Getting Started

* [Link to terraform page](https://registry.terraform.io/providers/fybe/fybe/latest)
* [Documentation link to terraform page](https://registry.terraform.io/providers/fybe/fybe/latest/docs)

1. Install [terraform cli](https://learn.hashicorp.com/tutorials/terraform/install-cli)
2. Copy the example `examples/main.tf.example` as `.tf` file to you project directory
3. Run terraform

    ```sh
    terraform init
    terraform plan
    # CAUTION:  with example main.tf you are about to order and pay an object storage
    terraform apply
    ```

## Local Development

1. Install [terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
2. Install `go` from here [Golang Page](https://go.dev/doc/install)
3. `git clone https://github.com/fybe/terraform.git`
4. Generate fybe api client client `rm -rf apiclient; docker run --rm -v "${PWD}:/local" --env JAVA_OPTS='-Dio.swagger.parser.util.RemoteUrl.trustAll=true -Dio.swagger.v3.parser.util.RemoteUrl.trustAll=true' openapitools/openapi-generator-cli:v5.2.1 generate --skip-validate-spec --input-spec 'https://api.fybe.com/api-v1.yaml' --generator-name go --output /local/apiclient`
5. Compile it `go mod tidy && go mod download && go build -gcflags="all=-N -l"  -o terraform-provider-fybe`
6. create `~/.terraformrc` with following content

    ```terraform
    provider_installation {

      dev_overrides {
        "fybe/fybe" = "/PATH_TO_YOUR/BINARY_BUILD"
      }

      direct {}
    }
    ```

7. Copy `./examples/main.tf.example` to `./main.tf` and fill in the provider config.
8. In the same directory execute

```sh
terraform plan
# CAUTION:  with example main.tf you are about to buy and pay an compute instance
terraform apply
```

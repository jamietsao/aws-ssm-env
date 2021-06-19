# aws-ssm-env

Simple utility to print parameters from Amazon EC2 Systems Manager (ssm) Parameter Store as environment variables. This is useful for injecting secure secrets into the environment of a docker container process.

## Environment Variables

### AWS_SSM_ENV_BACKOFF

Set `AWS_SSM_ENV_BACKOFF` to a Golang `time.Duration` string to set the automatic maximum exponential backoff for AWS API retries.

```bash
export AWS_SSM_ENV_BACKOFF="30s"
```

## Usage

Create secret parameters on AWS Parameter Store for your application using [hierarchies](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html#sysman-paramstore-su-organize):

```bash
aws ssm put-parameter --name /userservice/SECRET_1 --value "123456" --type SecureString
aws ssm put-parameter --name /accountservice/secret_2 --value "abcdef" --type SecureString
aws ssm put-parameter --name /accountservice/secret_3 --value "foobarbaz" --type SecureString
aws ssm put-parameter --name /database/production/password --value "productionpass" --type SecureString
aws ssm put-parameter --name /database/staging/password --value "stagingpass" --type SecureString
```

Add tags to categorize parameters in various ways:

```bash
aws ssm add-tags-to-resource --resource-type Parameter --resource-id /userservice/SECRET_1 --tags Key=userservice,Value=true Key=production,Value=true

aws ssm add-tags-to-resource --resource-type Parameter --resource-id /accountservice/secret_2 --tags Key=accountservice,Value=true Key=production,Value=true

aws ssm add-tags-to-resource --resource-type Parameter --resource-id /accountservice/secret_3 --tags Key=accountservice,Value=true Key=staging,Value=true

aws ssm add-tags-to-resource --resource-type Parameter --resource-id /database/production/password --tags Key=userservice,Value=true Key=accountservice,Value=true Key=production,Value=true

aws ssm add-tags-to-resource --resource-type Parameter --resource-id /database/staging/password --tags Key=userservice,Value=true Key=accountservice,Value=true Key=staging,Value=true
```

Retrieve parameters with `aws-ssm-env`:

```bash
# filter by 'userservice' parameters for 'production'
AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=userservice,production
SECRET_1=123456
PASSWORD=productionpass

# filter by 'accountservice' parameters for 'staging'
AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=accountservice,staging
SECRET_3=foobarbaz
PASSWORD=stagingpass

# filter by path (`/` will search all parameters)
AWS_REGION=<aws-region> aws-ssm-env --paths=/database
PASSWORD=productionpass
PASSWORD=stagingpass

# filter by path and tag
AWS_REGION=<aws-region> aws-ssm-env --paths=/database --tags=production
PASSWORD=productionpass
```

*Notice that parameter names are automatically capitalized.*

Use `export` with `aws-ssm-env` to inject secrets from Parameter Store into the environment.
Since `aws-ssm-env` may fail, it's recommended to capture the output and then export the variables:

```bash
VARS=$(AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=userservice,production)
if [[ $? -ne 0 ]]; then
    export "$VARS"
fi
```

If no error checking is needed, export environment variables directly for all parameters returned:

```
> export $(AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=userservice,production)
> env
...
...
SECRET_1=123456
PASSWORD=productionpass
```

## Installation

Install directly go get:

```bash
go get github.com/joberly/aws-ssm-env
```

Or download [binary](https://github.com/joberly/aws-ssm-env/releases/latest):

```bash
# Replace the value of VERSION with the version to download.
VERSION=v1.2.0
wget -O aws-ssm-env.zip https://github.com/joberly/aws-ssm-env/releases/download/${VERSION}/aws-ssm-env-${VERSION}-linux-amd64.zip
unzip aws-ssm-env.zip
chmod 755 aws-ssm-env
```

## Authors

John Oberly III

Jamie Tsao (original utility at https://github.com/jamietsao/aws-ssm-env)

## License

[See LICENSE file.](LICENSE)

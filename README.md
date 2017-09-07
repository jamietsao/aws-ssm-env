# aws-ssm-env
Simple utility to print parameters from Amazon EC2 Systems Manager Parameter Store as environment variables. This is useful for injecting secure secrets into the environment of a docker container process.

### Usage
Create secret parameters on AWS Parameter Store for your application using [hierarchies](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html#sysman-paramstore-su-organize)  :
```
> aws ssm put-parameter --name /production/app-1/SECRET_1 --value "123456" --type SecureString --key-id <your-key-id> --region <aws-region>
> aws ssm put-parameter --name /production/app-1/secret_2 --value "abcdef" --type SecureString --key-id <your-key-id> --region <aws-region>
```
Use `export` with `aws-ssm-env` to inject secrets from Parameter Store into the environment:
```
> export $(AWS_REGION=<aws-region> aws-ssm-env -paths /production/app-1/)
> env
...
...
SECRET_1=123456
SECRET_2=abcdef
```
Notice that parameter names are automatically capitalized

Multiple hierarchy paths can be passed in vis `-paths` (comma separated):
```
> aws ssm put-parameter --name /production/common/common_secret --value "foobarbaz" --type SecureString --key-id <your-key-id> --region <aws-region>
> export $(AWS_REGION=<aws-region> aws-ssm-env -paths /production/app-1/,/production/common/)
> env
...
...
SECRET_1=123456
SECRET_2=abcdef
COMMON_SECRET=foobarbaz
```

### Author
Jamie Tsao

### License
MIT

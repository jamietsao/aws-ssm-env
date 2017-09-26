# aws-ssm-env
Simple utility to print parameters from Amazon EC2 Systems Manager (ssm) Parameter Store as environment variables. This is useful for injecting secure secrets into the environment of a docker container process.

### Usage
Create secret parameters on AWS Parameter Store for your application using [hierarchies](http://docs.aws.amazon.com/systems-manager/latest/userguide/sysman-paramstore-working.html#sysman-paramstore-su-organize):
```
> aws ssm put-parameter --name /userservice/SECRET_1 --value "123456" --type SecureString
> aws ssm put-parameter --name /accountservice/secret_2 --value "abcdef" --type SecureString
> aws ssm put-parameter --name /accountservice/secret_3 --value "foobarbaz" --type SecureString
> aws ssm put-parameter --name /database/production/password --value "productionpass" --type SecureString
> aws ssm put-parameter --name /database/staging/password --value "stagingpass" --type SecureString
```
Add tags to categorize parameters in various ways:
```
> aws ssm add-tags-to-resource --resource-type Parameter --resource-id /userservice/SECRET_1 --tags Key=userservice,Value=true Key=production,Value=true
> aws ssm add-tags-to-resource --resource-type Parameter --resource-id /accountservice/secret_2 --tags Key=accountservice,Value=true Key=production,Value=true
> aws ssm add-tags-to-resource --resource-type Parameter --resource-id /accountservice/secret_3 --tags Key=accountservice,Value=true Key=staging,Value=true
> aws ssm add-tags-to-resource --resource-type Parameter --resource-id /database/production/password --tags Key=userservice,Value=true Key=accountservice,Value=true Key=production,Value=true
> aws ssm add-tags-to-resource --resource-type Parameter --resource-id /database/staging/password --tags Key=userservice,Value=true Key=accountservice,Value=true Key=staging,Value=true
```
Retrieve parameters with `aws-ssm-env`:
```
# filter by 'userservice' parameters for 'production'
> AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=userservice,production
SECRET_1=123456
PASSWORD=productionpass
# filter by 'accountservice' parameters for 'staging'
> AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=accountservice,staging
SECRET_3=foobarbaz
PASSWORD=stagingpass
# filter by path (`/` will search all parameters)
> AWS_REGION=<aws-region> aws-ssm-env --paths=/database
PASSWORD=productionpass
PASSWORD=stagingpass
# filter by path and tag
> AWS_REGION=<aws-region> aws-ssm-env --paths=/database --tags=production
PASSWORD=productionpass
```
*Notice that parameter names are automatically capitalized.*


Use `export` with `aws-ssm-env` to inject secrets from Parameter Store into the environment:
```
> export $(AWS_REGION=<aws-region> aws-ssm-env --paths=/ --tags=userservice,production)
> env
...
...
SECRET_1=123456
PASSWORD=productionpass
```

### Installation
go get:
```
> go get github.com/jamietsao/aws-ssm-env
```
Or download [binary](https://github.com/jamietsao/aws-ssm-env/releases/latest):
```
> wget -O aws-ssm-env.zip https://github.com/jamietsao/aws-ssm-env/releases/download/v1.0.0/aws-ssm-env-v1.0.0-linux-amd64.zip
> unzip aws-ssm-env.zip
> chmod 755 aws-ssm-env
```

### Author
Jamie Tsao

### License
MIT

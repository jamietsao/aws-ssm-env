# aws-ssm-env
Simple utility to print parameters from Amazon Systems Manager (ssm) Parameter Store as environment variables. This is useful for injecting secure secrets into the environment of a docker container process.

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

#### Retrieve and print parameters with `aws-ssm-env`:

*Notice that parameter names are automatically capitalized.*
```
# filter by tag: print 'userservice' parameters for 'production'
> AWS_REGION=<aws-region> aws-ssm-env --tags=userservice,production
SECRET_1=123456
PASSWORD=productionpass

# filter by tag: print 'accountservice' parameters for 'staging'
> AWS_REGION=<aws-region> aws-ssm-env --tags=accountservice,staging
SECRET_3=foobarbaz
PASSWORD=stagingpass

# filter by path: print all database parameters
> AWS_REGION=<aws-region> aws-ssm-env --paths=/database
PASSWORD=productionpass
PASSWORD=stagingpass

# filter by path and tag: print database parameters for 'production'
> AWS_REGION=<aws-region> aws-ssm-env --paths=/database --tags=production
PASSWORD=productionpass
```

**WARNING: Using `'/'` as a path (e.g. `--paths=/`) will recursively retrieve EVERY single parameter configured in Parameter Store.  This will increase the runtime of this script and could result in hitting SSM rate limits. Use of `'/'` as a path is highly discouraged.**


Use `export` with `aws-ssm-env` to inject secrets from Parameter Store into the environment:
```
> export $(AWS_REGION=<aws-region> aws-ssm-env --tags=userservice,production)
> env
...
...
SECRET_1=123456
PASSWORD=productionpass
```

#### Setting parameters into environment from application code
If you have a need to set SSM parameters as environment variables directly from application code:
```
// retrieve database parameters for production and set each as env variables
fetch.MustSetEnv([]string{"/database"}, []string{"production"}, true)
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

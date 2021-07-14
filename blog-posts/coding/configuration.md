# Configuration

## Motivation
In our previous team we worked on various environments, such as a `local`, a `dev` and `production` environment.
Obviously, we therefore have the need to configure certain environment specific parameters such as our database connection. 
To do so we use `.json` files to configure for environments. 
I know `.toml` is very popular among the Go community for its interoperability with Go, its verboseness, among other reasons. 
We chose `.json` over `.toml` since not all of our teams use Golang as their primary programming language.


## Usage 

### Default config

For parameters which are common across all environments we define a `default.go` 

```go
var defaultConfig = Config{
	Name: "Golang Service",
}
```

### Environment specific config
For environment specific parameters we have one config file per environment e.g. `dev.json`. Secrets such as the `{{ .db_name }}` will get injected during deployment. 

```json
{
  "environment": "dev",
  "persistence": {
    "dbName": "{{ .db_name }}",
    "dbHost": "go-postgres",
    "dbPort": 5432,
    "dbUsername": "{{ .db_user }}",
    "dbPassword": "{{ .db_password }}",
    "sslEnabled": false
  },
  "featureToggles": {
    "enableNewFeature": true
  }
}
```

### Local development

For environment specific parameters we have one config file per environment e.g. `local.json`

```json
{
  "environment": "local",
  "persistence": {
    "dbName": "golangservice",
    "dbHost": "go-postgres",
    "dbPort": 5432,
    "dbUsername": "postgres",
    "dbPassword": "",
    "sslEnabled": false
  }
}
```

### Injecting secrets
No matter if we're talking about env specific or local config, we can agree that we don't want to hardcode any secrets in our repository. 
In real life, we've got a `k8s` cluster setup, to learn more about this checkout our 

#### Injecting environment secrets
Secrets are injected during deployment in the following steps: 
1. fetch secrets from vault
2. run with given parameters per environment in your pipeline
 ```
 ./do deploy "$STAGE" "$VERSION" "$DBNAME" "$DBUSER" "$DBPASSWORD"
 ```

Eventually this `do-`script uses helm to inject the secrets into the configuration. 

- [ ] more info about helm 

## Injecting secrets for local development
A trick we do when developing locally is that we fetch all necessary secrets and inject them in our `local-temp.json` 
(of course this file is in our `.gitignore`!!). **ADVICE** use talisman to avoid pushing secrets by mistake

To test this out do 
```shell
./do generate-local-config
```

## Code

### How the config service is implemented

The way we implemented our config service we require 3 bits of information:
1. *configPath* - path to our configuration file, passed in as argument to our binary
2. *secretsDirectoryPath* - path to our secrets directory, passed in as argument to our binary
3. *secretsEnv* - environment variable 

If you look at our `Dockerfile` you will see how our `webservice` binary is started, there we pass in our config file as well as the path to our secrets. 

```
CMD ./webservice "/service/config/config.json" "/service/secrets"
```

- [ ] continue explanation

## Further Resources
- [ ] add resources
# Configuration

## Motivation
In our previous team we worked on various environments, such as a `local`, a `dev` and `production` environment.
Obviously, we therefore have the need to configure certain environment specific parameters such as our database connection. 
To do so we use `.json` files to configure for environments. 
I know `.toml` is very popular among the Go community for its interoperability with Go, its verboseness, among other reasons. 
We chose `.json` over `.toml` since not all of our teams use Golang as their primary programming language.


## Usage 

### Using the config

For parameters which are common across all environments we define a `default.go` 

```go
var defaultConfig = Config{
	Name: "Golang Service",
}

```

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

### Adding secrets

## Injecting environment secrets
- [ ] create a vault where we can fetch secrets from ?
- [ ] add text for secrets management

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

- [ ] continue explanation

## Further Resources
- [ ] add resources
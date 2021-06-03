# Integrating a Postgres database

If your service is trying to store information, you'll want to establish a connection to a database to persist this information over time. 

Make sure that you encapsulate all database logic in one place. 
This will allow you to move to a different database driver or another database management system altogether.

In our service we decided to go for a postgres database with the SQLBoiler ORM (https://github.com/volatiletech/sqlboiler). 
What makes it so special is the approach of first building the database structure and then generating the models from that. 
This ensures that the types defined in the Golang code match those defined in the database. 

For database migrations we decided to use Golang Migrate (https://github.com/golang-migrate/migrate).

## Starting up a postgres locally

In order to run our service against a local database, we start up a postgres docker container using the `./do.sh` script:
```
## run-db : start local postgres database
function task_run_db {
    green "Removing old database"
    docker rm -f local-db
    green "Pulling image"
    docker pull postgres:11

    docker run -d \
        -e POSTGRES_DB="postgres" \
        -e POSTGRES_USER="my-user" \
        -e POSTGRES_PASSWORD="my-users-password" \
        -e POSTGRES_HOST_AUTH_METHOD="trust" \
        --name="local-db" \
        -p 5432:5432 \
        -m 128m \
        postgres:11
    green "wait 2 seconds until database is up"
    sleep 2
}
```
You should now see this database in `docker ps` and you should be able to connect to it using 
```
psql "host=localhost port=5432 user=my-user dbname=postgres password=my-users-password sslmode=disable"
```
## Connecting to the postgres

In order to connect to the database, we have created a `postgres.go` file in `/persistence`, which takes the configuration and returns a driver.
```
func ConnectPostgres(config config.Postgres) (*sql.DB, error) {
	pgOptions := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", config.Host, config.Port, config.UserName, config.DBName, config.Password)
	if !config.SSLEnabled {
		pgOptions = pgOptions + " sslmode=disable"
	}

	db, err := sql.Open("postgres", pgOptions)
	if err != nil {
		return nil, err
	}

	log.Infof("PostgreSQL storage: connected to host %s:%s database %s with user %s", config.Host, config.Port, config.DBName, config.UserName)
	return db, nil
}
```

## Migration scripts
As a first step we define our initial database schema. To do so you can run: 
```
migrate create -ext sql -dir migrations -seq initial_setup  
```
This will create an `.up.sql` and `.down.sql` file in the `./migrations` folder. 
You can then use the `.up.sql` file to create the tables, and use the `.down.sql` to drop the tables again.

An example for a `.up.sql` file would be 
```
CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE,
   phone_number VARCHAR (20) UNIQUE
);
```
and then the `.down.sql` file would look like this
```
DROP TABLE IF EXISTS users;
```

Those can be run in the Command line, however in our example we have created another `main.go` in `./cmd/migration`. 
This takes all connection information as input parameters, then connects to the database and processes the up migrations defined in `./migrations`

```
func main() {
	initLogging()

	if len(os.Args) != 7 {
		log.Fatalf("the migrations command expects 6 command line parameters: dbHost dbPort dbName sslEnabled dbUsername dbPassword")
	}

	port, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("the migrations command expects the port to be a number: %s", os.Args[2])
	}

	pgConfig := config.Postgres{
		Host:       os.Args[1],
		Port:       int(port),
		DBName:     os.Args[3],
		SSLEnabled: os.Args[4] == "true",
		UserName:   os.Args[5],
		Password:   os.Args[6],
	}

	db, err := persistence.ConnectPostgres(pgConfig)
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: pgConfig.DBName,
	})
	if err != nil {
		log.Fatalf("could not create driver for DB migrations: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		pgConfig.DBName, driver)
	if err != nil {
		log.Fatalf("could not init DB migrations: %s", err.Error())
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("could not migrate the DB: %s", err.Error())
		} else {
			log.Info("No database changes necessary")
		}
	}
	log.Info("PostgreSQL storage: migrations finished")
}

func initLogging() {
	cLog := console.New(true)
	log.AddHandler(cLog, log.AllLevels...)
}
```

## Generating models using SQL Boiler
Now that we have tables in the database, we can generate models from those.
This is done using `sqlboiler`. First you need to specify the configuration in `sqlboiler.toml` the root directory:
```
pkgname  = "models"
output   = "persistence/models"
wipe     = true
no-tests = true

[psql]
  dbname = "postgres"
  host   = "localhost"
  port   = 5432
  user   = "my-user"
  pass   = "my-users-password"
  sslmode = "disable"
  blacklist = ["schema_migrations"]
```
and then the files can be generated using `sqlboiler psql` and they will be 
found in the directory specified in `output`.

Now you can use them in your `postgres.go` in the following way
```
func (p *Postgres) Add(ctx context.Context, user User) error {
	newUser := &models.User{
		Username:    user.UserName,
		Email:       null.StringFrom(user.Email),
		PhoneNumber: null.StringFrom(user.Telephone),
	}
	return newUser.Insert(ctx, p.db, boil.Infer())
}
```

## Further reading
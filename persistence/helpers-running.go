package persistence

import (
	"context"
	"database/sql"
	"time"
)

//FIXME: we ran into a bug where after 100 integration tests the database connection suddenly dropped bc it reached its max
// there should be the possibility to not open a connection if there's already one established
func EnsureConnected(ctx context.Context, connectionString string, pollingDelay time.Duration) error {
	done := make(chan error, 1)

	go func() {
		for {
			db, err := sql.Open("postgres", connectionString)
			if err == nil {
				err := db.Ping()
				if err == nil {
					done <- nil
				}
			}

			time.Sleep(pollingDelay)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

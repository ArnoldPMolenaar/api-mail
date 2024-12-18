package cache

import (
	"github.com/ArnoldPMolenaar/api-utils/cache"
	"github.com/valkey-io/valkey-go"
)

var Valkey valkey.Client

// OpenValkeyConnection Start a new valkey connection.
func OpenValkeyConnection() error {
	// Open connection to valkey.
	client, err := cache.ValkeyConnection()
	if err != nil {
		return err
	}

	// Set the global Valkey variable.
	Valkey = client

	return nil
}

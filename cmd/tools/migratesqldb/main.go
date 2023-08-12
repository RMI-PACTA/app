// Command migrate provides functionality for managing migration sets applied
// to a database, using the golang-migrate library. It currently takes
// configuration via a sops-encrypted file.
package main

import (
	"github.com/RMI/pacta/cmd/tools/migratesqldb/cmd"
)

func main() {
	cmd.Execute()
}

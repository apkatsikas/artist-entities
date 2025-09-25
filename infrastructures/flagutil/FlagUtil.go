package flagutil

import (
	"flag"
	"sync"
)

type FlagUtil struct {
	MigrateDB       bool
	MigrateUser     string
	MigratePassword string
	Secret          string
}

func (fu *FlagUtil) Setup() {
	flag.BoolVar(&fu.MigrateDB, "migrateDB", false, "Migrate the database")
	flag.StringVar(&fu.MigrateUser, "migrateUser", "", "User name to migrate")
	flag.StringVar(&fu.MigratePassword, "migratePassword", "", "Password for user to migrate")
	flag.StringVar(&fu.Secret, "secret", "", "Secret auth value")
	flag.Parse()
}

// Setup singleton
var (
	fu     *FlagUtil
	fuOnce sync.Once
)

func Get() *FlagUtil {
	if fu == nil {
		fuOnce.Do(func() {
			fu = &FlagUtil{}
		})
	}
	return fu
}

package infrastructures

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type SQLiteHandler struct {
    conn *gorm.DB
}

func (handler *SQLiteHandler) Connection() *gorm.DB {
    return handler.conn
}

func (handler *SQLiteHandler) ConnectSQLite(dsn string) error {
    db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
        // These logs do not go to a log file
        // Could be turned off on the local console with this line
        //Logger: logger.Default.LogMode(logger.Silent)
    })
    if err != nil {
        return err
    }
    handler.conn = db
    return nil
}

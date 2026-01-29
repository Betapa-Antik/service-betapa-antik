package utils

import "gorm.io/gorm"

// RunInTransaction starts a transaction, runs the provided function and commits/rolls back accordingly.
func RunInTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
    tx := db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    if err := fn(tx); err != nil {
        _ = tx.Rollback().Error
        return err
    }

    if err := tx.Commit().Error; err != nil {
        return err
    }
    return nil
}

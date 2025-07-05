package database

import (
	"gorm.io/gorm"
)

// UserIDContextKey is the key used to store the user ID in GORM's context.
const UserIDContextKey = "user_id"

// UpdateByCallback is a GORM plugin to automatically set CreatedBy, UpdatedBy, DeletedBy.
type UpdateByCallback struct{}

// Name returns the name of the plugin.
func (p *UpdateByCallback) Name() string {
	return "update_by_callback"
}

// Initialize registers the callbacks.
func (p *UpdateByCallback) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:before_create").Register("set_created_by", p.setCreatedBy)
	db.Callback().Update().Before("gorm:before_update").Register("set_updated_by", p.setUpdatedBy)
	db.Callback().Delete().Before("gorm:before_delete").Register("set_deleted_by", p.setDeletedBy)
	return nil
}

// setCreatedBy sets the CreatedBy field before creating a record.
func (p *UpdateByCallback) setCreatedBy(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if createdByField := db.Statement.Schema.LookUpField("CreatedBy"); createdByField != nil {
			if userID, ok := db.Statement.Context.Value(UserIDContextKey).(uint); ok && userID != 0 {
				if db.Statement.Changed(createdByField.Name) {
					return // Field already set, skip
				}
				db.Statement.SetColumn(createdByField.Name, userID)
			}
		}
	}
}

// setUpdatedBy sets the UpdatedBy field before updating a record.
func (p *UpdateByCallback) setUpdatedBy(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if updatedByField := db.Statement.Schema.LookUpField("UpdatedBy"); updatedByField != nil {
			if userID, ok := db.Statement.Context.Value(UserIDContextKey).(uint); ok && userID != 0 {
				if db.Statement.Changed(updatedByField.Name) {
					return // Field already set, skip
				}
				db.Statement.SetColumn(updatedByField.Name, userID)
			}
		}
	}
}

// setDeletedBy sets the DeletedBy field before deleting a record (soft delete).
func (p *UpdateByCallback) setDeletedBy(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if deletedByField := db.Statement.Schema.LookUpField("DeletedBy"); deletedByField != nil {
			if userID, ok := db.Statement.Context.Value(UserIDContextKey).(uint); ok && userID != 0 {
				if db.Statement.Changed(deletedByField.Name) {
					return // Field already set, skip
				}
				db.Statement.SetColumn(deletedByField.Name, userID)
			}
		}
	}
}

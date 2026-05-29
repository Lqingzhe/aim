package newerror

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func IsMysqlError(result *gorm.DB) error {
	err := result.Error
	if err != nil {
		if isContext, err2 := IsContextError(err); isContext {
			return err2
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return MakeError(http.StatusConflict, CodeResourceDuplicate, "The Key Already Exist", fmt.Errorf(`The Key Already Set, Should Use "Update"`), LevelWarn)
		}
		return MakeError(http.StatusInternalServerError, CodeDatabaseError, "Database Error", err, LevelError)
	}
	return nil
}

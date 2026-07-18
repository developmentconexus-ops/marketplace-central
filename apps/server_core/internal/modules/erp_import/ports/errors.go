package ports

import (
	"errors"
	"fmt"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

var ErrImportInProgress = errors.New("erp import in progress")

var ErrImportNotFound = errors.New("erp import not found")

type DuplicateFileError struct {
	ExistingID       domain.ImportID
	ExistingProtocol domain.Protocol
}

func (e *DuplicateFileError) Error() string {
	return fmt.Sprintf("duplicate ERP import file: existing import %s (%s)", e.ExistingID, e.ExistingProtocol)
}

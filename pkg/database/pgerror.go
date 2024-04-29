package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	// IntegrityConstraintViolationErrorCode when a value is missing or not associated with a valid foreign table.
	IntegrityConstraintViolationErrorCode = "23000"
	// NotNullViolationErrorCode when received an empty or null field.
	NotNullViolationErrorCode = "23502"
	// ForeignKeyViolationErrorCode FK constraint was created with a not valid value.
	ForeignKeyViolationErrorCode = "23503"
	// UniqueViolationErrorCode duplicate key value violates unique constraint.
	UniqueViolationErrorCode = "23505"
	// CheckViolationErrorCode when a value violates a specific requirement in a column.
	CheckViolationErrorCode = "23514"
	// InvalidTextRepresentation when a field is uuid but receive a different pattern.
	InvalidTextRepresentation = "22P02"
)

var (
	ErrStatementTimeout             = errors.New("statement timeout")
	ErrIntegrityConstraintViolation = errors.New("integrity constraint violation")
	ErrNotNullViolation             = errors.New("not null violation")
	ErrForeignKeyViolation          = errors.New("foreign key violation")
	ErrCheckViolation               = errors.New("check violation")
	ErrUniqueViolation              = errors.New("unique violation")
	ErrInternal                     = errors.New("internal server error")
	ErrNotFound                     = errors.New("it was not found")
	ErrInvalidTextRepresentation    = errors.New("invalid text representation")
	ErrNoRowsAffected               = errors.New("no rows affected")
)

func getPGError(err error) pgdriver.Error {
	var pgErr pgdriver.Error
	errors.As(err, &pgErr)
	return pgErr
}

func getPGErrorCode(err error) string {
	return getPGError(err).Field('C')
}

func GetPGErrorConstraint(err error) string {
	return getPGError(err).Field('n')
}

func GetPGErrorColumn(err error) string {
	return getPGError(err).Field('c')
}

func IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrNotFound) || err.Error() == ErrNotFound.Error()
}

func HandleAdapterError(err error) (bool, error) {
	var pgErr pgdriver.Error
	if errors.As(err, &pgErr) {
		if pgErr.StatementTimeout() {
			return true, fmt.Errorf("%w: %s", ErrStatementTimeout, pgErr.Error())
		}

		switch getPGErrorCode(err) {
		case InvalidTextRepresentation:
			return true, fmt.Errorf("%w: %s", ErrInvalidTextRepresentation, pgErr.Error())
		case IntegrityConstraintViolationErrorCode:
			return true, fmt.Errorf("%w: %s", ErrIntegrityConstraintViolation, pgErr.Error())
		case NotNullViolationErrorCode:
			return true, fmt.Errorf("%w: %s", ErrNotNullViolation, pgErr.Error())
		case ForeignKeyViolationErrorCode:
			return true, fmt.Errorf("%w: %s", ErrForeignKeyViolation, pgErr.Error())
		case CheckViolationErrorCode:
			return true, fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Error())
		case UniqueViolationErrorCode:
			return true, fmt.Errorf("%w: %s", ErrUniqueViolation, pgErr.Error())
		default:
			return true, pgErr
		}
	}

	if IsErrNoRows(err) {
		return true, fmt.Errorf("%w", ErrNotFound)
	}

	return false, fmt.Errorf("%w", ErrInternal)
}

func HandleResult(res sql.Result) error {
	affected, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		return rowsErr
	}

	if affected == int64(0) {
		return ErrNoRowsAffected
	}

	return nil
}

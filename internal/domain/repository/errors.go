package repository

type repositoryError string

const (
	ErrRecordNotFound  repositoryError = "RECORD_NOT_FOUND"
	ErrDuplicatedKey   repositoryError = "KEY_DUPLICATED"
	ErrSessionNotFound repositoryError = "SESSION_NOT_FOUND"
	ErrUndefined       repositoryError = "UNDEFINED_REPO_ERROR"
)

func (e repositoryError) Error() string {
	return string(e)
}

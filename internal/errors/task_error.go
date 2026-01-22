package taskerrors

type ErrorType string

const (
	ErrRetryable    ErrorType = "retryable"
	ErrNonRetryable ErrorType = "non_retryable"
)

type TaskError struct {
	Type ErrorType
	Err  error
}

func (t *TaskError) Error() string {
	return t.Err.Error()
}

func Retryable(err error) *TaskError {
	return &TaskError{Type: ErrRetryable, Err: err}
}

func NonRetryable(err error) *TaskError {
	return &TaskError{Type: ErrNonRetryable, Err: err}

}

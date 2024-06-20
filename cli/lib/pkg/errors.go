package cliadpt

import "fmt"

type (
	ErrInvalidConfig struct {
		error
	}

	ErrNilConfig struct{}

	ErrUseCaseNotFound struct {
		Args []string
		UseCaseName
	}
)

func NewInvalidConfigError(err error) error { return ErrInvalidConfig{error: err} }

func (e ErrNilConfig) Error() string { return "flagset.Config is nil" }

func (e ErrUseCaseNotFound) Error() string {
	return fmt.Sprintf("flagset.UseCase: [%s] not found to run args=%v", e.UseCaseName, e.Args)
}

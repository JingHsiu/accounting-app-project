package common

type Output interface {
	GetID() string
	GetExitCode() ExitCode
	GetMessage() string
}

type UseCaseOutput struct {
	ID       string
	ExitCode ExitCode
	Message  string
}

func (o UseCaseOutput) GetID() string         { return o.ID }
func (o UseCaseOutput) GetExitCode() ExitCode { return o.ExitCode }
func (o UseCaseOutput) GetMessage() string    { return o.Message }

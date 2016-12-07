package balancer

// InvalidWorkError is the error returned when a node receives invalid work.
type InvalidWorkError struct {
	Str string
}

func (err *InvalidWorkError) Error() string {
	return err.Str
}

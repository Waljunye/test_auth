package stores

func NewErrNoRowsWasAffected(msg string) ErrNoRowsWasAffected {
	return ErrNoRowsWasAffected{msg}
}

type ErrNoRowsWasAffected struct {
	message string
}

func (e ErrNoRowsWasAffected) Error() string {
	return "store error:" + " no rows was affected: " + e.message
}

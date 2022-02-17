package shterrors

var _ error = (*ErrorConflictSaveURL)(nil)

type ErrorConflictSaveURL struct {
	Err           error
	ExistShortURL string
}

func (*ErrorConflictSaveURL) Error() string {
	return "Сохраняемая ссылка уже существует"
}

func (e *ErrorConflictSaveURL) Is(tgt error) bool {
	_, ok := tgt.(*ErrorConflictSaveURL)
	return ok
}

package pss

// Error - общий тип константных ошибок пакета.
type Error string

// Error удовлетворяет интерфейсу error.
func (err Error) Error() string {
	return string(err)
}

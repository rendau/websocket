package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ServiceNA      = Err("service_not_available")
	NotAuthorized  = Err("not_authorized")
	ObjectNotFound = Err("object_not_found")
)

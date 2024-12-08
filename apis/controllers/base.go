package controllers

// ControllerBase define
type ControllerBase struct{}

func isError(val any) bool {
	if _, ok := val.(error); ok {
		return true
	} else {
		return false
	}
}

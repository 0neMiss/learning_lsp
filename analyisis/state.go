package analysis

type State struct {
	// Map of  file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}

}

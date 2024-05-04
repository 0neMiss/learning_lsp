package analysis

//If this were a real LSP for a real language you would talk with compliers to get information like errors, datatypes, keywords ect..
type State struct {
	// Map of  file names to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

package resolver

type stack struct {
	scopes []map[string]bool
}

func newStack() *stack {
	return &stack{make([]map[string]bool, 0)}
}

func (s *stack) size() int {
	return len(s.scopes)
}

func (s *stack) isEmpty() bool {
	return s.size() == 0
}

func (s *stack) get(i int) map[string]bool {
	return s.scopes[i]
}

func (s *stack) push() {
	s.scopes = append(s.scopes, make(map[string]bool))
}

func (s *stack) pop() {
	s.scopes = s.scopes[0 : s.size()-1]
}

func (s *stack) peek() map[string]bool {
	return s.get(s.size() - 1)
}

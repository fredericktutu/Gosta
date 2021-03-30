package gosta

type IntStack struct{
	stack []int
}


type SIDStack struct {
	stack []SID
}

type ActionStack struct {
	stack []Action
}

// begin push
func (s *IntStack) Push(i int) {
	s.stack = append(s.stack, i)
}

func (s *ActionStack) Push(i Action) {
	s.stack = append(s.stack, i)
}

func (s *SIDStack) Push(i SID) {
	s.stack = append(s.stack, i)
}
// end push


//begin pop
func (s *IntStack) Pop() {
	length := len(s.stack)
	s.stack = s.stack[:length - 1]
}

func (s *ActionStack) Pop() {
	length := len(s.stack)
	s.stack = s.stack[:length - 1]
}

func (s *SIDStack) Pop() {
	length := len(s.stack)
	s.stack = s.stack[:length - 1]
}
//end pop

//begin top
func (s *IntStack) Top() int {
	length := len(s.stack)
	return s.stack[length - 1]
}

func (s *ActionStack) Top() Action {
	length := len(s.stack)
	return s.stack[length - 1]
}

func (s *SIDStack) Top() SID {
	length := len(s.stack)
	return s.stack[length - 1]
}
//end top

//begin size
func (s *IntStack) Size() int {
	return len(s.stack)
}
func (s *ActionStack) Size() int {
	return len(s.stack)
}
func (s *SIDStack) Size() int {
	return len(s.stack)
}
//end size

//begin empty
func (s *IntStack) Empty() bool {
	return len(s.stack) == 0
}
func (s *ActionStack) Empty() bool {
	return len(s.stack) == 0
}
func (s *SIDStack) Empty() bool {
	return len(s.stack) == 0
}
//end empty

//begin get
func (s *IntStack) Get(n int) int{
	return s.stack[n]
}
func (s *ActionStack) Get(n int) Action{
	return s.stack[n]
}
func (s *SIDStack) Get(n int) SID{
	return s.stack[n]
}
//end get

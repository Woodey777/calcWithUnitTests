package commonStructs

import "errors"

type Stack []string

func (s *Stack) Push(value string) {
	*s = append(*s, value)
}

func (s *Stack) Pop() (string, error) {
	if len(*s) == 0 {
		return "", errors.New("stack is empty")
	}

	last_ind := len(*s) - 1
	value := (*s)[last_ind]
	*s = (*s)[:last_ind]
	return value, nil
}

func (s *Stack) Peek() (string, error) {
	if len(*s) == 0 {
		return "", errors.New("stack is empty")
	}
	return (*s)[len(*s)-1], nil
}

type FloatStack []float64

func (s *FloatStack) Push(value float64) {
	*s = append(*s, value)
}

func (s *FloatStack) Pop() (float64, error) {
	if len(*s) == 0 {
		return 0, errors.New("stack is empty")
	}

	last_ind := len(*s) - 1
	value := (*s)[last_ind]
	*s = (*s)[:last_ind]
	return value, nil
}

func (s *FloatStack) Peek() (float64, error) {
	if len(*s) == 0 {
		return 0, errors.New("stack is empty")
	}
	return (*s)[len(*s)-1], nil
}

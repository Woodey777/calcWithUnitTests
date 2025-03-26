package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Приоритет операторов
var precedence = map[string]int{
	"(": 0,
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"~": 4,
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	stack := make(Stack, 0, len(tokens))

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
		} else if tokenPrec, ok := precedence[token]; ok {
			stackTop, _ := stack.Peek()
			for len(stack) != 0 && precedence[stackTop] >= tokenPrec {
				output = append(output, stackTop)
				stack.Pop()
				stackTop, _ = stack.Peek()
			}
			stack.Push(token)
		} else {
			return nil, fmt.Errorf("invalid operation or number: %s", token)
		}
	}

	for len(stack) != 0 {
		stackTop, _ := stack.Pop()
		output = append(output, stackTop)
	}

	return output, nil
}

func evaluatePostfix(tokens []string) (float64, error) {
	numsStack := make(FloatStack, 0, len(tokens))

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			numsStack.Push(num)
		} else {
			if token == "~" {
				num, err := numsStack.Pop()
				if err != nil {
					return 0, fmt.Errorf("not enough operands for unary minus")
				}
				numsStack.Push(-num)
				continue
			}

			if len(numsStack) < 2 {
				return 0, fmt.Errorf("not enough operands for operation %s", token)
			}

			operand2, _ := numsStack.Pop()
			operand1, _ := numsStack.Pop()

			var result float64
			switch token {
			case "+":
				result, err = Add(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating addition: %v", err)
				}
			case "-":
				result, err = Sub(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating substraction: %v", err)
				}
			case "*":
				result, err = Mul(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating multiplication: %v", err)
				}
			case "/":
				if operand2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result, err = Div(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating division: %v", err)
				}
			default:
				panic("unknown operator")
			}
			numsStack.Push(result)
		}
	}

	if len(numsStack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	res, _ := numsStack.Pop()
	return res, nil
}

func tokenize(input string) []string {
	tokens := make([]string, 0, len(input))
	var currentToken strings.Builder

	for _, char := range input {
		if unicode.IsDigit(char) || char == '.' {
			currentToken.WriteRune(char)
			continue
		}

		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		}

		if !unicode.IsSpace(char) {
			if char == '-' {
				if len(tokens) == 0 {
					tokens = append(tokens, string('~'))
					continue
				} else {
					lasToken := []rune(tokens[len(tokens)-1])
					if lasToken[len(lasToken)-1] != '.' && !unicode.IsDigit(lasToken[len(lasToken)-1]) {
						tokens = append(tokens, string('~'))
						continue
					}
				}
			}

			tokens = append(tokens, string(char))
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

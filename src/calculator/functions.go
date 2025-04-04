package calculator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	structs "calcWithTests/src/commonstructs"
)

// Приоритет операторов
var precedence = map[string]int{
	"(": 0,
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"~": 3,
	"^": 4,
	"m": 4,
	"d": 4,
}

func Calculate(expression string) (float64, error) {
	tokens := tokenize(expression)
	fmt.Println(tokens)
	postfix, err := infixToPostfix(tokens)
	fmt.Println(postfix)
	if err != nil {
		return 0, fmt.Errorf("error while parsing: %v", err)
	}

	result, err := evaluatePostfix(postfix)
	if err != nil {
		return 0, fmt.Errorf("error while calculating: %v", err)
	}

	return result, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	stack := make(structs.Stack, 0, len(tokens))

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
		} else if token == "(" {
			stack.Push(token)
		} else if token == ")" {
			for len(stack) != 0 {
				stackTop, _ := stack.Pop()
				if stackTop == "(" {
					break
				}
				output = append(output, stackTop)
			}
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
	numsStack := make(structs.FloatStack, 0, len(tokens))

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
			case "*", "m":
				result, err = Mul(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating multiplication: %v", err)
				}
			case "/", "d":
				if operand2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result, err = Div(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating division: %v", err)
				}
			case "^":
				result, err = Pow(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating power: %v", err)
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

	runes := []rune(input)

	for i := 0; i < len(input); i++ {
		if unicode.IsDigit(runes[i]) || runes[i] == '.' {
			currentToken.WriteRune(runes[i])
			continue
		}

		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
			currentToken.Reset()
		}

		if runes[i] == 'e' && i+1 < len(input) && (runes[i+1] == '+' || runes[i+1] == '-') {
			eInd := i
			eTokens := make([]string, 0, 2)
			i++

			switch runes[i] {
			case '+':
				eTokens = append(eTokens, string('m'))
			case '-':
				eTokens = append(eTokens, string('d'))
			}
			i++

			for i < len(input) {
				if unicode.IsDigit(runes[i]) || runes[i] == '.' {
					currentToken.WriteRune(runes[i])
				} else {
					break
				}
				i++
			}
			i--

			powStr := currentToken.String()
			if currentToken.Len() > 0 {
				currentToken.Reset()
				pow, _ := strconv.ParseFloat(powStr, 64)
				num := math.Pow(10, pow)
				eTokens = append(eTokens, strconv.FormatFloat(num, 'f', -1, 64))
				tokens = append(tokens, eTokens...)
			} else {
				tokens = append(tokens, "e")
				i = eInd
			}
			continue
		}

		if !unicode.IsSpace(runes[i]) {
			if runes[i] == '-' {
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

			tokens = append(tokens, string(runes[i]))
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

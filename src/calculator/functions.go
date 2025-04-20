package calculator

import (
	structs "calcWithTests/src/commonStructs"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type CalculatorConfig struct {
	AngleUnits string
}

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
	"q": 4, // sqrt
	"l": 4, // ln
	"x": 4, // exp
	"s": 4, // sin
	"c": 4, // cos
	"t": 4, // tg
	"g": 4, // ctg
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Calculate(expression string, config CalculatorConfig) (float64, error) {
	tokens := tokenize(expression)
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, fmt.Errorf("error while parsing: %v", err)
	}

	result, err := evaluatePostfix(postfix, config.AngleUnits)
	if err != nil {
		return 0, fmt.Errorf("error while calculating: %v", err)
	}

	return result, nil
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	stack := make(structs.Stack, 0, len(tokens))

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil || token == "e" || token == "pi" {
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

func evaluatePostfix(tokens []string, angleUnits string) (float64, error) {
	numsStack := make(structs.FloatStack, 0, len(tokens))

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			numsStack.Push(num)
			continue
		}

		if token == "e" {
			numsStack.Push(math.E)
			continue
		}
		if token == "pi" {
			numsStack.Push(math.Pi)
			continue
		}

		var result float64
		var err error

		switch token {
		case "+", "-", "*", "/", "^", "m", "d":
			if len(numsStack) < 2 {
				return 0, fmt.Errorf("not enough operands for operation %s", token)
			}
			operand2, _ := numsStack.Pop()
			operand1, _ := numsStack.Pop()

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
				result, err = Div(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating division: %v", err)
				}
			case "^":
				result, err = Pow(operand1, operand2)
				if err != nil {
					return 0, fmt.Errorf("calculating power: %v", err)
				}
			}
		case "q", "l", "x", "s", "c", "t", "g", "~":
			if len(numsStack) < 1 {
				return 0, fmt.Errorf("not enough operands for operation %s", token)
			}
			operand, _ := numsStack.Pop()

			switch token {
			case "~":
				result = -operand
			case "q":
				result, err = Sqrt(operand)
				if err != nil {
					return 0, fmt.Errorf("calculating sqrt: %v", err)
				}
			case "l":
				result, err = Ln(operand)
				if err != nil {
					return 0, fmt.Errorf("calculating ln: %v", err)
				}
			case "x":
				result, err = Exp(operand)
				if err != nil {
					return 0, fmt.Errorf("calculating exp: %v", err)
				}
			case "s", "c", "t", "g":
				if angleUnits == "degree" {
					operand = degreesToRadians(operand)
				}
				switch token {
				case "s":
					result = Sin(operand)
				case "c":
					result = Cos(operand)
				case "t":
					result, err = Tg(operand)
					if err != nil {
						return 0, fmt.Errorf("calculating tg: %v", err)
					}
				case "g":
					result, err = Cot(operand)
					if err != nil {
						return 0, fmt.Errorf("calculating ctg: %v", err)
					}
				}
			}
		default:
			panic("unknown operator")
		}

		numsStack.Push(result)
	}

	if len(numsStack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	res, _ := numsStack.Pop()
	return res, nil
}

func tokenize(input string) []string {
	tokens := make([]string, 0, len(input))
	var currNumToken strings.Builder
	var currLetToken strings.Builder

	runes := []rune(input)

	for i := 0; i < len(input); i++ {
		if unicode.IsDigit(runes[i]) || runes[i] == '.' {
			if currLetToken.Len() > 0 {
				tokens = append(tokens, getLetToken(&currLetToken))
				currLetToken.Reset()
			}

			currNumToken.WriteRune(runes[i])
		} else if strings.Contains("(+-*/^)", string(runes[i])) {

			if currNumToken.Len() > 0 {
				tokens = append(tokens, currNumToken.String())
				currNumToken.Reset()
			}

			if currLetToken.Len() > 0 {
				tokens = append(tokens, getLetToken(&currLetToken))
				currLetToken.Reset()
			}

			if runes[i] == '-' {
				if len(tokens) == 0 {
					tokens = append(tokens, string('~'))
					continue
				}

				if lastToken := []rune(tokens[len(tokens)-1]); (lastToken[len(lastToken)-1] != '.' && !unicode.IsDigit(lastToken[len(lastToken)-1])) && !reflect.DeepEqual(lastToken, []rune("e")) {
					tokens = append(tokens, string('~'))
					continue
				}

				tokens = append(tokens, string('-'))
				continue
			}

			tokens = append(tokens, string(runes[i]))
		} else {

			if currNumToken.Len() > 0 {
				tokens = append(tokens, currNumToken.String())
				currNumToken.Reset()
			}

			if i-1 >= 0 && unicode.IsDigit(runes[i-1]) && runes[i] == 'e' &&
				i+1 < len(input) && (runes[i+1] == '+' || runes[i+1] == '-') {
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
						currNumToken.WriteRune(runes[i])
					} else {
						break
					}
					i++
				}
				i--

				if currNumToken.Len() > 0 {
					powStr := currNumToken.String()
					currNumToken.Reset()
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
				currLetToken.WriteRune(runes[i])
			}
		}
	}

	if currNumToken.Len() > 0 {
		tokens = append(tokens, currNumToken.String())
	}

	if currLetToken.Len() > 0 {
		tokens = append(tokens, currLetToken.String())
	}

	return tokens
}

func getLetToken(token *strings.Builder) string {
	var res string
	switch token.String() {
	case "sqrt":
		res = "q"
	case "ln":
		res = "l"
	case "exp":
		res = "x"
	case "sin":
		res = "s"
	case "cos":
		res = "c"
	case "tg":
		res = "t"
	case "ctg":
		res = "g"
	default:
		res = token.String()
	}

	return res
}

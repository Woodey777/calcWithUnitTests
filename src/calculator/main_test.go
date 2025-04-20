package calculator

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	type CaseTokenize struct {
		expr   string
		result []string
	}
	cases := []CaseTokenize{
		{expr: "1+2", result: []string{"1", "+", "2"}},
		{expr: "-2   + 3.5", result: []string{"~", "2", "+", "3.5"}},
		{expr: "2 * -23.3.5", result: []string{"2", "*", "~", "23.3.5"}},
		{expr: "2 *ab 5", result: []string{"2", "*", "ab", "5"}},
		{expr: "2   4 / 5", result: []string{"2", "4", "/", "5"}},
		{expr: "5e-02", result: []string{"5", "d", "100"}},
		{expr: "e-02", result: []string{"e", "-", "02"}},
		{expr: "2e+02.0", result: []string{"2", "m", "100"}},
		{expr: "2^3", result: []string{"2", "^", "3"}},
		{expr: "2e+02 +  28", result: []string{"2", "m", "100", "+", "28"}},
		{expr: "2 e+ 2", result: []string{"2", "e", "+", "2"}},
		{expr: "3 *(2+ 4)", result: []string{"3", "*", "(", "2", "+", "4", ")"}},
		{expr: "sqrt(  cos(3  +e-20) ) ", result: []string{"q", "(", "c", "(", "3", "+", "e", "-", "20", ")", ")"}},
		{expr: "e +tg (250) / pi", result: []string{"e", "+", "t", "(", "250", ")", "/", "pi"}},
		{expr: "etg(2)", result: []string{"etg", "(", "2", ")"}},
		{expr: "sqrt(ln(e))", result: []string{"q", "(", "l", "(", "e", ")", ")"}},
		{expr: "3e+sqrt(4)", result: []string{"3", "e", "+", "q", "(", "4", ")"}},
	}

	for _, c := range cases {
		result := tokenize(c.expr)
		if slices.Compare(result, c.result) != 0 {
			t.Errorf("tokenize(%q) = %q, expected %q", c.expr, result, c.result)
		}
	}
}

func TestInfixToPostfix(t *testing.T) {
	type CaseInfixToPostfix struct {
		tokens  []string
		result  []string
		isError bool
	}
	cases := []CaseInfixToPostfix{
		{tokens: []string{"1", "+", "2"}, result: []string{"1", "2", "+"}, isError: false},
		{tokens: []string{"~", "2", "+", "3.5"}, result: []string{"2", "~", "3.5", "+"}, isError: false},
		{tokens: []string{"2", "4", "/", "5"}, result: []string{"2", "4", "5", "/"}, isError: false},
		{tokens: []string{"3", "*", "(", "2", "+", "4", ")"}, result: []string{"3", "2", "4", "+", "*"}, isError: false},
		{tokens: []string{"2", "^", "3"}, result: []string{"2", "3", "^"}, isError: false},
		{tokens: []string{"q", "(", "c", "(", "3", "+", "e", "-", "20", ")", ")"}, result: []string{"3", "e", "+", "20", "-", "c", "q"}, isError: false},
		{tokens: []string{"e", "+", "t", "(", "250", ")", "/", "pi"}, result: []string{"e", "250", "t", "pi", "/", "+"}, isError: false},
		{tokens: []string{"q", "(", "l", "(", "e", ")", ")"}, result: []string{"e", "l", "q"}, isError: false},

		{tokens: []string{"2", "*", "~", "23.3.5"}, result: nil, isError: true},
		{tokens: []string{"2", "*", "ab", "5"}, result: nil, isError: true},
		{tokens: []string{"etg", "(", "2", ")"}, result: nil, isError: true},
	}
	for _, c := range cases {
		result, err := infixToPostfix(c.tokens)
		if c.isError && err == nil {
			t.Errorf("infixToPostfix(%q), expected error", c.tokens)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("infixToPostfix(%q), unexpected error: %s", c.tokens, err)
			continue
		}

		if slices.Compare(result, c.result) != 0 {
			t.Errorf("infixToPostfix(%q) = %q, expected %q", c.tokens, result, c.result)
		}
	}
}

func TestEvaluatePostfix(t *testing.T) {
	type CaseEvaluatePostfix struct {
		tokens     []string
		result     float64
		angleUnits string
		isError    bool
	}
	cases := []CaseEvaluatePostfix{
		{tokens: []string{"1", "2", "+"}, result: 3.0, angleUnits: "radians", isError: false},
		{tokens: []string{"2", "~", "3.5", "-"}, result: -5.5, angleUnits: "radians", isError: false},
		{tokens: []string{"3", "~", "4", "9", "*", "+", "20", "~", "0.5", "/", "-"}, result: 73, angleUnits: "radians", isError: false},
		{tokens: []string{"2", "3", "^"}, result: 8, angleUnits: "radians", isError: false},
		{tokens: []string{"e", "17", "+", "c", "q"}, result: 0.8036167306671608, angleUnits: "radians", isError: false},
		{tokens: []string{"e", "250", "t", "pi", "/", "+"}, result: 1.4363579480746753, angleUnits: "radians", isError: false},
		{tokens: []string{"e", "l", "q"}, result: 1, angleUnits: "radians", isError: false},

		{tokens: []string{"5", "0", "/"}, result: 0, angleUnits: "radians", isError: true},
		{tokens: []string{"1", "2", "+", "5"}, result: 0, angleUnits: "radians", isError: true},
		{tokens: []string{"~"}, result: 0, angleUnits: "radians", isError: true},
		{tokens: []string{"+"}, result: 0, angleUnits: "radians", isError: true},
		{tokens: []string{"pi", "2", "/", "t"}, result: 0, angleUnits: "radians", isError: true},
		{tokens: []string{"pi", "g"}, result: 0, angleUnits: "radians", isError: true},
	}

	for _, c := range cases {
		result, err := evaluatePostfix(c.tokens, c.angleUnits)
		if c.isError && err == nil {
			t.Errorf("evaluatePostfix(%q), expected error", c.tokens)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("evaluatePostfix(%q), unexpected error: %s", c.tokens, err)
			continue
		}

		require.Equal(t, c.result, result)
	}
}

func TestCalculate(t *testing.T) {
	type CaseCalculate struct {
		expression string
		result     float64
		angleUnits string
		isError    bool
	}
	cases := []CaseCalculate{
		{expression: "1+2", result: 3, angleUnits: "radian", isError: false},
		{expression: "3.375e+09^(1/3)", result: 1500, angleUnits: "radian", isError: false},
		{expression: "sqrt(cos(17+e) ) ", result: 0.8036167306671608, angleUnits: "radian", isError: false},
		{expression: "e+tg(250)/pi", result: 1.4363579480746753, angleUnits: "radian", isError: false},
		{expression: "e+tg(250)/pi", result: 3.592831053138179, angleUnits: "degree", isError: false},
		{expression: "sqrt(ln(e))", result: 1, angleUnits: "radian", isError: false},
		{expression: "exp(5)+ sin(pi) *ctg(3)", result: 148.4131591025766, angleUnits: "radian", isError: false},
		{expression: "sin(90)", result: 1, angleUnits: "degree", isError: false},
		{expression: "sin(pi/2)", result: 1, angleUnits: "radian", isError: false},
		{expression: "sqrt(2^2 * 5 + 1)", result: 4.58257569495584, angleUnits: "radian", isError: false},
		{expression: "ln(exp(2))", result: 2, angleUnits: "radian", isError: false},
		{expression: "ln(e^2)", result: 2, angleUnits: "radian", isError: false},

		{expression: "1e+300 * 1e+300", result: 0, angleUnits: "radian", isError: true},
		{expression: "1e+308 / 0.5", result: 0, angleUnits: "radian", isError: true},
		{expression: "1e+308 + 1e+308", result: 0, angleUnits: "radian", isError: true},
		{expression: "-1e+308 -1e+308", result: 0, angleUnits: "radian", isError: true},
		{expression: "1e+300^2", result: 0, angleUnits: "radian", isError: true},
		{expression: "2 * -23.3.5", result: 0, angleUnits: "radian", isError: true},
		{expression: "ln(-5)", result: 0, angleUnits: "radian", isError: true},
		{expression: "sqrt(-5)", result: 0, angleUnits: "radian", isError: true},
		{expression: "exp(2000)", result: 0, angleUnits: "radian", isError: true},
	}

	for _, c := range cases {
		result, err := Calculate(c.expression, CalculatorConfig{AngleUnits: c.angleUnits})
		if c.isError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, c.result, result)
	}
}

func TestBench(t *testing.T) {
	type CaseCalculate struct {
		expression string
		result     float64
		isError    bool
	}
	cases := []CaseCalculate{
		{expression: "1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1",
			result: 250, isError: false},
		{expression: "1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000+1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000+1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000+1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			result: 4e+249, isError: false},
		{expression: "1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1+1",
			result: 1026, isError: false},
		{expression: "2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2*2",
			result: 9.134385233318143e+46, isError: false},
		{expression: "((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((1+1)))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))))",
			result: 2, isError: false},
		{expression: "1+2*3-4/5+6*7-sqrt(1+1)^32+5+6+7+8+9+10+11+12+13+14+15+16+17+18+19+20+21+22+23+24+25+26+27+28+29+30+31+32+33+34+35+36+37+38+39+40+41+42+43+44+45+46+47+48+49+50+51+52+53+54+55+56+(4*3+e*3e+21+4*3+e+4*3+e+4*3+e+4*3+e)/2*3-4/5+6*7-sqrt(1+1)^32+5+6+7+8+9+10+11+12+13+14+15+16+17+18+19+20+21+22+23+24+25+26+27+28+29+30+31+32+33+34+35+36+37+38+39+40+41+42+43+44+45+46+47+48+49+50+51+52+53+54+5",
			result: 1.2232268228065703e+22, isError: false},
		{expression: "sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(sqrt(256))))))))))))))))))))+(2^100)/((3+4*5)^10-sqrt(10000))/(3+4*5)^8*3-4/5+6*7-sqrt(1+1)^32",
			result: 1.1067549351347717e+06, isError: false},

	}
	// "(((...(1+1)...)))" (100 уровней вложенности)
	

	for _, c := range cases {
		start := time.Now()
		result, err := Calculate(c.expression, CalculatorConfig{AngleUnits: "radian"})
		duration := time.Since(start)

		if c.isError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Less(t, duration, 200*time.Millisecond)
		require.Equal(t, c.result, result)
	}
}

func TestAdd(t *testing.T) {
	type CaseAdd struct {
		a, b    float64
		result  float64
		isError bool
	}
	cases := []CaseAdd{
		{a: 1, b: 2, result: 3, isError: false},
		{a: -1, b: -2, result: -3, isError: false},
		{a: math.MaxFloat64, b: math.MaxFloat64, result: 0, isError: true},
		{a: -math.MaxFloat64, b: -math.MaxFloat64, result: 0, isError: true},
	}

	for _, c := range cases {
		result, err := Add(c.a, c.b)
		if c.isError && err == nil {
			t.Errorf("Add(%f, %f), expected error", c.a, c.b)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("Add(%f, %f), unexpected error: %s", c.a, c.b, err)
			continue
		}

		if result != c.result {
			t.Errorf("Add(%f, %f) = %f, expected %f", c.a, c.b, result, c.result)
		}
	}
}

func TestMul(t *testing.T) {
	type CaseMul struct {
		a, b    float64
		result  float64
		isError bool
	}
	cases := []CaseMul{
		{a: 1, b: 2, result: 2, isError: false},
		{a: math.MaxFloat64, b: math.MaxFloat64, result: 0, isError: true},
		{a: -math.MaxFloat64, b: -math.MaxFloat64, result: 0, isError: true},
	}

	for _, c := range cases {
		result, err := Mul(c.a, c.b)
		if c.isError && err == nil {
			t.Errorf("Mul(%f, %f), expected error", c.a, c.b)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("Mul(%f, %f), unexpected error: %s", c.a, c.b, err)
			continue
		}

		if result != c.result {
			t.Errorf("Mul(%f, %f) = %f, expected %f", c.a, c.b, result, c.result)
		}
	}
}

func TestPow(t *testing.T) {
	type CasePow struct {
		a, b    float64
		result  float64
		isError bool
	}
	cases := []CasePow{
		{a: 5, b: 5, result: 3125, isError: false},
		{a: math.MaxFloat64, b: math.MaxFloat64, result: 0, isError: true},
	}

	for _, c := range cases {
		result, err := Pow(c.a, c.b)
		if c.isError && err == nil {
			t.Errorf("Pow(%f, %f), expected error", c.a, c.b)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("Pow(%f, %f), unexpected error: %s", c.a, c.b, err)
			continue
		}

		if result != c.result {
			t.Errorf("Pow(%f, %f) = %f, expected %f", c.a, c.b, result, c.result)
		}
	}
}

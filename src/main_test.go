package main

import (
	"math"
	"slices"
	"testing"
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
		{expr: "2 *ab 5", result: []string{"2", "*", "a", "b", "5"}},
		{expr: "2   4 / 5", result: []string{"2", "4", "/", "5"}},
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
		{tokens: []string{"2", "*", "~", "23.3.5"}, result: nil, isError: true},
		{tokens: []string{"2", "*", "a", "b", "5"}, result: nil, isError: true},
		{tokens: []string{"2", "4", "/", "5"}, result: []string{"2", "4", "5", "/"}, isError: false},
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
		tokens  []string
		result  float64
		isError bool
	}
	cases := []CaseEvaluatePostfix{
		{tokens: []string{"1", "2", "+"}, result: 3.0, isError: false},
		{tokens: []string{"2", "~", "3.5", "-"}, result: -5.5, isError: false},
		{tokens: []string{"3", "~", "4", "9", "*", "+", "20", "~", "0.5", "/", "-"}, result: 73, isError: false},
		{tokens: []string{"5", "0", "/"}, result: 0, isError: true},
		{tokens: []string{"1", "2", "+", "5"}, result: 0, isError: true},
		{tokens: []string{"~"}, result: 0, isError: true},
		{tokens: []string{"+"}, result: 0, isError: true},
		{tokens: []string{"+"}, result: 0, isError: true},
	}

	for _, c := range cases {
		result, err := evaluatePostfix(c.tokens)
		if c.isError && err == nil {
			t.Errorf("evaluatePostfix(%q), expected error", c.tokens)
			continue
		} else if !c.isError && err != nil {
			t.Errorf("evaluatePostfix(%q), unexpected error: %s", c.tokens, err)
			continue
		}

		if result != c.result {
			t.Errorf("evaluatePostfix(%q) = %f, expected %f", c.tokens, result, c.result)
		}
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

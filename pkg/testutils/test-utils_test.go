package testutils_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/paulcarlton-ww/goutils/pkg/testutils"
)

func TestContainsStringArray(t *testing.T) {
	tests := []struct {
		testNum  int
		one      []string
		two      []string
		expected bool
	}{
		{1, []string{"a", "b", "c", "d"}, []string{"b", "c", "d"}, true},
		{2, []string{"a", "b", "c", "d"}, []string{"a", "b", "c", "d"}, true},
		{3, []string{"a", "b", "c", "d"}, []string{"b", "d"}, false},
		{4, []string{"a", "b", "c", "d"}, []string{"a", "b"}, true},
		{5, []string{"a"}, []string{}, false},
		{6, []string{}, []string{}, true},
		{7, []string{"a"}, []string{"a", "b"}, false},
	}

	for _, test := range tests {
		result := testutils.ContainsStringArray(test.one, test.two, false)
		if result != test.expected {
			t.Errorf("\nTest: %d\narray one:\n%+v\narray two:\n%+v\nExpected: %t\nGot.....: %t",
				test.testNum, test.one, test.two, test.expected, result)
		}
	}
}

func TestReadBuf(t *testing.T) {
	tests := []struct {
		expected *[]string
	}{{&[]string{"b", "c", "d"}}}

	for _, test := range tests {
		buffer := &bytes.Buffer{}

		for _, input := range *test.expected {
			buffer.WriteString(fmt.Sprintf("%s\n", input))
		}

		result := testutils.ReadBuf(buffer)
		if !testutils.ContainsStringArray(*result, *test.expected, true) {
			t.Errorf("\nExpected:\n%+v\nGot.....:\n%+v", test.expected, result)
		}
	}
}

func TestDisplayStrings(t *testing.T) {
	tests := []struct {
		one      []string
		expected string
	}{
		{[]string{"a", "b", "c"}, "0 - a\n1 - b\n2 - c"},
	}

	for _, test := range tests {
		result := testutils.DisplayStrings(test.one)
		if result != test.expected {
			t.Errorf("\ninput:\n%+v\nExpected:\n%s\nGot.....:\n%s",
				test.one, test.expected, result)
		}
	}
}

func TestCompareWhereList(t *testing.T) { // nolint: funlen // ok
	tests := []struct {
		testNum  int
		one      []string
		two      []string
		expected bool
	}{
		{
			1,
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(34)"},
			true,
		},
		{
			2,
			[]string{"SomeFunction() - some_file.go(NN)", "NotThisFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(34)"},
			false,
		},
		{
			3,
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(123)", "AnotherFunction() - another_file.go(45)"},
			true,
		},
		{
			4,
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(12346)", "AnotherFunction() - another_file.go(7)"},
			true,
		},
		{
			5,
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(34)"},
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(34)"},
			true,
		},
		{
			6,
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go()"},
			true,
		},
		{
			7,
			[]string{"SomeFunction() - some_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			false,
		},
		{
			8,
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			[]string{"SomeFunction() - some_file.go(NN)"},
			false,
		},
		{
			9,
			[]string{"SomeFunction() - some_file.go(12)", "AnotherFunction() - another_file.go(34)"},
			[]string{"SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)"},
			true,
		},
	}

	for _, test := range tests {
		result := testutils.CompareWhereList(test.one, test.two)
		if result != test.expected {
			t.Errorf("\nTest: %d\nlist one:\n%+v\nlist two:\n%+v\nExpected: %t\nGot.....: %t",
				test.testNum, test.one, test.two, test.expected, result)
		}
	}
}

func TestCompareWhere(t *testing.T) {
	tests := []struct {
		testNum  int
		one      string
		two      string
		expected bool
	}{
		{1, "SomeFunction() - some_file.go(NN)", "AnotherFunction() - another_file.go(NN)", false},
		{2, "SomeFunction() - some_file.go(NN)", "SomeFunction() - some_file.go(1)", true},
		{3, "SomeFunction() - some_file.go(NN)", "SomeFunction() - some_file.go(12)", true},
		{4, "SomeFunction() - some_file.go(NN)", "SomeFunction() - some_file.go(123)", true},
		{5, "SomeFunction() - some_file.go(NN)", "SomeFunction() - some_file.go(12345)", true},
		{6, "SomeFunction() - some_file.go(12345)", "SomeFunction() - some_file.go(NN)", true},
		{7, "SomeFunction() - some_file.go(1234)", "SomeFunction() - some_file.go(1234)", true},
		{8, "SomeFunction() - some_file.go(1234)", "SomeFunction() - some_file.go(1235)", false},
		{9, "SomeFunction() - some_file.go(NN)", "SomeFunction() - some_file.go()", true},
	}

	for _, test := range tests {
		result := testutils.CompareWhere(test.one, test.two)
		if result != test.expected {
			t.Errorf("\nTest: %d\none:\n%s\ntwo:\n%s\nExpected: %t\nGot.....: %t",
				test.testNum, test.one, test.two, test.expected, result)
		}
	}
}

func TestCompareItems(t *testing.T) { // nolint: funlen // ok
	type TestInfo struct {
		testNum  int
		one      interface{}
		two      interface{}
		expected bool
	}

	type subData struct {
		S string
		A []int
	}

	type testData struct {
		B bool
		I int
		F float64
		X interface{}
		E subData
		subData
	}

	errAnError := errors.New("an error") // nolint:goerr113 // ?

	errDifferentError := errors.New("a different error") // nolint:goerr113 // ?

	tests := []TestInfo{
		{testNum: 1, one: 1, two: 2, expected: false},
		{testNum: 2, one: 1, two: 1, expected: true},
		{ // nolint: dupl // ok
			testNum: 5, one: testData{
				B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}},
			},
			two: testData{
				B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 11, 7}},
			}, expected: false,
		},
		{ // nolint: dupl // ok
			testNum: 6, one: testData{
				B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}},
			},
			two: testData{
				B: true, I: 1, F: 12.43,
				X:       subData{S: "interface"},
				E:       subData{S: "sub struct", A: []int{1, 2, 3}},
				subData: subData{S: "embedded", A: []int{9, 8, 7}},
			}, expected: true,
		},
		{testNum: 7, one: "one", two: "two", expected: false},
		{testNum: 8, one: "one", two: "one", expected: true},
		{testNum: 8, one: errAnError, two: errDifferentError, expected: false},
		{testNum: 10, one: errAnError, two: errAnError, expected: true},
	}

	for _, test := range tests {
		result := testutils.CompareItems(test.one, test.two)
		if result != test.expected {
			t.Errorf("Test: %d\nExpected:\n%t\nGot....:\n%t\nInput Data:\n%+v\n%+v\n",
				test.testNum, test.expected, result, test.one, test.two)
		}
	}
}

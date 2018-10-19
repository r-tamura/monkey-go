package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 * (2 + 10)", 60},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVMTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    "true",
			expected: true,
		},
		{
			input:    "false",
			expected: false,
		},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == false", false},
		{"true != false", true},
		{"true == true", true},
		{"true != true", false},
		{"false == true", false},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 10 })", true},
	}
	runVMTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVMTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}
	runVMTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1: 2, 2: 3}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 2}).HashKey(): 3,
		}},
		{"{1 + 1: 2 * 2, 3 + 3: 4 * 4}", map[object.HashKey]int64{
			(&object.Integer{Value: 2}).HashKey(): 4,
			(&object.Integer{Value: 6}).HashKey(): 16,
		}},
	}
	runVMTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVMTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) {10} else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
	}
	runVMTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVMTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
		let fivePlusTen = fn() { 5 + 10; };
		fivePlusTen()
		`,
			expected: 15,
		},
		{
			input: `
		let one = fn() { 1; };
		let two = fn() { 2; };
		one() + two()
		`,
			expected: 3,
		},
		{
			input: `
		let a = fn() { 1; };
		let b = fn() { a() + 1; };
		let c = fn() { b() + 1; };
		c()
		`,
			expected: 3,
		},
	}

	runVMTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() { };
			noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() { };
			let noReturnTwo = fn() { noReturn() };
			noReturn();
			noReturnTwo();
			`,
			expected: Null,
		},
	}
	runVMTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() { 1; };
			let returnOneReturner = fn() { returnsOne; };
			returnOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnOne
			}
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}
	runVMTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num;
			}
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}
	runVMTests(t, tests)
}

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
			return
		}
		if len(array.Elements) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
			return
		}
		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
			return
		}
		if len(hash.Pairs) != len(expected) {
			t.Errorf("hash hash wrong number of pairs. got=%d,want=%d", len(expected), len(hash.Pairs))
			return
		}
		for expectedKey, expectedValue := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}
			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null: %T (%+v)", actual, actual)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		// Memo %+v フィールド名付き
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if expected != result.Value {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if expected != result.Value {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if expected != result.Value {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

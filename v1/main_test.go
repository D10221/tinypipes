package main

import (
	"testing"
	"fmt"
)

// ClearScreen  does what it says
func ClearScreen() {
	Output.WriteString("\033[2J")
}


//ActionToSegmentAdapter Adapt func()void to segment
func ActionToSegmentAdapter(action func()) Segment {
	return func(x interface{}) (interface{}, error) {
		action()
		return x, nil
	}
}



// Test_ConnectInt tests IntSegment adapter
func Test_ConnectInt(t *testing.T) {

	intSegment := func(i int) (interface{}, error) {
		return i, nil
	}

	segment := func(x interface{}) (interface{}, error) {
		return 42, nil
	}

	pipe := NewPipe(
		segment,
		ConnectInt(intSegment),
	)

	x, e := pipe.Process()
	if e != nil {
		t.Error(e)
		return
	}
	if x == nil || x != 42 {
		t.Error("None or bad Value")
	}

}

//
func Test_ConnectBytes(t *testing.T){
	byteSegement := func(b []byte) (interface{}, error) {
		return b, nil
	}

	segment := func(x interface{}) (interface{}, error) {
		return []byte("42"), nil
	}

	pipe := NewPipe(
		segment,
		ConnectBytes(byteSegement),
	)

	x, e := pipe.Process()
	if e != nil {
		t.Error(e)
		return
	}
	if x == nil{
		t.Error("No Value")
	}
	 b,ok := x.([]byte)
	if !ok {
		t.Error("Not a []byte")
		return
	}
	if string(b) != "42" {
		t.Errorf("bad value: %v", b)
	}
}

func Test_ToFloat(t *testing.T) {

	t.Log("ToFloat")

	values := []interface{}{"3.3", 3.3 }

	for _, value := range values {
		if _, e := ConnectString(ParseFloat)(value); e != nil {
			t.Error(e);
			return
		}
	}

	values = []interface{}{true, 3 }

	for _, value := range values {
		if _, e := ConnectString(ParseFloat)(value); e == nil {
			t.Errorf("Shouldn't accept anything but %v", []KnownType{String, Float64});
			return
		}
	}
}

type testCase struct {
	Input    string
	expected float64
}

// Test_GetNumber
func Test_GetNumber(t *testing.T) {

	t.Log("Test_GetNumber")
	tests := []testCase{
		testCase{Input: "3.3", expected: 3.3 },
		testCase{Input: " 3.3", expected: 3.3 },
		testCase{Input: "\t3.3", expected: 3.3 },
		testCase{Input: "\t3.3\r", expected: 3.3 },
		testCase{Input: "\t3.3\n", expected: 3.3 },
		testCase{Input: "\t3.3 \r", expected: 3.3 },
		testCase{Input: "\t3.3 \r", expected: 3.3 },
		testCase{Input: "33", expected: 33 },
		testCase{Input: "3", expected: 3 },
	}

	for _, test := range tests {

		b := clean([]byte(test.Input))
		n, e := GetNumber(b)
		if e != nil {
			t.Error(e)
			return
		}
		t.Log(n)
		s, e := ConnectString(ParseFloat)(n)
		if e != nil {
			t.Error(e)
			return
		}
		f, ok := s.(float64)
		if !ok {
			t.Error("Bad conversion")
			return
		}
		if f != test.expected {
			t.Errorf("Expected: %v , got: %v", test.expected, f)
			return
		}
	}
}


// GetNumber from input?
func GetNumber(x interface{}) (interface{}, error) {
	value, ok := x.([]byte)
	if !ok {
		return x, fmt.Errorf("Expected %v", Byte.String())
	}
	outBytes := numberRegex.Find(clean(value))
	if outBytes == nil {
		return x, fmt.Errorf("\n No match for %v in '%v'", numberRegex.String(), clean(value))
	}
	return string(outBytes), nil
}


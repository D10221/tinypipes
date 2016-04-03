/*
main    its's a bit over engineered,
	and could be better with more files,
	but I know to know my limits
	and has to fit in one ..
	it would benefit from generics,
	but we have none,
	it would benefit code generation,
	but I coded none,
	I'm learning go :)
	mind the runtime checking ... ;]
	implementing sort of piping ...
		.. I'm gone
	NOTE: v2 will make it rhyme
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// main problem with this setup is that in Segment , things needs to get converted or type-checked all the time ,
	// only one value travels the pipe and it gets used transformed , copy  , etc
	// no access to anything like a context
	// pipe segments will need to access globals to modify state outside the parameter they received
	// the order of the pipes is important
	// having oddly shaped  pipes , needs to get adapted
	pipes := []*Pipe {
		NewPipe(
			GetInput("Enter integer value: "),
			ConnectBytes(GetNumberFromBytesAsString),
			ConnectString(ParseInt),
			Show,
		),
		NewPipe(
			//ActionToSegmentAdapter(ClearScreen),
			GetInput("Enter Decimal value: "),
			ConnectBytes(GetNumberFromBytesAsString),
			ConnectString(ParseFloat),
			Show,
		),
	}
	ProcessPipes(pipes)
}

// KnownType enum? with stringer
type KnownType int

const (
	Byte KnownType = iota
	String
	Int64
	Float64
)

const _KnownTypeName = "ByteStringInt64Float64"

var _KnownTypeIndex = [...]uint8{0, 4, 10, 15, 22}

func (i KnownType) String() string {
	if i < 0 || i >= KnownType(len(_KnownTypeIndex)-1) {
		return fmt.Sprintf("KnownType(%d)", i)
	}
	return _KnownTypeName[_KnownTypeIndex[i]:_KnownTypeIndex[i+1]]
}
// KnownType end

var (
	//Output to...
	Output = bufio.NewWriter(os.Stdout) // *bufio.Writer
	//Input  from...
	Input = *bufio.NewReader(os.Stdin)
)

// RESET_LINE Return cursor to start of line and clean it
const _Reset = "\r\033[K"

// Print prints parameters to Output
func Print(a ...interface{}) {
	fmt.Fprint(Output, _Reset)
	fmt.Fprint(Output, a...)
	Output.Flush()
}
/*Segments*/

//Segment to process in the tiny pipeline
type Segment func(interface{}) (interface{}, error);

type ByteSegment func([]byte)(interface{}, error);

type StringSegment func(string)(interface{}, error);

type IntSegment func(int)(interface{}, error);

type FloatSegment func(float64)(interface{}, error);

//Pipe TinyPipe
type Pipe struct {
	Segments []Segment
	Debug    bool
}

// NewPipe fty
func NewPipe(segments ...Segment) *Pipe {
	return &Pipe{Segments: segments[:]}
}

// Process all pipe's segments while no error
func (p *Pipe) Process() (interface{}, error) {
	var value interface{}
	var e error
	for i, segment := range p.Segments {
		value, e = segment(value)
		if e != nil {
			break
		}
		if p.Debug {
			Print(fmt.Sprintf("Segment: %v , out: %v", i, value))
		}
	}
	return value, e
}

//ProcessPipes Process all pipes while no error
func ProcessPipes(pipes []*Pipe) {
	for _, pipe := range pipes {
		_, e := Process(pipe)
		if e != nil {
			break
		}
	}
	Print("Goodbye!.\n")
}

// Process all pipe's segments while no error
func Process(pipe *Pipe) (interface{}, error) {
	x, err := pipe.Process()
	if err != nil {
		Print((fmt.Sprintf("Error: %s \n", err)))
	} else {
		Print(fmt.Sprintf("%v \n", x))
	}
	return x, err
}

/*
	Adapters
*/

// ConnectBytes adapts ByteSegment to Segment
func ConnectBytes(byteSegment ByteSegment) Segment{
	return func(x interface{}) (interface{}, error) {
		value, ok := x.([]byte)
		if !ok {
			return x, fmt.Errorf("Expected %v", Byte.String())
		}
		return byteSegment(value)
	}
}

// ConnectInt adapts intSegment to Segment
func ConnectInt(in IntSegment) Segment {
	return func(x interface{}) (interface{}, error) {
		value, ok := x.(int)
		if !ok {
			return x, fmt.Errorf("Expected %v", Int64.String())
		}
		return in(value)
	}
}

// ConnectString adapts intSegment to Segment
func ConnectString(in StringSegment) Segment {
	return func(x interface{}) (interface{}, error) {
		value, ok := x.(string)
		if !ok {
			return x, fmt.Errorf("Expected %v", String.String())
		}
		return in(value)
	}
}

// ConnectFloat connects a FloatSegment to a SegmentPipe
func ConnectFloat(fSegment FloatSegment) Segment {

	return func(x interface{}) (interface{}, error) {
		value, ok := x.(float64)
		if !ok {
			return x, fmt.Errorf("Expected %v , maybe %v", String.String(), Float64.String())
		}
		return fSegment(value)
	}
}

/*
	Segment implementations
	-----------------------
*/


var cleanRgx = regexp.MustCompile(`(\t|\r\|\n|\s)`)

func clean(b []byte) []byte {
	return cleanRgx.ReplaceAll(b, []byte(""))
}

// GetInput get user input
func GetInput(message string) Segment {
	return func(interface{}) (interface{}, error) {
		Print(message)
		x, e := Input.ReadBytes('\n')
		return x, e
	}
}

var numberRegex = regexp.MustCompile(`^(?:\s+)?\t?\r?(\d+|\d+\.?\d+?)\b$(?:\s+|\t|\r|\n)?`)

// GetNumber from input?
func GetNumberFromBytesAsString(value []byte) (interface{}, error) {
	outBytes := numberRegex.Find(clean(value))
	if outBytes == nil {
		return value, fmt.Errorf("\n No match for %v in '%v'", numberRegex.String(), clean(value))
	}
	return string(outBytes), nil
}

func ParseInt(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

func ParseFloat(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

func Show(x interface{}) (interface{}, error) {
	if x == nil {
		return x, fmt.Errorf("Expected %v", "something")
	}
	return fmt.Sprintf("You entered Number %v", x), nil

}

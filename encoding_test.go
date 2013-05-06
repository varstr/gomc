package gomc

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"text/template"
)

const (
	_TEST_STRUCT_TPL = `
    &TestStruct{
        "Str": {{.Str}}
        "Int": {{.Int}}
        "Slice": {
            {{range $_, $i := .Slice}} {{$i | printf "%q"}}, {{end}}
        },
        "Map": {
            {{range $k, $v := .Map}} {{$k | printf "%q"}}: {{$v | printf "%q"}}, {{end}}
        }
    }`
)

type TestStruct struct {
	Str   string            `json: str`
	Int   int               `json: int`
	Slice []string          `json: slice`
	Map   map[string]string `json: map`
}

func randomStr(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(33, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomStruct() *TestStruct {
	return &TestStruct{
		Str:   randomStr(10),
		Int:   randInt(0, 100),
		Slice: []string{randomStr(5), randomStr(5), randomStr(5)},
		Map:   map[string]string{randomStr(5): randomStr(10), randomStr(5): randomStr(10), randomStr(5): randomStr(10)},
	}
}

func (self *TestStruct) format() string {
	tpl := template.New("tpl")
	tpl, e := tpl.Parse(_TEST_STRUCT_TPL)
	if e != nil {
		panic(e)
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, self)
	return buf.String()
}

func equal(origin, restore interface{}) bool {
	switch restore.(type) {
	case *bool:
		return origin == *restore.(*bool)
	case *int:
		return origin == *restore.(*int)
	case *int8:
		return origin == *restore.(*int8)
	case *int16:
		return origin == *restore.(*int16)
	case *int32:
		return origin == *restore.(*int32)
	case *int64:
		return origin == *restore.(*int64)
	case *uint:
		return origin == *restore.(*uint)
	case *uint8:
		return origin == *restore.(*uint8)
	case *uint16:
		return origin == *restore.(*uint16)
	case *uint32:
		return origin == *restore.(*uint32)
	case *uint64:
		return origin == *restore.(*uint64)
	case *string:
		return origin == *restore.(*string)
	case *[]byte:
		return reflect.DeepEqual(origin, *restore.(*[]byte))
	case *TestStruct:
		return reflect.DeepEqual(origin, restore.(*TestStruct))
	}
	return false
}

func testBaseTypes(origin, restore interface{}, t *testing.T) {
	b, f, e := encode(origin, ENCODING_GOB)
	if e != nil {
		t.Error("Fail to encode:", e)
	} else if f != encodingFlag(ENCODING_DEFAULT) {
		t.Error("Error return flag:", f, ", expect:", encodingFlag(ENCODING_DEFAULT))
	} else if e = decode(b, encodingFlag(ENCODING_DEFAULT), restore); e != nil {
		t.Error("Fail to decode:", e)
	} else if !equal(origin, restore) {
		t.Error("Error restore:", restore, ", expect:", origin)
	}
}

func testStruct(origin, restore *TestStruct, encoding EncodingType, t *testing.T) {
	b, f, e := encode(origin, encoding)
	if e != nil {
		t.Error("Fail to encode:", e)
	} else if f != encodingFlag(encoding) {
		t.Error("Error return flag:", f, ", expect:", encodingFlag(encoding))
	} else if e = decode(b, encodingFlag(encoding), restore); e != nil {
		t.Error("Fail to decode:", e)
	} else if !equal(origin, restore) {
		t.Error("Error restore:", restore.format(), ", expect:", origin.format())
	}
}

func TestBool(t *testing.T) {
	origin := true
	restore := new(bool)
	testBaseTypes(origin, restore, t)
}

func TestInt(t *testing.T) {
	origin := 42
	restore := new(int)
	testBaseTypes(origin, restore, t)
}

func TestInt8(t *testing.T) {
	origin := int8(42)
	restore := new(int8)
	testBaseTypes(origin, restore, t)
}

func TestInt16(t *testing.T) {
	origin := int16(42)
	restore := new(int16)
	testBaseTypes(origin, restore, t)
}

func TestInt32(t *testing.T) {
	origin := int32(42)
	restore := new(int32)
	testBaseTypes(origin, restore, t)
}

func TestInt64(t *testing.T) {
	origin := int64(42)
	restore := new(int64)
	testBaseTypes(origin, restore, t)
}

func TestUint(t *testing.T) {
	origin := uint(42)
	restore := new(uint)
	testBaseTypes(origin, restore, t)
}

func TestUint8(t *testing.T) {
	origin := uint8(42)
	restore := new(uint8)
	testBaseTypes(origin, restore, t)
}

func TestUint16(t *testing.T) {
	origin := uint16(42)
	restore := new(uint16)
	testBaseTypes(origin, restore, t)
}

func TestUint32(t *testing.T) {
	origin := uint32(42)
	restore := new(uint32)
	testBaseTypes(origin, restore, t)
}

func TestUint64(t *testing.T) {
	origin := uint64(42)
	restore := new(uint64)
	testBaseTypes(origin, restore, t)
}

func TestByteSlice(t *testing.T) {
	origin := []byte{1, 1, 2, 3, 5, 8, 13}
	restore := new([]byte)
	testBaseTypes(origin, restore, t)
}

func TestStructGob(t *testing.T) {
	origin := randomStruct()
	restore := new(TestStruct)
	testStruct(origin, restore, ENCODING_GOB, t)
}

func TestStructJSON(t *testing.T) {
	origin := randomStruct()
	restore := new(TestStruct)
	testStruct(origin, restore, ENCODING_JSON, t)
}

func BenchmarkEncodeDefault(b *testing.B) {
	b.StopTimer()
	origin := randomStr(10)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		encode(origin, ENCODING_DEFAULT)
	}
}

func BenchmarkDecodeDefault(b *testing.B) {
	b.StopTimer()
	origin := randomStr(10)
	restore := new(string)
	buffer, flag, _ := encode(origin, ENCODING_DEFAULT)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		decode(buffer, flag, restore)
	}
}

func benchmarkEncode(b *testing.B, encoding EncodingType) {
	b.StopTimer()
	origin := randomStruct()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		encode(origin, encoding)
	}
}

func benchmarkDecode(b *testing.B, encoding EncodingType) {
	b.StopTimer()
	origin := randomStruct()
	restore := new(TestStruct)
	buffer, flag, _ := encode(origin, encoding)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		decode(buffer, flag, restore)
	}
}

func BenchmarkEncodeGob(b *testing.B) {
	benchmarkEncode(b, ENCODING_GOB)
}

func BenchmarkDecodeGob(b *testing.B) {
	benchmarkDecode(b, ENCODING_GOB)
}

func BenchmarkEncodeJSON(b *testing.B) {
	benchmarkEncode(b, ENCODING_JSON)
}

func BenchmarkDecodeJSON(b *testing.B) {
	benchmarkDecode(b, ENCODING_JSON)
}

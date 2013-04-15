package gomc

import (
    "reflect"
    "testing"
)

type Test struct {
    X int `json:"x"`
    Y int `json:"y"`
    Z int `json:"z"`
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
    case *Test:
        return reflect.DeepEqual(origin, restore.(*Test))
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
        t.Error("Error restore", restore, ", expect:", origin)
    }
}

func testStruct(origin, restore interface{}, encoding EncodingType, t *testing.T) {
    b, f, e := encode(origin, encoding)
    if e != nil {
        t.Error("Fail to encode:", e)
    } else if f != encodingFlag(encoding) {
        t.Error("Error return flag:", f, ", expect:", encodingFlag(encoding))
    } else if e = decode(b, encodingFlag(encoding), restore); e != nil {
        t.Error("Fail to decode:", e)
    } else if !equal(origin, restore) {
        t.Error("Error restore", restore, ", expect:", origin)
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
    origin := &Test{1,2,3}
    restore := new(Test)
    testStruct(origin, restore, ENCODING_GOB, t)
}

func TestStructJSON(t *testing.T) {
    origin := &Test{1,2,3}
    restore := new(Test)
    testStruct(origin, restore, ENCODING_JSON, t)
}

package gomc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"strconv"
)

type EncodingType uint
type EncodeFunc func(interface{}) ([]byte, error)
type DecodeFunc func([]byte, interface{}) error

const (
	_NUMERIC_BASE = 10
)

const (
	ENCODING_DEFAULT EncodingType = iota
	ENCODING_GOB
	ENCODING_JSON
)

var (
	encoders = map[EncodingType]EncodeFunc{
		ENCODING_DEFAULT: encodeDefault,
		ENCODING_GOB:     encodeGob,
		ENCODING_JSON:    json.Marshal,
	}

	decoders = map[EncodingType]DecodeFunc{
		ENCODING_DEFAULT: decodeDefault,
		ENCODING_GOB:     decodeGob,
		ENCODING_JSON:    json.Unmarshal,
	}
)

func encodingFlag(encoding EncodingType) uint32 {
	return 1 << encoding
}

func encodeDefault(object interface{}) (buffer []byte, err error) {
	switch object.(type) {
	case bool:
		buffer = strconv.AppendBool(buffer, object.(bool))
	case int:
		buffer = strconv.AppendInt(buffer, int64(object.(int)), _NUMERIC_BASE)
	case int8:
		buffer = strconv.AppendInt(buffer, int64(object.(int8)), _NUMERIC_BASE)
	case int16:
		buffer = strconv.AppendInt(buffer, int64(object.(int16)), _NUMERIC_BASE)
	case int32:
		buffer = strconv.AppendInt(buffer, int64(object.(int32)), _NUMERIC_BASE)
	case int64:
		buffer = strconv.AppendInt(buffer, object.(int64), _NUMERIC_BASE)
	case uint:
		buffer = strconv.AppendUint(buffer, uint64(object.(uint)), _NUMERIC_BASE)
	case uint8:
		buffer = strconv.AppendUint(buffer, uint64(object.(uint8)), _NUMERIC_BASE)
	case uint16:
		buffer = strconv.AppendUint(buffer, uint64(object.(uint16)), _NUMERIC_BASE)
	case uint32:
		buffer = strconv.AppendUint(buffer, uint64(object.(uint32)), _NUMERIC_BASE)
	case uint64:
		buffer = strconv.AppendUint(buffer, object.(uint64), _NUMERIC_BASE)
	case string:
		buffer = []byte(object.(string))
	case []byte:
		buffer = object.([]byte)
	default:
		err = errors.New("Invalid object for default encode")
	}
	return
}

func decodeDefault(buffer []byte, object interface{}) (err error) {
	str := string(buffer)
	switch object.(type) {
	case *bool:
		value := object.(*bool)
		*value, err = strconv.ParseBool(str)
	case *int:
		v, err := strconv.ParseInt(str, _NUMERIC_BASE, 0)
		if err == nil {
			value := object.(*int)
			*value = int(v)
		}
	case *int8:
		v, err := strconv.ParseInt(str, _NUMERIC_BASE, 8)
		if err == nil {
			value := object.(*int8)
			*value = int8(v)
		}
	case *int16:
		v, err := strconv.ParseInt(str, _NUMERIC_BASE, 16)
		if err == nil {
			value := object.(*int16)
			*value = int16(v)
		}
	case *int32:
		v, err := strconv.ParseInt(str, _NUMERIC_BASE, 32)
		if err == nil {
			value := object.(*int32)
			*value = int32(v)
		}
	case *int64:
		value := object.(*int64)
		*value, err = strconv.ParseInt(str, _NUMERIC_BASE, 64)
	case *uint:
		v, err := strconv.ParseUint(str, _NUMERIC_BASE, 0)
		if err == nil {
			value := object.(*uint)
			*value = uint(v)
		}
	case *uint8:
		v, err := strconv.ParseUint(str, _NUMERIC_BASE, 8)
		if err == nil {
			value := object.(*uint8)
			*value = uint8(v)
		}
	case *uint16:
		v, err := strconv.ParseUint(str, _NUMERIC_BASE, 16)
		if err == nil {
			value := object.(*uint16)
			*value = uint16(v)
		}
	case *uint32:
		v, err := strconv.ParseUint(str, _NUMERIC_BASE, 32)
		if err == nil {
			value := object.(*uint32)
			*value = uint32(v)
		}
	case *uint64:
		value := object.(*uint64)
		*value, err = strconv.ParseUint(str, _NUMERIC_BASE, 64)
	case *string:
		value := object.(*string)
		*value = str
	case *[]byte:
		value := object.(*[]byte)
		*value = buffer
	default:
		err = errors.New("Invalid object for default decode")
	}
	return
}

func encodeGob(object interface{}) (buffer []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err = encoder.Encode(object)
	if err != nil {
		return
	}
	buffer = buf.Bytes()
	return
}

func decodeGob(buffer []byte, object interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(buffer))
	return decoder.Decode(object)
}

func encode(object interface{}, encoding EncodingType) (buffer []byte, flag uint32, err error) {
	if buffer, err = encodeDefault(object); err == nil {
		flag = encodingFlag(ENCODING_DEFAULT)
	} else if encoder, ok := encoders[encoding]; ok {
		buffer, err = encoder(object)
		flag = encodingFlag(encoding)
	} else {
		err = errors.New("Unsupported encoding type")
	}
	return
}

func decode(buffer []byte, flag uint32, object interface{}) (err error) {
	for encoding, decoder := range decoders {
		if flag&encodingFlag(encoding) != 0 {
			return decoder(buffer, object)
		}
	}
	err = errors.New("Unsupported decoding flag")
	return
}

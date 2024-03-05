package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type dataType byte
type typeString rune
type Value struct {
	typ_  dataType
	str   string
	num   int
	bulk  string
	array []Value
}

const (
	STRING dataType = iota
	ERROR
	INTEGER
	ARRAY
	BULK
	NULL
	PARSERROR
)

const (
	stringType  typeString = '+'
	errorType   typeString = '-'
	integerType typeString = ':'
	bulkType    typeString = '$'
	arrayType   typeString = '*'
)

type Resp struct {
	reader *bufio.Reader
}

type Writer struct {
	writer io.Writer
}

func NewResp(r io.Reader) *Resp {
	return &Resp{bufio.NewReader(r)}
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (r *Resp) parseType(typ_ rune) (dataType, error) {
	switch typ_ {
	case '+':
		return STRING, nil
	case '-':
		return ERROR, nil
	case ':':
		return INTEGER, nil
	case '$':
		return BULK, nil
	case '*':
		return ARRAY, nil

	default:
		return PARSERROR, errors.New("unrecognized type")
	}
}

func (r *Resp) readLine() ([]byte, int, error) {
	var (
		line []byte
		n    int
	)

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (int, int, error) {
	l, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	readInt, err := strconv.ParseInt(string(l), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(readInt), n, nil
}

func (r *Resp) Read() (Value, error) {
	typ_, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	dt, err := r.parseType(rune(typ_))
	if err != nil {
		return Value{}, err
	}

	switch dt {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		return Value{}, errors.New("unrecognized type")
	}
}

func (r *Resp) readArray() (Value, error) {
	var val Value
	val.typ_ = ARRAY

	size, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	val.array = make([]Value, 0)
	for i := 0; i < size; i++ {
		v, err := r.Read()
		if err != nil {
			return v, err
		}

		val.array = append(val.array, v)
	}
	return val, nil
}

func (r *Resp) readBulk() (Value, error) {
	var val Value
	val.typ_ = BULK

	size, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	bulk := make([]byte, size)

	val.bulk = string(bulk)

	//skip clrf
	r.readLine()
	return val, nil
}

func (v Value) Marshal() []byte {
	switch v.typ_ {
	case ARRAY:
		return v.marshalArray()
	case BULK:
		return v.marshalBulk()
	case STRING:
		return v.marshalString()
	case ERROR:
		return v.marshalError()
	default:
		return v.marshalNull()

	}
}

func (v Value) marshalArray() []byte {
	var b []byte
	b = append(b, byte(arrayType))
	b = append(b, strconv.Itoa(len(v.array))...)
	b = append(b, "\r\n"...)

	for i := 0; i < len(v.array); i++ {
		b = append(b, v.array[i].Marshal()...)
	}

	return b
}

func (v Value) marshalBulk() []byte {
	var b []byte
	b = append(b, byte(bulkType))
	b = append(b, strconv.Itoa(len(v.bulk))...)
	b = append(b, "\r\n"...)
	b = append(b, v.bulk...)
	b = append(b, "\r\n"...)

	return b
}

func (v Value) marshalString() []byte {
	var b []byte
	b = append(b, byte(stringType))
	b = append(b, v.str...)
	b = append(b, "\r\n"...)

	return b
}

func (v Value) marshalError() []byte {
	var b []byte
	b = append(b, byte(errorType))
	b = append(b, v.str...)
	b = append(b, "\r\n"...)

	return b
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (w Writer) Write(v Value) error {
	_, err := w.writer.Write(v.Marshal())
	if err != nil {
		return err
	}

	return nil
}

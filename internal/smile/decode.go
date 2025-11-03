// Copyright (c) 2024 LeLuxNetpackage smile

// Licensed under the MIT License
// Original source: https://gitlab.com/LeLuxNet/X/-/blob/c09411c26dfb/encoding/smile/decode.go
//
// Package smile implements decoding of Smile as defined in
// https://github.com/FasterXML/smile-format-specification/blob/master/smile-specification.md
package smile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

const magic = ":)\n"

func Unmarshal(data []byte, v interface{}) error {
	if len(data) < 4 || string(data[:len(magic)]) != magic {
		return errors.New("smile: invalid header")
	}

	h := data[3]
	if ver := h >> 4; ver != 0 {
		return fmt.Errorf("smile: unsupported version: %d", ver)
	}

	d := decodeState{
		r:          bytes.NewReader(data[4:]),
		rawBinary:  h&4 != 0,
		sStringVal: h&2 != 0,
		sPropName:  h&1 != 0,
		buf:        make([]byte, 1),
	}
	return d.unmarshal(v)
}

type decodeState struct {
	r   io.Reader
	buf []byte

	rawBinary  bool
	sStringVal bool
	sPropName  bool

	sKeys shared
	sVals shared
}

func (d *decodeState) unmarshal(v interface{}) error {
	return d.decode(reflect.ValueOf(v).Elem())
}

const (
	emptyString = iota + 0x20
	null
	falseTok
	trueTok
	int32Tok
	int64Tok
	bigInt
	_
	float32Tok
	float64Tok
	bigDecimal
)

const (
	longAscii   = 0xe0
	longUnicode = 0xe4
	longSString = 0xec

	startArray  = 0xf8
	endArray    = 0xf9
	startObject = 0xfa
	endObject   = 0xfb

	endString = 0xfc
)

func (d *decodeState) ReadByte() (byte, error) {
	for {
		n, err := d.r.Read(d.buf)
		if err != nil {
			return 0, err
		}
		if n != 0 {
			return d.buf[0], nil
		}
	}
}

func (d *decodeState) decode(v reflect.Value) error {
	return d.value(v)
}

func (d *decodeState) value(v reflect.Value) error {
	b, err := d.ReadByte()
	if err != nil {
		return err
	}

	switch b & 0xe0 {
	case 0x00:
		return d.setString(v, d.sVals[b&0x1f-1])
	case 0x20:
		switch b {
		case emptyString:
			return d.setString(v, "")
		case null:
			switch v.Kind() {
			case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
				v.Set(reflect.Zero(v.Type()))
			}
			return nil
		case falseTok, trueTok:
			switch v.Kind() {
			case reflect.Bool:
				v.SetBool(b == trueTok)
			case reflect.Interface:
				if v.NumMethod() == 0 {
					v.Set(reflect.ValueOf(b == trueTok))
				}
			}
			return nil
		case int32Tok, int64Tok:
			n, err := d.int(true)
			if err != nil {
				return err
			}
			return d.setInt(v, n)
		case bigInt:
			n, err := d.bigInt()
			if err != nil {
				return err
			}

			switch v.Kind() {
			case reflect.String:
				v.SetString(n.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v.SetInt(n.Int64())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v.SetUint(n.Uint64())
			case reflect.Float32, reflect.Float64:
				v.SetFloat(float64(n.Int64()))
			case reflect.Interface:
				if v.NumMethod() == 0 {
					v.Set(reflect.ValueOf(n))
				}
			}
			return nil
		case float32Tok:
			n, err := d.float32()
			if err != nil {
				return err
			}

			switch v.Kind() {
			case reflect.String:
				v.SetString(strconv.FormatFloat(float64(n), 'f', -1, 32))
			case reflect.Float32, reflect.Float64:
				v.SetFloat(float64(n))
			case reflect.Interface:
				if v.NumMethod() == 0 {
					v.Set(reflect.ValueOf(n))
				}
			}
			return nil
		case float64Tok:
			n, err := d.float64()
			if err != nil {
				return err
			}

			switch v.Kind() {
			case reflect.String:
				v.SetString(strconv.FormatFloat(n, 'f', -1, 64))
			case reflect.Float32, reflect.Float64:
				v.SetFloat(n)
			case reflect.Interface:
				if v.NumMethod() == 0 {
					v.Set(reflect.ValueOf(n))
				}
			}
			return nil
		case bigDecimal:
			n, err := d.bigDecimal()
			if err != nil {
				return err
			}

			switch v.Kind() {
			case reflect.String:
				v.SetString(n.String())
			case reflect.Float32, reflect.Float64:
				f, _ := n.Float64()
				v.SetFloat(f)
			case reflect.Interface:
				if v.NumMethod() == 0 {
					v.Set(reflect.ValueOf(n))
				}
			}
			return nil
		}
	case 0x40:
		return d.string(v, b, 1, &d.sVals)
	case 0x60:
		return d.string(v, b, 1+32, &d.sVals)
	case 0x80:
		return d.string(v, b, 2, &d.sVals)
	case 0xa0:
		return d.string(v, b, 2+32, &d.sVals)
	case 0xc0:
		return d.setInt(v, d.smallInt(b))
	case 0xe0:
		switch b {
		case longAscii, longUnicode:
			s, err := d.longString()
			if err != nil {
				return err
			}
			return d.setString(v, s)
		case startArray:
			return d.array(v)
		case startObject:
			return d.object(v)
		default:
			if b&0xfc == longSString {
				s, err := d.longSharedString(b)
				if err != nil {
					return err
				}
				return d.setString(v, s)
			}
		}
	}
	return fmt.Errorf("smile: unexpected value type %x", b)
}

func (d *decodeState) valueInterface(b byte) (interface{}, error) {
	switch b & 0xe0 {
	case 0x00:
		return d.sVals[b&0x1f-1], nil
	case 0x20:
		switch b {
		case emptyString:
			return "", nil
		case null:
			return nil, nil
		case falseTok:
			return false, nil
		case trueTok:
			return true, nil
		case int32Tok, int64Tok:
			return d.int(true)
		case bigInt:
			return d.bigInt()
		case float32Tok:
			return d.float32()
		case float64Tok:
			return d.float64()
		case bigDecimal:
			return d.bigDecimal()
		}
	case 0x40:
		return d.stringInterface(b, 1, &d.sVals)
	case 0x60:
		return d.stringInterface(b, 1+32, &d.sVals)
	case 0x80:
		return d.stringInterface(b, 2, &d.sVals)
	case 0xa0:
		return d.stringInterface(b, 2+32, &d.sVals)
	case 0xc0:
		return d.smallInt(b), nil
	case 0xe0:
		switch b {
		case longAscii, longUnicode:
			return d.longString()
		case startArray:
			return d.arrayInterface()
		case startObject:
			return d.objectInterface()
		default:
			if b&0xfc == longSString {
				return d.longSharedString(b)
			}
		}
	}
	return nil, fmt.Errorf("smile: unexpected value type %x", b)
}

func (d *decodeState) array(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface:
		if v.NumMethod() == 0 {
			i, err := d.arrayInterface()
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(i))
		}
	}
	return nil
}

func (d *decodeState) arrayInterface() ([]interface{}, error) {
	var v = make([]interface{}, 0)
	for {
		b, err := d.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == endArray {
			return v, nil
		}

		val, err := d.valueInterface(b)
		if err != nil {
			return nil, err
		}

		v = append(v, val)
	}
}

func (d *decodeState) key(b byte) (string, error) {
	switch {
	case b == 0x20:
		return "", nil
	case b == 0x34:
		return d.longKeyString()
	case 0x30 <= b && b < 0x34:
		b2, err := d.ReadByte()
		if err != nil {
			return "", err
		}
		i := int(b&0x03)<<8 | int(b2)
		return d.sKeys[i], nil
	case 0x40 <= b && b < 0x80:
		return d.sKeys[b&0x3f], nil
	case 0x80 <= b && b < 0xc0:
		return d.stringInterface(b, 1, &d.sKeys)
	case 0xc0 <= b && b < 0xf8:
		return d.stringInterface(b, 2, &d.sKeys)
	}
	return "", fmt.Errorf("smile: unexpected key type %x", b)
}

func (d *decodeState) object(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface:
		if v.NumMethod() == 0 {
			i, err := d.objectInterface()
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(i))
		}
	}
	return nil
}

func (d *decodeState) objectInterface() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for {
		b, err := d.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == endObject {
			return m, nil
		}

		key, err := d.key(b)
		if err != nil {
			return nil, err
		}

		b, err = d.ReadByte()
		if err != nil {
			return nil, err
		}

		val, err := d.valueInterface(b)
		if err != nil {
			return nil, err
		}

		m[key] = val
	}
}

func (d *decodeState) setString(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Interface:
		if v.NumMethod() == 0 {
			v.Set(reflect.ValueOf(s))
		}
	}
	return nil
}

func (d *decodeState) string(v reflect.Value, b byte, add byte, share *shared) error {
	s, err := d.stringInterface(b, add, share)
	if err != nil {
		return err
	}
	d.setString(v, s)
	return nil
}

func (d *decodeState) stringInterface(b byte, add byte, share *shared) (string, error) {
	buf := make([]byte, b&0x1f+add)
	_, err := io.ReadFull(d.r, buf)
	if err != nil {
		return "", err
	}
	s := string(buf)
	share.add(s)
	return s, nil
}

func (d *decodeState) longString() (string, error) {
	var s strings.Builder
	for {
		b, err := d.ReadByte()
		if err != nil {
			return "", err
		}
		if b == endString {
			return s.String(), nil
		}
		s.WriteByte(b)
	}
}

func (d *decodeState) longSharedString(b byte) (string, error) {
	b2, err := d.ReadByte()
	if err != nil {
		return "", err
	}
	i := int(b&0x03)<<8 | int(b2)
	return d.sVals[i], nil
}

func (d *decodeState) longKeyString() (string, error) {
	return "", errors.New("smile: not implemented: long key string")
}

func (d *decodeState) setInt(v reflect.Value, n int64) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(strconv.FormatInt(n, 10))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(n))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(n))
	case reflect.Interface:
		if v.NumMethod() == 0 {
			v.Set(reflect.ValueOf(n))
		}
	}
	return nil
}

func (d *decodeState) smallInt(b byte) int64 {
	return zigZagDecode(int64(b & 0x1f))
}

func (d *decodeState) int(signed bool) (int64, error) {
	var v int64
	for {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}

		if b >= 0x80 {
			v <<= 6
			v |= int64(b & 0x3f)
			if signed {
				v = zigZagDecode(v)
			}
			return v, nil
		}

		v <<= 7
		v |= int64(b)
	}
}

func (d *decodeState) safeBytes() ([]byte, error) {
	l, err := d.int(false)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, 0, l)
	var scratch, scratchL byte

	for {
		b, err := d.ReadByte()
		if err != nil {
			return nil, err
		}

		// make sure highest bit is zero
		b &= 0x7f

		if len(bytes) == cap(bytes)-1 && scratchL >= 1 {
			return append(bytes, scratch|b), nil
		}

		switch scratchL {
		case 0:
			scratch = b << 1
			scratchL = 7
		case 1:
			bytes = append(bytes, b|scratch)
			scratchL = 0
		default:
			scratchL--
			bytes = append(bytes, scratch|b>>scratchL)
			scratch = b << (8 - scratchL)
		}
	}
}

func (d *decodeState) bigInt() (*big.Int, error) {
	bytes, err := d.safeBytes()
	if err != nil {
		return nil, err
	}

	n := new(big.Int)
	if bytes[0]&0b10000000 != 0 {
		for i, b := range bytes {
			bytes[i] = ^b
		}
		n = n.SetBytes(bytes)
		n = n.Not(n)
	} else {
		n = n.SetBytes(bytes)
	}

	return n, nil
}

func (d *decodeState) float32() (float32, error) {
	return 0, errors.New("smile: not implemented: float32")
}

func (d *decodeState) float64() (float64, error) {
	var bits uint64
	for i := 0; i < 10; i++ {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}

		bits <<= 7
		bits |= uint64(b)
	}
	return math.Float64frombits(bits), nil
}

func (d *decodeState) bigDecimal() (*big.Float, error) {
	return nil, errors.New("smile: not implemented: big decimal")
}

func zigZagDecode(n int64) int64 {
	return (n >> 1) ^ (-(n & 1))
}

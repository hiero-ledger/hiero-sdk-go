// SPDX-License-Identifier: Apache-2.0

package hiero

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Decode decodes the input with a given type
func Decode(t *Type, input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	val, _, err := decode(t, input)
	return val, err
}

// DecodeStruct decodes the input with a type to a struct
func DecodeStruct(t *Type, input []byte, out interface{}) error {
	val, err := Decode(t, input)
	if err != nil {
		return err
	}

	dc := &mapstructure.DecoderConfig{
		Result:           out,
		WeaklyTypedInput: true,
		TagName:          "abi",
	}
	ms, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}
	if err = ms.Decode(val); err != nil {
		return err
	}
	return nil
}

func decode(t *Type, input []byte) (interface{}, []byte, error) {
	var data []byte
	var length int
	var err error

	// safe check, input should be at least 32 bytes
	if len(input) < 32 {
		return nil, nil, fmt.Errorf("incorrect length")
	}

	if t.isVariableInput() {
		length, err = readLength(input)
		if err != nil {
			return nil, nil, err
		}
	} else {
		data = input[:32]
	}

	switch t.kind {
	case KindTuple:
		return decodeTuple(t, input)

	case KindSlice:
		return decodeArraySlice(t, input[32:], length)

	case KindArray:
		return decodeArraySlice(t, input, t.size)
	}

	var val interface{}
	switch t.kind {
	case KindBool:
		val, err = decodeBool(data)

	case KindInt, KindUInt:
		val = readInteger(t, data)

	case KindString:
		val = string(input[32 : 32+length])

	case KindBytes:
		val = input[32 : 32+length]

	case KindAddress:
		val, err = readAddr(data)

	case KindFixedBytes:
		val, err = readFixedBytes(t, data)

	case KindFunction:
		val, err = readFunctionType(t, data)

	default:
		return nil, nil, fmt.Errorf("decoding not available for type '%s'", t.kind)
	}

	return val, input[32:], err
}

var (
	maxUint256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(-1))
	maxInt256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(255), nil),
		big.NewInt(-1))
)

// Address is an Ethereum address
type Address [20]byte

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

// BytesToAddress converts bytes to an address object
func BytesToAddress(b []byte) Address {
	var a Address

	size := len(b)
	min := min(size, 20)

	copy(a[20-min:], b[len(b)-min:])
	return a
}

// Address implements the ethgo.Key interface Address method.
func (a Address) Address() Address {
	return a
}

// UnmarshalText implements the unmarshal interface
func (a *Address) UnmarshalText(b []byte) error {
	return unmarshalTextByte(a[:], b, 20)
}

// MarshalText implements the marshal interface
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// Bytes returns the bytes of the Address
func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) String() string {
	return a.checksumEncode()
}

func unmarshalTextByte(dst, src []byte, size int) error {
	str := string(src)

	str = strings.Trim(str, "\"")
	if !strings.HasPrefix(str, "0x") {
		return fmt.Errorf("0x prefix not found")
	}
	str = str[2:]
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != size {
		return fmt.Errorf("length %d is not correct, expected %d", len(b), size)
	}
	copy(dst, b)
	return nil
}

func (a Address) checksumEncode() string {
	address := strings.ToLower(hex.EncodeToString(a[:]))

	hash := hex.EncodeToString(Keccak256Hash([]byte(address)).Bytes())

	ret := "0x"
	for i := 0; i < len(address); i++ {
		character := string(address[i])

		num, _ := strconv.ParseInt(string(hash[i]), 16, 64)
		if num > 7 {
			ret += strings.ToUpper(character)
		} else {
			ret += character
		}
	}

	return ret
}

func readAddr(b []byte) (Address, error) {
	res := Address{}
	if len(b) != 32 {
		return res, fmt.Errorf("len is not correct")
	}
	copy(res[:], b[12:])
	return res, nil
}

func readInteger(t *Type, b []byte) interface{} {
	switch t.t.Kind() {
	case reflect.Uint8:
		return b[len(b)-1]

	case reflect.Uint16:
		return binary.BigEndian.Uint16(b[len(b)-2:])

	case reflect.Uint32:
		return binary.BigEndian.Uint32(b[len(b)-4:])

	case reflect.Uint64:
		return binary.BigEndian.Uint64(b[len(b)-8:])

	case reflect.Int8:
		return int8(b[len(b)-1])

	case reflect.Int16:
		return int16(binary.BigEndian.Uint16(b[len(b)-2:]))

	case reflect.Int32:
		return int32(binary.BigEndian.Uint32(b[len(b)-4:]))

	case reflect.Int64:
		return int64(binary.BigEndian.Uint64(b[len(b)-8:]))

	default:
		ret := new(big.Int).SetBytes(b)
		if t.kind == KindUInt {
			return ret
		}

		if ret.Cmp(maxInt256) > 0 {
			ret.Add(maxUint256, big.NewInt(0).Neg(ret))
			ret.Add(ret, big.NewInt(1))
			ret.Neg(ret)
		}
		return ret
	}
}

// nolint
func readFunctionType(t *Type, word []byte) ([24]byte, error) {
	res := [24]byte{}
	if !allZeros(word[24:32]) {
		return res, fmt.Errorf("function type expects the last 8 bytes to be empty but found: %b", word[24:32])
	}
	copy(res[:], word[0:24])
	return res, nil
}

// nolint
func readFixedBytes(t *Type, word []byte) (interface{}, error) {
	array := reflect.New(t.t).Elem()
	reflect.Copy(array, reflect.ValueOf(word[0:t.size]))
	return array.Interface(), nil
}

func decodeTuple(t *Type, data []byte) (interface{}, []byte, error) {
	res := make(map[string]interface{})

	orig := data
	origLen := len(orig)
	for indx, arg := range t.tuple {
		if len(data) < 32 {
			return nil, nil, fmt.Errorf("incorrect length")
		}

		entry := data
		if arg.Elem.isDynamicType() {
			offset, err := readOffset(data, origLen)
			if err != nil {
				return nil, nil, err
			}
			entry = orig[offset:]
		}

		val, tail, err := decode(arg.Elem, entry)
		if err != nil {
			return nil, nil, err
		}

		if !arg.Elem.isDynamicType() {
			data = tail
		} else {
			data = data[32:]
		}

		name := arg.Name
		if name == "" {
			name = strconv.Itoa(indx)
		}
		if _, ok := res[name]; !ok {
			res[name] = val
		} else {
			return nil, nil, fmt.Errorf("tuple with repeated values")
		}
	}
	return res, data, nil
}

func decodeArraySlice(t *Type, data []byte, size int) (interface{}, []byte, error) {
	if size < 0 {
		return nil, nil, fmt.Errorf("size is lower than zero")
	}
	if 32*size > len(data) {
		return nil, nil, fmt.Errorf("size is too big")
	}

	var res reflect.Value
	switch t.kind {
	case KindSlice:
		res = reflect.MakeSlice(t.t, size, size)
	case KindArray:
		res = reflect.New(t.t).Elem()
	}

	orig := data
	origLen := len(orig)
	for indx := 0; indx < size; indx++ {
		isDynamic := t.elem.isDynamicType()

		if len(data) < 32 {
			return nil, nil, fmt.Errorf("incorrect length")
		}

		entry := data
		if isDynamic {
			offset, err := readOffset(data, origLen)
			if err != nil {
				return nil, nil, err
			}
			entry = orig[offset:]
		}

		val, tail, err := decode(t.elem, entry)
		if err != nil {
			return nil, nil, err
		}

		if !isDynamic {
			data = tail
		} else {
			data = data[32:]
		}
		res.Index(indx).Set(reflect.ValueOf(val))
	}
	return res.Interface(), data, nil
}

func decodeBool(data []byte) (interface{}, error) {
	switch data[31] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("bad boolean")
	}
}

func readOffset(data []byte, len int) (int, error) {
	offsetBig := big.NewInt(0).SetBytes(data[0:32])
	if offsetBig.BitLen() > 63 {
		return 0, fmt.Errorf("offset larger than int64: %v", offsetBig.Int64())
	}
	offset := int(offsetBig.Int64())
	if offset > len {
		return 0, fmt.Errorf("offset insufficient %v require %v", len, offset)
	}
	return offset, nil
}

func readLength(data []byte) (int, error) {
	lengthBig := big.NewInt(0).SetBytes(data[0:32])
	if lengthBig.BitLen() > 63 {
		return 0, fmt.Errorf("length larger than int64: %v", lengthBig.Int64())
	}
	length := int(lengthBig.Uint64())

	// if we trim the length in the data there should be enough
	// bytes to cover the length
	if length > len(data)-32 {
		return 0, fmt.Errorf("length insufficient %v require %v", len(data), length)
	}
	return length, nil
}

func allZeros(b []byte) bool {
	for _, i := range b {
		if i != 0 {
			return false
		}
	}
	return true
}

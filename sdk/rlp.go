package hiero

import "encoding/hex"

// SPDX-License-Identifier: Apache-2.0

// RLPType represents the type of RLP item.
type RLPType int

const (
	VALUE_TYPE RLPType = iota
	LIST_TYPE
)

// RLPItem represents a single RLP item.
type RLPItem struct {
	itemType   RLPType    // Type of the RLP item (value or list)
	itemValue  []byte     // Holds the byte value for value-type items
	childItems []*RLPItem // Holds child items for list-type items
}

// NewRLPItem creates a new RLP item.
func NewRLPItem(typ RLPType) *RLPItem {
	return &RLPItem{itemType: typ}
}

// EncodeBinary encodes a number into a byte slice.
func encodeBinary(num uint64) []byte {
	var bytes []byte
	for num != 0 {
		bytes = append([]byte{byte(num & 0xFF)}, bytes...)
		num >>= 8
	}
	return bytes
}

// EncodeLength encodes the length of an item with an offset.
func encodeLength(num int, offset byte) []byte {
	if num < 56 {
		return []byte{offset + byte(num)}
	}
	encodedLength := encodeBinary(uint64(num))
	return append([]byte{byte(len(encodedLength) + int(offset) + 55)}, encodedLength...)
}

// Gets Child items for RLPItem
func (item *RLPItem) GetChildItems() []*RLPItem {
	return item.childItems
}

// Gets the item value for RLPItem
func (item *RLPItem) GetItemValue() []byte {
	return item.itemValue
}

// Assign methods to set values for the RLPItem
func (item *RLPItem) AssignValue(value []byte) *RLPItem {
	item.itemType = VALUE_TYPE
	if hex.EncodeToString(value) == "00" {
		item.itemValue = []byte{}
	} else {
		item.itemValue = value
	}
	return item
}

func (item *RLPItem) AssignBytes(value []byte) {
	item.itemType = VALUE_TYPE
	item.itemValue = value
}

func (item *RLPItem) AssignString(value string) {
	item.AssignBytes([]byte(value))
}

func (item *RLPItem) AssignList() {
	item.itemType = LIST_TYPE
}

// Clear resets the item values.
func (item *RLPItem) Clear() {
	item.itemValue = nil
	item.childItems = nil
}

// PushBack adds a value to a list item.
func (item *RLPItem) PushBack(child *RLPItem) {
	item.childItems = append(item.childItems, child)
}

// Size returns the size of the RLPItem.
func (item *RLPItem) Size() int {
	if item.itemType == VALUE_TYPE {
		return len(item.itemValue)
	}
	size := 0
	for _, child := range item.childItems {
		size += child.Size()
	}
	return size
}

// Write encodes the RLPItem to a byte slice.
func (item *RLPItem) Write() ([]byte, error) {
	if item.itemType == VALUE_TYPE {
		if len(item.itemValue) == 1 && item.itemValue[0] < 0x80 {
			return item.itemValue, nil
		}
		return append(encodeLength(len(item.itemValue), 0x80), item.itemValue...), nil
	}

	var bytes []byte
	for _, child := range item.childItems {
		childBytes, err := child.Write()
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, childBytes...)
	}
	return append(encodeLength(len(bytes), 0xC0), bytes...), nil
}

// Read decodes a byte slice into an RLPItem.
func (item *RLPItem) Read(bytes []byte) error {
	item.Clear()

	if len(bytes) == 0 {
		return nil
	}

	index := 0
	return item.decodeBytes(bytes, &index)
}

// decodeBytes decodes the bytes starting from the given index.
func (item *RLPItem) decodeBytes(bytes []byte, index *int) error {
	if *index >= len(bytes) {
		return errRLPUnexpectedEnd
	}
	prefix := bytes[*index]
	(*index)++

	// Single byte case
	if prefix < 0x80 {
		item.itemValue = []byte{prefix}
		item.itemType = VALUE_TYPE
		return nil
	}

	// Short string case
	if prefix < 0xB8 {
		stringLength := int(prefix) - 0x80
		if stringLength > len(bytes)-*index {
			return errRLPUnexpectedEnd
		}
		item.itemValue = bytes[*index : *index+stringLength]
		item.itemType = VALUE_TYPE
		*index += stringLength
		return nil
	}

	// Long string case
	if prefix < 0xC0 {
		stringLength, err := decodeLength(bytes, index, int(prefix)-0xB7)
		if err != nil {
			return err
		}
		item.itemValue = bytes[*index : *index+stringLength]
		item.itemType = VALUE_TYPE
		*index += stringLength
		return nil
	}

	// Short list case
	if prefix < 0xF8 {
		listLength := int(prefix) - 0xC0
		if listLength > len(bytes)-*index {
			return errRLPUnexpectedEnd
		}
		return item.decodeList(bytes, index, listLength)
	}

	// Long list case
	listLength, err := decodeLength(bytes, index, int(prefix)-0xF7)
	if err != nil {
		return err
	}
	return item.decodeList(bytes, index, listLength)
}

// decodeList decodes the child items of a list whose payload starts at the given index.
// Child items are not allowed to extend past the end of the list payload.
func (item *RLPItem) decodeList(bytes []byte, index *int, listLength int) error {
	end := *index + listLength
	for *index < end {
		childItem := NewRLPItem(LIST_TYPE)
		if err := childItem.decodeBytes(bytes[:end], index); err != nil {
			return err
		}
		item.PushBack(childItem)
	}
	item.itemType = LIST_TYPE
	return nil
}

// decodeLength reads a big-endian length of lengthLength bytes starting at the given index
// and checks that an item of that length fits in the remaining bytes.
func decodeLength(bytes []byte, index *int, lengthLength int) (int, error) {
	if lengthLength > len(bytes)-*index {
		return 0, errRLPUnexpectedEnd
	}
	var length uint64
	for i := 0; i < lengthLength; i++ {
		length = (length << 8) | uint64(bytes[*index])
		(*index)++
	}
	if length > uint64(len(bytes)-*index) {
		return 0, errRLPUnexpectedEnd
	}
	return int(length), nil
}

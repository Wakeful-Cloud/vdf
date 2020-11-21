package vdf

import (
	"encoding/binary"
	"errors"
)

//buffer represents a method of storing binary data
type buffer struct {
	//Data represents the raw binary data
	Data []byte

	//Position within the buffer
	Position uint32
}

//Map represents a VDF file map
type Map map[string]interface{}

const (
	//vdfMapStart type
	vdfMapStart byte = 0x00

	//vdfString type
	vdfString = 0x01

	//vdfNumber type
	vdfNumber = 0x02

	//vdfMapEnd type
	vdfMapEnd = 0x08
)

//mapItem represents a VDF map item
type mapItem struct {
	Type  byte
	Name  string
	Value interface{}
}

//ReadVdf reads a binary VDF file
func ReadVdf(data []byte) (Map, error) {
	//Initial buffer
	buffer := buffer{
		Data:     data,
		Position: 0,
	}

	return nextMap(&buffer)
}

//nextMap returns the next map
func nextMap(buffer *buffer) (Map, error) {
	contents := Map{}

	//Iterate over maps
	for {
		//Get map item
		mapItem, err := nextMapItem(buffer)

		if err != nil {
			return nil, err
		}

		if mapItem.Type == vdfMapEnd {
			break
		} else {
			contents[mapItem.Name] = mapItem.Value
		}
	}

	return contents, nil
}

//nextMapItem returns the next map item
func nextMapItem(buffer *buffer) (mapItem, error) {
	//Get item type
	typeByte := buffer.Data[buffer.Position]
	buffer.Position++

	//MapEnd
	if typeByte == vdfMapEnd {
		//Construct item
		item := mapItem{
			Type: vdfMapEnd,
		}

		return item, nil
	}

	//Get name
	name, err := nextStringZero(buffer)

	if err != nil {
		return mapItem{}, err
	}

	//Get value
	var value interface{}
	switch typeByte {
	case vdfMapStart:
		value, err = nextMap(buffer)

		if err != nil {
			return mapItem{}, err
		}
	case vdfString:
		value, err = nextStringZero(buffer)

		if err != nil {
			return mapItem{}, err
		}
	case vdfNumber:
		//Get slice
		positionInt := int(buffer.Position)
		slice := buffer.Data[positionInt : positionInt+4]

		//Convert to number
		value = binary.LittleEndian.Uint32(slice)
		buffer.Position += 4
	default:
		return mapItem{}, errors.New("unrecognized VDF type")
	}

	//Construct item
	item := mapItem{
		Type:  typeByte,
		Name:  name,
		Value: value,
	}

	return item, nil
}

//nextStringZero returns the next string (Up until the first 'NUL')
func nextStringZero(buffer *buffer) (string, error) {
	//Prevent buffer overflow
	if buffer.Position >= uint32(len(buffer.Data)) {
		return "", errors.New("aborted buffer overflow")
	}

	//Keep reference to start of string
	start := buffer.Position

	//Get the end position (NUL terminator)
	end := buffer.Position
	for {
		//Get current character
		char := buffer.Data[buffer.Position]
		buffer.Position++

		if char == 0 {
			break
		} else {
			end++
		}
	}

	return string(buffer.Data[start:end]), nil
}

//WriteVdf writes a binary VDF file
func WriteVdf(vdfMap Map) ([]byte, error) {
	return addMap(vdfMap)
}

//addMap adds a map
func addMap(vdfMap Map) ([]byte, error) {
	buffer := []byte{}

	//Iterate over keys
	for k, v := range vdfMap {
		switch v.(type) {
		case uint32:
			//Add KT
			kt, err := addKT(vdfNumber, k)

			if err != nil {
				return nil, err
			}

			buffer = append(buffer, kt...)

			//Convert to number
			value := v.(uint32)

			//Convert to bytes
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, value)

			buffer = append(buffer, bytes...)
		case string:
			//Add KT
			kt, err := addKT(vdfString, k)

			if err != nil {
				return nil, err
			}

			buffer = append(buffer, kt...)

			//Convert to string
			value := v.(string)

			//Convert to string byte array
			bytes, err := addString(value)

			if err != nil {
				return nil, err
			}

			buffer = append(buffer, bytes...)
		case Map:
			//Add KT
			kt, err := addKT(vdfMapStart, k)

			if err != nil {
				return nil, err
			}

			buffer = append(buffer, kt...)

			//Convert to VdfMap
			value := v.(Map)

			//Add map
			bytes, err := addMap(value)

			if err != nil {
				return nil, err
			}

			buffer = append(buffer, bytes...)
		default:
			return nil, errors.New("unrecognized go type")
		}
	}

	//Add map end
	buffer = append(buffer, vdfMapEnd)

	return buffer, nil
}

func addString(value string) ([]byte, error) {
	//Convert value to bytes
	bytes := []byte(value)

	//Ensure no nulls's in the string
	for _, v := range bytes {
		if v == 0 {
			return nil, errors.New("NUL terminator found in key")
		}
	}

	//Append null terminator
	bytes = append(bytes, 0)

	return bytes, nil
}

//addKT adds a key-type pair (Must manually add value)
func addKT(Type byte, key string) ([]byte, error) {
	buffer := []byte{}

	//Add type
	buffer = append(buffer, Type)

	//Add key
	keyBytes, err := addString(key)

	if err != nil {
		return nil, err
	}

	buffer = append(buffer, keyBytes...)

	return buffer, nil
}

package dht

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"reflect"

	"github.com/anacrolix/missinggo/v2/slices"
	"github.com/anacrolix/torrent/bencode"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type NodeInfo struct {
	ID   protocol.ID
	Addr NodeAddr
}

func (ni NodeInfo) String() string {
	return fmt.Sprintf("{%x at %s}", ni.ID, ni.Addr)
}

func RandomNodeInfo(ipLen int) (ni NodeInfo) {
	ni.ID = protocol.RandomNodeID()
	ni.Addr.Port = rand.Intn(math.MaxUint16 + 1)
	ni.Addr.IP = make(net.IP, ipLen)
	for i := 0; i < ipLen; i++ {
		ni.Addr.IP[i] = byte(rand.Intn(256))
	}
	return
}

var _ interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
} = (*NodeInfo)(nil)

func (ni NodeInfo) MarshalBinary() ([]byte, error) {
	var w bytes.Buffer
	w.Write(ni.ID[:])
	w.Write(ni.Addr.IP)
	if err := binary.Write(&w, binary.BigEndian, uint16(ni.Addr.Port)); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (ni *NodeInfo) UnmarshalBinary(b []byte) error {
	copy(ni.ID[:], b)
	return ni.Addr.UnmarshalBinary(b[20:])
}

type (
	CompactIPv4NodeInfo []NodeInfo
)

func (CompactIPv4NodeInfo) ElemSize() int {
	return 26
}

// func (me *CompactIPv4NodeInfo) Scrub() {
// 	slices.FilterInPlace(me, func(ni *NodeInfo) bool {
// 		ni.Addr.IP = ni.Addr.IP.To4()
// 		return ni.Addr.IP != nil
// 	})
// }

func (ni CompactIPv4NodeInfo) MarshalBinary() ([]byte, error) {
	return marshalBinarySlice(slices.Map(func(ni NodeInfo) NodeInfo {
		ni.Addr.IP = ni.Addr.IP.To4()
		return ni
	}, ni).(CompactIPv4NodeInfo))
}

func (ni CompactIPv4NodeInfo) MarshalBencode() ([]byte, error) {
	return bencodeBytesResult(ni.MarshalBinary())
}

func (ni *CompactIPv4NodeInfo) UnmarshalBinary(b []byte) error {
	return unmarshalBinarySlice(ni, b)
}

func (ni *CompactIPv4NodeInfo) UnmarshalBencode(b []byte) error {
	return unmarshalBencodedBinary(ni, b)
}

func unmarshalBencodedBinary(u encoding.BinaryUnmarshaler, b []byte) (err error) {
	var ub string
	err = bencode.Unmarshal(b, &ub)
	if err != nil {
		return
	}
	return u.UnmarshalBinary([]byte(ub))
}

type elemSizer interface {
	ElemSize() int
}

func unmarshalBinarySlice(slice elemSizer, b []byte) (err error) {
	sliceValue := reflect.ValueOf(slice).Elem()
	elemType := sliceValue.Type().Elem()
	bytesPerElem := slice.ElemSize()
	elem := reflect.New(elemType)
	for len(b) != 0 {
		if len(b) < bytesPerElem {
			err = fmt.Errorf("%d trailing bytes < %d required for element", len(b), bytesPerElem)
			break
		}
		if bu, ok := elem.Interface().(encoding.BinaryUnmarshaler); ok {
			err = bu.UnmarshalBinary(b[:bytesPerElem])
		} else if elem.Elem().Len() == bytesPerElem {
			reflect.Copy(elem.Elem(), reflect.ValueOf(b[:bytesPerElem]))
		} else {
			err = fmt.Errorf("can't unmarshal %v bytes into %v", bytesPerElem, elem.Type())
		}
		if err != nil {
			return
		}
		sliceValue.Set(reflect.Append(sliceValue, elem.Elem()))
		b = b[bytesPerElem:]
	}
	return
}

func marshalBinarySlice(slice elemSizer) (ret []byte, err error) {
	var elems []encoding.BinaryMarshaler
	makeInto(&elems, slice)
	for _, e := range elems {
		var b []byte
		b, err = e.MarshalBinary()
		if err != nil {
			return
		}
		if len(b) != slice.ElemSize() {
			panic(fmt.Sprintf("marshalled %d bytes, but expected %d", len(b), slice.ElemSize()))
		}
		ret = append(ret, b...)
	}
	return
}

func bencodeBytesResult(b []byte, err error) ([]byte, error) {
	if err != nil {
		return b, err
	}
	return bencode.Marshal(b)
}

// makes and sets a slice at *ptrTo, and type asserts all the elements from "from" to it.
func makeInto(ptrTo interface{}, from interface{}) {
	fromSliceValue := reflect.ValueOf(from)
	fromLen := fromSliceValue.Len()
	if fromLen == 0 {
		return
	}
	// Deref the pointer to slice.
	slicePtrValue := reflect.ValueOf(ptrTo)
	if slicePtrValue.Kind() != reflect.Ptr {
		panic("destination is not a pointer")
	}
	destSliceValue := slicePtrValue.Elem()
	// The type of the elements of the destination slice.
	destSliceElemType := destSliceValue.Type().Elem()
	destSliceValue.Set(reflect.MakeSlice(destSliceValue.Type(), fromLen, fromLen))
	for i := range make([]struct{}, fromSliceValue.Len()) {
		// The value inside the interface in the slice element.
		itemValue := fromSliceValue.Index(i)
		if itemValue.Kind() == reflect.Interface {
			itemValue = itemValue.Elem()
		}
		convertedItem := itemValue.Convert(destSliceElemType)
		destSliceValue.Index(i).Set(convertedItem)
	}
}

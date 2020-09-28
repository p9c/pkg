package simplebuffer

import (
	"encoding/binary"

	"github.com/p9c/pkg/app/slog"
)

type Serializer interface {
	// Encode returns the wire/storage form of the data
	Encode() []byte
	// Decode stores the decoded data from the head of the slice and returns the remainder
	Decode(b []byte) []byte
}

type Serializers []Serializer

type Container struct {
	Data []byte
}

// CreateContainer takes an array of serializer interface objects and renders the data into bytes
func (srs Serializers) CreateContainer(magic []byte) (out *Container) {
	if len(magic) != 4 {
		slog.Error("magic must be 4 bytes")
		return
	}
	out = &Container{}
	var offset uint32
	var length uint16
	var nodes []uint32
	for i := range srs {
		b := srs[i].Encode()
		// Debug(i, len(b), hex.EncodeToString(b))
		length++
		nodes = append(nodes, offset)
		offset += uint32(len(b))
		out.Data = append(out.Data, b...)
	}
	// L.Spew(out.Data)
	// Debug(offset, len(out.Data))
	nodeB := make([]byte, len(nodes)*4+2)
	start := uint32(len(nodeB) + 8)
	binary.BigEndian.PutUint16(nodeB[:2], length)
	for i := range nodes {
		b := nodeB[i*4+2 : i*4+4+2]
		binary.BigEndian.PutUint32(b, nodes[i]+start)
		// Debug(i, len(b), hex.EncodeToString(b))
	}
	// L.Spew(nodeB)
	out.Data = append(nodeB, out.Data...)
	size := offset + uint32(len(nodeB)) + 8
	// Debug("size", size, len(out.Data))
	sB := make([]byte, 4)
	binary.BigEndian.PutUint32(sB, size)
	out.Data = append(append(magic[:], sB...), out.Data...)
	return
}

func (c *Container) Count() uint16 {
	size := binary.BigEndian.Uint32(c.Data[4:8])
	// Debug("size", size)
	if len(c.Data) >= int(size) {
		// we won't touch it if it's not at least as big so we don't get
		// bounds errors
		return binary.BigEndian.Uint16(c.Data[8:10])
	}
	return 0
}

func (c *Container) GetMagic() (out []byte) {
	return c.Data[:4]
}

// Get returns the bytes that can be imported into an interface assuming the types are correct - field ordering is hard
// coded by the creation and identified by the magic. This is all read only and subslices so it should generate very
// little garbage or copy operations except as required for the output (we aren't going to go unsafe here, it isn't
// really necessary since already this library enables avoiding the decoding of values not being used from a message (or
// not used yet)
func (c *Container) Get(idx uint16) (out []byte) {
	length := c.Count()
	size := len(c.Data)
	if length > idx {
		// Debug("length", length, "idx", idx)
		if idx < length {
			offset := binary.BigEndian.Uint32(c.
				Data[10+idx*4 : 10+idx*4+4])
			// Debug("offset", offset)
			if idx < length-1 {
				nextOffset := binary.BigEndian.Uint32(c.
					Data[10+((idx+1)*4) : 10+((idx+1)*4)+4])
				// Debug("nextOffset", nextOffset)
				out = c.Data[offset:nextOffset]
			} else {
				nextOffset := len(c.Data)
				// Debug("last nextOffset", nextOffset)
				out = c.Data[offset:nextOffset]
			}
		}
	} else {
		slog.Error("size mismatch", length, size)
	}
	return
}

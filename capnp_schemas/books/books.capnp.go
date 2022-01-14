// Code generated by capnpc-go. DO NOT EDIT.

package books

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Book struct{ capnp.Struct }

// Book_TypeID is the unique identifier for the type Book.
const Book_TypeID = 0x8100cc88d7d4d47c

func NewBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Book{st}, err
}

func NewRootBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Book{st}, err
}

func ReadRootBook(msg *capnp.Message) (Book, error) {
	root, err := msg.Root()
	return Book{root.Struct()}, err
}

func (s Book) String() string {
	str, _ := text.Marshal(0x8100cc88d7d4d47c, s.Struct)
	return str
}

func (s Book) Title() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Book) HasTitle() bool {
	return s.Struct.HasPtr(0)
}

func (s Book) TitleBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Book) SetTitle(v string) error {
	return s.Struct.SetText(0, v)
}

func (s Book) PageCount() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Book) SetPageCount(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

// Book_List is a list of Book.
type Book_List struct{ capnp.List }

// NewBook creates a new list of Book.
func NewBook_List(s *capnp.Segment, sz int32) (Book_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Book_List{l}, err
}

func (s Book_List) At(i int) Book { return Book{s.List.Struct(i)} }

func (s Book_List) Set(i int, v Book) error { return s.List.SetStruct(i, v.Struct) }

func (s Book_List) String() string {
	str, _ := text.MarshalList(0x8100cc88d7d4d47c, s.List)
	return str
}

// Book_Future is a wrapper for a Book promised by a client call.
type Book_Future struct{ *capnp.Future }

func (p Book_Future) Struct() (Book, error) {
	s, err := p.Future.Struct()
	return Book{s}, err
}

const schema_85d3acc39d94e0f8 = "x\xda\x12Ht`2d\xdd\xcf\xc8\xc0\x10(\xc2\xca" +
	"\xb6\xbf\xe6\xca\x95\xeb\x1dg\x1a\x03\x15\x18\x19\xff\xffx" +
	"0e\xee\xe15\x97[\x19X\x19\xd9\x19\x18\x0c\x8fj" +
	"1\x0a^eg`\x10\xbcX\xce\xa0\xfb?9\xb1 " +
	"\xaf \xbe8\x999#57\xb1X?)??\x1b" +
	"J\xea\x81\xa5\xac\xf8\x9d\xf2\xf3\xb3\x03\x18\x19\x039\x98" +
	"Y\x18\x18X\x18\x19\x18\x045\x8d\x18\x18\x02U\x98\x19" +
	"\x03\x0d\x98\x18\x19\x19E\x18Ab\xbaA\x0c\x0c\x81:" +
	"\xcc\x8c\x81\x16L\x8c\xf2%\x99%9\xa9\x8c<\x0cL" +
	"\x8c<\x0c\x8c\xff\x0b\x12\xd3S\x9d\xf3K\xf3\x18\x18K" +
	"\x18Y\x18\x98\x18Y\x18\x18\x01\x01\x00\x00\xff\xff\x0aQ" +
	",w"

func init() {
	schemas.Register(schema_85d3acc39d94e0f8,
		0x8100cc88d7d4d47c)
}

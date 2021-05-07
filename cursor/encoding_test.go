package cursor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestEncoding(t *testing.T) {
	suite.Run(t, &encodingSuite{})
}

type encodingSuite struct {
	suite.Suite
}

/* bool */

type boolModel struct {
	Value    bool
	ValuePtr *bool
}

func (s *encodingSuite) TestBool() {
	c, _ := s.encodeValue(boolModel{Value: true})
	v, _ := s.decodeValue(boolModel{}, c)
	s.Equal(true, v)
}

func (s *encodingSuite) TestBoolPtr() {
	b := true
	c, _ := s.encodeValuePtr(boolModel{ValuePtr: &b})
	v, _ := s.decodeValuePtr(boolModel{}, c)
	s.Equal(true, *(v.(*bool)))
}

/* int */

type intModel struct {
	Value    int
	ValuePtr *int
}

func (s *encodingSuite) TestInt() {
	c, _ := s.encodeValue(intModel{Value: 123})
	v, _ := s.decodeValue(intModel{}, c)
	s.Equal(int(123), v)
}

func (s *encodingSuite) TestIntPtr() {
	i := 123
	c, _ := s.encodeValuePtr(intModel{ValuePtr: &i})
	v, _ := s.decodeValuePtr(intModel{}, c)
	s.Equal(int(123), *(v.(*int)))
}

/* uint */

type uintModel struct {
	Value    uint
	ValuePtr *uint
}

func (s *encodingSuite) TestUint() {
	c, _ := s.encodeValue(uintModel{Value: 123})
	v, _ := s.decodeValue(uintModel{}, c)
	s.Equal(uint(123), v)
}

func (s *encodingSuite) TestUintPtr() {
	ui := uint(123)
	c, _ := s.encodeValuePtr(uintModel{ValuePtr: &ui})
	v, _ := s.decodeValuePtr(uintModel{}, c)
	s.Equal(uint(123), *(v.(*uint)))
}

/* float */
type floatModel struct {
	Value    float64
	ValuePtr *float64
}

func (s *encodingSuite) TestFloat() {
	c, _ := s.encodeValue(floatModel{Value: 123.45})
	v, _ := s.decodeValue(floatModel{}, c)
	s.Equal(float64(123.45), v)
}

func (s *encodingSuite) TestFloatPtr() {
	f := 123.45
	c, _ := s.encodeValuePtr(floatModel{ValuePtr: &f})
	v, _ := s.decodeValuePtr(floatModel{}, c)
	s.Equal(float64(123.45), *(v.(*float64)))
}

/* string */

type stringModel struct {
	Value    string
	ValuePtr *string
}

func (s *encodingSuite) TestString() {
	c, _ := s.encodeValue(stringModel{Value: "hello"})
	v, _ := s.decodeValue(stringModel{}, c)
	s.Equal("hello", v)
}

func (s *encodingSuite) TestStringPtr() {
	str := "hello"
	c, _ := s.encodeValuePtr(stringModel{ValuePtr: &str})
	v, _ := s.decodeValuePtr(stringModel{}, c)
	s.Equal("hello", *(v.(*string)))
}

/* time */

type timeModel struct {
	Value    time.Time
	ValuePtr *time.Time
}

func (s *encodingSuite) TestTime() {
	t := time.Now()
	c, _ := s.encodeValue(timeModel{Value: t})
	v, _ := s.decodeValue(timeModel{}, c)
	s.Equal(t.Second(), v.(time.Time).Second())
}

func (s *encodingSuite) TestTimePtr() {
	t := time.Now()
	c, _ := s.encodeValuePtr(timeModel{ValuePtr: &t})
	v, _ := s.decodeValuePtr(timeModel{}, c)
	s.Equal(t.Second(), v.(*time.Time).Second())
}

/* struct */

type structModel struct {
	Value    structValue
	ValuePtr *structValue
}

type structValue struct {
	Value []byte
}

func (s *encodingSuite) TestStruct() {
	c, _ := s.encodeValue(structModel{
		Value: structValue{Value: []byte("123")},
	})
	v, _ := s.decodeValue(structModel{}, c)
	s.Equal(structValue{Value: []byte("123")}, v)
}

func (s *encodingSuite) TestStructPtr() {
	sv := structValue{Value: []byte("123")}
	c, _ := s.encodeValuePtr(structModel{ValuePtr: &sv})
	v, _ := s.decodeValuePtr(structModel{}, c)
	s.Equal(sv, *(v.(*structValue)))
}

/* multiple */

type multipleModel struct {
	ID        int
	Name      string
	CreatedAt *time.Time
}

func (multipleModel) Keys() []string {
	return []string{
		"ID",
		"Name",
		"CreatedAt",
	}
}

func (s *encodingSuite) TestMultipleFields() {
	keys := multipleModel{}.Keys()

	t := time.Now()
	c, err := NewEncoder(keys...).Encode(multipleModel{
		ID:        123,
		Name:      "Hello",
		CreatedAt: &t,
	})
	s.Nil(err)

	fields, err := NewDecoder(keys...).Decode(c, multipleModel{})
	s.Nil(err)

	s.Len(fields, 3)
	s.Equal(123, fields[0])
	s.Equal("Hello", fields[1])
	s.Equal(t.Second(), fields[2].(*time.Time).Second())
}

func (s *encoderSuite) TestMultipleFieldsWithZeroValue() {
	keys := multipleModel{}.Keys()

	c, err := NewEncoder(keys...).Encode(multipleModel{})
	s.Nil(err)

	fields, err := NewDecoder(keys...).Decode(c, multipleModel{})
	s.Nil(err)

	s.Equal(0, fields[0])
	s.Equal("", fields[1])
	s.Equal((*time.Time)(nil), fields[2])
}

/* decode struct */

func (s *encodingSuite) TestMultipleFieldsToStruct() {
	keys := multipleModel{}.Keys()

	t := time.Now()
	c, err := NewEncoder(keys...).Encode(multipleModel{
		ID:        123,
		Name:      "Hello",
		CreatedAt: &t,
	})
	s.Nil(err)

	var model multipleModel
	err = NewDecoder(keys...).DecodeStruct(c, &model)
	s.Nil(err)

	s.Equal(123, model.ID)
	s.Equal("Hello", model.Name)
	s.Equal(t.Second(), (*model.CreatedAt).Second())
}

func (s *encoderSuite) TestMultipleFieldsToStructWithZeroValue() {
	keys := multipleModel{}.Keys()

	c, err := NewEncoder(keys...).Encode(multipleModel{})
	s.Nil(err)

	var model multipleModel
	err = NewDecoder(keys...).DecodeStruct(c, &model)
	s.Nil(err)

	s.Equal(0, model.ID)
	s.Equal("", model.Name)
	s.Equal((*time.Time)(nil), model.CreatedAt)
}

func (s *encodingSuite) encodeValue(v interface{}) (string, error) {
	return NewEncoder("Value").Encode(v)
}

func (s *encodingSuite) encodeValuePtr(v interface{}) (string, error) {
	return NewEncoder("ValuePtr").Encode(v)
}

func (s *encodingSuite) decodeValue(m interface{}, c string) (interface{}, error) {
	fields, err := NewDecoder("Value").Decode(c, m)
	if err != nil {
		return nil, err
	}
	if len(fields) != 1 {
		s.FailNow("invalid value model: %v, fields %v", m, fields)
	}
	return fields[0], nil
}

func (s *encodingSuite) decodeValuePtr(m interface{}, c string) (interface{}, error) {
	fields, err := NewDecoder("ValuePtr").Decode(c, m)
	if err != nil {
		return nil, err
	}
	if len(fields) != 1 {
		s.FailNow("invalid value model: %v, fields %v", m, fields)
	}
	return fields[0], nil
}

package gorfh

import (
  "encoding/binary"
  "github.com/ibm-messaging/mq-golang/v5/ibmmq"
  "io"
)

const (
  StrucId     = "RFH "
  EncodingUTF8 = 1208
)

func OffsetRFH2(payload []byte) int {
  offset := uint32(0)
  var endian binary.ByteOrder
  if ibmmq.MQENC_NATIVE % 2 == 0 {
    endian = binary.LittleEndian
  } else {
    endian = binary.BigEndian
  }
  for string(payload[offset:4]) == StrucId {
    offset += endian.Uint32(payload[offset+8:offset+12])
  }
  return int(offset)
}

type RFH2 struct {
  StrucId string
  Version int32
  StrucLength int32
  Encoding int32
  CodedCharSetId int32
  Format string
  Flags int32
  NameValueCCSID int32
  NameValues     H
}

func NewMQRFH2() *RFH2 {
  // cmqc.h:5021
  return &RFH2{
    StrucId:        StrucId,
    Version:        ibmmq.MQRFH_VERSION_2,
    StrucLength:    ibmmq.MQRFH_STRUC_LENGTH_FIXED_2,
    Encoding:       ibmmq.MQENC_NATIVE,
    CodedCharSetId: ibmmq.MQCCSI_INHERIT,
    Format:         ibmmq.MQFMT_NONE,
    Flags:          ibmmq.MQRFH_NONE,
    NameValueCCSID: EncodingUTF8,
    NameValues: make(H),
  }
}

func (hdr *RFH2) getEndian() binary.ByteOrder {
  if hdr.Encoding % 2 == 0 {
    return binary.LittleEndian
  } else {
    return binary.BigEndian
  }
}

func (hdr *RFH2) MarshalBinary() ([]byte, error) {
  const space8 = "        "
  endian := hdr.getEndian()
  buf := &bufferWithLens{}

  _, _ = buf.Write([]byte(hdr.StrucId))
  _ = binary.Write(buf, endian, uint32(hdr.Version))
  _ = binary.Write(buf, endian, uint32(hdr.StrucLength))
  _ = binary.Write(buf, endian, uint32(hdr.Encoding))
  _ = binary.Write(buf, endian, uint32(hdr.CodedCharSetId))
  _, _ = buf.Write([]byte((hdr.Format+space8)[0:8]))
  _ = binary.Write(buf, endian, uint32(hdr.Flags))
  _ = binary.Write(buf, endian, uint32(hdr.NameValueCCSID))
  if hdr.NameValues != nil && len(hdr.NameValues) > 0 {
    _, err := hdr.NameValues.writeTo(buf)
    if err != nil {
      return nil, err
    }
  }

  offset := len(hdr.StrucId) + 4 // strucId + version
  buf.PutLen(offset, buf.Len())
  return buf.BytesOrdered(endian), nil
}

func (hdr *RFH2) ReadFrom(r io.Reader) (n int64, err error) {
  panic("implement me")
}

package gorfh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ibm-messaging/mq-golang/v5/ibmmq"
)

const (
	StrucId      = "RFH "
	EncodingUTF8 = 1208
)

var (
	ErrBinFormatRFH2 = errors.New("error decoding rfh2")
)

func OffsetRFH2(payload []byte) int {
	offset := uint32(0)
	var endian binary.ByteOrder
	if ibmmq.MQENC_NATIVE%2 == 0 {
		endian = binary.LittleEndian
	} else {
		endian = binary.BigEndian
	}
	for string(payload[offset:4]) == StrucId {
		offset += endian.Uint32(payload[offset+8 : offset+12])
	}
	return int(offset)
}

type RFH2 struct {
	StrucId        string
	Version        int32
	StrucLength    int32
	Encoding       int32
	CodedCharSetId int32
	Format         string
	Flags          int32
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
		NameValues:     make(H),
	}
}

func (hdr *RFH2) getEndian() binary.ByteOrder {
	if hdr.Encoding%2 == 0 {
		return binary.LittleEndian
	} else {
		return binary.BigEndian
	}
}

func (hdr *RFH2) UnmarshalBinary(data []byte) error {
	if !bytes.Equal([]byte(StrucId), data[:len(StrucId)]) {
		return ErrBinFormatRFH2
	} else {
		hdr.StrucId = StrucId
	}
	var n int
	var vi uint64
	offset := len(StrucId)
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		if n < 1 {
			return ErrBinFormatRFH2
		}
		hdr.Version = int32(vi)
		if hdr.Version != ibmmq.MQRFH_VERSION_2 {
			return ErrBinFormatRFH2
		}
		offset += 4
	}
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		if n < 1 {
			return ErrBinFormatRFH2
		}
		hdr.StrucLength = int32(vi)
		offset += 4
		if hdr.StrucLength < ibmmq.MQRFH_STRUC_LENGTH_FIXED_2 {
			return ErrBinFormatRFH2
		}
	}
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		if n < 1 {
			return ErrBinFormatRFH2
		}
		hdr.Encoding = int32(vi)
		offset += 4
	}
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		//if n < 1 {
		//  return ErrBinFormatRFH2
		//}
		hdr.CodedCharSetId = int32(vi)
		offset += 4
	}
	{
		hdr.Format = string(data[offset : offset+8])
		offset += 8
	}
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		//if n < 1 {
		//  return ErrBinFormatRFH2
		//}
		hdr.Flags = int32(vi)
		offset += 4
	}
	{
		vi, n = binary.Uvarint(data[offset : offset+4])
		//if n < 1 {
		//  return ErrBinFormatRFH2
		//}
		hdr.NameValueCCSID = int32(vi)
		offset += 4
		//allowed := []int32{
		//  EncodingUTF8,
		//  1200,
		//  13488,
		//  17584,
		//}
		//valid := false
		//for _, i := range allowed {
		//  if i == hdr.NameValueCCSID {
		//    valid = true
		//    break
		//  }
		//}
		//if !valid {
		//  return ErrBinFormatRFH2
		//}
	}

	hdr.NameValues = make(H)
	xmlPayload := data[ibmmq.MQRFH_STRUC_LENGTH_FIXED_2:hdr.StrucLength]
	for i := 0; i < len(xmlPayload); {
		length, n := binary.Uvarint(xmlPayload)
		if n < 1 {
			return ErrBinFormatRFH2
		}
		i += int(length) + 1
		holder, err := parseXml(xmlPayload[4 : i+4])
		if err != nil {
			return err
		}
		for k, v := range holder {
			hdr.NameValues[k] = v
		}
		xmlPayload = xmlPayload[i:]
	}
	return nil
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
	_, _ = buf.Write([]byte((hdr.Format + space8)[0:8]))
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

type RFH2Slice []*RFH2

func (s RFH2Slice) Len() (length int) {
	for _, hdr := range s {
		length += int(hdr.StrucLength)
	}
	return
}

package gorfh

import "bytes"

type UTF8Decoder struct {
	payload []byte
}

func NewDecoder(p []byte) *UTF8Decoder {
	return &UTF8Decoder{payload: p}
}

func (d *UTF8Decoder) DecodeAll() (out RFH2Slice, err error) {
	out = make(RFH2Slice, 0)
	for offset := 0; bytes.Equal(d.payload[offset:offset+len(StrucId)], []byte(StrucId)); {
		var hdr RFH2
		err = hdr.UnmarshalBinary(d.payload[offset:])
		if err != nil {
			return
		}
		out = append(out, &hdr)
		offset += int(hdr.StrucLength)
	}
	return
}

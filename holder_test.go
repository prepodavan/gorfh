package gorfh

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_parseXml(t *testing.T) {
	type args struct {
		p []byte
	}
	type want struct {
		d   H
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "default",
			args: args{p: []byte("<usr><filename>2</filename><operation>op</operation></usr>")},
			want: want{
				d: H{
					"usr": H{
						"filename":  int64(2),
						"operation": "op",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := parseXml(tt.args.p)
			if (err != nil) != tt.want.err {
				t.Errorf("parseXml() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(gotD, tt.want.d) {
				t.Errorf("parseXml() gotD = %v, want %v", gotD, tt.want.d)
			}
		})
	}
}

func Test_parse1(t *testing.T) {
	r := NewMQRFH2()
	r.NameValues["usr"] = H{"my": 1}
	r1 := NewMQRFH2()
	r1.NameValues["vtb"] = H{"bhive": 3}
	p1, err := r1.MarshalBinary()
	require.Nil(t, err)
	p, err := r.MarshalBinary()
	require.Nil(t, err)
	d := append(p, p1...)
	headers, err := NewDecoder(d).DecodeAll()
	require.Nil(t, err)

	jsonR, err := json.Marshal(r.NameValues)
	require.Nil(t, err)
	jsonH, err := json.Marshal(headers[0].NameValues)
	require.Nil(t, err)
	jsonR1, err := json.Marshal(r1.NameValues)
	require.Nil(t, err)
	jsonH1, err := json.Marshal(headers[1].NameValues)
	require.Equal(t, jsonR, jsonH)
	require.Equal(t, jsonR1, jsonH1)
}

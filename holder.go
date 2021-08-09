package gorfh

import (
	"io"
	"strconv"
)

type H map[string]interface{}

func (h H) writeTo(out *bufferWithLens) (n int, err error) {
	for folderName, folder := range h {
		n, err = out.WriteWithLen(func(dst *bufferWithLens) (n int, err error) {
			n, err = dst.Write([]byte("<" + folderName + ">"))
			if err != nil {
				return
			}
			var f int
			f, err = h.writeFolder(folder, dst)
			if err != nil {
				return
			}
			n += f
			save := n
			n ,err = dst.Write([]byte("</" + folderName + ">"))
			if err != nil {
				return
			}
			n += save
			if n % 4 != 0 {
				dst.Grow(((n / 4) + 1) * 4) // grow to next multiple of four
				for n % 4 != 0 { // and fill it with spaces
					err = dst.WriteByte(' ')
					if err != nil {
						return
					}
					n++
				}
			}
			return
		})
		if err != nil {
			return
		}
	}
	return
}

func (h H) writeFolder(folder interface{}, out io.Writer) (n int, err error) {
	switch v := folder.(type) {
	case []byte:
		n, err = out.Write(v)
	case string:
		n, err = out.Write([]byte(v))
	case bool:
		n, err = out.Write([]byte(strconv.FormatBool(v)))
	case uint:
		n, err = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
	case uint8:
		n, err = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
	case uint16:
		n, err = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
	case uint32:
		n, err = out.Write([]byte(strconv.FormatUint(uint64(v), 10)))
	case uint64:
		n, err = out.Write([]byte(strconv.FormatUint(v, 10)))
	case int:
		n, err = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
	case int8:
		n, err = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
	case int16:
		n, err = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
	case int32:
		n, err = out.Write([]byte(strconv.FormatInt(int64(v), 10)))
	case int64:
		n, err = out.Write([]byte(strconv.FormatInt(v, 10)))
	case float32:
		n, err = out.Write([]byte(strconv.FormatFloat(float64(v), 'b', -1, 10)))
	case float64:
		n, err = out.Write([]byte(strconv.FormatFloat(v, 'b', -1, 10)))
	case H:
		var k int
		for name, value := range v {
			k, err = out.Write([]byte("<" + name + ">"))
			if err != nil {
				return
			}
			n += k
			k, err = h.writeFolder(value, out)
			if err != nil {
				return
			}
			n += k
			k, err = out.Write([]byte("</" + name + ">"))
			if err != nil {
				return
			}
			n += k
		}
	}
	return
}

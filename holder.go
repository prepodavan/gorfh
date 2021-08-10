package gorfh

import (
	"bytes"
	"container/list"
	"errors"
	"github.com/prepodavan/gorfh/utils"
	"io"
	"strconv"
)

var (
	ErrXml = errors.New("invalid xml syntax")
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
			n, err = dst.Write([]byte("</" + folderName + ">"))
			if err != nil {
				return
			}
			n += save
			if n%4 != 0 {
				dst.Grow(((n / 4) + 1) * 4) // grow to next multiple of four
				for n%4 != 0 {              // and fill it with spaces
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

func parseXml(p []byte) (d H, err error) {
	type tag struct {
		start  int
		name   []byte
		holder H
	}
	i := 0
	stack := list.New()
	d = make(H)
	stack.PushBack(&tag{holder: d})
	for i < len(p) {
		if p[i] != '<' {
			i++
			continue
		}
		if i == len(p)-1 {
			err = ErrXml
			return
		}
		tagEnd := bytes.IndexByte(p[i+1:], '>')
		if tagEnd == -1 {
			err = ErrXml
			return
		}
		tagEnd += i + 1
		if p[i+1] != '/' { // if opening tag
			name := bytes.TrimSpace(p[i+1 : tagEnd])
			if !utils.IsName(name) {
				err = ErrXml
				return
			}
			stack.PushBack(&tag{name: name, holder: make(H), start: tagEnd + 1})
			i = tagEnd + 1
			continue
		}
		// if closing tag
		name := bytes.TrimSpace(p[i+2 : tagEnd])
		if !utils.IsName(name) {
			err = ErrXml
			return
		}
		t := stack.Back().Value.(*tag)
		if !bytes.Equal(name, t.name) {
			err = ErrXml
			return
		}
		stack.Remove(stack.Back())
		parent := stack.Back().Value.(*tag)
		if len(t.holder) > 0 { // it's a group
			parent.holder[string(t.name)] = t.holder
		} else { // it's an element
			buf := bytes.TrimSpace(p[t.start:i])
			if i, err := strconv.ParseInt(string(buf), 10, 64); err == nil {
				parent.holder[string(t.name)] = i
			} else if b, err := strconv.ParseBool(string(buf)); err == nil {
				parent.holder[string(t.name)] = b
			} else if i, err := strconv.ParseFloat(string(buf), 64); err == nil {
				parent.holder[string(t.name)] = i
			} else {
				parent.holder[string(t.name)] = string(buf)
			}
		}
		i = tagEnd + 1
	}
	return
}

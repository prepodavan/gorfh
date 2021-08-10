## Parsing
```go
headers, _ := gorfh.NewDecoder(payload).DecodeAll()
fmt.Println(headers[0]) // header
fmt.Println(payload[headers.Len():]) // msg
```
package rules

import "io"

type Interface interface {
	ConvertIoReaderToStruct(data io.Reader, model interface{})
}

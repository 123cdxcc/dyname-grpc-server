package protoutil

import (
	"encoding/json"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProtoAllMessage(t *testing.T) {
	a := assert.New(t)
	fileDescriptor, err := protoparse.Parser{
		ImportPaths: []string{
			"../../api/",
			".",
		},
	}.ParseFiles("test.proto", "grpc_helper.proto")
	a.NoError(err)
	definitions := GetProtoAllSchema(fileDescriptor)
	for name, schema := range definitions {
		fmt.Print(name, ": ")
		b, err := json.MarshalIndent(schema, "", "  ")
		a.NoError(err)
		fmt.Println(string(b))
	}
}

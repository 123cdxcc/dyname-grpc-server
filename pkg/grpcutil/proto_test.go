package grpcutil

import (
	"encoding/json"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProtoToJsonSchema(t *testing.T) {
	a := assert.New(t)
	fileDescriptor, err := protoparse.Parser{
		ImportPaths: []string{
			"../../api/",
		},
	}.ParseFiles("grpc_helper.proto")
	a.NoError(err)
	for _, descriptor := range fileDescriptor {
		for _, messageDescriptor := range descriptor.GetMessageTypes() {
			if messageDescriptor.GetName() == "TestMessage" {
				var m map[string]any
				err = json.Unmarshal([]byte(ProtoToJsonSchema(messageDescriptor)), &m)
				a.NoError(err)
				b, _ := json.MarshalIndent(m, "", "  ")
				fmt.Println(string(b))
			}
		}
	}
}

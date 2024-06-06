package grpcutil

import (
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
				ProtoToJsonSchema(messageDescriptor)
			}
		}
	}
}

func TestGetProtoAllMessage(t *testing.T) {
	a := assert.New(t)
	fileDescriptor, err := protoparse.Parser{
		ImportPaths: []string{
			"../../api/",
		},
	}.ParseFiles("grpc_helper.proto")
	a.NoError(err)
	m := GetProtoAllMessage(fileDescriptor)
	for k := range m {
		fmt.Println(k)
	}
}

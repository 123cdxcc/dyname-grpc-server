package grpcutil

import (
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

func ProtoToJsonSchema(descriptor *desc.MessageDescriptor) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for i, field := range descriptor.GetFields() {
		sb.WriteString("\"")
		sb.WriteString(field.GetName())
		sb.WriteString("\":\"")
		sb.WriteString(protoTypeToJsonType(field.GetType()))
		sb.WriteString("\"")
		if i < len(descriptor.GetFields())-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func protoTypeToJsonType(fieldType descriptorpb.FieldDescriptorProto_Type) string {
	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return "number"
	}
	return "null"
}

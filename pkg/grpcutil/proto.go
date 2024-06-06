package grpcutil

import (
	"dyname-grpc-server/pkg/jsonschema"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/liushuochen/gotable"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

func ProtoToJsonSchema(descriptor *desc.MessageDescriptor) string {
	table, err := gotable.Create(
		"Name",
		"type",
		"is map",
		"is array",
		"json id",
		"json type",
		"ref",
		"items",
		"AdditionalProperties",
	)
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return ""
	}
	store := make(map[string]jsonschema.Schema)
	for _, field := range descriptor.GetFields() {
		schema := ProtoFieldTypeToJsonSchemaType(field, store)
		items := make([]string, 0, len(schema.Items))
		for _, item := range schema.Items {
			items = append(items, item.Ref)
		}
		var AdditionalProperties string
		if schema.AdditionalProperties != nil {
			AdditionalProperties = schema.AdditionalProperties.Ref
		}
		table.AddRow([]string{
			field.GetName(),
			field.GetType().String(),
			fmt.Sprint(field.IsMap()),
			fmt.Sprint(field.IsRepeated()),
			schema.Id,
			string(schema.Type),
			schema.Ref,
			strings.Join(items, ","),
			AdditionalProperties,
		})
	}
	table.Align("Name", gotable.Left)
	table.Align("json id", gotable.Left)
	fmt.Println(table.String())
	return ""
}

func GetTypeID(field *desc.FieldDescriptor) string {
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return field.GetMessageType().GetFullyQualifiedName()
	} else {
		return field.GetType().String()
	}
}

func GetProtoAllMessage(fileDesc []*desc.FileDescriptor) map[string]*desc.MessageDescriptor {
	result := make(map[string]*desc.MessageDescriptor)
	for _, file := range fileDesc {
		m := GetAllMessage(file.GetMessageTypes())
		for s, descriptor := range m {
			result[s] = descriptor
		}
	}
	return result
}

func GetAllMessage(ms []*desc.MessageDescriptor) map[string]*desc.MessageDescriptor {
	result := make(map[string]*desc.MessageDescriptor)
	for _, message := range ms {
		if len(message.GetNestedMessageTypes()) > 0 {
			nms := GetAllMessage(message.GetNestedMessageTypes())
			for s, descriptor := range nms {
				result[s] = descriptor
			}
		}
		result[message.GetFullyQualifiedName()] = message
	}
	return result
}

func ProtoFieldTypeToJsonSchemaType(field *desc.FieldDescriptor, store map[string]jsonschema.Schema) (schema jsonschema.Schema) {
	defer func() {
		schema.Id = GetTypeID(field)
		if schema.Id != "" {
			store[schema.Id] = schema
		}
	}()
	fieldType := field.GetType()
	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return jsonschema.Schema{
			Type: jsonschema.IntegerType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return jsonschema.Schema{
			Type: jsonschema.NumberType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return jsonschema.Schema{
			Type: jsonschema.StringType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		if field.IsRepeated() {
			schema = jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []jsonschema.Schema{
				{
					Ref: GetTypeID(field),
				},
			}
			return schema
		}
		return jsonschema.Schema{
			Type: jsonschema.StringType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return jsonschema.Schema{
			Type:            jsonschema.StringType,
			ContentEncoding: "base64",
		}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return jsonschema.Schema{
			Type: jsonschema.BooleanType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum := make([]string, 0, len(field.GetEnumType().GetValues()))
		for _, item := range field.GetEnumType().GetValues() {
			enum = append(enum, item.GetName())
		}
		return jsonschema.Schema{
			Type: jsonschema.StringType,
			Enum: enum,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		schema = jsonschema.Schema{}
		if field.IsMap() {
			schema.Type = jsonschema.ObjectType
			id := GetTypeID(field.GetMapValueType())
			schema.AdditionalProperties = &jsonschema.Schema{
				Ref: id,
			}
		} else if field.IsRepeated() {
			schema.Type = jsonschema.ArrayType
			schema.Items = []jsonschema.Schema{
				{
					Ref: GetTypeID(field),
				},
			}
		} else if field.GetMessageType() != nil {
			schema.Type = jsonschema.ObjectType
			schema.Ref = GetTypeID(field)
		}
		return schema
	}
	return jsonschema.Schema{}
}

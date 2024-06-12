package protoutil

import (
	"dyname-grpc-server/pkg/jsonschema"
	"dyname-grpc-server/pkg/tools"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

func GetProtoAllMessage(fileDesc []*desc.FileDescriptor) map[string]*desc.MessageDescriptor {
	messageRegister := make(map[string]*desc.MessageDescriptor)
	for _, file := range fileDesc {
		for _, messageDescriptor := range file.GetMessageTypes() {
			getAllMessage(messageDescriptor, messageRegister)
		}
	}
	return messageRegister
}

func getAllMessage(message *desc.MessageDescriptor, messageRegister map[string]*desc.MessageDescriptor) {
	if len(message.GetNestedMessageTypes()) > 0 {
		for _, messageDescriptor := range message.GetNestedMessageTypes() {
			getAllMessage(messageDescriptor, messageRegister)
		}
	}
	messageRegister[message.GetFullyQualifiedName()] = message
}

type IGetFullyQualifiedName interface {
	GetFullyQualifiedName() string
}

func getName(i IGetFullyQualifiedName) string {
	name := i.GetFullyQualifiedName()
	return strings.Replace(name, ".", "_", -1)
}

func getRefName(i IGetFullyQualifiedName) string {
	return fmt.Sprintf("%s%s", jsonschema.DefinitionsPath, getName(i))
}

func ProtoField2JsonSchemaType(i *desc.FieldDescriptor) *jsonschema.Schema {
	switch i.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type: jsonschema.IntegerType,
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type: jsonschema.IntegerType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type: jsonschema.NumberType,
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type: jsonschema.NumberType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT64:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type: jsonschema.StringType,
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type: jsonschema.StringType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type: jsonschema.StringType,
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type: jsonschema.StringType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type:            jsonschema.StringType,
					ContentEncoding: "base64",
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type:            jsonschema.StringType,
			ContentEncoding: "base64",
		}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		if i.IsRepeated() {
			schema := &jsonschema.Schema{}
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Type: jsonschema.BooleanType,
				},
			}
			return schema
		}
		return &jsonschema.Schema{
			Type: jsonschema.BooleanType,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum := make([]string, 0, len(i.GetEnumType().GetValues()))
		for _, item := range i.GetEnumType().GetValues() {
			enum = append(enum, item.GetName())
		}
		return &jsonschema.Schema{
			Type: jsonschema.StringType,
			Enum: enum,
		}
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		schema := &jsonschema.Schema{}
		if i.IsMap() {
			schema.Type = jsonschema.ObjectType
			key := ProtoField2JsonSchemaType(i.GetMapKeyType())
			value := ProtoField2JsonSchemaType(i.GetMapValueType())
			schema.Items = []*jsonschema.Schema{
				{
					Properties: map[string]*jsonschema.Schema{
						"key":   key,
						"value": value,
					},
				},
			}
		} else if i.IsRepeated() {
			schema.Type = jsonschema.ArrayType
			schema.Items = []*jsonschema.Schema{
				{
					Ref: getRefName(i.GetMessageType()),
				},
			}
		} else if i.GetMessageType() != nil {
			schema.Type = jsonschema.ObjectType
			schema.Ref = getRefName(i.GetMessageType())
		}
		return schema
	}
	return nil
}

func GetRefs(schema *jsonschema.Schema) []string {
	refsMap := make(map[string]struct{})
	for _, v := range schema.Properties {
		switch v.Type {
		case jsonschema.ArrayType:
			for _, item := range v.Items {
				refsMap[item.Ref] = struct{}{}
			}
		case jsonschema.ObjectType:
			refsMap[v.Ref] = struct{}{}
		}
	}
	refs := make([]string, 0, len(refsMap))
	for k := range refsMap {
		if k == "" {
			continue
		}
		k = strings.Replace(k, jsonschema.DefinitionsPath, "", -1)
		refs = append(refs, k)
	}
	return refs
}

func GetProtoAllSchema(fileDesc []*desc.FileDescriptor) map[string]*jsonschema.Schema {
	messageRegister := GetProtoAllMessage(fileDesc)
	definitions := make(map[string]*jsonschema.Schema)
	for _, descriptor := range messageRegister {
		schema := &jsonschema.Schema{
			SchemaVersion: jsonschema.SchemaVersion,
			Id:            getName(descriptor),
			Type:          jsonschema.ObjectType,
			Definitions:   make(map[string]*jsonschema.Schema),
			Properties:    make(map[string]*jsonschema.Schema),
		}
		for _, fieldDescriptor := range descriptor.GetFields() {
			fieldType := ProtoField2JsonSchemaType(fieldDescriptor)
			schema.Properties[fieldDescriptor.GetName()] = fieldType
		}
		definitions[descriptor.GetFullyQualifiedName()] = schema
	}
	result := make(map[string]*jsonschema.Schema, len(definitions))
	for name, schema := range definitions {
		rSchema := tools.Copy(schema)
		rSchema.Definitions = make(map[string]*jsonschema.Schema, len(schema.Definitions))
		refs := GetRefs(schema)
		for _, ref := range refs {
			for _, v := range definitions {
				if v.Id == ref {
					vc := tools.Copy(v)
					vc.Id = ""
					vc.SchemaVersion = ""
					rSchema.Definitions[ref] = vc
					break
				}
			}
		}
		result[name] = rSchema
	}
	return result
}

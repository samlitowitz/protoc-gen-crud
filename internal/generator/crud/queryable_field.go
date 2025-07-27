package crud

import (
	"fmt"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	"google.golang.org/protobuf/types/descriptorpb"
)

type QueryableField struct {
	*descriptor.Field

	// IsInlined is true when the field associated with this column is inlined
	IsInlined bool
	// Parent is the field which the field associated with this column is derived from and is only set when IsInlined = true
	Parent *descriptor.Field
}

func QueryableFieldsFromFields(fields []*descriptor.Field) []*QueryableField {
	var qFields []*QueryableField

	for _, field := range fields {
		if !field.Inline {
			qFields = append(qFields, &QueryableField{Field: field})
			continue
		}
		if field.IsScalarGoType() {
			qFields = append(qFields, &QueryableField{Field: field})
			continue
		}
		// skip field, only handle in-line field types of message
		if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}
		// generator error
		if field.FieldMessage == nil {
			// TODO: panic or skip? panic for now
			panic(fmt.Errorf("generator error: %s: undefined FieldMessage", field.FQFN()))
		}

		// types with no CRUD definition
		if !field.FieldMessage.GenerateCRUD {
			for _, msgTypField := range field.FieldMessage.Fields {
				clone := shallowCopyField(msgTypField)
				qFields = append(qFields, &QueryableField{Field: clone, IsInlined: true, Parent: field})
			}
			continue
		}

		// types which have a CRUD definition to generate
		for _, msgTypField := range field.FieldMessage.PrimaryKey() {
			clone := shallowCopyField(msgTypField)
			qFields = append(qFields, &QueryableField{Field: clone, IsInlined: true, Parent: field})
		}
		for _, msgTypField := range field.FieldMessage.NonPrimeAttributes() {
			clone := shallowCopyField(msgTypField)
			qFields = append(qFields, &QueryableField{Field: clone, IsInlined: true, Parent: field})
		}
	}

	return qFields
}

func QueryableFieldsFromMessage(msg *descriptor.Message) []*QueryableField {
	return append(QueryableFieldsFromFields(msg.PrimaryKey()), QueryableFieldsFromFields(msg.NonPrimeAttributes())...)
}

func shallowCopyField(original *descriptor.Field) *descriptor.Field {
	return &descriptor.Field{
		FieldDescriptorProto: original.FieldDescriptorProto,
		Message:              original.Message,
		FieldEnum:            original.FieldEnum,
		FieldMessage:         original.FieldMessage,
		ForcePrefixedName:    original.ForcePrefixedName,
		Ignore:               original.Ignore,
		Inline:               original.Inline,
		Relationships:        original.Relationships,
		IsPrimeAttribute:     original.IsPrimeAttribute,
	}
}

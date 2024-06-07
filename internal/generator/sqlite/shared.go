package sqlite

import (
	"fmt"

	"github.com/samlitowitz/protoc-gen-crud/internal/generator/crud"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/types/descriptorpb"
)

func QuotedIdent(s string) string {
	return "\"" + Ident(s) + "\""
}

func Ident(s string) string {
	return strcase.ToSnake(s)
}

func ColumnsFromFields(fields []*crud.QueryableField) []*Column {
	var cols []*Column
	for _, field := range fields {
		cols = append(cols, &Column{QueryableField: field})
	}

	return cols
}

type Column struct {
	*crud.QueryableField
}

func (col *Column) GetName() string {
	if !col.IsInlined {
		return col.Field.GetName()
	}
	return col.Parent.GetName() + "_" + col.Field.GetName()
}

func (col *Column) GetComment() string {
	switch col.Field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return ""

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return fmt.Sprintf(
			" /* references %s.%s */",
			QuotedIdent(col.Field.FieldEnum.GetName()),
			QuotedIdent("id"),
		)

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	default:
		panic(fmt.Errorf("sqlite: sql: field %s: unsupported type %s", col.Field.GetName(), col.Field.GetType()))
	}
}

func (col *Column) GetType() string {
	switch col.Field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "REAL"

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "BLOB"

	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "TEXT"

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		// Enums will reference a look-up table
		return "INTEGER"

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	default:
		panic(fmt.Errorf("sqlite: sql: field %s: unsupported type %s", col.Field.GetName(), col.Field.GetType()))
	}
}

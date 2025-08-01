package descriptor

import (
	"fmt"

	relationshipOptions "github.com/samlitowitz/protoc-gen-crud/options/relationships"

	crudOptions "github.com/samlitowitz/protoc-gen-crud/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// loadCRUDs registers CRUDs from "targetFile" to "r"
// It must be called after loadFile is called for all files so that loadCRUDs
// can resolve names of message types and their fields
func (r *Registry) loadCRUDs(file *File) error {
	for _, msg := range file.Messages {
		msgOpts, err := extractMessageOptions(msg.DescriptorProto)
		if err != nil {
			return fmt.Errorf("%s: %v", msg.FQMN(), err)
		}
		// No CRUD message options defined, skip it
		if msgOpts == nil {
			continue
		}

		msg.GenerateCRUD = true

		err = assignMessageOptions(msg, msgOpts)
		if err != nil {
			return fmt.Errorf("%s: %v", msg.FQMN(), err)
		}

		for impl := range msg.Implementations {
			file.Implementations[impl] = struct{}{}
		}

		for _, field := range msg.Fields {
			_, isPrimeAttribute := msg.PrimaryKeyByFQFN[field.FQFN()]
			field.IsPrimeAttribute = isPrimeAttribute

			isFieldMaskField := msg.HasFieldMask() && msg.FieldMask.FQFN() == field.FQFN()
			if !isPrimeAttribute && !isFieldMaskField {
				msg.NonPrimeAttributesByFQFN[field.FQFN()] = field
			}

			fieldOpts, err := extractFieldOptions(field.FieldDescriptorProto)
			if err != nil {
				return fmt.Errorf("%s: %v", field.FQFN(), err)
			}
			if fieldOpts == nil {
				continue
			}
			err = assignFieldOptions(field, fieldOpts)
			if err != nil {
				return fmt.Errorf("%s: assign field options: %v", field.GetName(), err)
			}
			if field.Ignore && field.IsPrimeAttribute {
				return fmt.Errorf("%s: ignored field cannot be part of a primary key", field.FQFN())
			}
			if field.Inline && field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
				return fmt.Errorf("%s: inlined field must be of type message", field.FQFN())
			}

			isCreatedAt := msg.HasCreatedAt() && msg.CreatedAt.FQFN() == field.FQFN()
			if isCreatedAt && !field.AsTimestamp {
				return fmt.Errorf("%s: field designed as `createdAt` must be a timestamp", field.FQFN())
			}
			isUpdatedAt := msg.HasUpdatedAt() && msg.UpdatedAt.FQFN() == field.FQFN()
			if isUpdatedAt && !field.AsTimestamp {
				return fmt.Errorf("%s: field designed as `createdAt` must be a timestamp", field.FQFN())
			}

			err = assignRelationships(r, msg, field, fieldOpts)
			if err != nil {
				return fmt.Errorf("%s: assign relationship: %v", field.FQFN(), err)
			}
			file.Relationships = append(file.Relationships, field.Relationships...)
		}
	}
	return nil
}

func assignMessageOptions(msg *Message, msgOpts *crudOptions.MessageOptions) error {
	msg.Implementations = make(map[crudOptions.Implementation]struct{})
	msg.PrimaryKeyByFQFN = make(map[string]*Field)
	msg.NonPrimeAttributesByFQFN = make(map[string]*Field)

	if msgOpts.GetImplementations() != nil {
		for _, implementation := range msgOpts.GetImplementations() {
			msg.Implementations[implementation] = struct{}{}
		}
	}

	if msgOpts.GetFieldMask() != "" {
		field, err := msg.LookupField(msgOpts.GetFieldMask())
		if err != nil {
			return err
		}
		msg.FieldMask = field
	}

	if msgOpts.GetPrimaryKey() != nil {
		msg.PrimaryKeyByFQFN = make(map[string]*Field, len(msgOpts.GetPrimaryKey()))

		for _, fieldName := range msgOpts.GetPrimaryKey() {

			field, err := msg.LookupField(fieldName)
			if err != nil {
				return fmt.Errorf("%s: field `%s`: %v", msg.GetName(), fieldName, err)
			}
			switch *field.Type {
			case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				return fmt.Errorf(
					"%s: unsupported candidate key field type %s",
					field.FQFN(),
					field.GetType(),
				)
			}
			msg.PrimaryKeyByFQFN[field.FQFN()] = field
		}
	}

	if msgOpts.GetCreatedAt() != "" {
		field, err := msg.LookupField(msgOpts.GetCreatedAt())
		if err != nil {
			return err
		}
		msg.CreatedAt = field
	}
	if msgOpts.GetUpdatedAt() != "" {
		field, err := msg.LookupField(msgOpts.GetUpdatedAt())
		if err != nil {
			return err
		}
		msg.UpdatedAt = field
	}

	return nil
}

func assignRelationships(r *Registry, msg *Message, field *Field, fieldOpts *crudOptions.FieldOptions) error {
	if !fieldOpts.HasRelationship() {
		return nil
	}
	if field.Ignore {
		return fmt.Errorf("ignored field cannot be part of a relationship")
	}
	if field.FieldMessage == nil {
		return fmt.Errorf("only fields with a message type may be part of a relationship")
	}

	var err error
	var fieldType *Message
	fieldType, err = r.LookupMsg("", field.FieldMessage.FQMN())
	if err != nil {
		return err
	}

	switch fieldOpts.GetRelationship().GetType() {
	case relationshipOptions.Type_MANY_TO_ONE:
	case relationshipOptions.Type_MANY_TO_MANY:
	case relationshipOptions.Type_ONE_TO_MANY:
	case relationshipOptions.Type_ONE_TO_ONE:

	case relationshipOptions.Type_UNKNOWN_TYPE:
		fallthrough
	default:
		return fmt.Errorf("unsupported relationship type %s", fieldOpts.GetRelationship().GetType().String())
	}
	field.Relationships = append(field.Relationships, &Relationship{
		Relationship: fieldOpts.GetRelationship(),
		DefinedOn:    msg,
		With:         fieldType,
	})

	return nil
}

func assignFieldOptions(field *Field, fieldOpts *crudOptions.FieldOptions) error {
	field.Ignore = fieldOpts.GetIgnore()
	field.Inline = fieldOpts.GetInline()
	field.AsTimestamp = fieldOpts.GetAsTimestamp()
	return nil
}

func extractMessageOptions(msg *descriptorpb.DescriptorProto) (*crudOptions.MessageOptions, error) {
	if msg.GetOptions() == nil {
		return nil, nil
	}
	if !proto.HasExtension(msg.GetOptions(), crudOptions.E_CrudMessageOptions) {
		return nil, nil
	}
	ext := proto.GetExtension(msg.GetOptions(), crudOptions.E_CrudMessageOptions)
	opts, ok := ext.(*crudOptions.MessageOptions)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want MessageOptions", ext)
	}
	return opts, nil
}

func extractFieldOptions(fd *descriptorpb.FieldDescriptorProto) (*crudOptions.FieldOptions, error) {
	if fd.GetOptions() == nil {
		return nil, nil
	}
	if !proto.HasExtension(fd.GetOptions(), crudOptions.E_CrudFieldOptions) {
		return nil, nil
	}
	ext := proto.GetExtension(fd.GetOptions(), crudOptions.E_CrudFieldOptions)
	opts, ok := ext.(*crudOptions.FieldOptions)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want CRUD", ext)
	}
	return opts, nil
}

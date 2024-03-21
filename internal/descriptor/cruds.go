package descriptor

import (
	"fmt"
	"strings"

	"github.com/samlitowitz/protoc-gen-crud/options/relationships"

	"github.com/samlitowitz/protoc-gen-crud/options"
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
			return err
		}
		// No CRUD message options defined, skip it
		if msgOpts == nil {
			continue
		}

		msg.GenerateCRUD = true

		err = assignMessageOptions(msg, msgOpts)
		if err != nil {
			return err
		}

		for _, field := range msg.Fields {
			fieldOpts, err := extractFieldOptions(field.FieldDescriptorProto)
			if err != nil {
				return err
			}
			if fieldOpts == nil {
				continue
			}
			field.Ignore = fieldOpts.Ignore

			err = assignRelationships(r, msg, field, fieldOpts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func assignMessageOptions(msg *Message, msgOpts *options.MessageOptions) error {
	msg.Implementations = make(map[options.Implementation]struct{})
	msg.CandidateKey = make([]*Field, 0)

	if msgOpts.GetImplementations() != nil {
		for _, implementation := range msgOpts.Implementations {
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

	if msgOpts.GetCandidateKey() != nil {
		msg.CandidateKey = make([]*Field, 0, len(msgOpts.GetCandidateKey()))
		for _, fieldName := range msgOpts.GetCandidateKey() {
			field, err := msg.LookupField(fieldName)
			if err != nil {
				return err
			}
			switch *field.Type {
			case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
				fallthrough
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				return fmt.Errorf(
					"on message type %s: unsupported candidate key field type of %s from field %s",
					msg.GetName(),
					field.GetType(),
					field.GetName(),
				)
			}
			msg.CandidateKey = append(msg.CandidateKey, field)
		}
	}

	return nil
}

func assignRelationships(r *Registry, msg *Message, field *Field, fieldOpts *options.FieldOptions) error {
	if fieldOpts.Relationship == nil {
		return nil
	}

	typeName := field.GetTypeName()

	var err error
	var fieldType *Message
	switch typeName[0] {
	case '.':
		i := strings.LastIndex(typeName, ".")
		fieldType, err = r.LookupMsg(typeName[1:i-1], typeName[i+1:])
		if err != nil {
			return err
		}
	default:
		fieldType, err = r.LookupMsg("", typeName)
		if err != nil {
			return err
		}
	}

	if fieldOpts.Relationship.GetType() == relationships.Type_UNKNOWN_TYPE {
		return fmt.Errorf(
			"on message %s: unknown relationship type defined on field %s",
			msg.GetName(),
			field.GetName(),
		)
	}

	switch fieldOpts.Relationship.GetType() {
	case relationships.Type_MANY_TO_ONE:
	case relationships.Type_MANY_TO_MANY:
	case relationships.Type_ONE_TO_MANY:
	case relationships.Type_ONE_TO_ONE:

	default:
		return fmt.Errorf("unknown relationships type")
	}
	field.Relationship = &Relationship{
		Relationship: fieldOpts.Relationship,
		DefinedOn:    msg,
		With:         fieldType,
	}
	return nil
}

func extractMessageOptions(msg *descriptorpb.DescriptorProto) (*options.MessageOptions, error) {
	if msg.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(msg.Options, options.E_CrudMessageOptions) {
		return nil, nil
	}
	ext := proto.GetExtension(msg.Options, options.E_CrudMessageOptions)
	opts, ok := ext.(*options.MessageOptions)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want MessageOptions", ext)
	}
	return opts, nil
}

func extractFieldOptions(fd *descriptorpb.FieldDescriptorProto) (*options.FieldOptions, error) {
	if fd.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(fd.Options, options.E_CrudFieldOptions) {
		return nil, nil
	}
	ext := proto.GetExtension(fd.Options, options.E_CrudFieldOptions)
	opts, ok := ext.(*options.FieldOptions)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want CRUD", ext)
	}
	return opts, nil
}

package descriptor

import (
	"fmt"

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
		def := &CRUD{
			Message: msg,
		}
		assignMessageOptions(def, msgOpts)

		for _, field := range msg.Fields {
			fieldOpts, err := extractFieldOptions(field.FieldDescriptorProto)
			if err != nil {
				return err
			}
			if fieldOpts == nil {
				continue
			}
			err = assignUniqueIdentifiers(def, field, fieldOpts)
			if err != nil {
				return err
			}

			err = assignRelationships(r.msgs, field, fieldOpts)
			if err != nil {
				return err
			}
		}
		processFieldMaskField(def)
		msg.CRUD = def
		file.CRUDs = append(file.CRUDs, def)
		r.cruds[msg.FQMN()] = def
	}
	return nil
}

// fixupCRUDRelationships turns all relationships from CRUD stubs to registry CRUDs
func (r *Registry) fixupCRUDRelationships(crud *CRUD) error {
	for _, field := range crud.Fields {
		// can't fix a relationship that doesn't exist
		if !field.HasRelationship() {
			continue
		}
		target, ok := r.cruds[*field.TypeName]
		if !ok {
			return &UnsupportedTypeError{
				fmt.Sprintf("no CRUD defined: %s", *field.TypeName),
			}
		}
		if len(target.UniqueIdentifiers) == 0 {
			return &UnsupportedTypeError{
				fmt.Sprintf("CRUD must have at least one unique identifier defined: %s", *field.TypeName),
			}
		}
		if len(target.MinimalUIDFields()) != 1 {
			return &UnsupportedTypeError{
				fmt.Sprintf("CRUD must have unique identifier with exactly one field defined: %s", *field.TypeName),
			}
		}
		field.Relationship.CRUD = target
	}

	return nil
}

func assignMessageOptions(def *CRUD, msgOpts *options.MessageOptions) {
	if msgOpts == nil {
		def.Operations = make(map[options.Operation]struct{})
		def.Implementations = make(map[options.Implementation]struct{})
		return
	}

	def.Operations = make(map[options.Operation]struct{}, len(msgOpts.Operations))
	def.Implementations = make(map[options.Implementation]struct{}, len(msgOpts.Implementations))

	for _, operation := range msgOpts.Operations {
		def.Operations[operation] = struct{}{}
	}

	for _, implementation := range msgOpts.Implementations {
		def.Implementations[implementation] = struct{}{}
	}

	def.FieldMaskFieldName = msgOpts.FieldMaskField
}

func assignUniqueIdentifiers(def *CRUD, field *Field, fieldOpts *options.FieldOptions) error {
	if fieldOpts.Uids == nil {
		return nil
	}
	if len(fieldOpts.Uids) == 0 {
		return nil
	}

	if field.Type == nil && field.TypeName == nil {
		return &UnknownTypeError{}
	}
	if field.Type == nil {
		return &UnsupportedTypeError{TypName: *field.TypeName}
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
		return &UnsupportedTypeError{TypName: *field.TypeName}
	}

	if def.UniqueIdentifiers == nil {
		def.UniqueIdentifiers = make(map[string][]*Field)
	}
	for _, uid := range fieldOpts.Uids {
		if _, ok := def.UniqueIdentifiers[uid]; !ok {
			def.UniqueIdentifiers[uid] = make([]*Field, 0, 1)
		}
		def.UniqueIdentifiers[uid] = append(def.UniqueIdentifiers[uid], field)
		if len(def.UniqueIdentifiers[uid]) > 1 {
			panic(
				fmt.Errorf(
					"unique identifiers must have exactly one field: %d given",
					len(def.UniqueIdentifiers[uid]),
				),
			)
		}
	}
	if len(def.UniqueIdentifiers) > 1 {
		panic(
			fmt.Errorf(
				"at most one unique identifier per message: %d given",
				len(def.UniqueIdentifiers),
			),
		)
	}

	return nil
}

func assignRelationships(msgs map[string]*Message, field *Field, fieldOpts *options.FieldOptions) error {
	if fieldOpts.Relationship == nil {
		return nil
	}
	_, ok := msgs[*field.TypeName]
	if !ok {
		return &UnsupportedTypeError{*field.TypeName}
	}
	switch fieldOpts.Relationship.GetType() {
	case relationships.Type_MANY_TO_ONE:
		return fmt.Errorf("many to one relationships are not implemented")
	case relationships.Type_ONE_TO_MANY:
		return fmt.Errorf("one to many relationships are not implemented")

	case relationships.Type_MANY_TO_MANY:
	case relationships.Type_ONE_TO_ONE:

	default:
		return fmt.Errorf("unknown relationships type")
	}
	field.Relationship = &Relationship{
		Relationship: fieldOpts.Relationship,
		CRUD:         &CRUD{Message: msgs[*field.TypeName]},
	}
	return nil
}

func processFieldMaskField(def *CRUD) {
	if def.FieldMaskFieldName == "" {
		return
	}
	for _, field := range def.Fields {
		if field.GetName() != def.FieldMaskFieldName {
			continue
		}
		def.FieldMaskField = field
		field.IsFieldMaskField = true
		return
	}
	panic(fmt.Sprintf("unable to find specified `fieldMaskField`: %s", def.FieldMaskFieldName))
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

type UnknownTypeError struct{}

func (err UnknownTypeError) Error() string {
	return "unknown type"
}

type UnsupportedTypeError struct {
	TypName string
}

func (err UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type: %s", err.TypName)
}

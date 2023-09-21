package descriptor

import (
	"fmt"

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
		}
		file.CRUDs = append(file.CRUDs, def)
		r.cruds[msg.FQMN()] = def
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
}

func assignUniqueIdentifiers(def *CRUD, field *Field, fieldOpts *options.FieldOptions) error {
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
			def.UniqueIdentifiers[uid] = make([]*Field, 1)
		}
		def.UniqueIdentifiers[uid] = append(def.UniqueIdentifiers[uid], field)
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

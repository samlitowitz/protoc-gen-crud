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
		if msgOpts == nil {
			continue
		}
		def := &CRUD{
			Message:           msg,
			Operations:        make(map[options.Operation]struct{}),
			UniqueIdentifiers: make(map[string][]*Field),
		}
		for _, field := range msg.Fields {
			fieldOpts, err := extractFieldOptions(field.FieldDescriptorProto)
			if err != nil {
				return err
			}
			if fieldOpts == nil {
				continue
			}
			for _, uid := range fieldOpts.GetUniqueIdentifiers() {
				if _, ok := def.UniqueIdentifiers[uid.GetId()]; !ok {
					def.UniqueIdentifiers[uid.GetId()] = make([]*Field, 1)
				}
				def.UniqueIdentifiers[uid.GetId()] = append(def.UniqueIdentifiers[uid.GetId()], field)
			}
		}

		msg.CRUD = def
		r.cruds[msg.FQMN()] = def
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

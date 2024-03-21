// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/internal/descriptor/registry.go
package descriptor

import (
	"fmt"
	"sort"
	"strings"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Registry struct {
	// msgs is a mapping from fully-qualified message name to descriptor
	msgs map[string]*Message
	// enums is a mapping from fully-qualified enum name to descriptor
	enums map[string]*Enum
	// files is a mapping from file path to descriptor
	files map[string]*File
	// pkgAliases is a mapping from package aliases to package paths in go which are already taken.
	pkgAliases map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		msgs:       make(map[string]*Message),
		enums:      make(map[string]*Enum),
		files:      make(map[string]*File),
		pkgAliases: make(map[string]string),
	}
}

func (r *Registry) LoadFromPlugin(gen *protogen.Plugin) error {
	return r.load(gen)
}

func (r *Registry) load(gen *protogen.Plugin) error {
	filePaths := make([]string, 0, len(gen.FilesByPath))
	for filePath := range gen.FilesByPath {
		filePaths = append(filePaths, filePath)
	}
	sort.Strings(filePaths)

	for _, filePath := range filePaths {
		r.loadFile(filePath, gen.FilesByPath[filePath])
	}
	for _, filePath := range filePaths {
		if !gen.FilesByPath[filePath].Generate {
			continue
		}
		err := r.fixupFieldFieldMessage(r.files[filePath])
		if err != nil {
			return err
		}
	}

	for _, filePath := range filePaths {
		if !gen.FilesByPath[filePath].Generate {
			continue
		}
		file := r.files[filePath]
		if err := r.loadCRUDs(file); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) loadFile(filePath string, file *protogen.File) {
	pkg := GoPackage{
		Path: string(file.GoImportPath),
		Name: string(file.GoPackageName),
	}
	if err := r.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
		for i := 0; ; i++ {
			alias := fmt.Sprintf("%s_%d", pkg.Name, i)
			if err := r.ReserveGoPackageAlias(alias, pkg.Path); err == nil {
				pkg.Alias = alias
				break
			}
		}
	}
	f := &File{
		FileDescriptorProto:     file.Proto,
		GoPkg:                   pkg,
		GeneratedFilenamePrefix: file.GeneratedFilenamePrefix,
	}

	r.files[filePath] = f
	r.registerMsg(f, nil, file.Proto.MessageType)
	r.registerEnum(f, nil, file.Proto.EnumType)
}

func (r *Registry) registerMsg(file *File, outerPath []string, msgs []*descriptorpb.DescriptorProto) {
	for i, md := range msgs {
		m := &Message{
			DescriptorProto:   md,
			File:              file,
			Outers:            outerPath,
			Index:             i,
			ForcePrefixedName: false,
		}
		for _, fd := range md.GetField() {
			m.Fields = append(m.Fields, &Field{
				FieldDescriptorProto: fd,
				Message:              m,
				ForcePrefixedName:    false,
			})
		}
		file.Messages = append(file.Messages, m)
		r.msgs[m.FQMN()] = m
		var outers []string
		outers = append(outers, outerPath...)
		outers = append(outers, m.GetName())
		r.registerMsg(file, outers, m.GetNestedType())
		r.registerEnum(file, outers, m.GetEnumType())
	}
}

func (r *Registry) registerEnum(file *File, outerPath []string, enums []*descriptorpb.EnumDescriptorProto) {
	for i, ed := range enums {
		e := &Enum{
			EnumDescriptorProto: ed,
			File:                file,
			Outers:              outerPath,
			Index:               i,
			ForcePrefixedName:   false,
		}
		file.Enums = append(file.Enums, e)
		r.enums[e.FQEN()] = e
	}
}

func (r *Registry) fixupFieldFieldMessage(file *File) error {
	for _, msg := range file.Messages {
		for _, field := range msg.Fields {
			if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
				continue
			}
			if field.GetTypeName() == "" {
				continue
			}
			fieldMessage, err := r.LookupMsg("", field.GetTypeName())
			if err != nil {
				return fmt.Errorf(
					"failed to fix up field %s on message %s",
					*field.Name,
					*msg.Name,
				)
			}
			field.FieldMessage = fieldMessage
		}
	}
	return nil
}

func (r *Registry) LookupMsg(location, name string) (*Message, error) {
	if strings.HasPrefix(name, ".") {
		m, ok := r.msgs[name]
		if !ok {
			return nil, fmt.Errorf("no message found: %s", name)
		}
		return m, nil
	}
	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqmn := strings.Join(append(components, name), ".")
		if m, ok := r.msgs[fqmn]; ok {
			return m, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no message found: %s", name)
}

func (r *Registry) LookupFile(name string) (*File, error) {
	f, ok := r.files[name]
	if !ok {
		return nil, fmt.Errorf("no such file given: %s", name)
	}
	return f, nil
}

func (r *Registry) ReserveGoPackageAlias(alias, pkgpath string) error {
	if taken, ok := r.pkgAliases[alias]; ok {
		if taken == pkgpath {
			return nil
		}
		return fmt.Errorf("package name %s is already taken. use another alias", alias)
	}
	r.pkgAliases[alias] = pkgpath
	return nil
}

package gen_gen

import (
	"fmt"

	"github.com/samlitowitz/protoc-gen-crud/internal/descriptor"
	gen "github.com/samlitowitz/protoc-gen-crud/internal/generator"
)

type generator struct {
	gens []gen.Generator
}

func New(gens ...gen.Generator) gen.Generator {
	return &generator{
		gens: gens,
	}
}

func (g *generator) Generate(targets []*descriptor.File) ([]*descriptor.ResponseFile, error) {
	files := make([]*descriptor.ResponseFile, 0)
	var err error

	for _, genGen := range g.gens {
		genFiles, genErr := genGen.Generate(targets)
		if genErr == nil {
			files = append(files, genFiles...)
			continue
		}
		if err == nil {
			err = genErr
			continue
		}
		err = fmt.Errorf("%s: %w", err, genErr)
	}

	return files, err
}

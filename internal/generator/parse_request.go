package generator

import (
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// REFURL: https://github.com/grpc-ecosystem/grpc-gateway/blob/main/internal/codegenerator/parse_req.go

// ParseRequest parses a code generator request from a proto Message.
func ParseRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read code generator request: %w", err)
	}
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code generator request: %w", err)
	}
	return req, nil
}

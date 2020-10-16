package handlers

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2/code"
	lsctx "github.com/hashicorp/terraform-ls/internal/context"
	ilsp "github.com/hashicorp/terraform-ls/internal/lsp"
	lsp "github.com/sourcegraph/go-lsp"
)

const rootmodulesCommandResponseVersion = 0
const rootmodulesCommandURIArgNotFound code.Code = -32004

type rootmodulesCommandResponse struct {
	Version     int              `json:"version"`
	DoneLoading bool             `json:"doneLoading"`
	RootModules []rootModuleInfo `json:"rootModules"`
}

type rootModuleInfo struct {
	URI string `json:"uri"`
}

func executeCommandRootModulesHandler(ctx context.Context, args commandArgs) (interface{}, error) {
	walker, err := lsctx.RootModuleWalker(ctx)
	if err != nil {
		return nil, err
	}

	uri, ok := args.GetString("uri")
	if !ok || uri == "" {
		return nil, fmt.Errorf("%w: expected uri argument to be set", rootmodulesCommandURIArgNotFound.Err())
	}

	fh := ilsp.FileHandlerFromDocumentURI(lsp.DocumentURI(uri))

	cf, err := lsctx.RootModuleCandidateFinder(ctx)
	if err != nil {
		return nil, err
	}
	doneLoading := !walker.IsWalking()
	candidates := cf.RootModuleCandidatesByPath(fh.Dir())

	rootModules := make([]rootModuleInfo, len(candidates))
	for i, candidate := range candidates {
		rootModules[i] = rootModuleInfo{
			URI: candidate.URI(),
		}
	}
	return rootmodulesCommandResponse{
		Version:     rootmodulesCommandResponseVersion,
		DoneLoading: doneLoading,
		RootModules: rootModules,
	}, nil
}

package handlers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-ls/internal/terraform/rootmodule"
	"github.com/hashicorp/terraform-ls/langserver"
)

func TestLangServer_workspaceExecuteCommand_rootmodules_argumentError(t *testing.T) {
	tmpDir := TempDir(t)
	testFileURI := fmt.Sprintf("%s/main.tf", tmpDir.URI())
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		RootModules: map[string]*rootmodule.RootModuleMock{
			tmpDir.Dir(): {
				TfExecFactory: validTfMockCalls(),
			},
		},
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, tmpDir.URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"github\"\n\n}\n",
			"uri": %q
		}
	}`, testFileURI)})

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "workspace/executeCommand",
		ReqParams: `{
		"command": "rootmodules"
	}`}, rootmodulesCommandFileArgNotFound.Err())
}

func TestLangServer_workspaceExecuteCommand_rootmodules_basic(t *testing.T) {
	tmpDir := TempDir(t)
	testFileURI := fmt.Sprintf("%s/main.tf", tmpDir.URI())
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		RootModules: map[string]*rootmodule.RootModuleMock{
			tmpDir.Dir(): {
				TfExecFactory: validTfMockCalls(),
			},
		},
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, tmpDir.URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"github\"\n\n}\n",
			"uri": %q
		}
	}`, testFileURI)})

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "workspace/executeCommand",
		ReqParams: fmt.Sprintf(`{
		"command": "rootmodules",
		"arguments": ["file=%s"] 
	}`, testFileURI)}, fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 3,
		"result": {
			"version": 0,
			"doneLoading": true,
			"rootModules": [
				{
					"uri": %q
				}
			]
		}
	}`, tmpDir.URI()))
}

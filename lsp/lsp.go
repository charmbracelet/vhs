package lsp

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/vhs/lexer"
	"github.com/charmbracelet/vhs/parser"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

// Server is the language server protocol server.
type Server struct {
	cache  string
	lines  []string
	errors []parser.Error
}

const languageServerName = "vhs"

var version = "0.0.1"
var handler protocol.Handler

// Run runs the language server protocol server.
func Run(_ *cobra.Command, _ []string) error {
	f, _ := os.Create("debug.log")
	logger := log.New(f)
	logger.SetLevel(log.DebugLevel)
	log.SetDefault(logger)

	// This increases logging verbosity (optional)
	commonlog.Configure(1, nil)

	s := Server{}

	handler = protocol.Handler{
		Initialize:            initialize,
		Initialized:           initialized,
		Shutdown:              shutdown,
		SetTrace:              setTrace,
		TextDocumentDidChange: s.didChange,
		TextDocumentHover:     s.hover,
	}

	server := server.NewServer(&handler, languageServerName, false)

	return server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.TextDocumentSync = protocol.TextDocumentSyncOptions{
		Change:            ptr(protocol.TextDocumentSyncKindFull),
		OpenClose:         ptr(true),
		Save:              ptr(true),
		WillSave:          ptr(true),
		WillSaveWaitUntil: ptr(true),
	}
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    languageServerName,
			Version: &version,
		},
	}, nil
}

func ptr[T any](t T) *T {
	return &t
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *Server) didChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) <= 0 {
		return errors.New("invalid did change")
	}
	doc, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
	if !ok {
		return errors.New("invalid did change")
	}
	s.cache = doc.Text
	s.lines = strings.Split(doc.Text, "\n")
	l := lexer.New(doc.Text)
	p := parser.New(l)
	p.Parse()
	s.errors = p.Errors()
	log.Debug(s.errors)

	var diagnostics []protocol.Diagnostic
	for _, error := range s.errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(error.Token.Line - 1),
					Character: uint32(error.Token.Column - 1),
				},
				End: protocol.Position{
					Line:      uint32(error.Token.Line - 1),
					Character: uint32(error.Token.Column + runewidth.StringWidth(error.Token.Literal)),
				},
			},
			Severity: ptr(protocol.DiagnosticSeverityError),
			Source:   ptr("vhs"),
			Message:  error.Msg,
		})
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})

	return nil
}

func (s *Server) hover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	return &protocol.Hover{
		Contents: "Nice",
	}, nil
}

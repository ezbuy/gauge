package lang

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getgauge/gauge/formatter"
	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/gauge/parser"
	"github.com/getgauge/gauge/util"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func format(request *jsonrpc2.Request) (interface{}, error) {
	var params lsp.DocumentFormattingParams
	if err := json.Unmarshal(*request.Params, &params); err != nil {
		return nil, err
	}
	logger.APILog.Debugf("LangServer: request received : Type: Format Document URI: %s", params.TextDocument.URI)
	file := util.ConvertURItoFilePath(params.TextDocument.URI)
	if util.IsValidSpecExtension(file) {
		spec, parseResult := new(parser.SpecParser).Parse(getContent(params.TextDocument.URI), gauge.NewConceptDictionary(), file)
		if !parseResult.Ok {
			return nil, fmt.Errorf("failed to format document. Fix all the problems first")
		}
		newString := formatter.FormatSpecification(spec)
		return createTextEdit(getContent(params.TextDocument.URI), newString), nil
	}
	return nil, fmt.Errorf("failed to format document. %s is not a valid spec file", file)
}

func createTextEdit(oldContent string, newString string) []lsp.TextEdit {
	return []lsp.TextEdit{
		{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 0,
				},
				End: lsp.Position{
					Line:      len(strings.Split(oldContent, "\n")),
					Character: len(oldContent),
				},
			},
			NewText: newString,
		},
	}
}

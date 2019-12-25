package golinters

import (
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"honnef.co/go/tools/staticcheck"
)

func NewStaticcheck() *goanalysis.Linter {
	analyzers := analyzersMapToSlice(staticcheck.Analyzers)
	setAnalyzersGoVersion(analyzers)

	return goanalysis.NewLinter(
		"staticcheck",
		"Staticcheck is a go vet on steroids, applying a ton of static analysis checks",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

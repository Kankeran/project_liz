package domain

import (
	"github.com/pkg/errors"
	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/imports"
)

// CodeFormatter struct
type CodeFormatter struct {
}

// Format format specified data witch goImports tool
func (cf *CodeFormatter) Format(data string) (output []byte, err error) {
	output, err = returns.Process("", "", []byte(data), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	output, err = imports.Process("", output, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return output, nil
}

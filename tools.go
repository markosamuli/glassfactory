// +build tools

package tools

import (
  _ "golang.org/x/tools/cmd/stringer"
  _ "golang.org/x/tools/cmd/goimports"
  _ "github.com/uudashr/gopkgs/cmd/gopkgs"
  - "golang.org/x/lint/golint"
)

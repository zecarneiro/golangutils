package system

import (
	"golangutils/pkg/common"
	"golangutils/pkg/file"
)

func GenerateTempFile(prefixo string) string {
	tempName := common.GenerateTempName(prefixo)
	return file.JoinPath(TempDir(), tempName)
}

package cert

import "capten/pkg/util"

func generateCert() error {
	return util.OsExec("bash", "./generate.sh")
}

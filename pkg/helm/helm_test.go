package helm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadValues(t *testing.T) {

	assert := require.New(t)

	reader := strings.NewReader(
		`
apps:
  - "openebs-jiva"
  - "traefik"
  - "temporal"
  - "cert-manager"
  - "pre-install"
  - "vault-operator"
  - "vault-post"
  - "vault-secrets-webhook"
  - "loki"
  - "kubviz-client"
  - "kubviz-agent"
  - "signoz"
  - "kad"
  - "kyverno"
  - "policy-reporter"
  - "kubescape"
`,
	)

	list, err := getAppsList(reader, "apps")

	assert.Nil(err)
	assert.True(len(list) > 0)

}

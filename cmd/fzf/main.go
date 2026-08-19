package fzf

import (
	"github.com/koki-develop/go-fzf"
)
func FzfSelect(instanceIDs []string, instanceNames []string) (string, error) {
	f, err := fzf.New()
	if err != nil {
		return "", err
	}

	// display "id  name" in the fuzzy finder
	idxs, err := f.Find(instanceIDs, func(i int) string {
		if i < len(instanceNames) && instanceNames[i] != "" {
			return instanceIDs[i] + "  " + instanceNames[i]
		}
		return instanceIDs[i]
	})
	if err != nil {
		return "", err
	}

	if len(idxs) > 0 {
		return instanceIDs[idxs[0]], nil
	}
	return "", nil
}

package validation

import "errors"

func ExtractSingleArgIgnoringOthers(args []string, idx int) (string, error) {
	if len(args) < idx {
		return "", errors.New("no file parameter provided")
	}

	return args[idx-1], nil
}

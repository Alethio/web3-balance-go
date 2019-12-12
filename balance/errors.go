package balance

import "fmt"

// RequestError wraps a request with the error received
type RequestError struct {
	Request *Request
	Err     error
}

func (re RequestError) section() string {
	return fmt.Sprintf("    %s %s %s", re.Request.DefaultBlockParam, re.Request.Address, re.Request.Currency)
}

// CollectBalancesError wraps errors returned from the RPC requests
type CollectError struct {
	Errors []*RequestError
}

// Error message aggregates all unique errors
func (ce CollectError) Error() string {
	fullErrorMessage := renderFullErrorMessage(ce.Errors)
	return fmt.Sprintf("Unable to collect balances because of these errors:\n%s", fullErrorMessage)
}

func renderFullErrorMessage(errors []*RequestError) string {
	var fullErrorMessage string
	errorMessages := make(map[string]string)
	index := 0

	for _, reqError := range errors {
		errorMessageKey := reqError.Err.Error()
		sectionError, ok := errorMessages[errorMessageKey]
		if !ok {
			sectionError = fmt.Sprintf("[%d] %s\n    Requests:\n", index, errorMessageKey)
			index++
		}
		sectionError += reqError.section() + "\n"
		errorMessages[errorMessageKey] = sectionError
	}

	for _, errorMessage := range errorMessages {
		fullErrorMessage += errorMessage
	}

	return fullErrorMessage
}

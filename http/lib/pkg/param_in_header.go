package httpadpt

const (
	TagRequestHeader = "header"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(TagRequestHeader, getHeaderInParamValue)
}

func getHeaderInParamValue(input Request, headerName string) (any, error) {
	var err error
	if IsRequestHeaderNil(input, &err) {
		return nil, err
	}
	value, found := input.Header().GetValue(headerName)
	if !found {
		return nil, nil
	}
	return value, nil
}

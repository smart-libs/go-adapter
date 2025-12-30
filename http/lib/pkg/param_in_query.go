package httpadpt

const (
	TagQuery = "query"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(TagQuery, getQueryInParamValue)
}

func getQueryInParamValue(input Request, flagName string) (any, error) {
	var err error
	if IsRequestQueryNil(input, &err) {
		return nil, err
	}
	value, found := input.Query().GetValue(flagName)
	if !found {
		return nil, nil
	}
	return value, nil
}

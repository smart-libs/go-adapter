package httpadpt

const (
	TagPath = "path"
)

func init() {
	getInputParamSpecFactoryRegistry().AddOption2(TagPath, getPathInParamValue)
}

func getPathInParamValue(input Request, pathName string) (any, error) {
	var err error
	if IsRequestPathNil(input, &err) {
		return nil, err
	}
	value, found := input.Path().GetValue(pathName)
	if !found {
		return nil, nil
	}
	return value, nil
}

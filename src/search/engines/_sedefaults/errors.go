package _sedefaults

func ReadErrorChannel(numberOfErrors int, errChannel chan error) []error {
	retErrors := make([]error, 0)
	for range numberOfErrors {
		err := <-errChannel
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}
	return retErrors
}

package _sedefaults

func NonNilErrorsFromSlice(errs []error) []error {
	nonNilErrs := make([]error, 0)
	counter := 0
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
			counter++
		}
	}
	return nonNilErrs
}

package main

func RemoveIndex(s []interface{}, index int) []interface{} {
	var value interface{} = &s
	// Now do the removal:
	sp := value.(*[]interface{})

	return append((*sp)[:index], (*sp)[index+1:]...)
}

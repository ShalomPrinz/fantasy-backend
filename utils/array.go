package utils

// Returns new array and removes duplication of objects in array.
// Note: duplication is checked via id property on each object
func ConcatDeduplicate[T any](first []T, second []T) []T {
	var withoutDups []T

	for _, element := range second {
		id := GetStringFromStruct(element, "ID")

		if !arrayHasId(first, id) {
			withoutDups = append(withoutDups, element)
		}
	}

	return append(first, withoutDups...)
}

func arrayHasId[T any](arr []T, id string) bool {
	for _, element := range arr {
		elementId := GetStringFromStruct(element, "ID")
		if elementId == id {
			return true
		}
	}

	return false
}

type fn func(any) string

func ArrayContainsString[T any](arr []T, value string, extractString fn) bool {
	for _, element := range arr {
		extractedValue := extractString(element)
		if extractedValue == value {
			return true
		}
	}

	return false
}

/* Returns T array mapped by f function to M array*/
func Map[T, M any](arr []T, f func(T) M) []M {
	mapped := make([]M, len(arr))
	for index, element := range arr {
		mapped[index] = f(element)
	}
	return mapped
}

/* Returns T array filtered by f function. if the result of f is value on an item, the item will be removed */
func RemoveItemByValue[T any](arr []T, value string, f func(T) string) []T {
	for i := 0; i < len(arr); i++ {
		if value == f(arr[i]) {
			newArr := make([]T, 0)
			newArr = append(newArr, arr[:i]...)
			return append(newArr, arr[i+1:]...)
		}
	}

	return arr
}

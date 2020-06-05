package errcodes

import "fmt"

//ErrMissingKey When we get missing keys
type ErrMissingKey struct {
	Key string
}

func (e ErrMissingKey) Error() string {
	return fmt.Sprintf("datastore missing key: %s", e.Key)
}

//IsErrMissingKey Returns if passed error is a missing key error
func IsErrMissingKey(err error) bool {
	if _, ok := err.(*ErrMissingKey); ok {
		return true
	}

	return false
}

type ErrDataSourceDoesntExist struct {
	Name string
}

func (e ErrDataSourceDoesntExist) Error() string {
	return fmt.Sprintf("datasource doesn't exist: %s", e.Name)
}

//IsDataSourceDoesntExist Returns if passed error is a missing key error
func IsDataSourceDoesntExist(err error) bool {
	if _, ok := err.(*ErrDataSourceDoesntExist); ok {
		return true
	}

	return false
}

type ErrFileExists struct {
	Filename string
}

func (e ErrFileExists) Error() string {
	return fmt.Sprintf("file already exists! Name: %s", e.Filename)
}

//ErrFileExists Returns if passed error is a missing key error
func IsErrFileExists(err error) bool {
	if _, ok := err.(*ErrFileExists); ok {
		return true
	}

	return false
}

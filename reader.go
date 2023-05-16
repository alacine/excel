package excel

import (
	"github.com/xuri/excelize/v2"
	"reflect"
)

// IReader interface
// All readers must implement this interface
type IReader interface {
	Unmarshall() error
	SetColumnsTags(tags map[string]*Tags)
}

// Reader is the Excel reader
type Reader struct {
	file  *excelize.File
	Sheet Sheet
	Axis  Axis
}

// validate validates the reader
// It returns an error if :
// - the sheet is not valid
// - the axis is not valid
func (r *Reader) validate() error {
	if r.file == nil {
		return ErrFileIsNil
	}
	if !r.isSheetValid() {
		return ErrSheetNotValid
	}
	if !r.isAxisValid() {
		return ErrAxisNotValid
	}
	return nil
}

// newReader create the appropriate reader
func (r *Reader) newReader(container any) (IReader, error) {
	// The type of the reader depends on the Container
	v := reflect.ValueOf(container)
	t := reflect.Indirect(v).Type()

	// validate the container
	// It must be a pointer to a slice
	if v.Kind() != reflect.Pointer && t.Kind() != reflect.Slice {
		return nil, ErrContainerInvalid
	}

	// Get the element type of the container
	e := t.Elem()
	if e.Kind() == reflect.Pointer {
		e = e.Elem()
	}

	// create the reader according to the
	// type of element
	switch e.Kind() {
	case reflect.Struct:
		reader, err := newStructReader(r, v)
		return reader, err
	case reflect.Map:
		reader, err := newMapReader(r, v)
		return reader, err
	case reflect.Slice, reflect.Array:
		reader, err := newSliceReader(r, v)
		return reader, err
	}
	return nil, ErrNoReaderFound
}

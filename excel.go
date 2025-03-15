package excel

import (
	"github.com/xuri/excelize/v2"
)

// Excel is the main structure that provides functionality for reading from and writing to Excel files.
// It contains references to the underlying excelize.File, Reader, Writer, and Struct components.
type Excel struct {
	File   *excelize.File
	Reader *Reader
	Writer *Writer

	Struct *Struct
}

// NewReader creates a new Excel reader from an existing excelize.File.
// It returns an error if the file is nil.
// The returned Excel instance can be used to unmarshal Excel data into Go structures.
func NewReader(file *excelize.File) (*Excel, error) {
	if file == nil {
		return nil, ErrFileIsNil
	}
	r := &Reader{
		file: file,
	}
	e := &Excel{
		File:   file,
		Reader: r,
	}
	return e, nil
}

// Unmarshal reads the Excel file and unmarshals it into the provided container.
// The container must be a pointer to a slice of structs, maps, or slices.
// Optional tags can be provided to customize the unmarshaling process.
// It returns an error if the Excel configuration is invalid or if unmarshaling fails.
func (e *Excel) Unmarshal(container any, tags ...map[string]*Tags) error {
	// validate excel input
	err := e.validate()
	if err != nil {
		return err
	}

	// Create the reader
	reader, err := e.Reader.newReader(container)
	if err != nil {
		return err
	}

	// Set column tags
	if len(tags) > 0 {
		reader.SetColumnsTags(tags[0])
	}

	// Check if reader is a struct reader
	if _, ok := reader.(*StructReader); ok {
		e.Struct = reader.(*StructReader).Struct
	}

	// unmarshall
	e.Reader.Result, err = reader.Unmarshall()
	return err
}

// NewWriter creates a new Excel writer from an existing excelize.File.
// It returns an error if the file is nil.
// The returned Excel instance can be used to marshal Go structures into Excel data.
func NewWriter(file *excelize.File) (*Excel, error) {
	if file == nil {
		return nil, ErrFileIsNil
	}
	w := &Writer{
		file: file,
	}
	e := &Excel{
		File:   file,
		Writer: w,
	}
	return e, nil
}

// Marshal writes the provided container into the Excel file.
// The container must be a pointer to a slice of structs, maps, or slices.
// Optional tags can be provided to customize the marshaling process.
// It returns an error if the Excel configuration is invalid or if marshaling fails.
func (e *Excel) Marshal(container any, tags ...map[string]*Tags) error {
	// validate excel input
	err := e.validate()
	if err != nil {
		return err
	}

	// Create the writer
	writer, err := e.Writer.newWriter(container)
	if err != nil {
		return err
	}

	// Set column tags
	if len(tags) > 0 {
		writer.SetColumnsTags(tags[0])
	}

	// Check if writer is a struct writer
	if _, ok := writer.(*StructWriter); ok {
		e.Struct = writer.(*StructWriter).Struct
	}

	// unmarshall
	e.Writer.Result, err = writer.Marshall(container)
	return err
}

// validate validates the Excel configuration.
// It returns an error if:
// - the file is nil
// - the reader is not valid
// - the writer is not valid
func (e *Excel) validate() error {
	if e.File == nil {
		return ErrFileIsNil
	}
	if e.Reader != nil {
		return e.Reader.validate()
	} else if e.Writer != nil {
		return e.Writer.validate()
	}
	return ErrConfigNotValid
}

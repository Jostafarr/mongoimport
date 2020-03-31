package loaders

import (
	"io"

	"github.com/clbanning/mxj"
)

// XMLLoader ...
type XMLLoader struct {
	readers    io.Reader
	resultChan chan map[string]interface{}
	workerChan chan bool
	done       bool
	// entries
}

// DefaultXMLLoader ..
func DefaultXMLLoader() *XMLLoader {
	return &XMLLoader{}
}

// Start ...
func (xmll *XMLLoader) Start() error {
	// entry, err := mxj.NewMapXmlReader(xmll.readers)
	// if err != nil {
	// 	return err
	//}
	return nil

}

// Describe ...
func (xmll *XMLLoader) Describe() string {
	return "XML"
}

// Finish ...
func (xmll *XMLLoader) Finish() error {
	return nil
}

// Create ...
func (xmll XMLLoader) Create(readers io.Reader, skipSanitization bool) ImportLoader {
	// if len(readers) < 1 {
	// 	return nil, fmt.Errorf("XML reader needs at least one input file")
	// }
	return &XMLLoader{
		readers: readers,
	}
}

// Type ...
// func (xmll *XMLLoader) Type() InputType {
// 	return MultipleInput
//}

// Load ...
func (xmll *XMLLoader) Load() (map[string]interface{}, error) {
	entry, err := mxj.NewMapXmlReader(xmll.readers)
	if err != nil {
		return nil, err
	}
	return entry, io.EOF
}

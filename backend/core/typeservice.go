package core

import "fmt"

// RawToTextFunc provides raw data to text or html conversion.
type RawToTextFunc func(raw []byte) (string, error)

// TypeProvider provides a single point of reference
// for all type specific functionality.
type TypeProvider struct {
	TypeNum int
	Name    string
	ToHTML  RawToTextFunc
	ToPlain RawToTextFunc
}

// TypeService maps Entry table's RawType number into
// usable functionality based on that number.
type TypeService struct {
	byType map[int]*TypeProvider
}

// NewTypeService creates new service with no mappings.
func NewTypeService() *TypeService {
	return &TypeService{
		byType: make(map[int]*TypeProvider),
	}
}

// AddProvider registers new provider.
// If provider exists it does nothing (does not update).
func (ts *TypeService) AddProvider(tp *TypeProvider) {
	_, exists := ts.byType[tp.TypeNum]
	if !exists {
		ts.byType[tp.TypeNum] = tp
	} else {

	}
}

// Provider returns provider for specified type num.
// Returns error if provider is not registered.
func (ts TypeService) Provider(typeNum int) (tp *TypeProvider, err error) {
	var exists bool
	tp, exists = ts.byType[typeNum]
	if exists {
		return tp, nil
	}
	return nil, fmt.Errorf("Cannot find type provider for type number %v", typeNum)
}

// Initialize populates service with default type providers.
func (ts *TypeService) Initialize() {
	ts.AddProvider(plainTextProvider())
}

// plainTextProvider implements simple plain text provider.
func plainTextProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 1,
		Name:    "Plain Text",
		ToHTML: func(raw []byte) (string, error) {
			return fmt.Sprintf("<pre>%v</pre>", string(raw)), nil
		},
		ToPlain: func(raw []byte) (string, error) {
			return string(raw), nil
		},
	}
}

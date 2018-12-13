package rx

import (
	"fmt"
	"reflect"
	"strings"
)

// genericType represents a any reflect.Type.
type genericType int

var genericTyp = reflect.TypeOf(new(genericType)).Elem()

// genericFunc is a type used to validate and call dynamic functions.
type genericFunc struct {
	MethodName string
	ParamName  string
	FnValue    reflect.Value
	FnType     reflect.Type
	TypesIn    []reflect.Type
	TypesOut   []reflect.Type
}

// Call calls a dynamic function.
func (g *genericFunc) Call(params ...interface{}) interface{} {
	paramsIn := make([]reflect.Value, len(params))
	for i, param := range params {
		paramsIn[i] = reflect.ValueOf(param)
	}
	paramsOut := g.FnValue.Call(paramsIn)
	if len(paramsOut) >= 1 {
		return paramsOut[0].Interface()
	}
	return nil
}

// newGenericFunc instantiates a new genericFunc pointer
func newGenericFunc(methodName, paramName string, fn interface{}, validateFunc func(*genericFunc) error) (*genericFunc, error) {
	fnValue := reflect.ValueOf(fn)

	if fnValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("%s: parameter [%s] is not a function type. It is a '%s'", methodName, paramName, fnValue.Type())
	}

	genericFn := &genericFunc{
		FnValue:    fnValue,
		FnType:     fnValue.Type(),
		MethodName: methodName,
		ParamName:  paramName,
	}

	numTypesIn := genericFn.FnType.NumIn()
	genericFn.TypesIn = make([]reflect.Type, numTypesIn)
	for i := 0; i < numTypesIn; i++ {
		genericFn.TypesIn[i] = genericFn.FnType.In(i)
	}

	numTypesOut := genericFn.FnType.NumOut()
	genericFn.TypesOut = make([]reflect.Type, numTypesOut)
	for i := 0; i < numTypesOut; i++ {
		genericFn.TypesOut[i] = genericFn.FnType.Out(i)
	}
	if err := validateFunc(genericFn); err != nil {
		return nil, err
	}

	return genericFn, nil
}

// simpleParamValidator creates a function to validate genericFunc based in the
// In and Out function parameters.
func simpleParamValidator(In []reflect.Type, Out []reflect.Type) func(genericFn *genericFunc) error {
	return func(fn *genericFunc) error {
		var isValid = func() bool {
			if In != nil {
				if len(In) != len(fn.TypesIn) {
					return false
				}
				for i, paramIn := range In {
					if paramIn != genericTyp && paramIn != fn.TypesIn[i] {
						return false
					}
				}
			}
			if Out != nil {
				if len(Out) != len(fn.TypesOut) {
					return false
				}
				for i, paramOut := range Out {
					if paramOut != genericTyp && paramOut != fn.TypesOut[i] {
						return false
					}
				}
			}
			return true
		}

		if !isValid() {
			return fmt.Errorf("%s: parameter [%s] has a invalid function signature. Expected: '%s', actual: '%s'", fn.MethodName, fn.ParamName, formatFuncSignature(In, Out), formatFuncSignature(fn.TypesIn, fn.TypesOut))
		}
		return nil
	}
}

// newElemTypeSlice creates a slice of items elem types.
func newElemTypeSlice(items ...interface{}) []reflect.Type {
	typeList := make([]reflect.Type, len(items))
	for i, item := range items {
		typeItem := reflect.TypeOf(item)
		if typeItem.Kind() == reflect.Ptr {
			typeList[i] = typeItem.Elem()
		}
	}
	return typeList
}

// formatFuncSignature formats the func signature based in the parameters types.
func formatFuncSignature(In []reflect.Type, Out []reflect.Type) string {
	paramInNames := make([]string, len(In))
	for i, typeIn := range In {
		if typeIn == genericTyp {
			paramInNames[i] = "T"
		} else {
			paramInNames[i] = typeIn.String()
		}

	}
	paramOutNames := make([]string, len(Out))
	for i, typeOut := range Out {
		if typeOut == genericTyp {
			paramOutNames[i] = "T"
		} else {
			paramOutNames[i] = typeOut.String()
		}
	}
	return fmt.Sprintf("func(%s)%s", strings.Join(paramInNames, ","), strings.Join(paramOutNames, ","))
}

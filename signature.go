package jdwp

import (
	"bytes"
	"strings"
)

const (
	SIGNATURE_FUNC     = '('
	SIGNATURE_ENDFUNC  = ')'
	SIGNATURE_ENDCLASS = ';'
	TagArray           = 91 // '[' - an array object (objectID size).
	TagByte            = 66 // 'B' - a byte value (1 byte).
	TagChar            = 67 // 'C' - a character value (2 bytes).
	TagObject          = 76 // 'L' - an object (objectID size).
	TagFloat           = 70 // 'F' - a float value (4 bytes).
	TagDouble          = 68 // 'D' - a double value (8 bytes).
	TagInt             = 73 // 'I' - an int value (4 bytes).
	TagLong            = 74 // 'J' - a long value (8 bytes).
	TagShort           = 83 // 'S' - a short value (2 bytes).
	TagVoid            = 86 // 'V' - a void value (no bytes).
	TagBoolean         = 90 // 'Z' - a boolean value (1 byte).
)

type signatureParse struct {
	Signature    string
	TypeNames    []string
	CurrentIndex int
}

// TranslateSignatureToClassName SignatureTranslate 类名签名转化为类名
func TranslateSignatureToClassName(sign string) string {
	parse := signatureParse{}
	parse.Signature = sign
	typeNames := parse.getTypeNames()
	return typeNames[len(typeNames)-1]
}

func TranslateClassNameToSignature(className string) string {
	var result []byte
	tmpClassName := []byte(className)
	left := []byte("[")
	firstIndex := bytes.Index(tmpClassName, left)
	index := firstIndex
	for index != -1 {
		result = append(result, left...)
		index = bytes.Index(tmpClassName, tmpClassName[index+1:])
	}

	if firstIndex != -1 {
		tmpClassName = tmpClassName[0:firstIndex]
	}

	signature := string(tmpClassName)
	if signature == "boolean" {
		result = append(result, 'Z')
	} else if signature == "byte" {
		result = append(result, 'B')
	} else if signature == "char" {
		result = append(result, 'C')
	} else if signature == "short" {
		result = append(result, 'S')
	} else if signature == "int" {
		result = append(result, 'I')
	} else if signature == "long" {
		result = append(result, 'J')
	} else if signature == "float" {
		result = append(result, 'F')
	} else if signature == "double" {
		result = append(result, 'D')
	} else {
		result = append(result, 'L')
		result = append(result, []byte(strings.ReplaceAll(signature, ".", "/"))...)
		result = append(result, ';')
	}
	return string(result)
}

func (s *signatureParse) getTypeNames() []string {
	if s.TypeNames == nil {
		s.TypeNames = make([]string, 10)
	}
	var elem string
	s.CurrentIndex = 0
	for s.CurrentIndex < len(s.Signature) {
		elem = s.nextTypeName()
		s.TypeNames = append(s.TypeNames, elem)
	}
	if len(s.TypeNames) == 0 {
		panic("Invalid JNI signature'" + string(s.Signature) + "'")
	}
	return s.TypeNames
}
func (s *signatureParse) nextTypeName() string {
	key := s.Signature[s.CurrentIndex]
	s.CurrentIndex++
	var result string
	switch key {
	case TagArray:
		result = s.nextTypeName() + "[]"
	case TagByte:
		result = "byte"
	case TagChar:
		result = "char"
	case TagObject:
		endIndex := indexOf(s.Signature, string(SIGNATURE_ENDCLASS), s.CurrentIndex)
		retVal := s.Signature[s.CurrentIndex:endIndex]
		retVal = strings.ReplaceAll(retVal, "/", ".")
		s.CurrentIndex = endIndex + 1
		result = retVal
	case TagFloat:
		result = "float"
	case TagDouble:
		result = "double"
	case TagInt:
		result = "int"
	case TagLong:
		result = "long"
	case TagShort:
		result = "short"
	case TagVoid:
		result = "void"
	case TagBoolean:
		result = "boolean"
	case SIGNATURE_ENDFUNC:
		s.nextTypeName()
	case SIGNATURE_FUNC:
		s.nextTypeName()
	default:
		panic("Invalid JNI signature character '" + s.Signature + "'")
	}
	return result
}
func indexOf(str, substr string, start int) int {
	if start == 0 || start == len(str) {
		return -1
	}
	index := strings.Index(str[start:], substr)
	if index == -1 {
		return -1
	}
	return index + start
}

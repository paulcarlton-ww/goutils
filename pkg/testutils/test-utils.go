package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test utilities

// CheckFieldValue will check if a field value is equal to the expected value and set the test to failed if not.
func CheckFieldValue(u TestUtil, fieldName string, fieldInfo FieldInfo) bool {
	test := u.TestData()
	t := u.Testing()

	if actual := test.ObjStatus.GetField(t, test.ObjStatus.Object, fieldName); actual != fieldInfo.FieldValue || u.FailTests() {
		t.Errorf("\nTest: %d, %s\nField: %s\nGot.....: %s\nExpected: %s",
			test.Number, test.Description, fieldName, spew.Sdump(actual), spew.Sdump(fieldInfo.FieldValue))

		if !u.FailTests() {
			return false
		}
	}

	return true
}

// CheckFieldGetter will the getter function of a field and check the value is equal to the expected value and set the test to failed if not.
func CheckFieldGetter(u TestUtil, fieldName string, fieldInfo FieldInfo) bool {
	test := u.TestData()
	t := u.Testing()

	if len(fieldInfo.GetterMethod) > 0 {
		if results := test.ObjStatus.CallMethod(t, test.ObjStatus.Object, fieldInfo.GetterMethod, []interface{}{}); results[0] != fieldInfo.FieldValue || u.FailTests() {
			t.Errorf("\nTest: %d, %s\nField: %s, Getter function: %s\nGot.....: %s\nExpected: %s",
				test.Number, test.Description, fieldName, fieldInfo.GetterMethod, spew.Sdump(results), spew.Sdump([]interface{}{fieldInfo.FieldValue}))

			if !u.FailTests() {
				return false
			}
		}
	} else {
		t.Logf("field getter method not set, skipping check of field: %s getter", fieldName)
	}

	return true
}

// CheckFieldsValue will get the value of object fields and check the expected values are equal to the actual value and set the test to failed if not.
func CheckFieldsValue(u TestUtil) bool {
	test := u.TestData()
	t := u.Testing()

	if test.ObjStatus.GetField == nil {
		t.Logf("object Get Field function not set, skipping check of fields values")

		return true
	}

	for fieldName, fieldInfo := range test.ObjStatus.Fields {
		if !CheckFieldValue(u, fieldName, fieldInfo) {
			return false
		}
	}

	return true
}

// CheckFieldsGetter will call the getter functions on an object to check the expected values are equal to the actual value and set the test to failed if not.
func CheckFieldsGetter(u TestUtil) bool {
	test := u.TestData()
	t := u.Testing()

	if test.ObjStatus.CallMethod == nil {
		t.Logf("object call method function not set, skipping check of fields getter functions")

		return true
	}

	for fieldName, fieldInfo := range test.ObjStatus.Fields {
		if !CheckFieldGetter(u, fieldName, fieldInfo) {
			return false
		}
	}

	return true
}

// CastToStringList casts an interface to a []string.
func CastToStringList(t *testing.T, data interface{}) []string {
	strList, ok := data.([]string)
	if !ok {
		t.Fatalf("failed to cast interface to []string")

		return nil
	}

	return strList
}

// CastToMapStringString casts an interface to a map[string]string.
func CastToMapStringString(t *testing.T, data interface{}) map[string]string {
	strMap, ok := data.(map[string]string)
	if !ok {
		t.Fatalf("failed to cast interface to map[string]string")

		return nil
	}

	return strMap
}

// UnsetEnvs unsets a list of environmental variables.
func UnsetEnvs(t *testing.T, names []string) {
	for _, env := range names {
		UnsetEnv(t, env)

		if t.Failed() {
			return
		}
	}
}

// MapKeysToList returns the keys of a map[string]string as a list.
func MapKeysToList(envsMap map[string]string) []string {
	keys := make([]string, 0, len(envsMap))

	for k := range envsMap {
		keys = append(keys, k)
	}

	return keys
}

// SetEnv sets an environmental variable.
func SetEnv(t *testing.T, envName, envValue string) {
	if err := os.Setenv(envName, envValue); err != nil {
		t.Fatalf("failed to set environmental variable: %s to %s, %s", envName, envValue, err)

		return
	}

	t.Logf("environmental variable: %s, set to %s", envName, envValue)
}

// UnsetEnv unsets an environmental variable.
func UnsetEnv(t *testing.T, envName string) {
	if err := os.Unsetenv(envName); err != nil {
		t.Fatalf("failed to unset environmental variable: %s , %s", envName, err)

		return
	}
}

// SetEnvs sets a number of environmental variables.
func SetEnvs(t *testing.T, envSettings map[string]string) {
	for name, val := range envSettings {
		SetEnv(t, name, val)

		if t.Failed() {
			return
		}
	}
}

// RemoveBottom removes items from caller list including and prior to testing.tRunner().
func RemoveBottom(callers []string) []string {
	ourCallers := []string{}

	for _, caller := range callers {
		if strings.HasPrefix(caller, "testing.tRunner()") {
			break
		}

		ourCallers = append(ourCallers, caller)
	}

	return ourCallers
}

// CompareWhereList compares a list of strings returned by GetCaller or Callers but ignores line numbers.
func CompareWhereList(one, two []string) bool {
	if len(one) == 0 && len(two) == 0 { // Both empty return true
		return true
	}

	if len(two) != len(one) {
		return false
	}

	for index, field := range one {
		if !CompareWhere(field, two[index]) {
			return false
		}
	}

	return true
}

// CompareWhere compares strings returned by GetCaller or Callers but ignores line numbers.
func CompareWhere(one, two string) bool {
	if strings.HasSuffix(one, "(NN)") || strings.HasSuffix(two, "(NN)") {
		return one[:strings.LastIndex(one, "(")] == two[:strings.LastIndex(two, "(")]
	}

	return one == two
}

// CompareItems compares two items, returning true if both are nil.
func CompareItems(one, two interface{}) bool {
	if one == nil && two == nil {
		return true
	}

	if one == nil || two == nil {
		return false
	}

	return fmt.Sprintf("%+v", one) == fmt.Sprintf("%+v", two)
}

// ContainsStringArray checks items in two string arrays verifying that the second string array is
// contained in the first.
func ContainsStringArray(one, two []string, first bool) bool {
	if len(one) == 0 && len(two) == 0 { // Both empty return true
		return true
	}

	if len(two) == 0 {
		return false
	}

	offset := 0

	for index, field := range one {
		if field == two[offset] {
			offset++
			index++

			if index < len(one) && offset < len(two) {
				return ContainsStringArray(one[index:], two[offset:], true)
			}

			if offset == len(two) {
				return true
			}
		}

		if first {
			return false
		}
	}

	return false
}

// ReadBuf reads lines from a buffer into a string array.
func ReadBuf(out io.Reader) *[]string {
	output := []string{}

	for {
		p := make([]byte, 1024000)
		n, err := out.Read(p)

		if n > 0 {
			lines := strings.Split(string(p[:n]), "\n")
			output = append(output, lines...)
		}

		if err != nil {
			break
		}
	}

	// If last line is empty then remove it
	if len(output) > 0 && len(output[len(output)-1]) == 0 {
		output = output[:len(output)-1]
	}

	return &output
}

// DisplayStrings returns a []string as a single string, one element per line.
func DisplayStrings(strs []string) string {
	output := ""
	newline := ""

	for index, str := range strs {
		output = fmt.Sprintf("%s%s%d - %s", output, newline, index, str)
		newline = "\n"
	}

	return output
}

// ToJSON is used to convert a data structure into JSON format.
// nolint godox ToDo address circular depencancies with internal/common and remove this
func ToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, jsonData, "", "\t")
	if err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}

// GetTestJSON returns json representation of an interface or sets test error if it fails.
func GetTestJSON(t *testing.T, data interface{}) string {
	jsonResp, err := ToJSON(data)
	if err != nil {
		t.Errorf("failed to convert data: %+v to json, %s", data, err)
	}

	return jsonResp
}

func HandlePanic(t *testing.T) {
	if a := recover(); a != nil {
		t.Fatalf("%s", a)
	}
}

// CallMethod calls a method on an object and returns the results.
func CallMethod(t *testing.T, obj interface{}, methodName string, params []interface{}) []interface{} {
	defer HandlePanic(t)

	method := reflect.ValueOf(obj).MethodByName(methodName)
	p := []reflect.Value{}

	for _, param := range params {
		p = append(p, reflect.ValueOf(param))
	}

	values := method.Call(p)
	results := []interface{}{}

	for _, result := range values {
		results = append(results, result.Interface())
	}

	return results
}

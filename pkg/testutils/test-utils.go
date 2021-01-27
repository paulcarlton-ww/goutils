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
)

// Test utilities

// FailTests determines if tests are set to failed so expected and results are emitted even if tests pass.
// Set to true to test failure output.
var FailTests = false // nolint:gochecknoglobals // ok

// FieldInfo holds information about a field of a struct
type FieldInfo struct {
	GetterMethod string  `json:"getter,omitempty"` // The method to get the value, nil if no getter method
	SetterMethod string  `json:"setter,omitempty"` // The method to get the value, nil if no setter method
	FieldName    string  `json:"name"` // The name of the field, required value
	kind         reflect.Kind `json:"type,omitempty` // The type of the field, defaults to string
}

func fakeLogger() logr.Logger {
	return testlogr.NullLogger{}
}

func getBoolField(t *testing.T, obj interface{}, field string) bool {
	value, err := config.GetField(obj, field)
	if err != nil {
		t.Fatalf("failed to get field '%s' value, %s", field, err)

		return false
	}

	boolValue, ok := value.(bool)
	if !ok {
		t.Fatalf("failed to cast field '%s' value, %s to bool", field, err)

		return false
	}
	return boolValue
}

func getStringField(t *testing.T, obj interface{}, field string) string {
	value, err := obj.GetField(obj, field)
	if err != nil {
		t.Fatalf("failed to get field '%s' value %s", field, err)

		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		t.Fatalf("failed to cast field '%s' value, %s to string", field, err)

		return false
	}
	return boolValue
}

func getIntField(t *testing.T, obj interface{}, field string) int {
	value, err := config.GetField(obj, field)
	if err != nil {
		t.Fatalf("failed to get field value, %s", err)

		return 0
	}

	return int(value.Int())
}

func getDurationField(t *testing.T, obj interface{}, field string) time.Duration {
	value, err := config.GetField(obj, field)
	if err != nil {
		t.Fatalf("failed to get field value, %s", err)

		return 0
	}

	return time.Duration(value.Int())
}

func getFieldString(t *testing.T, obj interface{}, k string, v TypeValue) string { // nolint:deadcode,unused // ok
	switch v.name {
	case stringType:
		return getStringField(t, obj, k)
	case boolType:
		{
			return fmt.Sprintf("%t", getBoolField(t, obj, k))
		}
	case durationType:
		{
			return getDurationField(t, obj, k).String()
		}
	case intType:
		{
			return fmt.Sprintf("%d", getIntField(t, obj, k))
		}
	default:
		t.Fatalf("Invalid field: %s", k)

		return ""
	}
}

func getTypeString(t *testing.T, v TypeValue) string { // nolint:deadcode,unused // ok
	switch v.name {
	case stringType:
		str, ok := v.value.(string)
		if !ok {
			t.Fatalf("failed to cast value to string: %s", v)

			return ""
		}

		return str
	case boolType:
		{
			val, ok := v.value.(bool)
			if !ok {
				t.Fatalf("failed to cast value to bool: %s", v)

				return ""
			}

			return fmt.Sprintf("%t", val)
		}
	case durationType:
		{
			val, ok := v.value.(time.Duration)
			if !ok {
				t.Fatalf("failed to cast value to duration: %s", v)

				return ""
			}

			return val.String()
		}
	case intType:
		{
			val, ok := v.value.(int)
			if !ok {
				t.Fatalf("failed to cast value to int: %s", v)

				return ""
			}

			return fmt.Sprintf("%d", val)
		}
	default:
		t.Fatalf("Invalid field: %s", v)

		return ""
	}
}

func CheckField(t *testing.T, getFieldFunc GetFieldFunc, k string, v interface{}) bool {
	switch v.(type) {
	case string:
		actual := getFieldFunc(t, k)
		if actual != v {
			t.Fatalf("Field: %s, actual: %s, expected: %s", k, actual, v)

			return false
		}
	case bool:
		{
			actual := getBoolField(t, obj, k)
			if actual != v.value {
				t.Fatalf("Field: %s, actual: %t, expected: %t", k, actual, v.value)

				return false
			}
		}
	case reflect.:
		{
			actual := getDurationField(t, obj, k)
			if actual != v.value {
				t.Fatalf("Field: %s, actual: %s, expected: %s", k, actual, v.value)

				return false
			}
		}
	case reflect.Int:
		{
			actual := getIntField(t, obj, k)
			if actual != v.value {
				t.Fatalf("Field: %s, actual: %d, expected: %d", k, actual, v.value)

				return false
			}
		}
	default:
		t.Fatalf("Invalid field: %s", k)

		return false
	}

	return true
}

func CheckFields(t *testing.T, obj interface{}, fields map[string]interface{}) bool {
	for k, v := range fields {
		if !CheckField(t, obj, k, v) {
			return false
		}
	}

	return true
}

// CastToStringList casts an interface to a []string
func CastToStringList(t *testing.T, data interface{}) []string {
	strList, ok := data.([]string)
	if !ok {
		t.Fatalf("failed to cast interface to []string")

		return nil
	}
	retrun strList
}

// CastToMapStringString casts an interface to a map[string]string
func CastToMapStringString(t *testing.T, data interface{}) map[string]string {
	strMap, ok := data.(map[string]string)
	if !ok {
		t.Fatalf("failed to cast interface to map[string]string")

		return nil
	}
	retrun strMap
}

// UnsetEnvs unsets a list of environmental variables
func UnsetEnvs(t *testing.T, names []string) {
	for _, env := range names {
		UnsetEnv(t, env)
		if t.Failed() {
			return
		}
	}
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

// SetEnvs sets a number of environmental variables
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

func init() { // nolint:gochecknoinits // ok
	_, present := os.LookupEnv("FAILED_OUTPUT_TEST")

	if present {
		FailTests = true
	}
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

// CallMethod calls a method on an object and returns the results.
func CallMethod(obj interface{}, methodName string, params []interface{}) (results []interface{}, err error) {
	method := reflect.ValueOf(obj).MethodByName(methodName)
	p := []reflect.Value{}

	for _, param := range params {
		p = append(p, reflect.ValueOf(param))
	}

	values := method.Call(p)

	for _, result := range values {
		results = append(results, result.Interface())
	}

	return results, nil
}

// Package testutils provides a framework and helper functions for use during unit testing.
package testutils

import (
	"os"
	"reflect"
	"testing"
)

type (
	// PrepTestI defines function to be called before running a test.
	PrepTestI func(u TestUtil)
	// CheckTestI defines function to be called after test to check result.
	CheckTestI func(u TestUtil) bool

	// GetFieldFunc is the function to call to get the value of a field of an object.
	GetFieldFunc func(t *testing.T, obj interface{}, fieldName string) interface{}
	// SetFieldFunc is the function to call to set the value of a field of an object.
	SetFieldFunc func(t *testing.T, obj interface{}, fieldName string, value interface{})
	// CallMethodFunc is the function to call a method on an object.
	CallMethodFunc func(t *testing.T, obj interface{}, methodName string, values []interface{}) []interface{}

	// FieldInfo holds information about a field of a struct.
	FieldInfo struct {
		Comparer     func(t *testing.T, actual, expected interface{}) bool `json:"comparer,omitempty"` // Function to do field specific compare, nil if not set.
		GetterMethod string                                                `json:"getter,omitempty"`   // The method to get the value, nil if no getter method.
		SetterMethod string                                                `json:"setter,omitempty"`   // The method to get the value, nil if no setter method.
		FieldValue   interface{}                                           `json:"value"`              // The value to set or expected value to verify.
	}

	// Fields is a map of field names to information about the field.
	Fields map[string]FieldInfo

	// ObjectStatus hold details of the object under test.
	// This can be used to verify the internal state of an object after a test.
	ObjectStatus struct {
		Object     interface{}    // The object or interface under test, this needs to be set during test before calling post test actions.
		GetField   GetFieldFunc   // The function to call to get a field value.
		SetField   SetFieldFunc   // The function to call to set a field value.
		CallMethod CallMethodFunc // The function to call a method on an object.
		Fields     Fields         // The fields of an object.
	}

	// DefTest generic test data structure.
	DefTest struct {
		Number      int           // Test number.
		Description string        // Test description.
		Config      interface{}   // Test configuration information to be used by test function or custom pre/post test functions.
		Inputs      []interface{} // Test inputs.
		Expected    []interface{} // Test expected results.
		Results     []interface{} // Test results.
		ObjStatus   *ObjectStatus // Details of object under test including field names and expected values, used by CheckFunc to verify values.
		PrepFunc    PrepTestI     // Function to be called before a test.
		// leave unset to call default - which prints the test number and name.
		CheckFunc CheckTestI // Function to be called to check a test results.
		// leave unset to call default - which compares actual results with expected results and verifies object status.
	}

	// TestUtil the interface used to provide testing utilities.
	TestUtil interface {
		CallPrepFunc()                 // Call the custom or default test preparation function.
		CallCheckFunc() bool           // Call the custom or default test checking function.
		Testing() *testing.T           // testing object.
		SetFailTests(value bool)       // Set the fail test setting to verify test check reporting.
		FailTests() bool               // Get the fail test setting.
		SetVerbose(value bool)         // Set the verbose setting.
		Verbose() bool                 // Get the verbose setting.
		SetTestData(testData *DefTest) // Set the test data.
		TestData() *DefTest            // Get the test data.
	}

	// testUtil is used to hold configuration information for testing.
	testUtil struct {
		TestUtil             // TestUtil interface that operates on this object.
		t         *testing.T // Testing object.
		testData  *DefTest   // The definition of this test.
		failTests bool       // Set to make default test check function reported retrun false to test report function.
		verbose   bool       // Set to make testutils more verbose
	}
)

// NewTestUtil retruns a new TestUtil interface.
func NewTestUtil(t *testing.T, testData *DefTest) TestUtil {
	u := &testUtil{failTests: false}
	u.t = t
	u.testData = testData

	_, present := os.LookupEnv("TESTUTILS_FAIL")
	if present {
		u.failTests = true
	}

	_, present = os.LookupEnv("TESTUTILS_VERBOSE")
	if present {
		u.verbose = true
	}

	return u
}

// CallPrepFunc calls the pre test setup function.
func (u *testUtil) CallPrepFunc() {
	DefaultPrepFunc(u)

	if u.testData.PrepFunc != nil {
		u.testData.PrepFunc(u)
	}
}

// CallCheckTestsFunc calls the check test result function.
func (u *testUtil) CallCheckFunc() bool {
	if u.testData.CheckFunc == nil {
		return DefaultCheckFunc(u)
	}

	return u.testData.CheckFunc(u)
}

// Testing returns the testing object.
func (u *testUtil) Testing() *testing.T {
	return u.t
}

// SetVerbose sets the verbose flag.
func (u *testUtil) SetVerbose(value bool) {
	u.verbose = value
}

// Verbose gets the verbose flag.
func (u *testUtil) Verbose() bool {
	return u.verbose
}

// SetFailTests sets the fail tests flag.
func (u *testUtil) SetFailTests(value bool) {
	u.failTests = value
}

// FailTests returns the fail test setting.
func (u *testUtil) FailTests() bool {
	return u.failTests
}

// SetTestData sets the test data.
func (u *testUtil) SetTestData(testData *DefTest) {
	u.testData = testData
}

// TestData returns the test data.
func (u *testUtil) TestData() *DefTest {
	return u.testData
}

// DefaultPrepFunc is the default pre test function that prints the test number and name.
func DefaultPrepFunc(u TestUtil) {
	test := u.TestData()
	u.Testing().Logf("Test: %d, %s\n", test.Number, test.Description)
}

// DefaultCheckFunc is the default check test function that compares actual and expected.
func DefaultCheckFunc(u TestUtil) bool {
	return CheckCallResultsReflect(u) && CheckObjStatusFunc(u)
}

// CheckObjStatusFunc checks object fields values against expected and report if different.
func CheckObjStatusFunc(u TestUtil) bool {
	return CheckFieldsValue(u) && CheckFieldsGetter(u)
}

// CheckCallResultsReflect uses reflect.DeepEqual tochecks call results against expected and report if different .
func CheckCallResultsReflect(u TestUtil) bool {
	test := u.TestData()

	result := reflect.DeepEqual(test.Results, test.Expected)
	if !result || u.FailTests() {
		ReportSpew(u)
	}

	return result
}

// CheckCallResultsJSON uses json comparison to checks call results against expected and report if different .
func CheckCallResultsJSON(u TestUtil) bool {
	test := u.TestData()
	t := u.Testing()

	result := CompareAsJSON(t, test.Results, test.Expected)
	if !result || u.FailTests() {
		ReportJSON(u)
	}

	return result
}

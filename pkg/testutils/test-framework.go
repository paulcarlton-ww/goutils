package testutils

import (
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type (
	// PrepTestI defines function to be called before running a test.
	PrepTestI func(u TestUtil)
	// CheckTestI definesfunction to be called after test to check result.
	CheckTestI func(u TestUtil) bool

	// GetFieldFunc is the function to call to get the value of a field of an object.
	GetFieldFunc func(t *testing.T, obj interface{}, fieldName string) interface{}
	// SetFieldFunc is the function to call to set the value of a field of an object.
	SetFieldFunc func(t *testing.T, obj interface{}, fieldName string, value interface{})
	// CallMethodFunc is the function to call a method on an object.
	CallMethodFunc func(t *testing.T, obj interface{}, methodName string, values []interface{}) []interface{}

	// FieldInfo holds information about a field of a struct.
	FieldInfo struct {
		GetterMethod string      `json:"getter,omitempty"` // The method to get the value, nil if no getter method.
		SetterMethod string      `json:"setter,omitempty"` // The method to get the value, nil if no setter method.
		FieldValue   interface{} `json:"value"`            // The value to set or expected value to verify.
	}

	// Fields is a map of field names to information about the field.
	Fields map[string]FieldInfo

	// ObjectStatus hold details of the object under test, including the object, the functions to get and set fields and call methods.
	ObjectStatus struct {
		Object     interface{}    // The object or interface under test, this needs to be set during test before calling post test actions.
		GetField   GetFieldFunc   // The function to call to get a field value.
		SetField   SetFieldFunc   // The function to call to set a field value.
		CallMethod CallMethodFunc // The function to call a method on an object.
		Fields     Fields         // The fields of an object.
	}

	// DefTest generic tests data structure used by tests.
	DefTest struct {
		Number      int           // Test number.
		Description string        // Test description.
		Config      interface{}   // Test configuration, to be used by custom preTest Function.
		Inputs      []interface{} // Test inputs.
		Expected    []interface{} // Test expected results.
		Results     []interface{} // Test results.
		ObjStatus   ObjectStatus  // Details of object under test including field names and expected values, used by CheckFunc to verify values.
		PrepFunc    PrepTestI     // Function to be called before a test.
		// leave unset to call default - which prints the test number and name.
		CheckFunc CheckTestI // Function to be called to check a test results.
		// leave unset to call default - which compares actual and expected as strings.
	}

	// TestUtil the interface used to provide testing utilities.
	TestUtil interface {
		CallPrepFunc()       // Call the custom ot defsult test preparation function.
		CallCheckFunc() bool // Call the custom or default test checking function.
		Testing() *testing.T
		SetFailTests(value bool)
		FailTests() bool
		SetTestData(testData *DefTest)
		TestData() *DefTest
	}

	// testUtil is used to hold configuration information for testing.
	testUtil struct {
		TestUtil             // TestUtil interface that operates in this object.
		t         *testing.T // Testing object.
		testData  *DefTest   // The definition of this test.
		failTests bool       // Set to make default test check function reported retrun false to test report function.
	}
)

// NewTestUtil retruns a new TestUtil interface.
func NewTestUtil(t *testing.T, testData *DefTest) TestUtil {
	u := &testUtil{failTests: false}
	u.t = t
	u.testData = testData

	_, present := os.LookupEnv("FAILED_OUTPUT_TEST")
	if present {
		u.failTests = true
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

func (u *testUtil) Testing() *testing.T {
	return u.t
}

func (u *testUtil) SetFailTests(value bool) {
	u.failTests = value
}

func (u *testUtil) FailTests() bool {
	return u.failTests
}

func (u *testUtil) SetTestData(testData *DefTest) {
	u.testData = testData
}

func (u *testUtil) TestData() *DefTest {
	return u.testData
}

// DefaultPrepFunc is the default pre test function that prints the test number and name.
func DefaultPrepFunc(u TestUtil) {
	test := u.TestData()
	u.Testing().Logf("Test: %d, %s\n", test.Number, test.Description)
}

// DefaultCheckFunc is the default check test function that compares actual and expected as strings.
func DefaultCheckFunc(u TestUtil) bool {
	return CheckCallResultsFunc(u) && CheckObjStatusFunc(u)
}

// CheckObjStatusFunc checks object fields values against expected and report if different.
func CheckObjStatusFunc(u TestUtil) bool {
	return CheckFieldsValue(u) && CheckFieldsGetter(u)
}

// CheckCallResultsFunc checks call results against expected and report if different .
func CheckCallResultsFunc(u TestUtil) bool {
	test := u.TestData()
	t := u.Testing()

	if !reflect.DeepEqual(test.Results, test.Expected) || u.FailTests() {
		t.Errorf("\nTest: %d, %s\nInput...: %s\nGot.....: %s\nExpected: %s",
			test.Number, test.Description, spew.Sdump(test.Inputs), spew.Sdump(test.Results), spew.Sdump(test.Expected))

		if !u.FailTests() {
			return false
		}
	}

	return true
}

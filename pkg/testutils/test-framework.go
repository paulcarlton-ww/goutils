package testutils

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type (
	// PrepTestI defines function to be called before running a test.
	PrepTestI func(t *testing.T, test *DefTest)
	// CheckTestI definesfunction to be called after test to check result.
	CheckTestI func(t *testing.T, test *DefTest) bool
	// ReportTestI defines function to be called to report test results.
	ReportTestI func(t *testing.T, test *DefTest)

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
		ReportFunc ReportTestI // Function to be called to report test results.
		// leave unset to call default - which reports input, actual and expected as strings.
	}

	// TestUtil the interface used to provide testing utilities.
	TestUtil interface {
		CallPrepFunc()                          // Call the custom ot defsult test preparation function.
		CallCheckFunc() bool                    // Call the custom or default test checking function.
		CallReportFunc()                        // Call the custom or default test reporting function.
		CallPostTestActions() bool // Calls custom or default check and reporting functions.
		SetFailTests(value bool)
		GetFailTests() bool
		SetTestData(testData *DefTest)
		GetTestData() *DefTest
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

	return u
}

// CallPrepFunc calls the pre test setup function.
func (u *testUtil) CallPrepFunc() {
	if u.testData.PrepFunc == nil {
		DefaultPrepFunc(u.t, u.testData)

		return
	}

	u.testData.PrepFunc(u.t, u.testData)
}

// CallCheckTestsFunc calls the check test result function.
func (u *testUtil) CallCheckFunc() bool {
	if u.testData.CheckFunc == nil {
		return DefaultCheckFunc(u.t, u.testData)
	}

	return u.testData.CheckFunc(u.t, u.testData)
}

// CallReportFunc calls the report test results function.
func (u *testUtil) CallReportFunc() {
	if u.testData.ReportFunc == nil {
		DefaultReportFunc(u.t, u.testData)

		return
	}

	u.testData.ReportFunc(u.t, u.testData)
}

// PostTestActions call after test to call check function and report function if check fails.
func (u *testUtil) PostTestActions() bool {
	if !u.CallCheckFunc() {
		u.t.Logf("Test: %d, %s, failed", u.testData.Number, u.testData.Description)
		u.CallReportFunc()

		return false
	}

	return true
}

func (u *testUtil) SetFailTests(value bool) {
	u.failTests = value
}

func (u *testUtil) GetFailTests() bool {
	return u.failTests
}

func (u *testUtil) SetTestData(testData *DefTest) {
	u.testData = testData
}

func (u *testUtil) GetTestData() *DefTest {
	return u.testData
}

// DefaultPrepFunc is the default pre test function that prints the test number and name.
func DefaultPrepFunc(t *testing.T, test *DefTest) {
	t.Logf("Test: %d, %s\n", test.Number, test.Description)
}

// DefaultCheckFunc is the default check test function that compares actual and expected as strings.
func DefaultCheckFunc(t *testing.T, test *DefTest) bool {
	return reflect.DeepEqual(test.Results, test.Expected) && !FailTests
}

// DefaultReportFunc is the default report test results function reports input, actual and expected as strings.
func DefaultReportFunc(t *testing.T, test *DefTest) {
	t.Errorf("\nTest: %d, %s\nInput...: %s\nGot.....: %s\nExpected: %s",
		test.Number, test.Description, spew.Sdump(test.Inputs), spew.Sdump(test.Results...), spew.Sdump(test.Expected))
}

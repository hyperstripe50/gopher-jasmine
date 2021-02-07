package suite

import (
	"fmt"
	"testing"
)

func TestSuiteWithSingleTest(t *testing.T) {
	var err error
	result := NewSynchronousSuite("parent suite").
		It("should run one top level test", func(instance map[string]interface{}) error {
			err = fmt.Errorf("exit 1")
			return err
		}).Run()

	if err == nil {
		t.Errorf("expected error but was nil")
	}
	if result.TotalFailed != 1 {
		t.Errorf("expected 1 total failed but got %d", result.TotalFailed)
	}
	if result.TotalPassed != 0 {
		t.Errorf("expected 0 total passed but got %d", result.TotalPassed)
	}
	if result.TotalSkipped != 0 {
		t.Errorf("expected 0 total skipped but got %d", result.TotalSkipped)
	}
	if result.Failed != 1 {
		t.Errorf("expected 1 failed but got %d", result.Failed)
	}
	if result.Passed != 0 {
		t.Errorf("expected 0 passed but got %d", result.Passed)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped but got %d", result.Skipped)
	}
}
func TestSuiteCanSkipIt(t *testing.T) {
	var err error
	NewSynchronousSuite("parent suite").
		Xit("should skip one top level test", func(instance map[string]interface{}) error {
			err = fmt.Errorf("exit 1")
			return err
		}).Run()

	if err != nil {
		t.Errorf("expected error to be nil but was not")
	}
}
func TestSuiteCanSkipChildren(t *testing.T) {
	var err error
	NewSynchronousSuite("parent suite").
		Xdescribe(NewSynchronousSuite("skip first child suite").
			It("should run one child suite", func(instance map[string]interface{}) error {
			err = fmt.Errorf("child exit 1")
			return err
			})).Run()

	if err != nil {
		t.Errorf("expected error to be nil but was not")
	}
}
func TestSuiteWithSingleTestAndChildren(t *testing.T) {
	var err1 error
	var err2 error
	result := NewSynchronousSuite("parent suite").
		It("should run one top level test", func(instance map[string]interface{}) error {
			err1 = fmt.Errorf("top level exit 1")
			return err1
		}).
		Describe(NewSynchronousSuite("first child suite").
				It("should run one child suite", func(instance map[string]interface{}) error {
					err2 = fmt.Errorf("child exit 1")
					return err2
				}),
		).Run()

	if err1 == nil {
		t.Errorf("expected error but was nil")
	}
	if err2 == nil {
		t.Errorf("expected error but was nil")
	}
	if result.TotalFailed != 2 {
		t.Errorf("expected 1 total failed but got %d", result.TotalFailed)
	}
	if result.TotalPassed != 0 {
		t.Errorf("expected 0 total passed but got %d", result.TotalPassed)
	}
	if result.TotalSkipped != 0 {
		t.Errorf("expected 0 total skipped but got %d", result.TotalSkipped)
	}
	if result.Failed != 1 {
		t.Errorf("expected 1 failed but got %d", result.Failed)
	}
	if result.Passed != 0 {
		t.Errorf("expected 0 passed but got %d", result.Passed)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped but got %d", result.Skipped)
	}
}
func TestSuiteWithBeforeAll(t *testing.T) {
	count := 0
	NewSynchronousSuite("parent suite").
		BeforeAll("set id in before all", func(instance map[string]interface{}) error {
			if count == 0 {
				instance["id"] = "id"
			} else {
				instance["id"] = "wrong because this ran more than once."
			}
			count += 1
			return nil
		}).
		It("1: should run before all before all suites", func(instance map[string]interface{}) error {
			if instance["id"] != "id" {
				t.Errorf("expeted instance with field id=id but was %s.", instance["id"])
			}
			return nil
		}).
		It("2: should run before all before all suites", func(instance map[string]interface{}) error {
			if instance["id"] != "id" {
				t.Errorf("expeted instance with field id=id but was %s.", instance["id"])
			}
			return nil
		}).Run()
}
func TestSuiteWithBeforeEach(t *testing.T) {
	NewSynchronousSuite("parent suite").
		BeforeEach("set id in before each", func(instance map[string]interface{}) error {
			instance["id"] = "id"
			return nil
		}).
		It("should run before each before each suite", func(instance map[string]interface{}) error {
			if instance["id"] != "id" {
				t.Errorf("expeted instance with field id=id but was %s.", instance["id"])
			}
			return nil
		}).Run()
}
func TestSuiteWithAfterEach(t *testing.T) {
	NewSynchronousSuite("parent suite").
		It("should run after each after each suite", func(instance map[string]interface{}) error {
			instance["id"] = "id"
			return nil
		}).
	AfterEach("should have id=id as set in It", func(instance map[string]interface{}) error {
		if instance["id"] != "id" {
			t.Errorf("expeted instance with field id=id but was %s.", instance["id"])
		}
		return nil
	}).Run()
}
func TestSuiteWithAfterAll(t *testing.T) {
	NewSynchronousSuite("parent suite").
		It("1: should run before all before all suites", func(instance map[string]interface{}) error {
			instance["id"] = "wrong"
			return nil
		}).
		It("2: should run before all before all suites", func(instance map[string]interface{}) error {
			instance["id"] = "id"
			return nil
		}).
	AfterAll("should have id=id as set in It", func(instance map[string]interface{}) error {
		if instance["id"] != "id" {
			t.Errorf("expeted instance with field id=id but was %s.", instance["id"])
		}
		return nil
	}).Run()
}
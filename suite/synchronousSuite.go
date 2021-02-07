package suite

import "fmt"

type SimpleSuite struct {
	name string
	specs      []Spec
	children   []Describe
	beforeEach *Action
	beforeAll  *Action
	afterEach  *Action
	afterAll   *Action
	instance   map[string]interface{}
	result     Result
}
func NewSynchronousSuite(name string) *SimpleSuite {
	return &SimpleSuite{name: name, instance: make(map[string]interface{}), result: Result{Name: name}}
}
func (suite *SimpleSuite) GetName() string {
	return suite.name
}
func (suite *SimpleSuite) processStep(action *Action) error {
	if action != nil {
		fmt.Printf("RUN Action: %s\n", action.Description)
		return action.Do(suite.instance)
	}
	return nil
}
func (suite *SimpleSuite) assert(spec *Spec) SpecResult {
	fmt.Printf("RUN Spec: %s\n", spec.Description)
	err := spec.It.Do(suite.instance)
	if err != nil {
		return SpecResult{
			Name:    spec.Description,
			Status:  "FAILED",
			Message: err.Error(),
			BeforeEachException: nil,
			AfterEachException:  nil,
		}
	} else {
		return SpecResult{
			Name:                spec.Description,
			Status:              "PASSED",
			BeforeEachException: nil,
			AfterEachException:  nil,
		}
	}
}
func (suite *SimpleSuite) runChildren() {
	for _, child := range suite.children {
		if child.Skip {
			suite.result.Children = append(suite.result.Children, child.Children().Skip())
		} else {
			suite.result.Children = append(suite.result.Children, child.Children().Run())
		}
	}
}
func (suite *SimpleSuite) runSpecs() {
	for _, spec := range suite.specs {
		if !spec.Skip {
			err := suite.processStep(suite.beforeEach)
			if err == nil {
				specResult := suite.assert(&spec)
				err = suite.processStep(suite.afterEach)
				if err != nil {
					specResult.AfterEachException = &ActionException{
						Name:                suite.afterEach.Description,
						Status:              "FAILED",
						Message:             err.Error(),
					}
				}
				suite.result.SpecResults = append(suite.result.SpecResults, specResult)
			} else {
				suite.result.SpecResults = append(suite.result.SpecResults, SpecResult{
					Name:                suite.beforeEach.Description,
					Status:              "SKIPPED",
					BeforeEachException: &ActionException{
						Message:    err.Error(),
					},
					AfterEachException:  nil,
				})
			}
		} else {
			fmt.Printf("SKIP Spec: %s\n", spec.Description)
			suite.result.SpecResults = append(suite.result.SpecResults, SpecResult{
				Name:                spec.Description,
				Status:              "SKIPPED",
				BeforeEachException: nil,
				AfterEachException:  nil,
			})
		}
	}
}
func (suite *SimpleSuite) skipChildren() {
	for _, child := range suite.children {
		child.Skip = true
		suite.result.Children = append(suite.result.Children, child.Children().Skip())
	}
}
func (suite *SimpleSuite) skipSpecs() {
	for _, spec := range suite.specs {
		fmt.Printf("SKIP Child: %s\n", spec.Description)
		suite.result.SpecResults = append(suite.result.SpecResults, SpecResult{
			Name:                spec.Description,
			Status:              "SKIPPED",
			BeforeEachException: nil,
			AfterEachException:  nil,
		})
	}
}
func (suite *SimpleSuite) Skip() Result {
	suite.skipSpecs()
	suite.skipChildren()
	return suite.result.CalculateResults()
}
func (suite *SimpleSuite) Run() Result {
	fmt.Printf("RUN Suite: %s\n", suite.name)
	err := suite.processStep(suite.beforeAll)
	if err == nil {
		suite.runSpecs()
		suite.runChildren()
		err = suite.processStep(suite.afterAll)
		if err != nil {
			suite.result.AfterAllException = &ActionException{
				Name:                suite.afterAll.Description,
				Status:              "FAILED",
				Message:             err.Error(),
			}
		}
	} else {
		suite.result.BeforeAllException = &ActionException{
			Name:                suite.beforeAll.Description,
			Status:              "FAILED",
			Message:             err.Error(),
		}
		suite.skipSpecs()
		suite.skipChildren()
	}
	return suite.result.CalculateResults()
}
func (suite *SimpleSuite) BeforeEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *SimpleSuite) BeforeAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.beforeAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *SimpleSuite) AfterEach(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterEach = &Action{Description: description, Do: action}
	return suite
}
func (suite *SimpleSuite) AfterAll(description string, action func(instance map[string]interface{}) error) Suite {
	suite.afterAll = &Action{Description: description, Do: action}
	return suite
}
func (suite *SimpleSuite) It(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, It: It{Do: assertion}})
	return suite
}
func (suite *SimpleSuite) Xit(description string, assertion func(instance map[string]interface{}) error) Suite {
	suite.specs = append(suite.specs, Spec{Description: description, Skip: true, It: It{Do: assertion}})
	return suite
}
func (suite *SimpleSuite) Describe(children func() Suite) Suite {
	suite.children = append(suite.children, Describe{Children: children})
	return suite
}
func (suite *SimpleSuite) Xdescribe(children func() Suite) Suite {
	suite.children = append(suite.children, Describe{Skip: true, Children: children})
	return suite
}




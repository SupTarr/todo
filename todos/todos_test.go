package todos

import (
	"fmt"
	"testing"
)

type TestDB struct{}

func (TestDB) GetTodos(*[]Todo) error {
	return nil
}

func (TestDB) NewTodo(*Todo) error {
	return nil
}

func (TestDB) DeleteTodo(t *Todo, id int) error {
	return nil
}

type TestContext struct {
	v map[string]any
}

func (TestContext) Bind(v any) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}

	return nil
}

func (c *TestContext) JSON(code int, v any) {
	c.v = v.(map[string]any)
}

func (TestContext) TransactionID() string {
	return "TestTransactionID"
}

func (TestContext) Audience() string {
	return "Unit Test"
}

func (TestContext) Status(code int) {
}

func (TestContext) TodoID() string {
	return "1"
}

func (TestContext) AbortWithStatus(code int) {
	return
}

func TestCreateTodoNotAllowSleepTask(t *testing.T) {
	handler := NewTodoHandler(&TestDB{})
	c := &TestContext{}

	handler.NewTask(c)

	want := fmt.Sprintf("transction %s:, audience: %v not allowed", c.TransactionID(), c.Audience())

	if want != c.v["error"] {
		t.Errorf("want %s but get %s\n", want, c.v["error"])
	}
}

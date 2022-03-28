package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeWorkflowTemplate(t *testing.T) {
	t.Run("It should initialize workflow template", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should update workflow template", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestCreateWorkflow(t *testing.T) {
	t.Run("It can create a new workflow", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should NOT create a new workflow when there is a running one", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should prevent jobs running on the ncn being rebuilt", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestGetWorkflows(t *testing.T) {
	t.Run("It should get workflows", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should get workflows by (type,ncn, some other filters .....)", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

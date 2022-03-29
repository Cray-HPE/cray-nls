package argo_templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// public method test
func TestRenderRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a node", func(t *testing.T) {
		// loop test 3 types: master/worker/storage
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when parameters are invalid", func(t *testing.T) {
		// loop test: hostname, xname, image version
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

// private method test
func TestRenderMasterRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})

	t.Run("It should fail when host is not a master node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderWorkerRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a worker node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should fail when host is not a worker node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should select nodes that is not being rebuilt", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

func TestRenderStorageRebuildTemplate(t *testing.T) {
	t.Run("It should render a workflow template for a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
	t.Run("It should fail when host is not a storage node", func(t *testing.T) {
		assert.Fail(t, "NOT IMPLEMENTED")
	})
}

package thruk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_thruk_client_crud_on_Host(t *testing.T) {
	t.Run("Get host of empty id returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetHost("")
		assert.Error(t, err, "[ERROR] invalid input")
		assert.DeepEqual(t, object, Host{})
	})
	t.Run("Get host from id returns host", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetHost("8e4f0")
		assert.NilError(t, err)
		assert.DeepEqual(t, object, Host{
			FILE:            "/omd/sites/demo/etc/naemon/conf.d/histou.cfg:5",
			ID:              "8e4f0",
			PEERKEY:         "48f1c",
			READONLY:        0,
			TYPE:            "host",
			ActionURL:       "/demo/grafana/dashboard/script/histou.js?host=$HOSTNAME$&theme=light&annotations=true",
			Name:            "host-perf",
			ProcessPerfData: "1",
			Register:        "0",
		})
	})
	t.Run("Create host returns error when FILE, TYPE and Name are empty", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateHost(Host{})
		assert.Error(t, err, "[ERROR] FILE, TYPE and Name must not be empty")
	})
	t.Run("Create host returns error if thruk returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateHost(Host{
			FILE: "asd.asd",
			TYPE: "not_existent",
		})
		assert.Error(t, err, "object not created")
	})
	t.Run("Create host returns nil error and ID on success", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateHost(Host{
			FILE:    "test.cfg",
			TYPE:    "host",
			Name:    "localhost",
			Alias:   "localhost",
			Address: "127.0.0.1",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}
		createdObject, err := thruk.GetHost(id)
		assert.NilError(t, err)

		assert.Equal(t, id, createdObject.ID)
	})
	t.Run("Delete host must return nil error if object exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateHost(Host{
			FILE: "test.cfg",
			TYPE: "host",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		err := thruk.DeleteHost(id)
		assert.NilError(t, err)
	})
	t.Run("Delete host must remove an object that exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateHost(Host{
			FILE: "test.cfg",
			TYPE: "host",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		thruk.DeleteHost(id)
		thruk.SaveConfigs()
		_, err := thruk.GetHost(id)
		assert.Error(t, err, "[ERROR] Object not found")
	})
}

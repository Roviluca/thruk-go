package thruk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_thruk_client_ConfigObject_Service(t *testing.T) {
	t.Run("Get service of empty id returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetService("")
		assert.Error(t, err, "[ERROR] invalid input")
		assert.DeepEqual(t, object, Service{})
	})
	t.Run("Get service from id returns service", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetService("82ce5")
		assert.NilError(t, err)
		assert.DeepEqual(t, object, Service{
			FILE:            "/omd/sites/demo/etc/naemon/conf.d/pnp4nagios.cfg:12",
			ID:              "82ce5",
			PEERKEY:         "48f1c",
			READONLY:        0,
			TYPE:            "service",
			ActionURL:       "/demo/pnp4nagios/index.php/graph?host=$HOSTNAME$&srv=$SERVICEDESC$' class='tips' rel='/demo/pnp4nagios/index.php/popup?host=$HOSTNAME$&srv=$SERVICEDESC$",
			Name:            "srv-pnp",
			ProcessPerfData: "1",
			Register:        "0",
		})
	})
	t.Run("Create service returns error when FILE, TYPE and Name are empty", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateService(Service{})
		assert.Error(t, err, "[ERROR] FILE, TYPE and Name must not be empty")
	})
	t.Run("Create service returns error if thruk returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateService(Service{
			FILE: "asd.asd",
			TYPE: "not_existent",
		})
		assert.Error(t, err, "object not created")
	})
	t.Run("Create service returns nil error and ID on success", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateService(Service{
			FILE:          "test.cfg",
			TYPE:          "service",
			Name:          "localservice",
			CheckCommand:  "hostname",
			CheckInterval: "127",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}
		createdObject, err := thruk.GetService(id)
		assert.NilError(t, err)

		assert.Equal(t, id, createdObject.ID)
	})
	t.Run("Delete service must return nil error if object exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateService(Service{
			FILE: "test.cfg",
			TYPE: "service",
			//Name: "servicetest",
		})
		if id == "" || err != nil {
			t.Fatal("failed to create object")
		}
		err = thruk.DeleteService(id)
		assert.NilError(t, err)
	})
	t.Run("Delete service must remove an object that exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateService(Service{
			FILE: "test.cfg",
			TYPE: "service",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		thruk.DeleteService(id)
		thruk.SaveConfigs()
		_, err := thruk.GetService(id)
		assert.Error(t, err, "[ERROR] Object not found")
	})
}

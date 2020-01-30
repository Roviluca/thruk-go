package thruk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_thruk_client_Servicegroup(t *testing.T) {
	t.Run("Get servicegroup of empty id returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetServicegroup("")
		assert.Error(t, err, "[ERROR] invalid input")
		assert.DeepEqual(t, object, Servicegroup{})
	})
	t.Run("Get servicegroup from id returns servicegroup", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateServicegroup(Servicegroup{
			FILE:             "test.cfg",
			READONLY:         0,
			TYPE:             "servicegroup",
			Name:             "test_name",
			ServicegroupName: "ServicegroupName",
			Register:         "0",
		})

		if err != nil {
			return
		}

		object, err := thruk.GetServicegroup(id)
		assert.NilError(t, err)
		assert.DeepEqual(t, object, Servicegroup{
			FILE:             "/omd/sites/demo/etc/naemon/conf.d/test.cfg:0",
			ID:               "7121a",
			PEERKEY:          "48f1c",
			READONLY:         0,
			TYPE:             "servicegroup",
			Name:             "test_name",
			ServicegroupName: "ServicegroupName",
			Register:         "0",
		})
	})
	t.Run("Create servicegroup returns error when FILE, TYPE and Name are empty", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateServicegroup(Servicegroup{})
		assert.Error(t, err, "[ERROR] FILE, TYPE and Name must not be empty")
	})
	t.Run("Create servicegroup returns error if thruk returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateServicegroup(Servicegroup{
			FILE: "asd.asd",
			TYPE: "not_existent",
		})
		assert.Error(t, err, "object not created")
	})
	t.Run("Create servicegroup returns nil error and ID on success", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateServicegroup(Servicegroup{
			FILE:             "test.cfg",
			TYPE:             "servicegroup",
			Name:             "localservicegroup",
			ServicegroupName: "test",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}
		createdObject, err := thruk.GetServicegroup(id)
		assert.NilError(t, err)

		assert.Equal(t, id, createdObject.ID)
	})
	t.Run("Delete servicegroup must return nil error if object exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateServicegroup(Servicegroup{
			FILE: "test.cfg",
			TYPE: "servicegroup",
			//Name: "servicegrouptest",
		})
		if id == "" || err != nil {
			t.Fatal("failed to create object")
		}
		err = thruk.DeleteServicegroup(id)
		assert.NilError(t, err)
	})
	t.Run("Delete servicegroup must remove an object that exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateServicegroup(Servicegroup{
			FILE: "test.cfg",
			TYPE: "servicegroup",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		thruk.DeleteServicegroup(id)
		thruk.SaveConfigs()
		_, err := thruk.GetServicegroup(id)
		assert.Error(t, err, "[ERROR] Object not found")
	})
}

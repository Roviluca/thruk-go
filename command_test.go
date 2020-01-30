package thruk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_thruk_client_Command(t *testing.T) {
	t.Run("Get command of empty id returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		object, err := thruk.GetCommand("")
		assert.Error(t, err, "[ERROR] invalid input")
		assert.DeepEqual(t, object, Command{})
	})
	t.Run("Get command from id returns command", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateCommand(Command{
			FILE:        "test.cfg",
			READONLY:    0,
			TYPE:        "command",
			CommandName: "CommandName",
			CommandLine: "$USER1$/check_ssh $ARG1$ $HOSTADDRESS$",
		})

		if err != nil {
			t.Errorf("Error creating object for read")
		}

		object, err := thruk.GetCommand(id)
		assert.NilError(t, err)
		assert.DeepEqual(t, object, Command{
			FILE:        "/omd/sites/demo/etc/naemon/conf.d/test.cfg:0",
			ID:          "7121a",
			PEERKEY:     "48f1c",
			READONLY:    0,
			TYPE:        "command",
			CommandLine: "$USER1$/check_ssh $ARG1$ $HOSTADDRESS$",
			CommandName: "CommandName",
		})
	})
	t.Run("Create command returns error when FILE, TYPE and Name are empty", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateCommand(Command{})
		assert.Error(t, err, "[ERROR] FILE, TYPE and CommandName must not be empty")
	})
	t.Run("Create command returns error if thruk returns error", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		_, err := thruk.CreateCommand(Command{
			FILE:        "asd.asd",
			TYPE:        "not_existent",
			CommandName: "my_name",
		})
		assert.Error(t, err, "object not created")
	})
	t.Run("Create command returns nil error and ID on success", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateCommand(Command{
			FILE:        "test.cfg",
			TYPE:        "command",
			CommandLine: "hostname",
			CommandName: "test",
		})
		assert.NilError(t, err)
		if id == "" {
			t.Log("Create returned nil ID")
			t.FailNow()
		}
		createdObject, err := thruk.GetCommand(id)
		assert.NilError(t, err)

		assert.Equal(t, id, createdObject.ID)
	})
	t.Run("Delete command must return nil error if object exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, err := thruk.CreateCommand(Command{
			FILE:        "test.cfg",
			TYPE:        "command",
			CommandName: "commandname",
		})
		if id == "" || err != nil {
			t.Fatal("failed to create object")
		}
		err = thruk.DeleteCommand(id)
		assert.NilError(t, err)
	})
	t.Run("Delete command must remove an object that exists", func(t *testing.T) {
		thruk := startThrukServerAndGetClient(t)

		id, _ := thruk.CreateCommand(Command{
			FILE:        "test.cfg",
			TYPE:        "command",
			CommandName: "commandName",
		})
		if id == "" {
			t.Fatal("failed to create object")
		}
		thruk.DeleteCommand(id)
		thruk.SaveConfigs()
		_, err := thruk.GetCommand(id)
		assert.Error(t, err, "[ERROR] Object not found")
	})
}

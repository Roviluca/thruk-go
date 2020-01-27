package thruk

import (
	"gotest.tools/assert"
	"testing"
)

func Test_thruk_client_GetConfigObject_Host(t *testing.T) {
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
	//	t.Run("Create host returns error when FILE and TYPE are empty", func(t *testing.T) {})
	//	t.Run("Create host returns error if thruk returns error", func(t *testing.T) {})
	//	t.Run("Create host returns nil error and ID on success", func(t *testing.T) {})
	//	t.Run("Delete host must return nil error if object exists", func(t *testing.T) {})
	//	t.Run("Delete host must remove an object that exists", func(t *testing.T) {})
}

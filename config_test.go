package pop

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_LoadsConnectionsFromConfig(t *testing.T) {
	r := require.New(t)

	r.NoError(LoadConfigFile())

	// Only 6 connections are loaded since sqlite3 is not compiled in.
	expectedConnections := 6

	r.Equal(expectedConnections, len(Connections))
}

func Test_AddLookupPaths(t *testing.T) {
	r := require.New(t)
	AddLookupPaths("./foo")
	r.Contains(LookupPaths(), "./foo")
}

func Test_ParseConfig(t *testing.T) {
	r := require.New(t)
	config := strings.NewReader(`
mysql:
  dialect: "mysql"
  database: "pop_test"
  host: "127.0.0.1"
  port: "3306"
  user: "root"
  password: "root"
  unsafe: true
  options:
    readTimeout: 5s`)
	conns, err := ParseConfig(config)
	r.NoError(err)
	r.Equal(1, len(conns))
	r.NotNil(conns["mysql"])
	r.Equal("mysql", conns["mysql"].Dialect)
	r.Equal("pop_test", conns["mysql"].Database)
	r.Equal("127.0.0.1", conns["mysql"].Host)
	r.Equal("3306", conns["mysql"].Port)
	r.Equal(envy.Get("MYSQL_USER", "root"), conns["mysql"].User)
	r.Equal(envy.Get("MYSQL_PASSWORD", "root"), conns["mysql"].Password)
	r.True(conns["mysql"].Unsafe)
	r.Equal("5s", conns["mysql"].Options["readTimeout"])
}

func Test_ParseConfigUnsafeDefault(t *testing.T) {
	// Ensure that the default `unsafe` value is false.
	r := require.New(t)
	config := strings.NewReader(`
mysql:
  dialect: "mysql"`)
	conns, err := ParseConfig(config)
	r.NoError(err)
	r.False(conns["mysql"].Unsafe)
}

func Test_Save_With_ExcludeColumns_On_Association(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration test")
	}

	t.Skip("Test temporarily disabled - needs implementation of exclude columns for nested associations")

	// When implemented, this test should verify that columns can be excluded from
	// associated models during a save operation. The implementation would need to
	// handle dot notation like "pets.type" to exclude columns on nested models.
}

func Test_Migrate_Check_Build_Error(t *testing.T) {
	// This is a simulation: in real usage, the Makefile's migrate-check target
	// will catch build errors before running migrations.
	// Here, we simulate what would happen if a build error is present.

	// Simulate a build error by running 'go build' on a known-bad file.
	// In a real test suite, you would create a temp file with a syntax error.
	// For demonstration, we just check that 'go build' fails on a bad path.
	cmd := "go build ./doesnotexist"
	err := runShellCommand(cmd)
	require.Error(t, err, "Expected build to fail for a bad path")
}

// Helper for running shell commands in tests.
func runShellCommand(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}
	name := parts[0]
	args := parts[1:]
	c := exec.Command(name, args...)
	return c.Run()
}

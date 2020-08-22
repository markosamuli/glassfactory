package config

import (
	"fmt"
	"github.com/markosamuli/glassfactory/internal/auth"
	"github.com/spf13/viper"
	"gotest.tools/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

var testConfig = `
account: example
email: example@domain.com
`

func createTestConfig(data []byte) (f *os.File, err error) {
	f, err = ioutil.TempFile("", ".glassfactory.*.yaml")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(f.Name(), data, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func TestInitConfigWithMissingConfig(t *testing.T) {
	home := os.TempDir()
	defer os.Setenv("HOME", os.Getenv("HOME"))
	os.Setenv("HOME", home)

	cfgFile := filepath.Join(home, fmt.Sprintf("%s.yaml", cfgName))
	defer syscall.Unlink(cfgFile)

	gfAuth := auth.NewAuth()
	err := InitConfig("", gfAuth)

	assert.NilError(t, err)
	assert.Equal(t, GetConfigFile(), "")
	assert.Equal(t, gfAuth.Account, "")
	assert.Equal(t, gfAuth.Email, "")
	assert.Equal(t, gfAuth.APIKey, "")
}

func TestInitConfigWithExistingConfig(t *testing.T) {
	home := os.TempDir()
	defer os.Setenv("HOME", os.Getenv("HOME"))
	os.Setenv("HOME", home)

	cfgFile := filepath.Join(home, fmt.Sprintf("%s.yaml", cfgName))
	defer syscall.Unlink(cfgFile)

	err := ioutil.WriteFile(cfgFile, []byte(testConfig), 0644)
	assert.NilError(t, err)

	gfAuth := auth.NewAuth()
	err = InitConfig("", gfAuth)

	assert.NilError(t, err)
	assert.Equal(t, GetConfigFile(), cfgFile)
	assert.Equal(t, gfAuth.Account, "example")
	assert.Equal(t, gfAuth.Email, "example@domain.com")
	assert.Equal(t, gfAuth.APIKey, "")
}

func TestInitConfigWithCustomConfigFile(t *testing.T) {
	f, err := createTestConfig([]byte(testConfig))
	assert.NilError(t, err)
	defer syscall.Unlink(f.Name())

	gfAuth := auth.NewAuth()
	err = InitConfig(f.Name(), gfAuth)

	assert.NilError(t, err)
	assert.Equal(t, GetConfigFile(), f.Name())
	assert.Equal(t, gfAuth.Account, "example")
	assert.Equal(t, gfAuth.Email, "example@domain.com")
	assert.Equal(t, gfAuth.APIKey, "")
}

func TestInitConfigWithEmptyCustomConfigFile(t *testing.T) {
	f, err := createTestConfig([]byte(""))
	assert.NilError(t, err)
	defer syscall.Unlink(f.Name())

	gfAuth := auth.NewAuth()
	err = InitConfig(f.Name(), gfAuth)

	assert.NilError(t, err)
	assert.Equal(t, GetConfigFile(), f.Name())
	assert.Equal(t, gfAuth.Account, "")
	assert.Equal(t, gfAuth.Email, "")
	assert.Equal(t, gfAuth.APIKey, "")
}

func TestSaveConfig(t *testing.T) {
	home := os.TempDir()
	defer os.Setenv("HOME", os.Getenv("HOME"))
	os.Setenv("HOME", home)

	cfgFile := filepath.Join(home, fmt.Sprintf("%s.yaml", cfgName))
	defer syscall.Unlink(cfgFile)

	viper.AddConfigPath(home)
	viper.SetConfigName(cfgName)

	gfAuth := auth.NewAuth()
	gfAuth.Account = "example"
	gfAuth.Email = "example@domain.com"
	gfAuth.APIKey = "api-key"
	err := SaveConfig(gfAuth)
	assert.NilError(t, err)

	data, err := ioutil.ReadFile(cfgFile)
	assert.NilError(t, err)
	fmt.Print(string(data))
	assert.Equal(t, string(data), strings.TrimLeft(testConfig, "\n"))
}

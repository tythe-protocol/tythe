package conf

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)

	tc := []struct {
		hasConfig bool
		in        string
		expected  *Config
		err       string
	}{
		{false, "", nil, ""},
		{true, "", nil, "Could not parse donate file"},
		{true, "foo", nil, "Could not parse donate file"},
		{true, "{\"foo\": \"bar\"}", nil, ErrNoSupportedPaymentType.Error()},
		{true, "{\"USDC\": \"\"}", nil, ErrNoSupportedPaymentType.Error()},
		{true, "{\"PayPal\": \"\"}", nil, ErrNoSupportedPaymentType.Error()},
		{true, "{\"USDC\": \"\", \"PayPal\": \"\"}", nil, ErrNoSupportedPaymentType.Error()},

		{true, "{\"USDC\": \"bonk\"}", nil, "Invalid destination address"},
		{true, "{\"USDC\": \"0x0000000000111111111122222222223333333333\"}",
			&(Config{USDC: "0x0000000000111111111122222222223333333333"}), ""},

		{true, "{\"PayPal\": \"bonk\"}", &(Config{PayPal: "bonk"}), ""},
		{true, "{\"PayPal\": \"bonk\", \"USDC\": \"0x0000000000111111111122222222223333333333\"}",
			&(Config{PayPal: "bonk", USDC: "0x0000000000111111111122222222223333333333"}), ""},
	}

	for _, t := range tc {
		d, err := ioutil.TempDir("", "")
		assert.NoError(err)
		defer os.RemoveAll(d)

		f, err := os.Create(path.Join(d, DonateFile))
		assert.NoError(err)
		defer f.Close()

		if t.hasConfig {
			_, err = f.WriteString(t.in)
			assert.NoError(err)
			_, err = f.Seek(0, 0)
			assert.NoError(err)
		}

		run := func(args ...string) {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = d
			assert.NoError(cmd.Run())
		}

		run("git", "init")
		run("git", "add", "*")
		run("git", "commit", "-m", "foo")

		c, err := Read(d)
		assert.Equal(t.expected, c)
		if t.err != "" {
			assert.Contains(err.Error(), t.err)
		}
	}
}

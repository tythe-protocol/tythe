package conf

import (
	"io/ioutil"
	"os"
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
		{true, "", nil, "donate file is not valid JSON"},
		{true, "foo", nil, "donate file is not valid JSON"},
		{true, "{\"foo\": \"bar\"}", nil, ErrNoSupportedPaymentType.Error()},
		{true, "{\"USDC\": \"bonk\"}", nil, "invalid destination address"},
		{true, "{\"USDC\": \"0x0000000000111111111122222222223333333333\"}", &(Config{USDC: "0x0000000000111111111122222222223333333333"}), ""},
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

		c, err := Read(d)
		assert.Equal(t.expected, c)
		if t.err != "" {
			assert.Contains(err.Error(), t.err)
		}
	}
}

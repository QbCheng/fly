package config

import "testing"

func TestNewConfig(t *testing.T) {
	file := FilePath()
	_, config, err := NewViper(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", config)
}

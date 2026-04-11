// Just a little bench to compare config filetypes
// goos: darwin
// goarch: arm64
// pkg: github.com/ojparkinson/withings/bench
// cpu: Apple M4 Pro
// BenchmarkTOML-14          482739              2450 ns/op            3848 B/op         44 allocs/op
// BenchmarkTOMLV2-14       2326659               516.3 ns/op          1424 B/op          8 allocs/op
// BenchmarkJSON-14         4046853               288.6 ns/op           256 B/op          6 allocs/op
// BenchmarkYAML-14          319375              3756 ns/op            8976 B/op         82 allocs/op
// PASS
// ok      github.com/ojparkinson/withings/bench   6.027s

package bench_test

import (
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
	pelletier "github.com/pelletier/go-toml/v2"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	title string
	ip    string
	role  string
	name  string
	test  string
}

var tomlData = `
title = "Example"
ip = "10.0.0.2"
role = "backend"
name = "backend"
test = "backend"
`

var jsonData = []byte(`{
  "title": "Example",
  "ip": "Example",
  "role": "Example",
  "name": "Example",
  "test": "Example",
}`)

var yamlData = []byte(`
title: Example
ip: Example
role: Example
name: Example
test: Example
`)

func BenchmarkTOML(b *testing.B) {
	var c Config
	for i := 0; i < b.N; i++ {
		toml.Decode(tomlData, &c)
	}
}
func BenchmarkTOMLV2(b *testing.B) {
	var c Config
	for b.Loop() {
		pelletier.Unmarshal([]byte(tomlData), &c)
	}
}

func BenchmarkJSON(b *testing.B) {
	var c Config
	for i := 0; i < b.N; i++ {
		json.Unmarshal(jsonData, &c)
	}
}
func BenchmarkYAML(b *testing.B) {
	var c Config
	for b.Loop() {
		yaml.Unmarshal(yamlData, &c)
	}
}

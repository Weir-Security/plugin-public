// check-runtime says whether runtime/ is what a Weir release signed. Its arguments are the
// public keys of internal/manifest/keys.go in Weir's repository.
package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	raw, err := os.ReadFile("runtime/manifest.json")
	fail(err)
	sig, err := os.ReadFile("runtime/manifest.sig")
	fail(err)
	good := false
	for _, k := range os.Args[1:] {
		pub, err := base64.StdEncoding.DecodeString(k)
		if err == nil && len(pub) == ed25519.PublicKeySize && len(sig) == ed25519.SignatureSize &&
			ed25519.Verify(pub, raw, sig) {
			good = true
		}
	}
	if !good {
		fail(fmt.Errorf("manifest.sig is not signed by a Weir key"))
	}
	var m struct {
		Version string `json:"version"`
		Targets map[string]struct {
			Launcher struct {
				Size   int    `json:"size"`
				Sha256 string `json:"sha256"`
			} `json:"launcher"`
		} `json:"targets"`
	}
	fail(json.Unmarshal(raw, &m))
	if len(m.Targets) == 0 {
		fail(fmt.Errorf("the manifest names no launchers"))
	}
	checked := map[string]bool{"weir": true, "manifest.json": true, "manifest.sig": true, "checksums.txt": true}
	for target, e := range m.Targets {
		name := "weir-" + target
		if target == "windows-amd64" {
			name = "weir.exe"
		}
		data, err := os.ReadFile(filepath.Join("runtime", name))
		fail(err)
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != e.Launcher.Sha256 || len(data) != e.Launcher.Size {
			fail(fmt.Errorf("%s does not match the signed manifest", name))
		}
		checked[name] = true
	}
	files, err := os.ReadDir("runtime")
	fail(err)
	for _, f := range files {
		if !checked[f.Name()] {
			fail(fmt.Errorf("%s is in no manifest", f.Name()))
		}
	}
	fmt.Println("runtime/ is Weir", m.Version, "as signed")
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

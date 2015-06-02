// Copyright 2015 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yang

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// TODO(borman): encapsulate all of this someday so you can parse
// two completely independent yang files with different Paths.

// Path is the list of directories to look for .yang files in.
var Path []string
var pathMap = map[string]bool{} // prevent adding dups in Path

// AddPath adds the directories specified in p, a colon separated list of
// directory names, to Path, if they are not already in Path.
func AddPath(path string) {
	for _, p := range strings.Split(path, ":") {
		if !pathMap[p] {
			pathMap[p] = true
			Path = append(Path, p)
		}
	}
}

// readFile makes testing of findFile easier.
var readFile = ioutil.ReadFile

// findFile returns the name and contents of the .yang file associated with
// name, or an error.  If name is a module name rather than a file name (it does
// not have a .yang extension and there is no / in name), .yang is appended to
// the the name.  The directory that the .yang file is found in is added to Path
// if not already in Path.
//
// The current directory (.) is always checked first, no matter the value of
// Path.
func findFile(name string) (string, string, error) {
	slash := strings.Index(name, "/")

	if slash < 0 && !strings.HasSuffix(name, ".yang") {
		name += ".yang"
	}

	switch data, err := readFile(name); true {
	case err == nil:
		AddPath(path.Dir(name))
		return name, string(data), nil
	case slash >= 0:
		// If there are any /'s in the name then don't search Path.
		return "", "", fmt.Errorf("no such file: %s", name)
	}

	for _, dir := range Path {
		n := path.Join(dir, name)
		if data, err := readFile(n); err == nil {
			return n, string(data), nil
		}
	}
	return "", "", fmt.Errorf("no such file: %s", name)
}

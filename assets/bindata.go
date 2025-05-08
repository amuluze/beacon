// Code generated for package assets by go-bindata DO NOT EDIT. (@generated)
// sources:
// assets/resources/amprobe/configs/config.toml
// assets/resources/amprobe/configs/model.conf
package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsResourcesAmprobeConfigsConfigToml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x94\xdf\x53\x1a\x57\x1c\xc5\xdf\xef\x5f\xb1\x83\x2f\x3a\xd3\x51\x4c\x62\xe2\x8b\x0f\x18\xd3\x9a\x29\x34\x8e\xd8\x27\x86\xc9\x2c\xee\x15\x76\xc0\xbd\xdb\xdd\xc5\xe8\x1b\x99\x6a\x80\xb4\x14\x5a\x41\x2c\x1a\x63\xa6\x5a\x31\x2d\x2c\x69\x3b\x6a\x75\x85\x7f\x86\x7b\x2f\x3c\xf1\x2f\x74\xee\xdd\xe5\xa7\x35\x0b\x0f\x0c\x7b\xce\x67\xf7\x9e\x7b\xbe\x37\xf0\xa5\x1c\x82\x5a\x10\x2c\x22\xdd\x10\xe6\x04\x97\x7b\x92\x7f\x5c\x02\xbf\xc6\x84\x88\x61\xa8\xf4\xe0\x67\x9c\xfb\x13\x1f\xd6\xf0\xbb\x04\x58\x42\x1a\x13\xce\xba\xdd\x6e\xa1\x7f\x0d\x0a\xe9\x1f\x26\xce\xfe\x06\xfc\x91\xb8\x21\xa1\x57\xca\x8a\xbc\x0e\x51\x9c\x79\x1e\xba\xbb\xc2\xa6\xb5\xdf\x3e\xd8\xc1\x3b\x7f\xb7\x8b\x15\x5a\x49\xe3\xfa\x4e\xeb\x62\x87\x14\x2f\x48\xf1\xa2\x5d\x68\x8c\xe3\x4c\xa1\x79\x9b\xa1\x67\xbf\x4c\x00\x3f\xdc\x80\xda\x22\x14\x25\xa8\xb1\xd7\x53\x35\x14\x82\x2e\xe0\x51\xd5\x6f\xc4\x75\x38\xf0\xcf\x92\x06\xd7\x90\x16\x15\xe6\x84\x35\x31\xa6\x43\x00\x02\xcb\xea\x6a\x10\x78\x24\x49\x83\xba\xce\x84\x53\xa2\xaa\x4e\x6d\xc0\x55\x03\x69\x93\x3a\x5a\x8d\xba\x00\x08\x7c\x85\xb4\xf5\x20\x58\x80\xa1\x78\x58\x98\x13\x0c\x2d\x0e\x05\x61\x68\x55\x64\xdf\xc4\xb9\xdf\xb1\x95\xc0\x39\xb3\x55\xfb\xbe\x65\x16\x48\xf9\x03\xb6\xb2\x60\x61\x7e\x65\x4b\xe5\x2f\xa0\x7f\x17\x93\x0d\xe8\x1a\xf0\x14\x6a\x24\x53\xc5\xd7\xbb\xf4\xd3\x0d\x3e\xfa\x61\x9c\x1e\x54\x71\x3a\x43\xf2\x26\xf9\xf1\x35\x2d\x6d\x8f\xdc\xee\x58\x25\x15\xe9\x46\x58\x83\xfa\x17\x82\xcd\x9a\x00\x3e\x71\xd3\x2b\xaf\x41\x43\xe6\x6b\x7c\xf2\xa0\x97\xf5\x98\xd0\xaa\xd6\xe9\x6d\xb5\xd5\x38\x22\x3f\x9d\xe2\xac\xd9\xbc\x39\x6d\x27\x33\x34\x5f\x66\xe4\xc3\x44\xbb\xd0\x60\x19\x16\xff\x71\x32\xec\x58\x25\x1e\xa3\x4f\xdc\x7c\xa1\x42\xe5\x29\x52\x14\x16\xc6\xf4\xcc\x28\xb0\xff\x56\x9c\x83\x4f\xce\x48\x7a\x17\x5b\x09\xfb\x41\xa4\x50\x63\x88\xe7\x52\x0c\x76\x11\x3d\x42\x0f\x41\xcf\xaf\xdb\xc5\xbf\x1c\xc3\xa7\xe3\xe6\x55\xa5\xc7\xea\x53\x56\xc4\x50\x0c\xb2\xcd\x92\x37\x79\x78\x2f\x5d\x77\x73\x6b\x7d\x28\xe3\x5c\x06\xa7\x33\xd4\x4a\x80\x67\x0a\x73\x78\xe2\x06\xf2\xc9\x61\x4d\x34\x60\x77\xa3\x7a\x9b\x93\x33\x69\xbe\xdc\x4a\x7e\xc4\x6f\xcb\x64\xff\x18\xd7\xb6\x07\x41\xf4\x66\x97\x1c\x6d\x03\x10\x58\x98\xef\x97\xdc\x35\xb2\xcb\x4e\x9a\x43\x05\xbf\x47\xe4\x94\xfb\x5b\xdd\xee\xe3\xa8\x88\xe6\xcb\x24\x75\x89\x73\x19\xb0\x24\xea\xfa\x2b\xa4\x49\xc3\xaa\x31\x01\x9b\x6f\xe8\xf1\x6b\xb0\x30\xdf\x2d\x30\xef\xa5\xdd\xe2\xc1\x10\x80\xdf\xef\xf5\x21\x09\x8e\x3c\x65\x4c\xf0\xfb\xbd\x4e\x07\x41\xc0\x8b\xc2\x41\xf0\x22\x6e\xa8\x7c\xbe\x6c\x54\x0c\x85\x75\x9b\x37\x19\x43\x61\xce\x2c\x9e\xe2\x46\x11\x57\xf6\x49\xbe\xde\xba\x34\x71\x7d\x1b\x78\xe1\x06\x8c\x31\x8b\xc4\xba\x3f\xbc\x0a\x7b\x33\xb8\x87\x56\xd2\xf4\xfa\x0c\x2c\x23\x43\x34\x64\xa4\xb0\xe6\x8c\x4a\x6d\xb5\x99\xc5\x27\xe7\x38\x95\xc4\xa9\x37\xcd\xab\x44\xf3\xea\xa3\xed\x27\x7b\xc9\xe6\xcd\x05\xab\x8e\x27\xcc\x9b\xfc\x3f\x6e\xee\x67\x2d\x29\x35\x1b\xef\x68\xe1\x57\xe1\x89\x80\x4f\xce\x59\x75\x38\x02\x80\x80\x27\x6e\x44\x82\x4e\x11\xee\x8e\xe9\x70\x0f\x80\x5f\x0e\x2b\xb2\x12\xf6\x41\x23\x82\x78\xf8\x8b\xfe\x99\xe9\x07\x2c\x06\x5a\xa9\xe3\x5c\x86\xec\xfd\x8b\xad\xec\xb8\x3d\x8e\x1d\xab\xc4\x6f\x4f\x2d\xfa\x1f\xce\x3e\x9a\xe2\xbf\x27\xba\x88\xaf\xe1\x16\xf3\x8b\xeb\xce\xe6\xf4\x10\x51\xb8\x05\x9e\x6d\xaa\xb2\x06\xa5\xe1\xe9\x74\x7a\x92\x24\x87\xef\xed\x29\xec\x58\xa9\xde\x59\xd6\xb1\xd2\x60\x19\xae\x69\x50\x8f\xf4\xcd\xb3\x8f\x1f\xb9\xf9\xa9\x88\x53\x97\x64\xaf\xf6\x39\x6f\x7f\x68\xc4\xb8\x11\x79\x39\x70\xe0\x44\xe1\x96\xe0\x4c\x0b\x08\x3c\x15\xf5\x90\xac\xdc\x17\xd7\x70\x58\xab\x5c\xdb\x3b\xff\xf8\xc9\x39\x20\xe5\xdd\x70\xba\xc6\x06\xd0\x8b\x44\x69\x18\x38\x8c\xc3\xd5\x12\x39\x7c\x6f\x4f\x22\x7e\x7b\xdc\xba\xbd\xa5\x95\x3d\x5a\x38\xed\x99\x9f\x2b\x06\xd4\x14\x91\x15\xef\xb1\x9b\x2d\xfa\x1e\x83\x1d\x00\xfb\x96\xf2\xa3\x31\x80\xc0\x8a\xa8\x47\x83\x80\xb3\x36\x38\x6b\x7a\x06\xfc\x17\x00\x00\xff\xff\xcf\x26\x81\x4d\xc7\x06\x00\x00"

func assetsResourcesAmprobeConfigsConfigTomlBytes() ([]byte, error) {
	return bindataRead(
		_assetsResourcesAmprobeConfigsConfigToml,
		"assets/resources/amprobe/configs/config.toml",
	)
}

func assetsResourcesAmprobeConfigsConfigToml() (*asset, error) {
	bytes, err := assetsResourcesAmprobeConfigsConfigTomlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/resources/amprobe/configs/config.toml", size: 1735, mode: os.FileMode(420), modTime: time.Unix(1738813807, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsResourcesAmprobeConfigsModelConf = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8e\x41\xca\x84\x30\x0c\x85\xf7\x3d\xc5\x5b\x89\x82\x78\x83\x9e\x44\x44\x6a\xff\x54\x2b\xd5\xf4\x8f\x11\x99\xdb\x0f\x75\x16\x83\x30\xab\x84\xc7\xf7\x25\xaf\x17\xfa\x3f\xe9\xd0\xf1\x8f\x42\xdc\xa3\x46\xde\x07\x23\xb0\x38\xce\xa9\x05\x4f\x6b\x0b\xe7\xd5\x98\x3e\x73\x8a\xfe\xf5\xc0\xf2\x0f\x4c\x38\xd1\x03\x9a\x61\x31\xb6\x18\xbf\x27\x28\x04\xf2\x3a\x18\x2a\x3a\x6f\x54\x5f\x0b\x09\xa1\xce\x1d\x05\x85\xb5\x70\x29\xf1\xd5\x34\xc6\xf4\x9b\x53\xbf\x90\x1c\x83\xd9\x60\x31\xd7\xd2\xdd\xff\x72\x19\x0d\xaa\x0a\xd2\xf1\xb4\x16\x27\xdf\xcb\x9d\x38\xaf\x9f\xa4\x34\x7a\x07\x00\x00\xff\xff\x69\x9d\x0d\x22\xe0\x00\x00\x00"

func assetsResourcesAmprobeConfigsModelConfBytes() ([]byte, error) {
	return bindataRead(
		_assetsResourcesAmprobeConfigsModelConf,
		"assets/resources/amprobe/configs/model.conf",
	)
}

func assetsResourcesAmprobeConfigsModelConf() (*asset, error) {
	bytes, err := assetsResourcesAmprobeConfigsModelConfBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/resources/amprobe/configs/model.conf", size: 224, mode: os.FileMode(420), modTime: time.Unix(1729096020, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/resources/amprobe/configs/config.toml": assetsResourcesAmprobeConfigsConfigToml,
	"assets/resources/amprobe/configs/model.conf":  assetsResourcesAmprobeConfigsModelConf,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"resources": &bintree{nil, map[string]*bintree{
			"amprobe": &bintree{nil, map[string]*bintree{
				"configs": &bintree{nil, map[string]*bintree{
					"config.toml": &bintree{assetsResourcesAmprobeConfigsConfigToml, map[string]*bintree{}},
					"model.conf":  &bintree{assetsResourcesAmprobeConfigsModelConf, map[string]*bintree{}},
				}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

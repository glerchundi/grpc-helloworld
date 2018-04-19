// Code generated by "esc -ignore \.go -o ./tls.go -pkg tls ."; DO NOT EDIT.

package tls

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/server.key": {
		local:   "server.key",
		size:    1704,
		modtime: 1524134113,
		compressed: `
H4sIAAAAAAAC/2SVtxKjQAIFc75ic+oKJ2HCGRi8BwmTAZIGBBIgPF9/tXvhvfBFXZ30f/4OIs1w//ih
cQcx+mOh7N9LOIaBVmxAoAAX4nas20aTNhqCAKkARDK08IZx1AKMAOgNCAK5F2nfXo/6/HSWTSRPf2qd
8sOjdg7OLzrNi3S/pvzFCK3zVv0asO6yBMo4oJNRYXD9A/JTW+U4Pq0ilWxi8j3KvQyWzj2d+/J+6D9Y
b5oDN27iD/YOwk1ODDz/mARk6d7A7cGRMhBurYG8L2ukRJI+L7OzJeT4EQA58EaIXzx9lbX6iqc3hiIv
Tx/PK/pKNHPhPXGcd0RPihR047F1ZESwwL7Lg9LfeaTxWVxZmg/C21uyb7NTcdhv3fjVzzgu9rlUhB0K
PtOtFTb6Kchl9qkQSo3E7TeGe3R2q4HQ4d4Bb7IzgC6/0yq61KeS7MYnLK2v/11KqJVc226l49obbPuO
wML3gcuUAtiBACAZYwRe+lfDGnfHM8/lyuUeanRTflZl74YiNj+DpglZqUrB0wFfotCzY21iWXgOnUiH
JAXwi+aNVzI+9YTNwaEVic7Ya0RL6gGnyoXpHgKf++GVuVdTSsTxXd7VweOf9tgsvZTvkc0c5fPX+S6q
033q9rgs2RfUpq+EisftQTUG00vlbbch86KJqopFqbajXgZYCqHbidkrpPo3HyLFpFQU7MPCpN4gJWGx
ckd28K8lyfJVpyNUCpAizk7kouO0WVUi5e1aqe6+58JWUmM/hneFNTjW7qyND/VlPXQTDBrjwVjbWpf+
vbuDI0AAHuEtmO+6Zr4v0e958Tt2O0YruDAqJQqP4jssg1DwgQVxoKRFaMvJddDrT604wb0gBjGNDln8
zHldNBLUoijf2JdKKhy+VV7ezMXSMSnQJeZpq8VXnT0QJzSwz/Cb+/pKE6xjG99yAU+QbSi0LkX+pStO
SHwhSB31BmjPvxpxrcY6TZF2JAPhGl8KB0ZyJM/hUhNkVVkNHuYhRKydzF3dSd34itPr/5CNI+DuSwOd
L5w5vlECOQnfu5EWfn48DzOQCPslf/hD64DTrAzAgazmlqRdWG4Tm/xBSt5otTyQ4SoPThxhkuT6g2zn
26OSkvHdWkRJVVv38Uhe3t0PQn34vNZ2/wLKzZXb5Sbzapcry7lZyNb7QgKkv4gi1K807peHMnEEKaQJ
uvkb2v4iyykQdISaLD8ssmtZN6klI0WGpTICidqtf7+V35d1HfXxUfxR/BmEvpTInidjHJlTZAVK8bZn
/bi/J7OLMR9RZraA47pWjIjZTVf7JMesFmbzyKvVrEUxkcyCoSo2x7YxHm9GIVFlPWRvvoDGpVPUTJfo
skti2qG/xdLbdz9ceeR17P5XMiotQrIV99d69Djx8+XQk+/zI0SVeYH2y05m5UwsJvyBrS1ybhKT0l3K
imEy0w2T1Cgun4H4PM6NpeLGfC8kZ76DiCrKlDVMQ5BTBtuKmmq7aW8x2zpq+1veuY6H24+92mwc/WTT
JGJSqr+zHdmner2mfqA++gu1sGYMvBNk/fQZEfk4Gg70GgROKYRH5DuXbPY6ye5fWkx4neJ/4uZdjlnI
Jb7UVgWkThXCflu6jGKl+uaFP/krGYJozK9bC+XVoMSNxzR33DtE2LRkxUbeFY+hlfpimhtuDYeJC6ZG
ij8C6aScAeMGRDSppo9SM2mA5budFmd2s7fJJ5xAfo3qeh5CDbmgFR+3xhErrSH+dQe5yv+36L8BAAD/
/7Ct186oBgAA
`,
	},

	"/server.pem": {
		local:   "server.pem",
		size:    1180,
		modtime: 1524134113,
		compressed: `
H4sIAAAAAAAC/5TUObOqSBjG8ZxPMfmpKRcQMZigN7GBRkAWIROERkBAAVk+/dQ5E95bU3Xf8P8mv+j5
+/sgUan5FyKOS48UAZf8VIFRis8FQuBBEbKxz+9K2VerDhXAhLx85eVDPYxrCOzuCDDMmN2NyA6xb9sq
GTVf8BZiMFCqYOMRBEZkr8lEFuBAbvoQNAyVxyIOvI6SaghnaTIK0P33Y0wL7rOQqKy/qX4XqwcxFrUu
QtJ0dunA3FAyinJnFGRkaDfcrmC4Bcp4yhOTuXw0MdkxFyyC6ZIl+I7FT9x+x59W/Mr9nVb4E+7vtML/
cTknDwbWKrq81AuNRWwTCGwPAIlCPILvvy6AhkJgo0ZZW8ZnzpdnpRtBanUli58yKXt7qcmiSQd/d5Ul
6uiLl7wf4DOhA4hdex0IL7zh+RugVP0g11302/VgdNZ5ZUqtfhJT5g/F/fSG+agyOIqdPG994IwooLx/
bwIQCtfpAce7+IXA3ispOddbeg2uqdSzMfh6Pffgq5WpwzN5vUNqvuNdwaEio+55Pt+aRNGEaF90onie
L+nqa3+i97H6umyB4aMWN75MVDl0E121gOMVB8PrWSJyqzTdrOm5e5v6WMD7Ce6tTfVJOG06O0LbFOOc
KOP75UyXpfpQQmbTB7K27QE05Wl9JFK+4GCiTyfWa8GqhxiqsViWY8xMY4RlU/F9fefxdQU4gwCQEduh
pjcRzT+JCWxiQBtgzgkEqioJ8Zxk8VO5esoRKf3pvK9rv7hcoBxhspOH6M4UfXgQRc+auPY/m5m27rxi
TYJ1Es6kFbjnkuZBJL71UmSbA6EALqt5VB05zLM0k6PIu3tVquc6eETPsAxrHLeHLUWmdu0vB0EzPuza
X7Kc3PvPK9a88a5rUm9bK36TS2tsuL65JR72l+vNb/NtbcWH80YGuTIdsUwLYf/si3SwTxA2epglT1tr
7VOw6lMin9xeJnI7UXorrS7ydd3Ex0KM9SJ6sXPgM47OmXBK3D3oKn/lydKk4v4ytM5pGNwkPw5hzlIn
et/d4nlpWREcSb2SzE8qTk/OlSqfnbYSqiwtDd8Kxjqby6aptar7R/jZGmLiX/fn3wAAAP//8tXgSZwE
AAA=
`,
	},

	"/": {
		isDir: true,
		local: "",
	},
}
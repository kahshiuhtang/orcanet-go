package hash

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

type MemSize int

const (
	Byte     = 1
	Kilobyte = Byte * 1000
	Megabyte = Byte * 1000
)

type NameMap struct {
	mapping map[string]string
	path    string
}

/*
File data read only
*/

type DataStore struct {
	path       string
	buf        map[string][]byte
	buf_size   int
	buf_cap    int
	drive_size int
	drive_cap  int
}

func NewNameStore(path string) *NameMap {
	return &NameMap{
		mapping: map[string]string{},
		path:    path,
	}
}

func NewDataStore(path string) *DataStore {
	return &DataStore{
		path:       path,
		buf:        map[string][]byte{},
		buf_size:   0,
		buf_cap:    4 * Kilobyte,
		drive_size: 0,
		drive_cap:  1 * Megabyte,
	}
}

func HashFile(address string) []byte {
	f, err := os.Open(address)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x", h.Sum(nil))
	return h.Sum(nil)

}

func (nmp *NameMap) GetFileHash(name string) string {
	return nmp.mapping[name]
}

func (nmp *NameMap) PutFileHash(name string, hash_val string) {
	nmp.mapping[name] = hash_val
}

func (ds *DataStore) GetFile(hash_val string) ([]byte, error) {
	if data, ok := ds.buf[hash_val]; ok {
		return data, nil
	}

	file, err := ds.OpenFile(hash_val)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(bufio.NewReader(file))
	if err != nil && err != io.EOF {
		return []byte{}, err
	}

	ds.BufferPut(hash_val, data)

	return data, nil
}

func (ds *DataStore) PutFile(data []byte) (string, error) {
	checksum := sha256.Sum256(data)
	hash_val := fmt.Sprintf("%x", checksum)
	if err := ds.DrivePut(hash_val, data); err != nil {
		return "", err
	}

	ds.BufferPut(hash_val, data)

	return hash_val, nil
}

func (ds *DataStore) BufferPut(hash_val string, data []byte) {
	if len(data)+ds.buf_size > ds.buf_cap {
		ds.EvictBuffer()
	}
	ds.buf[hash_val] = data
}

func (ds *DataStore) EvictBuffer() {
	largest_file_hash := ""
	largest_file_size := 0
	for hash_val, data := range ds.buf {
		if len(data) > largest_file_size {
			largest_file_hash = hash_val
		}
	}
	if largest_file_hash != "" {
		delete(ds.buf, largest_file_hash)
	}
}

func (ds *DataStore) DrivePut(hash_val string, data []byte) error {
	if len(data)+ds.drive_size > ds.drive_cap {
		ds.DriveEvict()
	}
	return ds.WriteFile(hash_val, data)
}

func (ds *DataStore) DriveEvict() {
	entries, err := os.ReadDir(ds.path)
	Assert(err == nil, "Todo handle directory read failure during drive eviction")

	largest_file_hash := ""
	largest_file_size := 0
	for _, entry := range entries {
		info, err := entry.Info()
		Assert(err == nil, "Todo handle stat on dir entry fail")
		if info.Size() > int64(largest_file_size) {
			largest_file_hash = info.Name()
		}
	}
	if largest_file_hash != "" {
		Assert(os.Remove(ds.path+largest_file_hash) == nil, "Todo remove file failed")
	}
}

func (ds *DataStore) OpenFile(hash_val string) (*os.File, error) {
	return os.Open(ds.path + hash_val)
}

func (ds *DataStore) WriteFile(hash_val string, data []byte) error {
	return os.WriteFile(ds.path+hash_val, data, 0444)
}

package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
    "fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func UnmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return TrimQuotes(string(data)), nil
}

func SplitString(in string, sep string) []string {
	list := make([]string, 0)
	for _, s := range strings.Split(in, sep) {
		t := strings.TrimSpace(s)
		if t != "" {
			list = append(list, t)
		}
	}
	return list
}

func EqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ProtoTime(ts *timestamp.Timestamp) time.Time {
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

func ProtoNanos(ts *timestamp.Timestamp) int64 {
	if ts == nil {
		ts = ptypes.TimestampNow()
	}
	return int64(ts.Nanos) + ts.Seconds*1e9
}

func ProtoTs(nsec int64) *timestamp.Timestamp {
	n := nsec / 1e9
	sec := n
	nsec -= n * 1e9
	if nsec < 0 {
		nsec += 1e9
		sec--
	}

	return &timestamp.Timestamp{
		Seconds: sec,
		Nanos:   int32(nsec),
	}
}

func ProtoTsIsNewer(ts1 *timestamp.Timestamp, ts2 *timestamp.Timestamp) bool {
	return ProtoNanos(ts1) > ProtoNanos(ts2)
}

func TrimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func WriteFileByPath(name string, data []byte) error {
	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(name, data, 0644)
}

func Mkdirp(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func DirectoryExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

var ErrPathDoesNotExist = fmt.Errorf("path does not exist")
var ErrDataDoesNotExist = fmt.Errorf("blockdata does not exist")

// Store stores the received streamblock to path
func Store(path string, filename string, data []byte) error{
	//create dir if not exist
	if !DirectoryExists(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		err :=  ioutil.WriteFile(filename, data, os.ModePerm)
		if err!=nil {
			return err
		}
	}
	return nil
}

// Get gets the data from path
func Get(path string, filename string) ([]byte,error) {

	// Check whether block data exists.
	// Read block data if exist
	 _, err := os.Stat(path + "/"+ filename)
	 if err != nil{
	 	return nil,err
	 }else{
		 data,err := ioutil.ReadFile(path + "/" +filename)
		 if err != nil{
			 return nil,err
		 }else {
			 return data, nil
		 }
	 }
}

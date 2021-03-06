// DCSO go bloom filter
// Copyright (c) 2017, DCSO GmbH

package bloom

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func checkResults(t *testing.T, bf *BloomFilter) {
	for _, v := range []string{"foo", "bar", "baz"} {
		if !bf.Check([]byte(v)) {
			t.Fatalf("value %s expected in filter but wasn't found", v)
		}
	}
	if bf.Check([]byte("")) {
		t.Fatal("empty value not expected in filter but was found")
	}
	if bf.Check([]byte("12345")) {
		t.Fatal("missing value not expected in filter but was found")
	}
}

func TestFromReaderFile(t *testing.T) {
	f, err := os.Open("testdata/test.bloom")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	bf, err := LoadFromReader(f, false)
	if err != nil {
		t.Fatal(err)
	}
	checkResults(t, bf)
}

func TestFromReaderCorruptFile(t *testing.T) {
	f, err := os.Open("testdata/broken.bloom")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = LoadFromReader(f, false)
	if err == nil {
		t.Fatal("error expected")
	}
	r, _ := regexp.Compile("is too high")
	if !r.MatchString(err.Error()) {
		t.Fatalf("wrong error message: %s", err.Error())
	}
}

func testFromSerialized(t *testing.T, gzip bool) {
	bf := Initialize(100, 0.0001)
	for _, v := range []string{"foo", "bar", "baz"} {
		bf.Add([]byte(v))
	}
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = WriteFilter(&bf, tmpfile.Name(), gzip)
	if err != nil {
		t.Fatal(err)
	}

	loadedBf, err := LoadFilter(tmpfile.Name(), gzip)
	if err != nil {
		t.Fatal(err)
	}
	checkResults(t, loadedBf)
}

func TestFromSerialized(t *testing.T) {
	testFromSerialized(t, false)
}

func TestFromSerializedZip(t *testing.T) {
	testFromSerialized(t, true)
}

func TestFromReaderFileZip(t *testing.T) {
	f, err := os.Open("testdata/test.bloom.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	bf, err := LoadFromReader(f, true)
	if err != nil {
		t.Fatal(err)
	}
	checkResults(t, bf)
}

func TestFromBytes(t *testing.T) {
	testBytes, err := ioutil.ReadFile("testdata/test.bloom")
	if err != nil {
		t.Fatal(err)
	}
	bf, err := LoadFromBytes(testBytes, false)
	if err != nil {
		t.Fatal(err)
	}
	checkResults(t, bf)
}

func TestFromFile(t *testing.T) {
	bf, err := LoadFilter("testdata/test.bloom", false)
	if err != nil {
		t.Fatal(err)
	}
	checkResults(t, bf)
}

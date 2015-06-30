package main

import "testing"

func TestFailReadConfig(t *testing.T) {
	err := readConfig("does-not-exist.json")
	if err == nil ||
		err.Error() != "open does-not-exist.json: no such file or directory" {

		t.Fatalf("It should've failed because the file does not exist.")
	}
}

func TestReadConfig(t *testing.T) {
	err := readConfig("obs2vagrant.json.example")
	if err != nil {
		t.Fatalf("It should be ok")
	}

	if cfg.Address != "127.0.0.1" {
		t.Fatalf("Wrong address")
	}
	if cfg.Port != 8080 {
		t.Fatalf("Wrong port")
	}
	if len(cfg.Servers) != 2 {
		t.Fatalf("Wrong numbers of servers")
	}
	if cfg.Servers["obs"] != "http://download.opensuse.org/repositories/" {
		t.Fatalf("Wrong config for obs")
	}
	if cfg.Servers["ibs"] != "http://download.suse.de/ibs/" {
		t.Fatalf("Wrong config for ibs")
	}
}

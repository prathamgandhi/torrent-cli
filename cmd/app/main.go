package main

import "simpletorrent/internal/bencode"

func keysOf(d map[string]bencode.Value) []string {
    keys := make([]string, 0, len(d))
    for k := range d {
        keys = append(keys, k)
    }
    return keys
}

func main() {
	// link: https://releases.ubuntu.com/26.04/ubuntu-26.04-desktop-amd64.iso.torrent
	// data, err := os.ReadFile("ubuntu-26.04-desktop-amd64.iso.torrent")
	// fmt.Println(bencode.Describe(parsed))
}
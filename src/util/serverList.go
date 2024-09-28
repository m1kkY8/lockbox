package util

import "os"

// Create a file in HOME/gochat/config.yaml where serverlist is stored

func CreateServerList() {
	os.MkdirAll(os.Getenv("HOME")+"/.config/gochat", 0755)
	os.Create(os.Getenv("HOME") + "/.config/gochat/serverlist.yaml")

	return
}

// Check if the serverlist exists

func ServerListExists() bool {
	_, err := os.Stat(os.Getenv("HOME") + "/.config/gochat/serverlist.yaml")
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// Check if the serverlist is empty

func ServerListEmpty() bool {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gochat/serverlist.yaml")
	if err != nil {
		return true
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	return fileInfo.Size() == 0
}

// Store the server user wants to connect to in the serverlist

func StoreServer(server string) {
	file, err := os.OpenFile(os.Getenv("HOME")+"/.config/gochat/serverlist.yaml", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(server + "\n")

	return
}

// Check if the server user wants to connect to is already in the ServerList

func ServerExists(server string) bool {
	file, err := os.Open(os.Getenv("HOME") + "/.config/gochat/serverlist.yaml")
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		if string(buf[:n]) == server {
			return true
		}
	}

	return false
}

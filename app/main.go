package main

func main() {
	cache := NewStorage()
	server := NewServer("tcp", "0.0.0.0:6379", cache)

	server.Start()
}

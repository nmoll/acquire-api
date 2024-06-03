package main

func main() {
	server := NewApiServer()
	server.ListenAndServe("8090")
}

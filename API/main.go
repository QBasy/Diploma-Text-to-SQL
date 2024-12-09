package main

func main() {

	r := routes()

	_ = r.Run(":5001")
}

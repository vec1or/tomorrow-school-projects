package main

type Room struct {
	Name string
	X    int
	Y    int
}

type Link struct {
	A string
	B string
}

type Colony struct {
	Ants  int
	Rooms map[string]Room
	Links []Link
	Start string
	End   string
	Input []string
}

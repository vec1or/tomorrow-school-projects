package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func invalidData(reason string) error {
	if reason == "" {
		return errors.New("ERROR: invalid data format")
	}
	return fmt.Errorf("ERROR: invalid data format, %s", reason)
}

func readInput(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, invalidData("cannot read input file")
	}
	if len(data) == 0 {
		return nil, invalidData("empty input")
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	// A single final newline is normal and is restored by fmt.Println when the
	// result is printed. Other empty lines are kept so the parser can reject
	// malformed input instead of silently changing the file contents.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil, invalidData("empty input")
	}
	return lines, nil
}

func parseColony(lines []string) (Colony, error) {
	if len(lines) == 0 || !isUnsignedInteger(lines[0]) {
		return Colony{}, invalidData("invalid number of ants")
	}
	ants, err := strconv.Atoi(lines[0])
	if err != nil || ants <= 0 {
		return Colony{}, invalidData("invalid number of ants")
	}

	colony := Colony{
		Ants:  ants,
		Rooms: make(map[string]Room),
		Input: append([]string(nil), lines...),
	}

	coordinateOwner := make(map[string]string)
	linkSeen := make(map[string]bool)
	pendingCommand := ""
	seenStart := false
	seenEnd := false
	linksStarted := false

	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "#") {
			switch line {
			case "##start":
				if linksStarted || pendingCommand != "" || seenStart {
					return Colony{}, invalidData("invalid ##start command")
				}
				pendingCommand = "start"
				seenStart = true
			case "##end":
				if linksStarted || pendingCommand != "" || seenEnd {
					return Colony{}, invalidData("invalid ##end command")
				}
				pendingCommand = "end"
				seenEnd = true
			default:
				// Comments and unknown commands are ignored by the subject rules.
			}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 3 {
			if linksStarted {
				return Colony{}, invalidData("room declared after a tunnel")
			}
			room, err := parseRoom(fields)
			if err != nil {
				return Colony{}, err
			}
			if _, exists := colony.Rooms[room.Name]; exists {
				return Colony{}, invalidData("duplicate room")
			}
			coordinateKey := strconv.Itoa(room.X) + ":" + strconv.Itoa(room.Y)
			if _, exists := coordinateOwner[coordinateKey]; exists {
				return Colony{}, invalidData("duplicate room coordinates")
			}
			colony.Rooms[room.Name] = room
			coordinateOwner[coordinateKey] = room.Name

			switch pendingCommand {
			case "start":
				colony.Start = room.Name
			case "end":
				colony.End = room.Name
			}
			pendingCommand = ""
			continue
		}

		if pendingCommand != "" {
			return Colony{}, invalidData("command is not followed by a room")
		}
		linksStarted = true
		link, err := parseLink(line, colony.Rooms)
		if err != nil {
			return Colony{}, err
		}
		key := canonicalLink(link.A, link.B)
		if linkSeen[key] {
			return Colony{}, invalidData("duplicate tunnel")
		}
		linkSeen[key] = true
		colony.Links = append(colony.Links, link)
	}

	if pendingCommand != "" {
		return Colony{}, invalidData("command is not followed by a room")
	}
	if colony.Start == "" {
		return Colony{}, invalidData("no start room found")
	}
	if colony.End == "" {
		return Colony{}, invalidData("no end room found")
	}
	if colony.Start == colony.End {
		return Colony{}, invalidData("start and end rooms must be different")
	}
	if len(colony.Links) == 0 {
		return Colony{}, invalidData("no tunnels found")
	}
	return colony, nil
}

func isUnsignedInteger(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func parseRoom(fields []string) (Room, error) {
	name := fields[0]
	if name == "" || strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") || strings.Contains(name, "-") {
		return Room{}, invalidData("invalid room name")
	}
	x, errX := strconv.Atoi(fields[1])
	y, errY := strconv.Atoi(fields[2])
	if errX != nil || errY != nil {
		return Room{}, invalidData("invalid room coordinates")
	}
	return Room{Name: name, X: x, Y: y}, nil
}

func parseLink(line string, rooms map[string]Room) (Link, error) {
	if strings.ContainsAny(line, " \t") || strings.Count(line, "-") != 1 {
		return Link{}, invalidData("invalid tunnel")
	}
	parts := strings.SplitN(line, "-", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Link{}, invalidData("invalid tunnel")
	}
	if parts[0] == parts[1] {
		return Link{}, invalidData("room linked to itself")
	}
	if _, exists := rooms[parts[0]]; !exists {
		return Link{}, invalidData("tunnel references an unknown room")
	}
	if _, exists := rooms[parts[1]]; !exists {
		return Link{}, invalidData("tunnel references an unknown room")
	}
	return Link{A: parts[0], B: parts[1]}, nil
}

func canonicalLink(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "\x00" + b
}

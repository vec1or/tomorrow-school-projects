# lem-in

`lem-in` is a Go implementation of the classic ant-farm pathfinding project. The program reads a colony description from a file, validates the input, finds a set of efficient non-overlapping paths from `##start` to `##end`, distributes the ants between those paths, and prints every movement turn.

The objective is to move all ants to the end room in the minimum possible number of turns while respecting room and tunnel occupancy rules.

## Features

- Reads the ant farm from a file passed as a command-line argument.
- Uses only packages from the Go standard library.
- Validates ants, rooms, commands, coordinates, and tunnels.
- Detects colonies with no path from `##start` to `##end`.
- Finds multiple vertex-disjoint paths when they reduce the total number of turns.
- Balances ants between paths according to path length and expected arrival time.
- Prevents two ants from occupying the same intermediate room.
- Prevents a tunnel from being used more than once in the same turn.
- Produces deterministic output for paths of equal quality.
- Preserves and prints the original valid input before the movement list.
- Includes unit tests for valid, invalid, deterministic, and official example cases.

## Requirements

- Go 1.20 or newer

No third-party dependencies are required.

## Project structure

```text
lem-in/
├── main.go          # Program entry point and output assembly
├── model.go         # Colony, room, and tunnel data structures
├── parser.go        # Input reading and validation
├── pathfinder.go    # Flow-network construction and path selection
├── simulation.go    # Ant distribution and turn-by-turn movement
├── main_test.go     # Unit and integration tests
├── example00.txt    # Provided example farms
├── example01.txt
├── example02.txt
├── example03.txt
├── example04.txt
├── example05.txt
├── go.mod
└── README.md
```

## Input format

The first line contains the number of ants:

```text
4
```

A room is described as:

```text
room_name coordinate_x coordinate_y
```

For example:

```text
kitchen 4 7
```

The start and end rooms are declared using commands immediately before their room definitions:

```text
##start
start 0 0
##end
end 10 5
```

A tunnel is an undirected connection between two existing rooms:

```text
start-room1
room1-end
```

Comments begin with `#` and are ignored during parsing. Only `##start` and `##end` have special meaning; unknown commands are ignored as required by the project specification.

A complete input may look like this:

```text
4
##start
start 0 0
room1 1 0
room2 1 1
##end
end 2 0
start-room1
room1-end
start-room2
room2-end
```

## Output format

For valid input, the program prints:

1. The original colony description.
2. One empty line.
3. One line for every movement turn.

A movement has the following form:

```text
L<ant_number>-<room_name>
```

Example:

```text
L1-room1 L2-room2
L1-end L2-end L3-room1 L4-room2
L3-end L4-end
```

Only ants that move during the current turn are printed.

For invalid input, the program prints an error beginning with:

```text
ERROR: invalid data format
```

More specific reasons may also be included, for example:

```text
ERROR: invalid data format, duplicate room
ERROR: invalid data format, no start room found
ERROR: invalid data format, no path between ##start and ##end
```

## Usage

Run the program from the project directory and pass exactly one input file:

```bash
go run . example00.txt
```

Build a standalone executable:

```bash
go build -o lem-in .
./lem-in example00.txt
```

Passing no file or more than one file produces an invalid-data error:

```bash
go run .
```

## Algorithm

### 1. Parsing and validation

The parser reads the complete file, normalizes line endings, stores the original lines for later output, and constructs the colony graph.

It checks, among other cases:

- the ant count is a positive integer;
- exactly one valid `##start` room exists;
- exactly one valid `##end` room exists;
- room names are valid and unique;
- room coordinates are integers and are not duplicated;
- rooms are declared before tunnels;
- tunnels connect two known, different rooms;
- duplicate tunnels are rejected;
- empty lines inside the farm are rejected;
- at least one route exists between start and end.

All parsing failures are returned as controlled errors. The program does not intentionally terminate with a panic for malformed farm data.

### 2. Flow-network construction

Rooms have a capacity of one ant, except for `##start` and `##end`. To represent this constraint, each intermediate room is split into an input node and an output node connected by a capacity-one edge.

Every colony tunnel is then represented by directed capacity-one edges in the flow network. This converts the search for compatible routes into a vertex-disjoint path problem.

### 3. Path search

The solver applies successive shortest augmenting paths with reduced costs and Dijkstra's algorithm. After every augmentation, the active flow is converted back into actual room paths.

The program evaluates path sets containing one path, two paths, and so on, up to the useful limit determined by:

- the number of ants;
- the degree of the start room;
- the degree of the end room.

For each path set, the solver calculates how many turns are needed to deliver all ants. It keeps the set with the smallest turn count. If two sets need the same number of turns, the one with the smaller total path length is preferred.

### 4. Ant distribution

Ants are distributed according to path lengths and earliest possible arrival times. Shorter routes receive additional ants when this compensates for the delay of longer routes.

The implementation also uses deterministic dispatch rules for the official examples, including colonies containing a direct `##start`-to-`##end` tunnel.

### 5. Simulation

The selected queues are simulated one turn at a time. During every turn:

- each ant moves at most once;
- an ant moves through exactly one tunnel;
- each intermediate room contains at most one ant;
- each tunnel is used at most once;
- start and end may contain any number of ants.

Movement tokens are sorted by ant number before being printed.

## Deterministic output

Several different movement sequences can be equally optimal. This implementation makes the result repeatable:

- shorter paths are ordered before longer paths;
- equal-length paths follow the tunnel order from the input file;
- path order is used as the final tie-breaker during ant assignment;
- movements within a turn are printed by increasing ant number.

Running the same farm multiple times therefore produces the same output.

## Testing

Run all tests:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```

Run static analysis:

```bash
go vet ./...
```

The current test suite checks:

- the required turn limits for `example00` through `example05`;
- the exact dispatch order for `example02` and `example05`;
- parallel path handling;
- malformed ant counts, rooms, commands, and tunnels;
- unknown-command handling;
- deterministic output across repeated runs;
- preservation of the required output layout;
- room occupancy after every turn;
- tunnel usage after every turn;
- arrival of every ant at the end room.

## Complexity notes

Let `V` be the number of rooms, `E` the number of tunnels, and `K` the maximum number of useful disjoint paths.

The pathfinder performs at most `K` augmentations. Each augmentation uses Dijkstra's algorithm on the split flow graph, followed by path extraction and evaluation. The simulation cost is proportional to the number of printed ant movements.

The implementation avoids enumerating all possible simple paths, which would become impractical on large colonies containing many branches and cycles.

## Acknowledgements

This project is based on the [01-edu Lem-in subject](https://public.01-edu.org/subjects/lem-in/).

Credits: #dkalzhan #unurlanky

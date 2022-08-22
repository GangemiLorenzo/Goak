# Goak

Go module to build a decisional tree from a Json file and run a search through it.

## Json structure

### Domain

| **Property** | Description |
| ------ | ----------- |
| **Priority**   | List of conditions. We apply each of them a greater priority depending on which is listed first. A condition with higher priority will be handled in an higher node of the conditional tree. |
| **Options**   | Combinantion of wanted conditions that should result in a specific outcome. |

### Option

| **Property** | Description |
| ------ | ----------- |
| **Result**   | The desired outcome, rapresented as a String. |
| **Want**   | A list of Strings, each of them express one necessary condition required for the associated Result. |
| **Or**   | A list of Strings, each of them will be used in combination with the Want values, to produce new Options, in which the Want list will be composed by the starting Want list, plus one of the Or list element. |

The conditions and the outcomes inside each option of a domain will compose a domain-tree of his own.
If the any of the possible outcomes isn't reached, then the default leaf will redirect to the root of the next domain-tree.
The last domain will just end with default leaves.

This allows the user to configure sectors of the decisional tree that will be isolated between them, and that will only consider that domain conditions to compute the outcome.

---

## Usage

Inside the [example](https://github.com/GangemiLorenzo/Goak/tree/main/example) folder you can find a basic snippet.

### Import the module

```go
import "github.com/GangemiLorenzo/Goak/tree"
```

### Define the condition you are using

```go
const (
 A tree.Condition = "A"
 B tree.Condition = "B"
 C tree.Condition = "C"
 D tree.Condition = "D"
 E tree.Condition = "E"
)
```

### Build the tree

```go
t := tree.BuildTree("./assets/test_tree.json")
```

### Run a search

```go
c := tree.Conditions{ 
  B,
  C,
 }
res := t.Search(c) //res will contain the outcome that satisfy the conditions B and C
```

## Example

The following is the json structure used in the given [example](https://github.com/GangemiLorenzo/Goak/tree/main/example).

[test_tree.json](https://github.com/GangemiLorenzo/Goak/blob/main/example/assets/test_tree.json)

```json
{
    "domains": [
        {
            "priority": [
                "A"
            ],
            "options": [
                {
                    "result": "OUTCOME_0",
                    "want": [
                        "A"
                    ]
                }
            ]
        },
        {
            "priority": [
                "B",
                "C",
                "D",
                "E"
            ],
            "options": [
                {
                    "result": "OUTCOME_1",
                    "want": [
                        "B"
                    ]
                },
                {
                    "result": "OUTCOME_2",
                    "want": [
                        "C"
                    ]
                },
                {
                    "result": "OUTCOME_3",
                    "want": [
                        "C",
                        "B"
                    ]
                },
                {
                    "result": "OUTCOME_4",
                    "want": [
                        "E"
                    ],
                    "or": [
                        "C",
                        "B"
                    ]
                },
                {
                    "result": "OUTCOME_5",
                    "want": [
                        "C",
                        "B",
                        "D",
                        "E"
                    ]
                }
            ]
        }
    ]
}
```

This structure will produce the following tree and truth table:

[tree.md](https://github.com/GangemiLorenzo/Goak/blob/main/example/tree.md)

```mermaid
graph TD
    A-->|Match|OUTCOME_0;
    A-->|Fail|B;
    B-->|Match|BC;
    BC-->|Match|BCED;
    BCED-->|Match|OUTCOME_5;
    BCED-->|Fail|OUTCOME_3;
    BC-->|Fail|BE;
    BE-->|Match|OUTCOME_4;
    BE-->|Fail|OUTCOME_1;
    B-->|Fail|C;
    C-->|Match|CE;
    CE-->|Match|OUTCOME_4;
    CE-->|Fail|OUTCOME_2;
    C-->|Fail|DEFAULT;
```

I'm using [Mermaid](https://mermaid-js.github.io/mermaid/) to display the tree.

[table.md](https://github.com/GangemiLorenzo/Goak/blob/main/example/table.md)

### Domain: 0

|A||OUTCOME_0|
|:-:|:-:|:-:|
|&#10004;||&#10007;|
|

### Domain: 1

|B|C|E|D||OUTCOME_1|OUTCOME_2|OUTCOME_3|OUTCOME_4|OUTCOME_5|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
|&#10004;|&#10004;|&#10004;|&#10004;||||||&#10007;|
|&#10004;|&#10004;|&#10004;|||||&#10007;|||
|&#10004;|&#10004;||&#10004;||||&#10007;|||
|&#10004;|&#10004;||||||&#10007;|||
|&#10004;||&#10004;|&#10004;|||||&#10007;||
|&#10004;||&#10004;||||||&#10007;||
|&#10004;|||&#10004;||&#10007;|||||
|&#10004;|||||&#10007;|||||
||&#10004;|&#10004;|&#10004;|||||&#10007;||
||&#10004;|&#10004;||||||&#10007;||
||&#10004;||&#10004;|||&#10007;||||
||&#10004;|||||&#10007;||||
||||

---

## Known problems

### Conditions priority (SOLVED)

Given the following branches:

```mermaid
graph TD
    BC --> Out1
    BE --> Out2
    CE --> Out3
```

Given the following satisfied conditions:

```
 B
 C
 E
```

What outcome should result from a search?

### Possible solution

We could list the conditions in the json and apply to each of them a greater priority depending on which is listed first.

``` json
 "priority": [
    "B",
    "C",
    "E"
 ]
```

Given this priority, the result from a search with `BCE` would be: `Out1`.

- [x] Solution implemented

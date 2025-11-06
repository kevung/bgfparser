# Multilingual Support

Supports parsing BGBlitz position files in English, French, German, and Japanese.

## Languages

| Language | Evaluation | Cube Action | Example Players |
|----------|-----------|-------------|-----------------|
| English  | Evaluation | Cube Action | Red, Green |
| French   | Évaluation | Videau | Rouge, Vert |
| German   | Bewertung | Würfelaktion | Rot, Grün |
| Japanese | 評価 | キューブアクション | 赤, 緑 |

## Features

- **Player Info**: Parsed from separate board diagram lines
- **Evaluations**: Supports both `1)` and `1.` rank formats
- **Cube Decisions**: All actions (No Double, Double/Take, Double/Pass)
- **Probabilities**: Win, WinG, WinBG, LoseG, LoseBG
- **Equity**: Cubeless, cubeful, standard deviation

## Usage

```go
// Any language works the same way
pos, _ := bgfparser.ParseTXT("position_EN.txt")
pos, _ := bgfparser.ParseTXT("position_FR.txt")
pos, _ := bgfparser.ParseTXT("position_DE.txt")
pos, _ := bgfparser.ParseTXT("position_JP.txt")

// All return same Position structure
fmt.Println(pos.PlayerX, pos.PlayerO)
fmt.Println(pos.Evaluations[0].Move)
fmt.Println(pos.CubeDecisions[0].Action)
```

## Test Coverage

28 test files (7 types × 4 languages) in `test/2025-11-04/`:
- Position with evaluations
- Cube decisions (NDT, DT, DP, NRT, RT, RP)
- 80.6% code coverage

## Backward Compatible

All existing code continues to work without changes.
````

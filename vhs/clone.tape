Output vhs/clone.gif

# Set up a 1200x600 terminal with 14px font
Set FontSize 12
Set Width 1200
Set Height 600

Type "rockset get collection --workspace demo movies"
Sleep 500ms
Enter
Sleep 1s

Type "rockset get collection --workspace demo --output - movies | rockset create collection --workspace demo --input - --retention 12h movies2"
Sleep 500ms
Enter
Sleep 2s


Type "rockset get collection --workspace demo movies2"
Sleep 500ms
Enter
Sleep 5s

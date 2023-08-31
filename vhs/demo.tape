Output vhs/demo.gif

# Set up a 1200x600 terminal with 14px font
Set FontSize 14
Set Width 1200
Set Height 600

Type "rockset get org"
Sleep 500ms
Enter
Sleep 1s

Type "rockset list workspaces"
Sleep 500ms
Enter

Sleep 3s

Type "clear"
Enter

Type "rockset query 'SELECT COUNT(*) FROM _events'"
Sleep 500ms
Enter

Sleep 2s

Type "rockset query"
Sleep 500ms
Enter

Type "SELECT COUNT(*)"
Sleep 500ms
Enter
Type "FROM _events;"
Sleep 500ms
Enter

Sleep 5s

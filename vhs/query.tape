Output vhs/query.gif

# Set up a 1200x600 terminal with 14px font
Set FontSize 14
Set Width 1200
Set Height 600

Type "rockset query 'SELECT COUNT(*) FROM demo.movies'"
Sleep 500ms
Enter

Sleep 2s

Type "rockset query"
Sleep 500ms
Enter

Type "SELECT COUNT(*)"
Sleep 500ms
Enter
Type "FORM demo.movies;"
Sleep 500ms
Ctrl+A
Sleep 300ms
Right@300ms 3
Backspace@300ms 2
Type@300ms "RO"
Ctrl+E
Sleep 500ms
Enter
Ctrl+D

Sleep 2s

Type "echo 'SELECT COUNT(*) FROM demo.movies;' | rockset query"
Sleep 500ms
Enter

Sleep 5s

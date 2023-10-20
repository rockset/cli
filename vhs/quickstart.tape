Output vhs/quickstart.gif

# Set up a 1200x600 terminal with 14px font
Set FontSize 14
Set Width 1200
Set Height 600

Type "rockset version"
Sleep 500ms
Enter
Sleep 2s

Type "rockset get org"
Sleep 500ms
Enter
Sleep 2s

Type "rockset list workspaces"
Sleep 500ms
Enter

Sleep 2s

Type "rockset create workspace demo --description 'for demos'"
Sleep 500ms
Enter

Sleep 5s

Type "clear"
Enter

Type "cat vhs/movie_releases_ingest_transform.sql"
Sleep 2s
Enter

# create movie collection with
Type "rockset create sample collection --dataset movies --workspace demo --retention 1h --ingest-transformation-file vhs/movie_releases_ingest_transform.sql movies"
Sleep 500ms
Enter

Sleep 5s

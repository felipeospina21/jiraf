# VHS demo tape for jiraf
# Run: vhs demo.tape

Output docs/assets/demo.gif

Set Shell "bash"
Set FontFamily "JetBrainsMono Nerd Font"
Set FontSize 14
Set Width 1600
Set Height 900
Set Padding 10

# Launch in demo mode (mocked data, no API calls)
Type "jiraf -demo"
Enter
Sleep 2s

# Select a board (Project Board)
Sleep 500ms
Enter
Sleep 2s

# Browse tickets
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Sleep 2s

# Open filter
Type "/"
Sleep 1s

# Navigate to Priority section and toggle "High"
Tab
Sleep 500ms
Space
Sleep 500ms

# Apply filter
Enter
Sleep 2s

# Open ticket details
Enter
Sleep 2s

# Scroll details
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Sleep 1s

# Create a comment
Shift+C
Sleep 1s
Type "Looks good, let's ship it"
Sleep 2s

# Submit comment (opens confirm)
Ctrl+S
Sleep 1s

# Confirm submission
Type "y"
Sleep 4s

# Close details panel
Escape
Sleep 1s

# Toggle help
Type "?"
Sleep 2s
Escape
Sleep 1s

# Quit
Ctrl+C

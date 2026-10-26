# Findings are Indicators or Matches, never a malware verdict

The GUI must explain why a file was flagged. Collapsing everything into one “finding” would hide whether HexHound saw a structural oddity or a detection method firing. Indicators are suspicious or structurally unusual conditions. Matches are named detection methods whose criteria were met (for example a hash list). Neither is presented as proof of malware. No findings means current methods reported nothing, not that the file is clean.

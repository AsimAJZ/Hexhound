# Identification produces format hypotheses, not identity

Extension, magic bytes, and similar cheap signals are untrusted. They yield zero or more format hypotheses on a scanned file. They do not prove the file is that format and they are not Match findings. Format analyzers parse untrusted bytes and emit Indicators (including extension/magic disagreement and polyglot / extra hypotheses). The identification labels MATCH/MISMATCH are retired so they cannot be confused with a Match.

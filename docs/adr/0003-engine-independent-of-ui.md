# The analysis engine is independent of CLI and GUI

HexHound’s product surface will be a desktop GUI, with a CLI kept as a development harness. The engine therefore accepts scan inputs and returns structured results (scanned files, findings, risk ratings). It must not print, hardcode corpus paths, or depend on a terminal or windowing toolkit. CLI and GUI are presentations of the same engine.

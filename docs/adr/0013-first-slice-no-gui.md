# First implementation phase: engine plus CLI, not the GUI

The product remains a desktop GUI later (ADR 0002). This implementation phase does not build the GUI. The first engine slice is the hint layer (extension, actual prefix magic, conflict Indicators) plus the PNG format analyzer. No Match methods, ZIP unpacking, PE/PDF/JPEG analyzers, hash detection, or required trailing-data/appended-PE hunting. PNG trailing-byte checks may come later and are not required now. The CLI must present the engine’s structured results and must not imply broader Coverage than is implemented.

# HexHound is a static risk analyzer, not an AV engine

HexHound’s job is to report risk-relevant conditions in untrusted files without executing them. Format validation exists to supply evidence for that analysis, not as a conformance product. HexHound must not automatically claim that a file is malware. Hash lists and other detection methods may be added later, but they remain distinct from suspicious indicators: a match is “this detection method fired,” not “this is malware.”

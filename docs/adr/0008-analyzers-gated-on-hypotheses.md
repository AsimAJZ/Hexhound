# A format analyzer runs only when its format is hypothesized

HexHound must not parse every file as PNG (or as any one format). An analyzer runs only if the scanned file has at least one hypothesis for that format. The hypothesis is still not proof of validity: the analyzer must check structure and treat all bytes as untrusted. Missing analyzers for other hypotheses do not mark analysis Incomplete; they are a Coverage concern, not Analysis status.

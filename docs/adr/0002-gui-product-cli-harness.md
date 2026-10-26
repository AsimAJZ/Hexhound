# Desktop GUI is the product; CLI is a development harness

The long-term product is a desktop application: the operator selects a folder, runs a scan, and inspects findings, explanations, and an overall risk rating. A CLI is acceptable during development because it is easier to test. Analysis must not be designed as a terminal program with a GUI bolted on later; presentation (CLI or GUI) sits on top of structured scan results.

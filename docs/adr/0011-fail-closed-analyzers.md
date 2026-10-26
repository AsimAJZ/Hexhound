# Analyzers fail closed; one file must not abort the scan

When required prerequisites are unusable, an analyzer must not continue into unsafe parsing or decompression. It emits Indicators supported by the evidence and stops according to its rules. Complete means in-scope methods finished by those rules; Incomplete means intended work could not finish because of an error or limit. Policy limits that stop work are both Incomplete and explained as HexHound policy. A hostile file must never abort the rest of the directory scan.

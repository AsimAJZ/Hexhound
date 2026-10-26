# HexHound

HexHound is a static risk analyzer. It inspects untrusted files without executing them and reports conditions that raise risk. It does not declare that a file is malware.

## Language

**Operator**:
The trusted person who runs HexHound and chooses what to inspect.
_Avoid_: user, client, customer

**Scanned file**:
The untrusted bytes HexHound inspects. Choosing the path is the operator’s act; the contents of that path are not trusted.
_Avoid_: sample, specimen, input, artifact, payload

**Static analysis**:
Inspection of a scanned file’s bytes without executing them, without following URLs found in them, and without making network requests because of them.
_Avoid_: sandbox, detonation, dynamic analysis, AV engine

**Scan**:
One operator-initiated static analysis of a chosen directory tree. Each ordinary filesystem file in that tree is a scanned file. A scan does not unpack archives and does not follow symlinks.
_Avoid_: job, run, session

**Skipped entry**:
A directory-tree entry HexHound did not inspect, recorded with its path and reason (this phase: symlinks). It is not a scanned file, not an Incomplete result, and not an Indicator by itself.
_Avoid_: ignored file, filtered file, unscanned file (as the type name)

**Finding**:
A reported condition on a scanned file, with an explanation. A finding is either an Indicator or a Match. It is not a claim that the file is malware.
_Avoid_: alert, verdict, detection, malware

**Explanation**:
The operator-facing reason a finding exists. It states whether the condition is a file-format specification violation or a HexHound resource-policy limit.
_Avoid_: message, log line, reason (as the type name)

**Indicator**:
A finding that the scanned file is suspicious or structurally unusual under static analysis.
_Avoid_: anomaly, suspicion, warning (as a type name)

**Match**:
A finding that a named detection method’s criteria were met. It is not a claim that the file is malware.
_Avoid_: signature hit, virus name, IoC (as a type name)

**Risk rating**:
An ordinal summary of findings: Low, Medium, High, or Critical. It is not a probability of malware. A scanned file with no findings is Low; that is not Clean, Safe, or Clear. Incomplete does not by itself raise the rating. Malformed structure does not by itself raise a file to High or Critical. A scan’s overall risk rating is the maximum among scanned files (Critical > High > Medium > Low). Skipped entries do not participate. Mapping from individual finding kinds to levels is not defined yet.
_Avoid_: score, confidence, threat level, malware probability, Clean, Safe, Clear

**Attribution**:
HexHound does not state that a condition is malicious or that it is merely accidental corruption. Findings state the observed condition and the rule (format specification or HexHound policy).
_Avoid_: malware, corrupted, benign, false positive (as HexHound verdicts)

**Hint layer**:
Cheap identification from extension and prefix magic bytes actually read from the scanned file. It produces format hypotheses and may emit conflict Indicators. It does not parse format structure, and it does not see embedded or appended formats.
_Avoid_: sniffer, MIME sniff, file(1)

**Format hypothesis**:
A cheap, untrusted guess at a scanned file’s format, taken from extension or HexHound’s magic-byte table. Only bytes actually read from the scanned file count as evidence. Magic hypotheses are prefix-only. A scanned file may have several hypotheses. A hypothesis is not a finding and not proof of format.
_Avoid_: detected format, MIME type, file type, MATCH

**Format analyzer**:
Untrusted parsing of a hypothesized format that emits Indicators. It must handle malformed input without executing content. Running an analyzer does not mean the file is valid.
_Avoid_: parser as a product name, validator (as if conformance were the goal)

**In-scope methods**:
The hint-layer checks, implemented format analyzers for this file’s hypotheses, and enabled Match methods that this version of HexHound applies to a scanned file.
_Avoid_: full scan, deep scan (undefined)

**Coverage**:
How much of HexHound’s possible analysis applied to a scanned file. Coverage is not Analysis status. Complete does not mean exhaustive. Coverage reporting is deferred.
_Avoid_: using Complete to mean “every conceivable check ran”

**Analysis status**:
Whether in-scope methods finished for a scanned file: Complete or Incomplete. Incomplete carries a reason and is not a finding. A hypothesis with no implemented analyzer does not make analysis Incomplete. Indicators found before a stop still stand.
_Avoid_: error, exception, crash (as the status name)

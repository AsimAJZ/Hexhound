# Do not follow symlinks; report them as skipped entries

Following directory or file symlinks can escape the folder the operator chose and can cycle. This phase does not follow symlinks. Each symlink is a skipped entry (path + reason), not a scanned file, not Incomplete, and not an Indicator merely for being a symlink. The operator must still be able to see that HexHound did not inspect those paths.

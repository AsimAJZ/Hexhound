# A scan covers filesystem files; it does not unpack archives

A scan is a walk of the operator-chosen directory tree. Each ordinary filesystem file is one scanned file. Unpacking ZIP/Office/other containers is deferred: it enlarges the trust and denial-of-service surface (zip bombs, path traversal, nested members). Magic-byte labels such as “zip” remain identification hints, not extraction.
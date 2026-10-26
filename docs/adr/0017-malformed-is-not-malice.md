# Malformed structure is evidence of a spec violation, not of malice

If PNG (or later format) prerequisites fail, HexHound must not decompress further, must emit Indicators for the evidence, and may mark that method Complete because it followed stop rules. Those Indicators explain the specific specification or policy rule. They must not automatically map to High or Critical. HexHound must not claim the file is malware, and must not claim it is “just corrupted,” unless some other defined finding supplies that evidence (none does today).

# Every scanned file has Complete or Incomplete analysis status

If analysis cannot finish (unreadable path, policy limit, parse abort, walk cut off), the operator must see that. Incomplete is not a third finding type and not a Match. Findings already produced still stand. A blank finding list plus Incomplete must not be read as “nothing to see.” Resource limits may also emit an Indicator that the file demanded excessive resources.

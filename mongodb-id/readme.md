# MongoDB ObjectId

- 12 bytes: `{timestamp}{random}{incrementing counter}`
    - 4 bytes timestamp
    - 5 bytes random
    - 3 bytes incrementing counter
- Feature
    - Roughly sortable: epoch to the left most
    - Unique: timestamp (always move forward) + random (to avoid collision) + counter (tie-breaker for same timestamp)

```sh
# create new ObjectId (hex: 1 char = 4 bits -> 2 chars = 1 byte)
test> ObjectId()
id = ObjectId('66eebd16e04acbceb6c76a8b')
# breakdown: 66eebd16.e04acbceb6.c76a8b
# 4 bytes timestamp: 66eebd16 -> 1726922006 -> 2024-09-21T12:33:26.000Z
# 5 bytes random: e04acbceb6 -> 963327545014
# 3 bytes incrementing counter: c76a8b -> 13068939

# get timestamp
test> id.getTimestamp()
ISODate('2024-09-21T12:33:54.000Z')

# get timestamp in epoch milliseconds
test> id.getTimestamp().getTime()
1726922034000

# object id is comparable
test> past = ObjectId()
ObjectId('66efdb489c18bdc04d1681ed')
test> now = ObjectId()
ObjectId('66efdb509c18bdc04d1681ee')
test> now > past
true
```

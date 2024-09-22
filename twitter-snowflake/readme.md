# Twitter Snowflake breakdown

- 64-bit integer: `0` + `41-bit timestamp` + `10-bit workerId` + `12-bit sequenceId`
- 1 bit unused
- 41 bits timestamp (milliseconds since epoch)
- 10 bits worker id
- 12 bits sequence number (within worker)

# Features

- Roughly sortable (with millisecond precision), as timestamp is on the left
- Unique

# Others
- https://github.com/sony/sonyflake
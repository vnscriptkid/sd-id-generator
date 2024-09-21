# UUID

- 128 bit ~ 16 bytes (vs 4 bytes for int, 4x increase)
    - pros: almost no collision
    - cons: insufficient for large scale (bloated index, can't fit in RAM)

# MySQL UUID

```sh
# create new UUID
mysql> select uuid();
# +--------------------------------------+
# | uuid()                               |
# +--------------------------------------+
# | 1295b4ac-7819-11ef-a589-0242ac130002 |
# +--------------------------------------+
```

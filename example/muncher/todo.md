### Monsters
With a several monsters should not try to yield its turn, it should instead try to take an alternative route. (This is currently a panic condition)

### Map
Clean up these:
This still happens when the bottom of the map scrolls of when going up. Repro by going far down, then back up.
 * 2017/12/07 23:17:02 Out of bounds Requested invalid position (16,-3) on board of dimensions 100x30
 * 2017/12/07 23:17:02 Out of bounds Requested invalid position (17,-3) on board of dimensions 100x30
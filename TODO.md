# TODO

###### Started todo on 11-26-2026

---

## Fixes

### Not Started
- [ ] Clean up the hovering outline left over from moving into a new chunk

### In-Progress
- [ ] Hovering, Block breaking & placing should work for a certain distance
- [ ] (maybe) Change `WorldX` and `WorldY` of structs to type `int`

### Done
- [x] Fix jump collisions

---

## Features

### Not Started
- [ ] `Queue generation`: Push missing neighbor requests onto a small queue and generate 1–2 chunks per tick (or spread work across frames) so movement stays smooth.

### In-Progress

### Done
- [x] Add vertical collisions
- [x] Add horizontal collisions

---

## Optimizations

### Not Started

### In-Progress

### Done
- [x] UpdateScreenPos for the current chunk and its surrounding chunks instead of all chunks

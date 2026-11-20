# testdata

`sample.log` is a real six-second capture from a sailing vessel's NMEA 2000
bus, recorded with `candump -L`. It holds ~1,500 frames across ~28 PGNs:
heading, attitude, rate of turn, rudder, water depth, boat speed, wind,
temperature, battery status, GNSS satellite/DOP data, autopilot status,
heartbeats, and an address claim — including fast-packet PGNs that span
multiple frames.

Frames whose PGNs carry vessel position, routes, or identity (GNSS position,
AIS, DSC, waypoint data) were removed for privacy, along with a few
high-rate proprietary PGNs that add volume without adding variety. Everything
else is byte-for-byte as it appeared on the wire.

`h5000-tables.raw` contains six assembled PGN 130822 messages from a B&G H5000
CPU. They come from `analyzer/tests/nested-repeating.in` at upstream revision
`17a4eb3436a7903c8e5cfc52b16a9822dd9b62e6`. Section 15 reaches the 223-byte
fast-packet limit during its second table.

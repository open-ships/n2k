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

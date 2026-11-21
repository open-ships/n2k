# Upgrading to the v10.0.0 PGN database in n2k v1.12.0

If you consume typed messages or their JSON, check the fields and lookup names
below. Re-decoding stored wire payloads applies the new signedness, masks,
and unavailable-value rules automatically.

## B&G H5000 settings and calibration

A type switch that handled `*pgn.NavicoUdbDatabaseObjectDump` for command 6
of PGN 130822 will now receive a more specific type for these sections.

| Section | Go message | How to read it |
| --- | --- | --- |
| 10 | `NavicoUdbDatabaseSettingsDump` | Each `Repeating1` entry carries `Flag`, `Length`, and raw `Value` bytes. The setting is identified by the record’s position. Values can have different widths. |
| 11–15, 18 | `NavicoUdbDatabaseTableDump11`, `12`, `13`, `14`, `15`, `18` (same type-name prefix) | Each `Repeating1` entry is one table. Its `Repeating2` slice contains the length and `*float32` value of each received record. |

For example, a section 14 message exposes its values in table order.

```go
if dump, ok := msg.(*pgn.NavicoUdbDatabaseTableDump14); ok {
    for i, table := range dump.Repeating1 {
        for j, record := range table.Repeating2 {
            if record.Value != nil {
                fmt.Printf("table %d, value %d = %g\n", i, j, *record.Value)
            }
        }
    }
}
```

JSON uses `repeating1[i].repeating2[j].value`. On input, a table’s `Count`
records what the device declared. Iterate the slice for values actually
received. In the 223-byte section 15 capture the second table declares 22
values and contains 19. Forwarding an untouched message preserves the exact
packet. Editing a table or constructing a new message derives each table’s
count from its slice length. Set a new float record’s `Length` to 4. Editing a
section 10 setting derives its length from its byte slice.

## Go fields and JSON values

The field names and JSON keys in this table stay the same. Fields that changed
from byte slices previously marshaled to base64 strings. Their new non-nil
values marshal as numbers. Numeric fields with an unavailable value decode to
nil and their raw key is omitted. Bitmaps retain all bits, including unnamed ones.

| Message and field | New Go type | Meaning |
| --- | --- | --- |
| `FurunoNavpilotStatus.RudderAngle` | `*int64` (was `*uint64`) | Signed ticks of 0.0001 rad. Use `RudderAngleValue()` for radians. The signed unavailable value is 0x7fff. |
| `Seatalk1PilotMode.PilotModeData` | `*int64` (was `[]uint8`) | Signed rudder position in degrees. Use `PilotModeDataValue()` or the JSON `physical.pilotModeData` measurement. |
| `SeatalkAlarm.AlarmPriority`, `SeatalkPilotMode.SubMode` | `*uint64` (was `[]uint8`) | Unsigned 16-bit values. |
| `WindlassControlStatus.SpeedControl` | `*uint64` (was `[]uint8`) | Unsigned 8-bit value. |
| `MercuryEngineStatus.Helm`, `MercuryEngineStatus.SubHelm` | `*uint64` (was `[]uint8`) | Upper and lower four bits of the helm byte, respectively. |
| `MercuryEngineStatus.Capabilities` | `*uint64` (was `[]uint8`) | An 8-bit mask. Use `MercuryEngineCapabilityBeacon` and `MercuryEngineCapabilityCruiseControl` to test individual bits. |
| `BepMarineCzoneCircuitStatus.Bitmap`, `BepMarineCzone65301.StatusBitmap` | `*uint64` (was `[]uint8`) | Little-endian 32-bit masks. The database leaves individual bits unnamed. |

## Lookup names

Replace these constants in Go. If you store the text returned by `String()`,
update the affected labels too. Numeric codes are stable.

| Previous constant | Replacement | New label |
| --- | --- | --- |
| `CzoneAlarmTypeHighTemperatureValue51` | `CzoneAlarmTypeHighTemperature51` | High Temperature (51) |
| `Seatalk1CommandAlarmONOFFForGuardValue171` | `Seatalk1CommandAlarmONOFFForGuard171` | Alarm ON/OFF for Guard (171) |
| `SeatalkCommandAutoTurnValue17` | `SeatalkCommandAutoTurn17` | Auto Turn (17) |
| `SimnetAlarmIdTrueWindShiftValue124` | `SimnetAlarmIdTrueWindShift124` | True wind shift (124) |
| `SimnetAlarmIdEVCComErrorValue247` | `SimnetAlarmIdEVCComError247` | EVC Com Error (247) |

## Payload sizes and measurements

Newly encoded `FurunoSensorSetup` messages grow from 28 to 30 bytes,
`FusionVersions` from 11 to 12, and `FusionTrackPosition` from 8 to 9. The
additional bytes are spare and encode as zero. An untouched older capture
retains its original length and padding. `FurunoSensorSetup.AirPressureOffset`
now has a physical range of −327,680 to 327,640 Pa, in steps of 10 Pa.

The NMEA 2000 input uses schema format 2.9.0 at revision
`17a4eb3436a7903c8e5cfc52b16a9822dd9b62e6`, with original SHA-256
`396888a761aaf47743a387f616bddbe5e2ec1717068fb837ff41adf458b7592b`.
The manifest covers 647 variants over 348 PGNs. It marks 110 complete for
both encoding and decoding, and zero as hardware-verified. All seven new
Navico definitions are marked incomplete by the source database.

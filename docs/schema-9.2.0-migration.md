# PGN database v9.2.0 migration (n2k v1.10.0)

The offline input is `cmd/pgngen/schema.json`, pinned to revision
`5add9cefe5711d3d3ab600ea4b37caa95a17f2f6` and original-source SHA-256
`35002ef57b0793654fc39a254df808b84deb8ab51df1b4e43b98bf3f0452bce8`.
The database version is 9.2.0, its schema format remains 2.8.0, and its protocol
is `nmea2000`. This upgrade includes v9.1.0 and addresses issues #41 and #42.
Coverage grows from 632 to 640 typed variants over the same 348 PGNs; 109
variants are marked complete for encode/decode and none are hardware-verified.
Refresh with `PGN_REFRESH_SOURCE=1 just pgn-sync`; ordinary `just pgn-sync`
works offline. Generation retains the full source attribution and license.

This release intentionally refines the evolving v1 API. Update affected Go
fields, JSON schemas, and lookup constants before upgrading.

## DSC addresses and dispatch

`DscCallInformation` and `DscDistressCallInformation` fields `DscMessageAddress`
and `MmsiOfShipInDistress` change from `*uint64` to `string`. The new
`DscDistressAcknowledgementOrRelay` uses the same string representation.
The old numeric values were little-endian binary interpretations of decimal
pair bytes, not usable MMSIs. Re-decode stored payloads; do not convert the old
integer to base-ten text. `PhysicalValue` now rejects these non-measurement
fields instead of returning their old binary interpretation.

For example, wire bytes `00 22 38 4e 5a` become `"003456789"`. Each byte
encodes two decimal digits as a binary number from 0 through 99, not BCD.
Leading zeros survive. A final zero is omitted as MMSI padding, except for
`DscCallInformation.DscMessageAddress` with format 102 (geographical area):
that same payload becomes `"0034567890"` with all ten digits retained.
The selector is retained in source/runtime metadata and the manifest as
`TrailingZeroIsDigitField`/`TrailingZeroIsDigitValue` (camelCase in JSON).

Use nine digits for an MMSI or ten for a full address. Nine-digit values append
one padding zero on encode. Ten-digit values ending in zero are accepted only
for geographical addresses, avoiding ambiguous padding; geographical values
must have ten digits. Empty strings represent unavailable data and encode as
five `ff` bytes. Any byte above 99 or a truncated decimal field decodes as
unavailable. JSON values are strings and unavailable fields are omitted.
Unchanged decoded messages, including malformed/unavailable fields, re-encode
to the original bytes.

Dispatch selects `DscDistressCallInformation` when format is 112, then
`DscDistressAcknowledgementOrRelay` when category is 112, and otherwise
`DscCallInformation`. Format wins if both are 112. Switch on the new relay type
where the old decoder used `DscDistressCallInformation` for all distress
categories. Distress calls now expose category as an ordinary lookup field.

## Other field and wire changes

| API | Migration |
| --- | --- |
| `BepMarineCzoneStatusExtended.Records` | Replace the opaque bytes with `OutputCurrent0`–`7`, `OutputLevel0`–`7`, and `OutputUnknownFlags0`–`7`. Current scales by 0.2 A, level by 0.1%; flags are six-bit byte slices. `Dipswitch` keeps its Go/JSON name but describes the ZCF module ID. |
| `BepMarineProprietaryPgn130822.Data` | Replace with `ModuleId`, `Page`, and the same eight current/level/flags groups. Current scales by 0.1 A and level by 0.1%. The payload is 28 bytes, down from the old opaque 223-byte definition. |
| `BepMarineProprietaryPgn130825.Data` | Replace with `ModuleId`, `Page`, `UnknownHeaderByte`, and `OutputCurrentRaw0`–`7` / `OutputLevelRaw0`–`7`. These are packed 11-bit raw values with no asserted physical scale. The payload is 27 bytes. |
| `BepMarineProprietaryPgn65294.Data` | The first old data byte is now `Brightness` (with physical accessors). `Data` contains only the remaining five bytes. |
| `FurunoMultiSatsInViewExtendedRepeating1.Azimuth` | Change `*int64` to `*uint64`. Values span 0 to 2π radians with 0.0001 rad ticks; `ffff` is unavailable. Values above π now decode correctly. |
| `Heave` | Accept new `Delay` and `DelaySource` fields and `physical.delay` JSON. Delay is in 0.01 s ticks; use `DelayValue` / `SetDelayValue`. `DelaySourceConst` names the four source choices. |
| `SimnetEvent` | Accept new `CommandType` and `Event` fields/JSON keys, replacing formerly reserved bytes. |

Re-decode opaque CZone/BEP payloads instead of copying old `Records`/`Data`
bytes into a single new field. CZone Go type names remain unchanged; their
metadata descriptions now identify AC Output Interface, DC Output Interface,
and Control X Plus. Physical accessors and JSON expose the decoded AC/DC
current and level values.

`LumishoreProprietary`, `MaretronProprietaryConfiguration`,
`NavicoBoatSpeedPolarTable`, and
`Message0x1ff000x1ffffManufacturerSpecificFastPacketNonAddressed` now have
223-byte payloads. Their binary fields shrink to 221, 216, 216, and 221 bytes,
respectively. Update buffers and authored payloads that used the old nominal
225/230-byte layouts. Untouched decoded payloads retain exact replay.

Seven added types cover Furuno MCU-006 remote messages: `FurunoCommand`,
`FurunoRemoteKeyboardReport`, `FurunoRemotePointerReport`,
`FurunoRemoteAnnounce`, `FurunoRemoteBacklight`, `FurunoRemoteInfo`, and
`FurunoRemoteHost`. These provide wire codecs, not an automatic pairing or
remote-control session.

## B&G lookup constants

The table lists every removed constant or existing constant whose value
changed. Select the replacement for the intended channel, using its corrected
meaning. The final column maps the *old wire key* to its current name.

| Previous constant | Old key | Current name for that key |
| --- | --- | --- |
| `BandgKeyValueOutsideTemperature` | 28 | `BandgKeyValueOutsideTemperatureFChannel` |
| `BandgKeyValueOutsideTemperatureValue29` | 29 | `BandgKeyValueOutsideTemperature` |
| `BandgKeyValueOutsideTemperatureValue30` | 30 | `BandgKeyValueWaterTemperatureFChannel` |
| `BandgKeyValueMagneticVariation` | 52 | `BandgKeyValueAttitudeRoll` |
| `BandgKeyValueCurrentSet` | 73 | `BandgKeyValueHeading` |
| `BandgKeyValueWindSpeedTrue` | 79 | `BandgKeyValueWindSpeedApparentKnotsChannel` |
| `BandgKeyValueWindAngleTrue` | 81 | `BandgKeyValueWindAngleApparent` |
| `BandgKeyValueWindSpeedTrueValue85` | 85 | `BandgKeyValueWindSpeedTrue` |
| `BandgKeyValueWaterTemperatureValue86` | 86 | `BandgKeyValueWindSpeedTrue2` |
| `BandgKeyValueTrueWindDirection` | 89 | `BandgKeyValueWindAngleTrue` |
| `BandgKeyValueWindDirection` | 109 | `BandgKeyValueTrueWindDirection` |
| `BandgKeyValueCurrentSetValue132` | 132 | `BandgKeyValueCurrentSet` |
| `BandgKeyValueAttitudeRoll` | 155 | `BandgKeyValueAttitudePitch` |
| `BandgKeyValueDepthValue195` | 195 | `BandgKeyValueDepthFathomsChannel` |
| `BandgKeyValueWaterSpeedValue235` | 235 | `BandgKeyValueSpeedOverGround` |
| `BandgKeyValueLaylineDistance` | 258 | `BandgKeyValueLaylineDistance1` |
| `BandgKeyValueLaylineDistanceValue259` | 259 | `BandgKeyValueLaylineDistance2` |
| `BandgKeyValueStartLineLatitude` | 263 | `BandgKeyValueMOBDRLatitude` |
| `BandgKeyValueStartLineLongitude` | 264 | `BandgKeyValueMOBDRLongitude` |

Several names remain but now select different keys: `OutsideTemperature`
28→29, `CurrentSet` 73→132, `WindSpeedTrue` 79→85, `WindAngleTrue` 81→89,
`TrueWindDirection` 89→109, and `AttitudeRoll` 155→52 (all with the
`BandgKeyValue` prefix). Audit these references even when compilation succeeds.
Lookup `String()` labels also follow the corrected channel names.

This upgrade imports the NMEA 2000 database and its wire semantics. Upstream
CLI/server features, including mDNS and command-line flags, are separate from
n2k's Go API.

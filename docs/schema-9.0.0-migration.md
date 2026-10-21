# PGN database v9.0.0 migration (n2k v1.9.0)

The offline input is `cmd/pgngen/schema.json`, pinned to revision
`5cf36705eef9fd0228ec63d428b182ed968df978` and original-source SHA-256
`726a222fb7e675a487441c7f9ee917bbaa9716e4bf1b07164d62a44e1a79abd0`.
The database version is 9.0.0, its schema format is 2.8.0, and its declared
protocol is `nmea2000`. This upgrades directly from the v7.1.0 database in
n2k v1.8.1, with 632 typed variants over 348 distinct PGNs. Refresh with
`PGN_REFRESH_SOURCE=1 just pgn-sync`; ordinary `just pgn-sync` works offline.
Generation retains the complete source attribution and license in the JSON.

This release intentionally refines the evolving v1 API. Regenerated Go types
and JSON keys follow the new definitions. Update affected field references,
pointer types, JSON schemas, and lookup constants before upgrading.
Go 1.26.10 or newer is required for the standard-library security fixes.

## Struct fields

| API | Migration |
| --- | --- |
| `FusionAuxGain.Gain`, `FusionBalance.Value` | Change `*uint64` to `*int64`; negative values are signed 8-bit ticks. |
| `AirmarDeviceInformation.SupplyVoltage` | Change `*uint64` to `*int64`; signed 16-bit ticks scale by 0.01 V. Use `SupplyVoltageValue` or `SetSupplyVoltageValue` for physical volts. |
| `UtilityPhaseBAcReactivePower.ReactivePower`, `UtilityPhaseCAcReactivePower.ReactivePower` | Change `*uint64` to `*int64`. The field is now 32 bits, with unsigned raw ticks and a -2,000,000,000 VAR physical offset. Re-decode stored payloads; the following fields move by 16 bits. |
| `AirmarPost.TestId`, `.TestResult` | Iterate `Repeating1 []AirmarPostRepeating1` and access each entry's `TestId` and `TestResult`. JSON moves `testId`/`testResult` into `repeating1` entries. Encoding derives `NumberOfIdTestResultPairsToFollow` from the slice length. A zero-result payload is four bytes. |
| `FusionZoneName.Number` | Use `Zone` (`number` becomes `zone` in JSON). |
| `FusionMenuLockId.LockId`, `.Flags` | The old 32-bit `LockId` field is `RequestToken`; the old 8-bit `Flags` field is `LockId`. The response carries both a request token and a lock ID. Update their meaning as well as their names. |
| `MaretronAnnunciator.Field4`, `.Field5`, `.Field6`, `.Field8` | Use `AnnunciatorInstance`, `AnnunciatorState`, `Pattern`, `AlertId`, respectively. `Field7` remains. JSON keys change accordingly. |
| `FurunoSixDegreesOfFreedomMovement.E`, `.F`, `.G`, `.H` | Replace the four unknown fields with signed 32-bit `RollRate`, `PitchRate`, `YawRate`. The old widths do not map one-to-one; re-decode stored payloads. Use the new `*Value`/`Set*Value` accessors for rad/s. `I` remains the final 16-bit field. |
| `YanmarEngineDataA.Data` | Use `UnknownSelector`, `EngineInstance`, `ThrottlePosition`, `TransmissionGear`, `EngineSpeed`, and `UnknownData`. `UnknownData` contains only the final byte. Re-decode stored payloads instead of assigning the old opaque six-byte payload to it. Throttle and engine-speed accessors expose percent and rpm. |

`AisClassAStaticAndVoyageRelatedData`, `AisClassBPositionReport`, and
`AisUtcAndDateReport` add `SequenceId`; `GnssSatsInViewRepeating1` adds
`GnssSystem`. Strict JSON consumers must accept the new keys. AIS UTC/date
reports include the newly described trailing spare/reserved fields.

## Wire and value semantics

- All 24 offset-encoded utility/generator power fields use unsigned 32-bit
  wire ticks stored in `*int64`. Apply `raw - 2,000,000,000` once to obtain W,
  VA, or VAR, or use the typed physical accessors. For example, raw
  3,000,000,000 represents +1,000,000,000, not a negative two's-complement
  integer. The valid raw range is 0–4,294,967,292; raw 2,147,483,647 is valid.
  4,294,967,293/294 remain raw sentinel values with unavailable physical
  measurements; 4,294,967,295 decodes to nil. Existing signed raw pointer
  types stay the same except for the two newly corrected fields above.
- `PowerFactorValue` and `SetPowerFactorValue` on `UtilityPhaseAAcReactivePower`,
  `UtilityPhaseBAcReactivePower`, `UtilityPhaseCAcReactivePower`,
  `UtilityTotalAcReactivePower`, `GeneratorPhaseAAcReactivePower`,
  `GeneratorPhaseBAcReactivePower`, `GeneratorPhaseCAcReactivePower`, and
  `GeneratorTotalAcReactivePower` now apply `raw / 16384 - 1`. A physical
  factor of 0.5 is raw 24576. Recompute physical JSON and value snapshots.
- Corrected signed minima include the full negative two's-complement endpoint;
  explicit reserved/error/unknown sentinel metadata now controls measurement
  availability consistently. Physical setters and fresh encoding enforce the
  corrected ranges. Unchanged decoded payloads retain their original bytes.
- `AirmarMasterReset` and `AirmarResetEeprom` add typed variants of PGN 126720;
  `AirmarPost` now decodes every result pair in PGN 130944.
- `ADDRESS` fields in `IsoCommandedAddress` and Simnet PGNs preserve the
  full 0–255 range; 255 is a non-nil value. This does not make reserved,
  cannot-claim, or broadcast addresses claimable by a Client.
- `FusionTuner.Track` is UTF-8 text converted from the declared `RDS_G0`
  repertoire when the wire text is not already valid UTF-8. Existing NUL and
  other string bytes retain their current treatment. Unmodified messages
  retain their original payload. New or edited text uses G0 when it decodes
  unambiguously, otherwise UTF-8, following the source decoder's precedence.
  Invalid UTF-8 Go strings are rejected on encoding this field. The 254-byte
  STRING_LZ limit applies after conversion to wire bytes.
- `FieldDescriptor.Encoding`, `SourceFieldDefinition.Encoding`, and the
  optional manifest field `encoding` retain the string repertoire.
- `HondaEngineData`, `HondaEngineStatus`, and `HondaEngineAlerts` now match
  manufacturer code 257 rather than 175. Update explicit manufacturer fields
  and test payloads; omitted match fields use the corrected value automatically.
- `FusionCapabilities` and `FusionMenuItem.Flags` are bit lookups;
  `FusionLowPassFilter.Filter` is a lookup; `FusionSetting.Value` and
  `FusionSettingsRepeating1.Value` use the setting-key-dependent lookup.
  These raw fields keep their Go types. Use the corresponding lookup enums
  and do not treat valid all-ones bitfields as missing numeric values.
- `SimnetCommandApStandby`, `SimnetCommandApNodrift`, `SimnetCommandApWind`,
  `SimnetCommandApNav`, and `SimnetCommandApHeading` now contain 32 trailing
  reserved bits. `SimnetApCommand` drops the two obsolete reserved fields.
  Newly encoded lengths follow the corrected layouts; unchanged decoded
  messages retain their original bytes.
- `GeneratorPhaseCAcPower.ApparentPowerValue` reports VA rather than VAR.
  Angular bounds and binary-fraction scales now use the more precise source
  values, including navigation/almanac, reactive power-factor, Maretron
  voltage, and Garmin turn-angle accessors. Recompute physical-value snapshots
  from raw fields rather than assuming the old rounded coefficients.

This upgrade imports the NMEA 2000 database. n2k's physical JSON continues to
use each field's schema units and `raw * Resolution + Offset`; upstream CLI
formatting and optional SI-output conversions are separate from this schema.

## Lookup constants

The tables below list every removed constant name and every existing constant
whose numeric value changed from v1.8.1. New constants are additional APIs.
Lookup `String()` output also follows updated source labels; it is display
text and should not be used as a stable identifier.

### Renamed constants

| Previous name | Replacement (same value) |
| --- | --- |
| `CzoneAlarmTypeShuntMistmatch` | `CzoneAlarmTypeShuntMismatch` |
| `FusionMessageIdDABUpdtateCommand` | `FusionMessageIdDABUpdateCommand` |
| `FusionMessageIdDvdPressIrKey` | `FusionMessageIdDVDPressIRKey` |
| `FusionMessageIdFapiGetServerInfo` | `FusionMessageIdGetServerInfo` |
| `FusionMessageIdFapiSetSourceEnabled` | `FusionMessageIdSetSourceEnabled` |
| `FusionMessageIdFapiSetSourceMultiroomEnabled` | `FusionMessageIdSetSourceMultiroomEnabled` |
| `FusionMessageIdFapiSetSourceName` | `FusionMessageIdSetSourceName` |
| `FusionMessageIdRequestHeadUnitDspSettings` | `FusionMessageIdRequestHeadUnitDSPSettings` |
| `FusionMessageIdScanMultirooomDevices` | `FusionMessageIdScanMultiroomDevices` |
| `FusionMessageIdSendDvdCommand` | `FusionMessageIdSendDVDCommand` |
| `FusionMessageIdSetLowPassFiler` | `FusionMessageIdSetLowPassFilter` |
| `FusionSettingIdAccessorySource` | `FusionSettingIDAccessorySource` |
| `FusionSourceTypeIpod` | `FusionSourceTypeIPod` |
| `FusionStatusMessageIdSiriusXmSignal` | `FusionStatusMessageIdSiriusXMSignal` |
| `ManufacturerCodeARKSEnterprisesInc` | `ManufacturerCodeARKSEnterprises` |
| `ManufacturerCodeAcrossOceansSystemsLtd` | `ManufacturerCodeAcrossOceansSystems` |
| `ManufacturerCodeAlltekMarineElectronicsCorp` | `ManufacturerCodeAlltekMarineElectronics` |
| `ManufacturerCodeAlpsAlpineCoLtd` | `ManufacturerCodeAlpsAlpine` |
| `ManufacturerCodeAuElectronicsGroup` | `ManufacturerCodeAEMPerformanceElectronics` |
| `ManufacturerCodeAventicsGmbH` | `ManufacturerCodeAventics` |
| `ManufacturerCodeBEPMarine2` | `ManufacturerCodeBEPMarine` |
| `ManufacturerCodeBSBArtificialIntelligenceGmbH` | `ManufacturerCodeBSBArtificialIntelligence` |
| `ManufacturerCodeBalmarLLC` | `ManufacturerCodeBalmar` |
| `ManufacturerCodeBarnacleSystemsInc` | `ManufacturerCodeBarnacleSystems` |
| `ManufacturerCodeBavariaYacts` | `ManufacturerCodeBavariaYachts` |
| `ManufacturerCodeBoeningAutomationstechnologieGmbHCoKG` | `ManufacturerCodeBoeningAutomationstechnologie` |
| `ManufacturerCodeCPACSystemsAB` | `ManufacturerCodeCPACSystems` |
| `ManufacturerCodeChetcoDigitialInstruments` | `ManufacturerCodeChetcoDigitalInstruments` |
| `ManufacturerCodeDataPanelCorp` | `ManufacturerCodeDataPanel` |
| `ManufacturerCodeEMMINETWORKSL` | `ManufacturerCodeEMMINETWORK` |
| `ManufacturerCodeEVEXCoLTD` | `ManufacturerCodeEVEX` |
| `ManufacturerCodeEgersundMarineElectronicsAS` | `ManufacturerCodeEgersundMarineElectronics` |
| `ManufacturerCodeElettromediaSpa` | `ManufacturerCodeElettromedia` |
| `ManufacturerCodeEnergySolutionsUKLTD` | `ManufacturerCodeEnergySolutionsUK` |
| `ManufacturerCodeFloscanInstrumentCoInc` | `ManufacturerCodeFloscanInstrument` |
| `ManufacturerCodeGMEAkaStandardCommunicationsPtyLTD` | `ManufacturerCodeGMEAkaStandardCommunications` |
| `ManufacturerCodeGobiusSensorTechnologyAB` | `ManufacturerCodeGobiusSensorTechnology` |
| `ManufacturerCodeHemisphereGPSInc` | `ManufacturerCodeHemisphereGPS` |
| `ManufacturerCodeHondaElectronicsCoLTD` | `ManufacturerCodeHondaElectronics` |
| `ManufacturerCodeHondaMotorCompanyLTD` | `ManufacturerCodeHondaMotor` |
| `ManufacturerCodeHumminbirdMarineElectronics2` | `ManufacturerCodeHumminbirdMarineElectronics` |
| `ManufacturerCodeHumphreeAB` | `ManufacturerCodeHumphree` |
| `ManufacturerCodeISSPROInc` | `ManufacturerCodeISSPRO` |
| `ManufacturerCodeITCInc` | `ManufacturerCodeITC` |
| `ManufacturerCodeInteractiveTechnologiesInc` | `ManufacturerCodeInteractiveTechnologies` |
| `ManufacturerCodeJLMarineYstems` | `ManufacturerCodeJLMarineSystems` |
| `ManufacturerCodeJapanRadioCo` | `ManufacturerCodeJapanRadio` |
| `ManufacturerCodeKobeltManufacturingCoLtd` | `ManufacturerCodeKobeltManufacturing` |
| `ManufacturerCodeKodenElectronicsCoLtd` | `ManufacturerCodeKodenElectronics` |
| `ManufacturerCodeKvasarAB` | `ManufacturerCodeKvasar` |
| `ManufacturerCodeLencoMarineInc` | `ManufacturerCodeLencoMarine` |
| `ManufacturerCodeLewmarInc` | `ManufacturerCodeLewmar` |
| `ManufacturerCodeLithionicsBatteryLLC` | `ManufacturerCodeLithionicsBattery` |
| `ManufacturerCodeMGEnergySystemsBV` | `ManufacturerCodeMGEnergySystems` |
| `ManufacturerCodeMarinesoftCoLTD` | `ManufacturerCodeMarinesoft` |
| `ManufacturerCodeMarlinTechnologiesInc` | `ManufacturerCodeMarlinTechnologies` |
| `ManufacturerCodeMcMurdoGroupAkaOroliaLTD` | `ManufacturerCodeMcMurdoGroupAkaOrolia` |
| `ManufacturerCodeMicroAirLLC` | `ManufacturerCodeMicroAir` |
| `ManufacturerCodeNKEMarineElecronics` | `ManufacturerCodeNKEMarineElectronics` |
| `ManufacturerCodeNanjingSandemarineInformationTechnologyCoLtd` | `ManufacturerCodeNanjingSandemarineInformationTechnology` |
| `ManufacturerCodeNaocontrolSL` | `ManufacturerCodeNaocontrol` |
| `ManufacturerCodeNautibusElectronicGmbH` | `ManufacturerCodeNautibusElectronic` |
| `ManufacturerCodeNobletec` | `ManufacturerCodeNobeltec` |
| `ManufacturerCodeOceanLEDMarineLtd` | `ManufacturerCodeOceanLEDMarine` |
| `ManufacturerCodeOceanSatBV` | `ManufacturerCodeOceanSat` |
| `ManufacturerCodeOffshoreSystemsUKLtd` | `ManufacturerCodeOffshoreSystemsUK` |
| `ManufacturerCodeOrcaTechnologoesAS` | `ManufacturerCodeOrcaTechnologies` |
| `ManufacturerCodePredictWindLtd` | `ManufacturerCodePredictWind` |
| `ManufacturerCodeProChargingSystemsLLC` | `ManufacturerCodeProChargingSystems` |
| `ManufacturerCodeQuickTeckElectronicsLtd` | `ManufacturerCodeQuickTeckElectronics` |
| `ManufacturerCodeRadianIoTInc` | `ManufacturerCodeRadianIoT` |
| `ManufacturerCodeRaritanEngineeringCompanyInc` | `ManufacturerCodeRaritanEngineering` |
| `ManufacturerCodeRideControllerLLC` | `ManufacturerCodeRideController` |
| `ManufacturerCodeRockfordCorp` | `ManufacturerCodeRockford` |
| `ManufacturerCodeSailmonBV` | `ManufacturerCodeSailmon` |
| `ManufacturerCodeSailormadeMarineTelemetryTetraTechnologyLTD` | `ManufacturerCodeSailormadeMarineTelemetryTetraTechnology` |
| `ManufacturerCodeSeaCrossMarineAB` | `ManufacturerCodeSeaCrossMarine` |
| `ManufacturerCodeSeaMacineRoboticsInc` | `ManufacturerCodeSeaMachineRobotics` |
| `ManufacturerCodeSeekeeper` | `ManufacturerCodeSeakeeper` |
| `ManufacturerCodeShadowCasterLEDLightingLLC` | `ManufacturerCodeShadowCasterLEDLighting` |
| `ManufacturerCodeSleipnerMotorAS` | `ManufacturerCodeSleipnerMotor` |
| `ManufacturerCodeSonicCorporation` | `ManufacturerCodeSonic` |
| `ManufacturerCodeSuzukiMotorCorporation` | `ManufacturerCodeSuzukiMotor` |
| `ManufacturerCodeTBSElectronicsBV` | `ManufacturerCodeTBSElectronics` |
| `ManufacturerCodeTheMarineGuardianLLC` | `ManufacturerCodeTheMarineGuardian` |
| `ManufacturerCodeTriskelMarineLtd` | `ManufacturerCodeTriskelMarine` |
| `ManufacturerCodeTrueHeadingAB` | `ManufacturerCodeTrueHeading` |
| `ManufacturerCodeUltraflexSpA` | `ManufacturerCodeUltraflex` |
| `ManufacturerCodeVesperMarineLtd` | `ManufacturerCodeVesperMarine` |
| `ManufacturerCodeWATTFuelCellCorp` | `ManufacturerCodeWATTFuelCell` |
| `ManufacturerCodeWetSoundsLLC` | `ManufacturerCodeWetSounds` |
| `ManufacturerCodeXantrexTechnologyInc` | `ManufacturerCodeXantrexTechnology` |
| `NavicoDataTypeGeoidalSeperation` | `NavicoDataTypeGeoidalSeparation` |
| `Seatalk1CommandDeviceIndentification` | `Seatalk1CommandDeviceIdentification` |

### Removed manufacturer codes

| Previous constant | Previous code | Migration |
| --- | --- | --- |
| `ManufacturerCodeFischerPanda` | 311 | Use the verified device NAME/manufacturer; retained definitions are `ManufacturerCodeFischerPandaGenerators` (356) and `ManufacturerCodeFischerPandaDE` (785). |
| `ManufacturerCodeGarmin2` | 645 | Use the verified device NAME/manufacturer; the retained Garmin code is `ManufacturerCodeGarmin` (229). |
| `ManufacturerCodeHondaMarine` | 175 | Use `ManufacturerCodeHondaMotor` (257), the retained Honda definition. |
| `ManufacturerCodeHondaMarine2` | 200 | Use `ManufacturerCodeHondaMotor` (257), the retained Honda definition. |
| `ManufacturerCodeHondaMarine3` | 225 | Use `ManufacturerCodeHondaMotor` (257), the retained Honda definition. |
| `ManufacturerCodeHondaMarine4` | 250 | Use `ManufacturerCodeHondaMotor` (257), the retained Honda definition. |

### Removed backlight constants

`SimnetBacklightLevelDayMode` (1) and `SimnetBacklightLevelNightMode` (4) are
removed. These bytes describe backlight intensity, not a day/night mode.
Choose the intended intensity: the retained table spans raw 0 (10%) through
99 (100%), with a step of 11 for each 10% increase. There is no day/night
replacement constant.

### Corrected numeric values

| Constant | Previous value | New value |
| --- | --- | --- |
| `FusionSettingBTAutoConnect` | 76 | 77 |
| `FusionSettingMaxVolumeStart` | 75 | 76 |
| `FusionSettingNullSetting` | 77 | 1000 |
| `ManufacturerCodeBEPMarine` | 116 | 295 |
| `ManufacturerCodeHumminbirdMarineElectronics` | 467 | 476 |

Select the constant for the intended device or setting using the corrected
source meaning. Do not preserve a previous numeric value merely to keep a
snapshot unchanged. In particular, the BEP Marine and Humminbird manufacturer
codes above change even though their constant names remain available.

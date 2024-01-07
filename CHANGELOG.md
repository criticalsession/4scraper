# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- New parallel download functionality, throttled at 20 concurrent threads at a time. Enable/disable parallel downloads in the config file. Default: `true`.
- Find functionality allows users to search a specific board for given keywords.

## [v1.3] - 2023-12-25

A literal Christmas 🎄 release. What am I doing with my life?

### Added

- New flag `--output` to define what directory the downloader should throw downloaded files into

### Fixed

- When enabling `OriginalFileNames=true` there's an extra `.` before the extension (e.g. `a78f817..jpg`)

## [v1.2] - 2023-12-23

### Added

- Added ability for users to set up a config file to manage basic settings

## [v1.1] - 2023-12-23

### Added

- Flags to allow users to run the program in silent mode

## [v1.0] - 2023-12-22

- Initial Release
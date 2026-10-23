# Changelog

All notable changes to this project will be documented in this file.

## [1.1.6](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.5...v1.1.6) (2026-10-23)


### Bug Fixes

* **text:** support content-type parameters, trim whitespace, and add nil receiver guards ([95698f8](https://github.com/microsoft/kiota-serialization-text-go/commit/95698f863021d9f7affbf5955051e6025b007af4))
* **text:** support content-type parameters, trim whitespace, and add nil receiver guards ([0bff70b](https://github.com/microsoft/kiota-serialization-text-go/commit/0bff70b83636fd73ca0719a208a57e566d286713))

## [1.1.5](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.4...v1.1.5) (2026-10-20)


### Bug Fixes

* **text:** parse byte as unsigned and preserve integer precision for int32/int64 ([2a8641f](https://github.com/microsoft/kiota-serialization-text-go/commit/2a8641f68c76bcd07382c212d270e19ab975dc47))
* **text:** parse byte as unsigned and preserve integer precision for int32/int64 ([2313caf](https://github.com/microsoft/kiota-serialization-text-go/commit/2313caf290edbef4d56103910a536ab23f7258df))
* **text:** preserve decimal and exponent parsing ([4447d28](https://github.com/microsoft/kiota-serialization-text-go/commit/4447d28bfa901c0db07d0e3bde47f48e4e8e1fe8))
* **text:** preserve leading plus byte values ([2837b8e](https://github.com/microsoft/kiota-serialization-text-go/commit/2837b8e1ceea5ed2e7612576d256334c8e0e4fff))

## [1.1.4](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.3...v1.1.4) (2026-10-16)


### Bug Fixes

* serialize text timestamps in RFC3339 format ([0287efb](https://github.com/microsoft/kiota-serialization-text-go/commit/0287efb8ab5c2a42812b95b803c81ddfc47dc1da))
* serialize text timestamps in RFC3339 format ([887e203](https://github.com/microsoft/kiota-serialization-text-go/commit/887e2031da0039d2d1d57b3d7b7ddd5487c30d5f))

## [1.1.3](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.2...v1.1.3) (2025-09-22)


### Bug Fixes

* extraneous assert dependency in the package ([c57dacc](https://github.com/microsoft/kiota-serialization-text-go/commit/c57dacccc3b89ae29a9a34ae63e460dfa23c6e4f))

## [1.1.2](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.1...v1.1.2) (2025-04-02)


### Bug Fixes

* removes common go dependency ([ba74296](https://github.com/microsoft/kiota-serialization-text-go/commit/ba74296bef18a77d0c367dda4a3908c04c38f9ed))
* removes common go dependency ([b2692be](https://github.com/microsoft/kiota-serialization-text-go/commit/b2692be0cef923152a0697665f94753653051f89))

## [1.1.1](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.1.0...v1.1.1) (2025-03-24)


### Bug Fixes

* upgrades common go dependency to solve triming issues ([da7894d](https://github.com/microsoft/kiota-serialization-text-go/commit/da7894dfc2fc351c8337c3570f50548994624098))
* upgrades common go dependency to solve triming issues ([c5b463c](https://github.com/microsoft/kiota-serialization-text-go/commit/c5b463c8f5308d4aea15fc520dacb6c9e6b89b7d))

## [1.1.0](https://github.com/microsoft/kiota-serialization-text-go/compare/v1.0.0...v1.1.0) (2025-03-13)


### Features

* upgrades required go version from go1.18 to go 1.22 ([c0886ad](https://github.com/microsoft/kiota-serialization-text-go/commit/c0886ad112bff2cf3fbea56dcd45922685d2e565))

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.7.1] - 2023-04-17

### Changed

- Implement serialization and parseNode listeners.

## [0.7.0] - 2023-01-23

### Added

- Added support for backing store.

## [0.6.0] - 2022-09-22

### Added

- Implement additional serialization method `WriteAnyValues` and parse method `GetRawValue` introduced in abstractions.

## [0.5.0] - 2022-09-02

### Added

- Added support for composed types serialization.

## [0.4.2] - 2022-09-02

### Changed

- Updated kiota abstractions library.

## [0.4.1] - 2022-06-07

### Changed

- Upgraded abstractions and yaml dependencies.

## [0.4.0] - 2022-05-26

### Changed

- Updated reference to kiota abstractions for enum response support.

## [0.3.0] - 2022-05-19

### Changed

- Updated kiota abstractions library.

## [0.2.0] - 2022-04-19

### Changed

- Upgraded abstractions to 0.4.0.
- Upgraded to go 18.

## [0.1.0] - 2022-03-30

### Added

- Initial tagged release of the library.

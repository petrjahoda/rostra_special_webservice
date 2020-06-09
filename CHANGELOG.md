# Rostra Special WebService Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.

## [2020.2.3.9] - 2020-06-9

### Added
- pseudocode for first controls

### Changed
- new handling for operation and workplace check

## [2020.2.3.8] - 2020-06-8

### Changed
- new handling for user and order check
- removed other checks (to do next)


## [2020.2.2.12] - 2020-05-12

### Added
- saving NOK fails to Zapsi
- mn1, mn2, mn3 controls

### Changed
- better handling enabling and disabling buttons, inputs, etc.


## [2020.2.2.11] - 2020-05-11

### Added
- radio button for clovek-stroj-serizeni
- proper handling and checking ok and nok pcs

### Removed
- radio button for ok and nok pcs

## [2020.2.2.4] - 2020-05-04

### Fixed
- proper handling first controls

### Added
- message parametr to displaying information from program to page

## [2020.2.1.29] - 2020-04-29

### Added
- added controls before enabling buttons, based on first_controls.png

## [2020.2.1.23] - 2020-04-23

### Added
- better username handling
- saving user syteline number as user login to zapsi

## [2020.2.1.22] - 2020-04-22

### Added
- software as service
- creating user, order and product in zapsi... if not present
- starting and ending terminalInputOrder in zapsi

## [2020.2.1.20] - 2020-04-20

### Added
- added structs for user, order, operation and workplaces
- fully working user input with checking in syteline
- updated rostra.html file 
- checking which button was pressed


## [2020.2.1.16] - 2020-04-16

### Added
- checking operations and workplaces


## [2020.2.1.15] - 2020-04-15

### Removed
- everything about workplacegroup was removed, because order+operation get list of workplaces from syteline

### Added
- displaying information about Syteline communication problem

## [2020.2.1.14] - 2020-04-14

### Added
- proper and better handling users and orders
- focus
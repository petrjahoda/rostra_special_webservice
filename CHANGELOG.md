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

## [2020.3.1.14] - 2020-07-14

### Fixed
- number of fixes in behavior

## [2020.2.3.29] - 2020-06-29

### Changed
- when closing order, additional checks are made: sytelineWorkplace.priznak_mn_1 == "0" || (sytelineWorkplace.priznak_mn_1 == "1" && countFromUser == (countFromZapsi - countFromSyteline))


## [2020.2.3.23] - 2020-06-23

### Changed
- when starting order, first checking sameorder and then sameuser
- when closing order, close is allowed only when amount from user is the same as difference between zapsi amount and syteline amount

### Added
- handling time_divisor (saving to zapsi2.device.Setting)
    - increasing with another open order added
    - reseting to 1 when no open order
- html table with open orders and calculated data

## [2020.2.3.22] - 2020-06-22

### Changed
- displaying order with additional data
- displaying operation with additional data
- displaying workplace with additional data
- GUI for OK, NOK, and NOK type


### Added
- suffix from order is trimmed from leading zeros
- when saving order to zapsi, type is saved to note
- when loading order from zapsi, loading from note, what to display on radio

### Fixed
- proper saving failtype


## [2020.2.3.12] - 2020-06-12

### Fixed
- when updating terminal_input_order in zapsi, updating just one, not all
- when inserting record to syteline, ANSI WARNING OFF has to be set


## [2020.2.3.11] - 2020-06-11

### Added
- complete saving to syteline

## [2020.2.3.10] - 2020-06-10

### Added
- complete first controls
- complete second controls
- complete saving to zapsi

### Changed
- removed table from html
- code formatted to more files 

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
fastlane documentation
================
# Installation

Make sure you have the latest version of the Xcode command line tools installed:

```
xcode-select --install
```

Install _fastlane_ using
```
[sudo] gem install fastlane -NV
```
or alternatively using `brew cask install fastlane`

# Available Actions
## Android
### android test
```
fastlane android test
```
Runs all the tests
### android analytics
```
fastlane android analytics
```
Submit a new Build to Crashlytics Beta
### android createscreenshots
```
fastlane android createscreenshots
```
Add new screen shots
### android alpha
```
fastlane android alpha
```
Deploy new alpha version to the Google Play
### android internal
```
fastlane android internal
```
Deploy new version to internal testers on Google Play

----

This README.md is auto-generated and will be re-generated every time [fastlane](https://fastlane.tools) is run.
More information about fastlane can be found on [fastlane.tools](https://fastlane.tools).
The documentation of fastlane can be found on [docs.fastlane.tools](https://docs.fastlane.tools).

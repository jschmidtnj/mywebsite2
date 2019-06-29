# joshuaschmidt

A new Flutter project.

## get emulator working

```bash
$ flutter emulators --create --name Pixel-test
Emulator 'Pixel-test' created successfully.
$ flutter emulators --launch Pixel-test
$ flutter run
```

## generate

### launch icons

`flutter pub run flutter_launcher_icons:main`

### signed key file

- `keytool -delete -alias key -keystore android/key/key.jks`
- `keytool -genkey -v -alias key -keystore android/key/key.jks -keyalg RSA -keysize 2048 -validity 10000`

### create ssh key for git allowing fastlane ios to work with match (ios)

see [this](https://help.github.com/en/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) for instructions

### encrypt secret files android - run in linux!

- `tar cvf secrets.tar lib/config.dart android/key.properties android/fastlane/serviceAccount.json android/key/key.jks`
- `gem install travis`
- `travis encrypt-file secrets.tar`

### create release version

- `flutter build apk --split-per-abi`
- `flutter build appbundle`
- see [this](https://flutter.dev/docs/deployment/android) for more info

### use fastlane CI/CD

see [this](https://flutter.dev/docs/deployment/cd#other-services) for config info

## Getting Started

This project is a starting point for a Flutter application.

A few resources to get you started if this is your first Flutter project:

- [Lab: Write your first Flutter app](https://flutter.dev/docs/get-started/codelab)
- [Cookbook: Useful Flutter samples](https://flutter.dev/docs/cookbook)

For help getting started with Flutter, view our [online documentation](https://flutter.dev/docs), which offers tutorials, samples, guidance on mobile development, and a full API reference.

## replace image html tag with Image in markdown

use regex: `<img\s.*?src=(?:'|")([^'">]+)(?:'|")`: [here](https://stackoverflow.com/a/1028370/8623391)
markdown: [here](https://github.com/flutter/flutter_markdown)

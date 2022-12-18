Test Instructions

First of all, activate firestore emulator by calling
    **firebase emulators:start --import=./data**

Note that you need to do it from ./tests/firebase,
as this directory contains all the info about the emulator

Then you can actually test everything by calling
    **go test**
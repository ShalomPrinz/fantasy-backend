# Test Instructions

**All commands should be executed from tests folder.**

## Emulator

First of all, activate firestore emulator:
    
    bash start_emulator.sh

Note that you need to do it from ./tests/firebase, as this directory contains all the info about the emulator.

Some tests might fail without the startup data, which happens if you forget adding the ***--import=./data*** part to the command

## Test

Then you can actually test everything:
    
    go test

## Coverage

    bash coverage.sh

#!/bin/sh

./bin/blendbase db:migrate
./bin/blendbase db:seed
./bin/blendbase server

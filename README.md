## Overview
This is Midsouth Makers 2nd generation door access software. Software is written in Go.
Door control software main purpose is to authenticate users. This is done by taking serial input from the
Arduino and verifying it with a local SQLite database. This database is checked every 30 seconds against
a separate web service that manages all of Midsouth Makers.

Webservice Repository:
https://github.com/svpernova09/midsouthmakers-rfid

## Hardware
- Wiegand RFID reader and keypad
- Arduino Nano that acts as a Wiegand to serial bridge
- Raspberry pi (or other linux machine) running the software.

## Wiegand to Serial
Arduino code for a Arduino Nano is located in the Arduino directory.

## Enviroment Variables
You will also need to set some environment variables. (I like using a tool like
[autoenv](https://github.com/kennethreitz/autoenv) or [direnv](https://direnv.net/)
to use while in development.  My `.env` file looks something like the following:)

export RFID_MQTTSERVER="tcp://192.168.10.5:1883"
export RFID_MQTTUSERNAME=""
export RFID_MQTTPASSWORD=""
export RFID_DBFILE="./rfid.sqlite"
export RFID_APIKEY="REALLYLONGSTRING"

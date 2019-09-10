## Overview
This is Midsouth Makers 2nd generation door access software. Software is written in Go.
Door control software main purpose is to authenticate users. This is done by taking serial input from the
Arduino and verifying it with a local SQLite database. This database is checked every 30 seconds against
a separate web service that manages all of Midsouth Makers.

Webservice Repository:
https://github.com/svpernova09/midsouthmakers-rfid

## Hardware
- Wiegand RFID reader and keypad, only tested with wiegand 34bit readers. 26bit wiegand could pose an issue.
- Arduino Nano that acts as a Wiegand to serial bridge
- Raspberry pi (or other linux machine) running the software.
- Relays to drive the door latch, buzzer, led, and doorbell

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


## Install Notes at Midsouthmakers
### Hardware
- Arudino nano 168p was used, but any atmega nano would work.
- Replacements for every element of the system are on hand. Extra keypad, set of relays, arduino, powersupply, etc. Contact Ben if there is a possible part needing replaced.
- 5v usb diode has been removed from the arduino nano to disable being powered over USB.

### Software Install
- IP Address: 192.168.1.11
- This repo was directly compiled on the local machine @ /usr/pi/go
- Enviroment varables were added to pi's .profile file
- .profile is called from crontab at reboot
- binary is called from crontab at reboot (/usr/pi/door_control)
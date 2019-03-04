## Development
You will also need to set some environment variables. (I like using a tool like
[autoenv](https://github.com/kennethreitz/autoenv) or [direnv](https://direnv.net/)
to manage these per-project.  My `.env` file looks something like the following:)

    export RFID_APIKEY='SOMEREALLYLONGSTRING!'
    export RFID_MQTTSERVER='tcp://192.168.10.5:1883'
    export RFID_MQTTUSERNAME='bseishen'
    export RFID_MQTTPASSWORD='ssshItsASecret'
    export RFID_SERIALPORT='/dev/ttyUSB0'
    export RFID_DBFILE='./rfid.sqlite'

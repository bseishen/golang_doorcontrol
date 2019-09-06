/***************
 * wigand to serial -- program that takes wigand data and sends it over uart
 *
 * Author: Ben S. Eishen
 *
 * Usage:
 * UART is 9600N1
 * RFID badge scan will send an R(for 34bit wiegand) or S(for 26bit wiegand) folowed by the ASCII conversion of the number ending in a newline ie "R1234567\n"
 * Button presses will be sent over as they are pressed in ASCII
 * Enter(or and escape will send \n and \esc\n respectively.
 *
 * If an 'O' is received, it will toggle the door relay, key buzzer and led for 3 seconds and release.
 *
 * If an 'E' is received, it will rapidly flash the the led 5 times indicating there was an error.
 *  
 * If an 'B' is received, ring the bell two short times
 * 
 * This project uses a wiegand library located at https://github.com/paulo-raca/YetAnotherArduinoWiegandLibrary
 * and has been tested with version 2.0
 * 
 * 
 */
#include <Wiegand.h>

//#define DEBUG 1

#define DOOR_RELAY 7
#define KEY_BUZZER 8
#define KEY_LED    5
#define BELL       6
#define PIN_D0     2
#define PIN_D1     3

//WIGAND Values for enter and esc (4bit)
#define ESC_BUTTON    0x0A //esc
#define ENTER_BUTTON  0x0B //enter


// The object that handles the wiegand protocol
Wiegand wiegand;
char inbyte;
byte escCount; //Pressed 3 times, ring the bell

void toggleDoor(){
  digitalWrite(DOOR_RELAY, LOW);
  digitalWrite(KEY_BUZZER, LOW);
  digitalWrite(KEY_LED,    LOW);

  delay(3000);

  digitalWrite(DOOR_RELAY, HIGH);
  digitalWrite(KEY_BUZZER, HIGH);
  digitalWrite(KEY_LED,    HIGH);
}

void showError(){
  byte count = 6;

  while(count){
    count--;
    digitalWrite(KEY_LED,    LOW);
    delay(100);
    digitalWrite(KEY_LED,    HIGH);
    delay(100);
  }
}

void ringBell(){
  byte count = 2;

  while(count){
    count--;
    digitalWrite(BELL,    LOW);
    delay(300);
    digitalWrite(BELL,    HIGH);
    delay(500);
  }
}

// Initialize Wiegand reader
void setup() {
  pinMode(DOOR_RELAY, OUTPUT);
  pinMode(KEY_BUZZER, OUTPUT);
  pinMode(KEY_LED, OUTPUT);
  pinMode(BELL, OUTPUT);

  digitalWrite(DOOR_RELAY, HIGH);
  digitalWrite(KEY_BUZZER, HIGH);
  digitalWrite(KEY_LED,    HIGH);
  digitalWrite(BELL,       HIGH);
  
  Serial.begin(9600);

  //Install listeners and initialize Wiegand reader
  wiegand.onReceive(receivedData, "Card readed: ");
  wiegand.onReceiveError(receivedDataError, "Card read error: ");
  wiegand.onStateChange(stateChanged, "State changed: ");
  wiegand.begin(Wiegand::LENGTH_ANY, true);

  //initialize pins as INPUT
  pinMode(PIN_D0, INPUT);
  pinMode(PIN_D1, INPUT);
}

// Continuously checks for pending messages and polls updates from the wiegand inputs
void loop() {
  // Checks for pending messages
  wiegand.flush();

  // Check for changes on the the wiegand input pins
  wiegand.setPin0State(digitalRead(PIN_D0));
  wiegand.setPin1State(digitalRead(PIN_D1));

  if(Serial.available()){
    inbyte = Serial.read();
    if(inbyte == 'O'){
      toggleDoor();
    }
    if(inbyte == 'E'){
      showError();
    }
    if(inbyte == 'B'){
      ringBell();
    }
  }
}

// Notifies when a reader has been connected or disconnected.
// Instead of a message, the seconds parameter can be anything you want -- Whatever you specify on `wiegand.onStateChange()`
void stateChanged(bool plugged, const char* message) {
#ifdef DEBUG
    Serial.print(message);
    Serial.println(plugged ? "CONNECTED" : "DISCONNECTED");
#endif
}

// Notifies when a card was read.
// Instead of a message, the seconds parameter can be anything you want -- Whatever you specify on `wiegand.onReceive()`
void receivedData(uint8_t* data, uint8_t bits, const char* message) {
    uint32_t t;

    //Check for rfid tag data, if there is data preface it with 'R' if its wiegand 34
    if(bits==32){
      Serial.print("R");
      t = ((uint32_t)data[3]&0xFF) | (((uint32_t)data[2]&0xFF) << 8) | (((uint32_t)data[1]&0xFF) << 16) | (((uint32_t)data[0]& 0xFF )<< 24);
      Serial.print(t);
      Serial.print("\n");
      escCount = 0;
    }
    //Check for rfid tag data, if there is data preface it with 'S' if its wiegand 26
    else if(bits==24){
      Serial.print("S");
      t = ((uint32_t)data[2]&0xFF) | (((uint32_t)data[1]&0xFF) << 8) | (((uint32_t)data[0]&0xFF) << 16);
      Serial.print(t);
      Serial.print("\n");
      escCount = 0;
    }

    //Keypresses are going to be 4 bit (usually).
    else if(bits==4){
      t = (data[0] & 0x0000000F);
      switch(t){
        //ESC
        case ESC_BUTTON:
          Serial.write(0x1B);
          Serial.print("\n");
          escCount++;
          break;
        //Enter
        case ENTER_BUTTON:
          Serial.print("\n");
          escCount = 0;
          break;
        //All other numbers
        default:
          Serial.print(t);
          escCount = 0;
      }
    }


    if(escCount > 2){
      escCount = 0;
      ringBell();
    }

#ifdef DEBUG
      Serial.print(message);
      Serial.print(bits);
      Serial.print("bits / ");
      //Print value in HEX
      uint8_t bytes = (bits+7)/8;
      for (int i=0; i<bytes; i++) {
          Serial.print(data[i] >> 4, 16);
          Serial.print(data[i] & 0xF, 16);
      }
      Serial.println();
#endif
   
}

// Notifies when an invalid transmission is detected
void receivedDataError(Wiegand::DataError error, uint8_t* rawData, uint8_t rawBits, const char* message) {
#ifdef DEBUG
    Serial.print(message);
    Serial.print(Wiegand::DataErrorStr(error));
    Serial.print(" - Raw data: ");
    Serial.print(rawBits);
    Serial.print("bits / ");

    //Print value in HEX
    uint8_t bytes = (rawBits+7)/8;
    for (int i=0; i<bytes; i++) {
        Serial.print(rawData[i] >> 4, 16);
        Serial.print(rawData[i] & 0xF, 16);
    }
    Serial.println();
#endif
}

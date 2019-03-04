/***************
 * wigand to serial -- program that takes wigand data and sends it over uart
 *
 * Author: Ben S. Eishen
 *
 * Usage:
 * UART is 9600N1
 * RFID badge scan will send an R folowed by the ASCII conversion of the number ending in a newline ie "R1234567\n"
 * Button presses will be sent over as they are pressed in ASCII
 * Enter and escape will send \n and \esc\n respectively.
 *
 * If an 'O' is received, it will toggle the door relay, key buzzer and led for 3 seconds and release.
 *
 * If an 'E' is received, it will rapidly flash the the led 5 times indicating there was an error.
 */

#define DOOR_RELAY 7
#define KEY_BUZZER 4
#define KEY_LED    5

#include <Wiegand.h>

WIEGAND wg;
char inbyte;

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

void setup() {
  pinMode(DOOR_RELAY, OUTPUT);
  pinMode(KEY_BUZZER, OUTPUT);
  pinMode(KEY_LED, OUTPUT);

  digitalWrite(DOOR_RELAY, HIGH);
  digitalWrite(KEY_BUZZER, HIGH);
  digitalWrite(KEY_LED,    HIGH);

	Serial.begin(9600);
	wg.begin();
}

void loop() {
  //delay(200);

  //Wiegand Check and Send
  if(wg.available()){

    //Check for rfid tag data, if there is data preface it with 'R'
    if((wg.getWiegandType()==26)||(wg.getWiegandType()==34)){
      Serial.print("R");
      Serial.print(wg.getCode());
      Serial.print("\n");

    //Keypresses are going to be 4 bit (usually).
    }else{
      switch(wg.getCode()){
        //ESC
        case 0x1B:
          Serial.write(wg.getCode());
          Serial.print("\n");
          break;
        //Enter
        case 0x0D:
          Serial.print("\n");
          break;
        //All other numbers
        default:
          Serial.print(wg.getCode());
      }
    }
  }

  if(Serial.available()){
    inbyte = Serial.read();
    if(inbyte == 'O'){
      toggleDoor();
    }
    if(inbyte == 'E'){
      showError();
    }
  }
}

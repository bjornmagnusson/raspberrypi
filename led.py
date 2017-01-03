import RPi.GPIO as GPIO
import time

pins = [4, 17, 27]

GPIO.setmode(GPIO.BCM)
for pin in pins:
	GPIO.setup(pin, GPIO.OUT)

num_blinks = input("Number of blinks?")

for i in range(num_blinks):
	print "Lighting" + str(i)
	for pin in pins:
		GPIO.output(pin,True)
		time.sleep(1)
		GPIO.output(pin,False)
		time.sleep(1)

print "Done"

GPIO.cleanup()

import RPi.GPIO as GPIO
import time

GPIO.setmode(GPIO.BOARD)
GPIO.setup(7, GPIO.OUT)

num_blinks = 0

num_blinks = input("Number of blinks?")

for i in range(count):
	print "Lighting" + str(i)
	GPIO.output(7,True)
	time.sleep(1)
	GPIO.output(7,False)
	time.sleep(1)

print "Done"

GPIO.cleanup()

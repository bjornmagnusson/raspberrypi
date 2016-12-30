import RPi.GPIO as GPIO
import time

GPIO.setmode(GPIO.BOARD)
GPIO.setup(7, GPIO.OUT)

for i in range(100):
	print "Lighting" + str(i)
	GPIO.output(7,True)
	time.sleep(1)
	GPIO.output(7,False)
	time.sleep(1)

print "Done"

GPIO.cleanup()
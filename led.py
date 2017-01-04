from gpiozero import LED
import time

leds = [LED(4), LED(17), LED(27)]

num_blinks = input("Number of blinks?")

for i in range(num_blinks):
	print "Lighting" + str(i)
	for led in leds:
		led.on()
		time.sleep(1)
		led.off()
		time.sleep(1)

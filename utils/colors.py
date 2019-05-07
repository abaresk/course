


class Color():
	def __init__(self, red=0, green=0, blue=0):
		self.red = red
		self.green = green
		self.blue = blue

	def asString(self):
		return "#" + hex(self.red)[2:] + hex(self.green)[2:] + hex(self.blue)[2:]

	def fromHex(self, string):
		self.red = int(string[1:3], 16)	
		self.green = int(string[3:5], 16)
		self.blue = int(string[5:], 16)

	def __add__(self, other):
		return Color((self.red + other.red) // 2, (self.green + other.green) // 2, (self.blue + other.blue) // 2)

	def blend(self, other, alpha):
		assert 0 <= alpha <= 1
		return Color((self.red + alpha * other.red) // 2, (self.green + alpha * other.green) // 2, (self.blue + alpha * other.blue) // 2)

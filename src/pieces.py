DIR_UP = 0
DIR_LEFT = 1
DIR_DOWN = 2
DIR_RIGHT = 3

NODE_POINT = 1
NODE_CURVE = 2
NODE_3WAY = 3
NODE_4WAY_REG = 4
NODE_4WAY_ICE = 5
NODE_4WAY_2PATH = 6

class Point():
	def __init__(self, x=None, y=None):
		self.x = x
		self.y = y

	def __add__(self, p2):
		return Point(self.x + p2.x, self.y + p2.y)

	def __sub__(self, p2):
		return Point(self.x - p2.x, self.y - p2.y)

	def __eq__(self, other):
		return self.x == other.x and self.y == other.y

	def __hash__(self):
		return hash((self.x, self.y))

	def __repr__(self):
		return "(" + str(self.x) + ", " + str(self.y) + ")"

	def distSquared(self):
		return self.x * self.x + self.y * self.y

	def nthPoint(self, direction, length):
		if direction == DIR_UP:
			return Point(self.x, self.y + length)
		elif direction == DIR_LEFT:
			return Point(self.x - length, self.y)
		elif direction == DIR_DOWN:
			return Point(self.x, self.y - length)
		else:
			return Point(self.x + length, self.y)

# Returns direction from p1 to p2
def vectorDir(p1, p2):
	diff = p2 - p1
	if diff.x == 0 and diff.y > 0:
		return DIR_UP
	elif diff.y == 0 and diff.x < 0:
		return DIR_LEFT
	elif diff.x == 0 and diff.y < 0:
		return DIR_DOWN
	elif diff.y == 0 and diff.x > 0:
		return DIR_RIGHT
	else:
		return None

def nextDir(direction, turns):
	return (direction + turns) % 4

def areLinked(node1, node2):
	return node1 in node2.nexts and node2 in node1.nexts





# class Track():
# 	def __init__(self, pos, orient):
# 		self.pos = pos,
# 		self.orient = orient 
# 		self.nexts = [None for i in range(2)]





class Node():
	def __init__(self, pos=Point(0, 0)):
		self.pos = pos
		self.nexts = [None for i in range(4)]

	def getNext(self, direction):
		return self.nexts[direction]

	def move(self, direction, length):
		self.pos = self.pos.nthPoint(direction, length)

	def allNullnodes(self):
		return [node for node in self.nexts if type(node) is NullNode]

	def validNewTrackDirections(self):
		return [vectorDir(self.pos, node.pos) for node in self.allNullnodes()]

	def __repr__(self):
		return str(id(self))[-4:] + " " + repr(self.pos)

# Invariant: NullNode should have exactly 1 link to a non-NullNode
class NullNode(Node):
	RADIUS = 0
	def __init__(self, pos):
		super(NullNode, self).__init__(pos)

	def allPoints(self):
		return [self.pos]

	def rotate(self):
		return

	def getParent(self):
		return [node for node in self.nexts if node is not None][0]

	def allPorts(self):
		return [self.pos.nthPoint(i, NullNode.RADIUS + 1) for i in range(4)]

class PointNode(Node):
	RADIUS = 0
	def __init__(self, pos):
		super(PointNode, self).__init__(pos)

	def allPoints(self):
		return [self.pos]

	def rotate(self):
		return

	def allPorts(self):
		return [self.pos.nthPoint(i, PointNode.RADIUS + 1) for i in range(4)]

class CurveNode(Node):
	RADIUS = 1
	def __init__(self, pos, orient):
		super(CurveNode, self).__init__(pos)
		self.orient = orient	# which direction is rightmost port

	def allPoints(self):
		pass

	def rotate(self):
		return

	def allPorts(self):
		return [self.pos.nthPoint(nextDir(self.orient, i), CurveNode.RADIUS + 1) for i in range(2)]

class ThreewayNode(Node):
	RADIUS = 1
	def __init__(self, pos, orient, passState, default):
		super(ThreewayNode, self).__init__(pos)
		self.orient = orient		# which direction is rightmost port
		self.passState = passState 	# can be either switched or blocked
		self.default = default 		# should be either left or right

	def allPoints(self):
		pass

	def rotate(self):
		pass

	def allPorts(self):
		return [self.pos.nthPoint(nextDir(self.orient, i), ThreewayNode.RADIUS + 1) for i in range(2)]

class FourwayNode(Node):
	RADIUS = 1
	def __init__(self, pos, passState):
		super(FourwayNode, self).__init__(pos)
		self.passState = passState

	def allPoints(self):
		pass

	def rotate(self):
		return

	def allPorts(self):
		return [self.pos.nthPoint(i, FourwayNode.RADIUS + 1) for i in range(4)]

class FourwayRegularNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(FourwayRegularNode, self).__init__(pos, passState)
		self.orient = orient

class FourwayIcyNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(FourwayIcyNode, self).__init__(pos, passState)
		self.orient = orient

class Fourway2PathsNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(Fourway2PathsNode, self).__init__(pos, passState)
		self.orient = orient

class Edge():
	def __init__(self, node1, node2):
		self.bridge = set([node1, node2])

	def has1NullNode(self):
		return len([type(n) is NullNode for n in self.bridge]) == 1

	def allNullnodes(self):
		return [type(n) is NullNode for n in self.bridge]

	def getNonNullNode(self):
		if self.has1NullNode():
			return[type(n) is not NullNode for n in self.bridge][0]
		return None

	def getDirToNullNode(self):
		if self.has1NullNode():
			nullnode = self.allNullnodes()[0]
			mainnode = self.getNonNullNode()
			return vectorDir(mainnode.pos, nullnode.pos)
		return None

	def getNext(self, direction):
		pass
			
	# TODO: this won't work
	def validNewTrackDirections(self):
		return [vectorDir(self.node1.pos, node.pos) for node in self.allNullnodes()]


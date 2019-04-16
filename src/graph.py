'''
Graph class -- stores course information

API:
	Modify graph:
		addTrackFromNode(node, direction)
		addTrackFromEdge(edge)
		addNodeFromEdge(edge, newNode)
		rotateNode(node)
		addNewNode(point)
		deleteTrackAtPoint(edge, point)
		deleteTrackAlongEdge(edge)
		deleteNode(node)
		translateRegion(piece) 	# should take node or edge (node1, node2)
	Getters:
		...
'''
from collections import defaultdict

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
	if diff.x == 0 and diff.x > 0:
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

class Node():
	def __init__(self, pos):
		self.pos = pos
		self.nexts = [None for i in range(4)]

	def getNext(self, direction):
		return self.nexts[direction]

	def move(self, direction, length):
		self.pos = self.pos.nthPoint(direction, length)

# Invariant: NullNode should have exactly 1 link to a non-NullNode
class NullNode(Node):
	RADIUS = 0
	def __init__(self, pos):
		super(NullNode, self).__init__(self, pos)

	def allPoints(self):
		return [self.pos]

	def rotate(self):
		return

	def getParent(self):
		pass

	def allPorts(self):
		return [self.pos.nthPoint(i, NullNode.RADIUS + 1) for i in range(4)]

class PointNode(Node):
	RADIUS = 0
	def __init__(self, pos):
		super(PointNode, self).__init__(self, pos)

	def allPoints(self):
		return [self.pos]

	def rotate(self):
		return

	def allPorts(self):
		return [self.pos.nthPoint(i, PointNode.RADIUS + 1) for i in range(4)]

class CurveNode(Node):
	RADIUS = 1
	def __init__(self, pos, orient):
		super(CurveNode, self).__init__(self, pos)
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
		super(ThreewayNode, self).__init__(self, pos)
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
		super(FourwayNode, self).__init__(self, pos)
		self.passState = passState

	def allPoints(self):
		pass

	def rotate(self):
		return

	def allPorts(self):
		return [self.pos.nthPoint(i, FourwayNode.RADIUS + 1) for i in range(4)]

class FourwayRegularNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(FourwayRegularNode, self).__init__(self, pos, passState)
		self.orient = orient

class FourwayIcyNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(FourwayIcyNode, self).__init__(self, pos, passState)
		self.orient = orient

class Fourway2PathsNode(FourwayNode):
	def __init__(self, pos, orient, passState):
		super(Fourway2PathsNode, self).__init__(self, pos, passState)
		self.orient = orient

class Edge():
	def __init__(self, node1, node2):
		self.bridge = set([node1, node2])

	def has1NullNode(self):
		return len([type(n) is NullNode for n in self.bridge]) == 1

	def getNullNode(self):
		if self.has1NullNode():
			return [type(n) is NullNode for n in self.bridge][0]
		return None

	def getNonNullNode(self):
		if self.has1NullNode():
			return[type(n) is not NullNode for n in self.bridge][0]
		return None

	def getDirToNullNode(self):
		if self.has1NullNode():
			nullnode = self.getNullNode()
			mainnode = self.getNonNullNode()
			return vectorDir(mainnode.pos, nullnode.pos)
		return None


class Graph():
	def __init__(self):
		self.pointmap = defaultdict(list)
		self.nodes = set()
		# TODO: set the first node

	def addTrackFromNode(self, node, direction):
		assert(node.validDirection(direction))
		assert(node.getNext(direction) is NullNode)
		self._addTrack(node.getNext(direction), direction, Edge(node, node.getNext(direction)))

	def addTrackFromEdge(self, edge):
		assert(edge.has1NullNode())
		self._addTrack(edge.getNullNode(), edge.getDirToNullNode(), edge)

	def _addTrack(self, nullnode, direction, edge):
		oldpoint, newpoint = nullnode.pos, nullnode.nthPoint(direction, 1)
		
		if self._checkCollision(nullnode, oldpoint, newpoint):
			return

		nullnode.move(direction, 1)

		if self._mergeable(nullnode):
			self._mergeNodes(nullnode)
		else:
			self.pointmap[oldpoint] = edge
			self.pointmap[newpoint] = nullnode

	def addNodeFromEdge(self, edge, nodeType):
		assert(edge.has1NullNode())

		newNode = self._makeNode(nodeType)
		point = edge.getNullNode().nthPoint(edge.getDirToNullNode(), type(newNode).RADIUS + 1)
		newNode.pos = point

		if self._addNode(newNode, point):
			# Link new node with the non-null
			self._replaceNullNode(edge, newNode)

	def rotateNode(self, node, clockwise=True):
		self._removePointsFromPointmap(node)
		node.rotate(clockwise)
		self._addPointsToPointmap(node)

	def addNewNode(self, point, nodeType):
		newNode = self._makeNode(nodeType)
		newNode.pos = point

		self._addNode(newNode, point)

	def _addNode(self, node, point):
		if self._checkCollision(node, None, point):
			return False

		self._addPointsToPointmap(node)
		self._addNullNodes(node)
		
		for port in node.allPorts():
			if self._mergeable(port):
				self._mergeNodes(port)
			else:
				self.pointmap[port.pos] = Edge(node, port)
		return True

	def _mergeable(self, node):
		'''
		Input: NullNode
		Do: merge with another NullNode (cancel each other out)
		Output: True if NullNode was able to be merged
		'''
		assert(type(node) is NullNode)
		# use NullNode's getParent method
		parent = node.getParent()
		direction = vectorDir(parent, node)

		if NullNode not in self.pointmap[node.pos.nthPoint(direction, -1)]:
			return False

		adjNode = [type(n) is NullNode for n in self.pointmap[node.pos.nthPoint(direction, -1)]][0]
		adjParent = adjNode.getParent()

		if vectorDir(adjParent,adjNode) != nextDir(direction, 2):
			return False

		# They can be merged!
		return True

	def _mergeNodes(self, node):
		assert(self._mergeable(node))

		parent = node.getParent()
		adjNode = [type(n) is NullNode for n in self.pointmap[node.pos.nthPoint(direction, -1)]][0]
		adjParent = adjNode.getParent()
		direction = vectorDir(parent, node)

		parent.nexts[direction] = adjParent
		adjParent.nexts[nextDir(direction, 2)] = parent

		# TODO: change any edge in between the two new nodes, and add any new edges
		newEdge = Edge(parent, adjParent)
		for point in self._allIntermediatePoints(parent, adjParent):
			self._removeEdgeWith(point, node)
			self._removeEdgeWith(point, adjNode)
			self.pointmap[point].append(newEdge)

		# Remove nodes from graph
		self._wipeNode(node)
		self._wipeNode(adjNode)

	def _wipeNode(self, node):
		assert(node in self.nodes)
		assert(node in self.pointmap[node.pos])
		self.pointmap[node.pos].remove(node)
		self.nodes.remove(node)

	def _checkCollision(self, node, oldpos, newpos):
		'''
		No track should be able to pass through a node's territory.
		Only NullNodes are allowed in the territory, and even they
		should only be in the boundary layer of the node
		'''
		if type(node) is NullNode:
			return self._NullNodeInvasion(oldpos, newpos)
		return any([type(value) is not NullNode for point in node.allPoints() for value in self.pointmap[point]])

	# A NullNode can overlap with one other type of node. But it cannot
	# overlap with that same node the next turn
	def _NullNodeInvasion(self, oldpos, newpos):
		if not self._nonNullNodesAtPoint(oldpos) or not self._nonNullNodesAtPoint(newpos):
			return False
		if len(self._nonNullNodesAtPoint(oldpos)) > 1 or len(self._nonNullNodesAtPoint(newpos)) > 1:
			return True
		return self._nonNullNodesAtPoint(oldpos)[0] is self._nonNullNodesAtPoint(newpos)[0]

	def _nonNullNodesAtPoint(self, point):
		return [type(n) is Node and type(n is not NullNode for n in self.pointmap[point])]

	def _makeNode(self, nodeType):
		if nodeType == NODE_POINT:
			return PointNode()
		elif nodeType = NODE_CURVE:
			return CurveNode()
		elif nodeType == NODE_3WAY:
			return ThreewayNode()
		elif nodeType == NODE_4WAY_REG:
			return FourwayRegularNode()
		elif nodeType == NODE_4WAY_ICE:
			return FourwayIcyNode()
		elif nodeType == NODE_4WAY_2PATH:
			return Fourway2PathsNode()
		return None

	def _replaceNullNode(self, edge, newNode):
		assert(edge.has1NullNode())

		# Link nodes together
		edge.getNonNullNode().nexts[direction] = newNode
		newNode.nexts[nextDir(direction, 2)] = edge.getNonNullNode()

		# Update graph
		self.nodes.remove(edge.getNullNode())
		self.nodes.add(newNode)

		# Update edge
		edge.bridge.remove(edge.getNullNode())
		edge.bridge.add(newNode)

	def _addPointsToPointmap(self, node):
		for point in node.allPoints():
			self.pointmap[point].append(node)

	def _removePointsFromPointmap(self, node):
		for point in node.allPoints():
			self.pointmap[point].remove(node)

	def _allIntermediatePoints(self, node1, node2):
		assert(areLinked(node1, node2))
		direction = vectorDir(node1, node2)

		output = []
		point = node1.pos
		while point != node2.pos:
			if node1 not in self.pointmap[point] and node2 not in self.pointmap[point]:
				output.append(Point(point.x, point.y))
			point = point.nthPoint(direction, 1)
		return output

	def _removeEdgeWith(self, point, node):
		for item in self.pointmap[point]:
			if type(item) is Edge and node in item.bridge:
				self.pointmap[point].remove(item)

	# Create a halo of NullNodes around a node
	def _addNullNodes(self, node):
		for point in node.allPorts():
			direction = vectorDir(node.pos, point)
			nullnode = NullNode(point)

			# Link nodes
			node.nexts[direction] = nullnode
			nullnode.nexts[nextDir(direction, 2)] = node

			# Add to pointmap
			self.pointmap[point].append(nullnode)


from collections import defaultdict

DIR_UP = 0
DIR_LEFT = 1
DIR_DOWN = 2
DIR_RIGHT = 3

NODE_NONSWITCH = 0
NODE_SWITCH = 1
NODE_BLOCK = 2

NODE_3WAY = 0
NODE_4WAY = 1

CURVE_NONE = 0
CURVE_LEFT = 1
CURVE_RIGHT = -1

def nextDir(direction, turns):
	return (direction + spaces) % 4

class Point():
	def __init__(self, x=None, y=None):
		self.x = x
		self.y = y

	def __add__(self, p2):
		return Point(self.x + p2.x, self.y + p2.y)

	def nthPoint(self, direction, length):
		if direction == DIR_UP:
			return Point(self.x, self.y + length)
		elif direction == DIR_LEFT:
			return Point(self.x - length, self.y)
		elif direction == DIR_DOWN:
			return Point(self.x, self.y - length)
		else:
			return Point(self.x + length, self.y)

# This is a linked list node for each piece
class Piece():
	def __init__(self, pos, grid):
		self.next = [None for i in range(4)]
		self.pos = pos

	def getNext(self, direction):
		return self.next[direction]

	def getNextNode(self, direction):
		piece = self.next[direction]
		while piece is not None:
			if type(piece) is Node:
				return piece
			piece = piece.next[direction]
		return None

class Edge(Piece):
	def __init__(self, pos, grid):
		Piece.__init__(pos, grid)

	def addCellsToGrid(self, grid):
		grid[pos].append(self)

class Node(Piece):
	def __init__(self, pos, grid, kind=NODE_NONSWITCH, nEntries=NODE_3WAY, orient=DIR_UP):
		Piece.__init__(pos, grid)
		self.kind = kind
		self.nEntries = nEntries
		self.orient = orient

	def addCellsToGrid(self, grid):
		grid[pos].append(self)
		for i in range(4):
			if not(self.nEntries == 3 and i == n.orient):
				grid[pos.nthPoint(DIR_UP + i)].append(self)

	def validNode(self, grid):
		if grid.get(self.pos):
			return False
		valid = True
		for i in range(4):
			if self.nEntries == 4 or i != self.orient:
				valid = valid and not grid.get(self.pos.nthPoint(i, 1))
		return valid

	def validOrientation(self, direction, grid):
		self.orient = nextDir(direction - 1)
		# keep rotating until there's a good one
		while self.orient != nextDir(direction, 2):
			if self.nEntries == 3 and self.kind == NODE_NONSWITCH and self.orient == direction:
				continue
			if self.validNode(grid):
				return True
			self.orient = nextDir(self.orient, 1)
		return False

	def rotate(self, grid, clockwise=True):
		if self.nEntries == 4:
			return

		turnInc = 1 if clockwise else -1

		first = self.orient
		self.orient = nextDir(self.orient, turnInc)
		# TODO: check validity
		# while self.orient != first:
		# 	if self.kind == NODE_NONSWITCH and self.orient == 
		# 	self.orient = nextDir(self.orient, turnInc)


class Graph():
	def __init__(self, pos):
		self.currPos = pos
		self.grid = defaultdict(list)
		self.currPiece = Node(self.pos, self.grid)
		self.currDirection = DIR_RIGHT
		self.addTrack(DIR_RIGHT).addUnit(DIR_RIGHT)

	def addTrack(self, direction, curve=CURVE_NONE):
		# make sure you aren't adding space b/w 2 nodes
		assert(self.currPiece.getNext(self.currDirection) is None), "Can't add piece on closed track"

		if curve:
			self.addCurve(curve)
			return

		posNext = self.currPos.nthPoint(self.currDirection, 1)

		# if there's a node there already...
		space = self.grid.get(posNext)
		if len(space) == 1 and type(space[0]) is Node:
			assert(space.pos != posNext), "Can't place track on top of node"
			# merge it with the node
			self.linkPiece(space)
		else:
			self.currPiece.next[self.currDirection] = Edge(posNext, self.grid)

		self.goToNext()

	def linkPiece(self, nextPiece):
		self.currPiece.next[self.currDirection] = nextPiece
		nextPiece.next[nextDir(self.currDirection, 2)] = self

	def goToNext(self):
		assert(self.currPiece.getNext(self.currDirection) is not None), "You can't move to an empty piece"
		self.currPiece = self.currPiece.getNext(self.currDirection)
		self.currPos = self.currPiece.pos

		# TODO: figure out what next direction would be
		if type(self.currPiece) is Edge:
			return
		self.currDirection = DIR_UP

	def addCurve(self, curve):
		assert self.curveValidityCheck(curve), "Curve cannot be placed here"
		# Create new node
		nodePos = self.currPos.nthPoint(self.currDirection, 2)
		curveNode = Node(nodePos, self.grid, NODE_NONSWITCH)

		# Add curve node's spaces to the grid
		self.grid[self.currPos.nthPoint(self.currDirection, 1)].add(curveNode)
		self.grid[nodePos.nthPoint(nextDir(self.currDirection, curve), 1)].add(curveNode)

		# Link curve node to current node
		self.currPiece.next[self.currDirection] = curveNode

		self.goToNext()

	def curveValidityCheck(self, curve):
		posNext = self.currPos.nthPoint(self.currDirection, 1)
		posNode = posNext.nthPoint(self.currDirection, 1)
		posNextNext = posNode.nthPoint(nextDir(self.currDirection, curve), 1)
		
		spaceNext, spaceNode, spaceNextNext = self.grid.get(posNext), self.grid.get(posNode), self.grid.get(posNextNext)
		return not spaceNext and not spaceNode and not spaceNextNext
		# valid2 = not spaceNextNext or (type(spaceNextNext) is Edge and spaceNextNext.nearestNode[(self.currDirection + curve)%4])
		# return valid1 and valid2

	def addNode(self, nodeKind, nEntries):
		n = Node(self.currPos.nthPoint(self.currDirection, 2), nodeKind, nEntries)

		# find a good initial orientation
		assert n.validOrientation(), "This type of node cannot be placed here"

		self.linkPiece(n)
		n.addCellsToGrid(self.grid)

		self.goToNext()
		

	def rotateNode(self, node, clockwise=True):
		pass


from src.maps import *

import struct

MAX_INT = 2**32-1

# class FileNode():
# 	def __init__(self):
# 		self.pos = None
# 		self.mainType = None	# gate, curve, 3way, 4way
# 		self.passType = None	# set, switch, block
# 		self.orient = None		# cross, topleft, topright, etc.
# 		self.trackType = None	# 1 path (no ice), 1 path (ice), 2 paths
# 		self.nextNodes = []		# next node in each direction

class FileObject():
	def __init__(self):
		self.type = None
		self.pos = None 	# position of head
		self.destination = None 	# for portal objects
		self.pairedNode = None 		# for track switches
		self.size = None 
		self.nextDir = None 
		self.pathState = None
		self.path = []		# list of nodes (should be circular)
		self.speed = None 

class SaveFileWriter():	
	def __init__(self, nodes, objects, filename):
		self.nodes = nodes
		self.objects = objects
		self.filename = filename

		self.nodeData, self.objData, self.loopData = [], [], []
		self.loops = []

		self.header = struct.pack(">ccc", b'L', b'V', b'L')

		self.anchor = self.findAnchor()

	def write(self):
		'''
		Given a list of track nodes, objects, and filename, write a binary
		file representing the course.
		'''
		with open(self.filename + ".lvl", "wb") as file:
			self.getNodes(nodes)
			self.getObjects(objects, nodes)
			self.getLoops()

			self.writeHeader(file)
			self.writeNodes(file)
			self.writeObjects(file)
			self.writeLoops(file)

	def findAnchor(self):
		for node in self.nodes:
			if node.mainType == GATE_START:
				return node.pos
		return Point(0, 0)

	def getNodes(self):
		fmt = '>2i4I4B'

		for node in self.nodes:
			# First pack base position, then direction indices, then node_type
			posx = (node.pos - self.anchor).x
			posy = (node.pos - self.anchor).y
			nextNodes = getNextNodes(node)
			subTypes = getNodeType(node)
			
			self.nodeData.append(struct.pack(fmt, \
				posx, posy, nextNodes, nextNodes[0], nextNodes[1], nextNodes[2], nextNodes[3], \
				subTypes[0], subTypes[1], subTypes[2], subTypes[3]))

		self.nodeHeader = struct.pack(">4cI", b'N', b'O', b'D', b'E', len(nodes) * struct.calcsize(fmt))

	def getNextNodes(self, node):
		return [MAX_INT if n not in self.nodes for n in node.nextNodes else self.nodes.index(n)]

	def getNodeType(self, node):
		subType1, subType2, subType3, subType4 = 0, 0, 0, 0
		if node.mainType == GATE_START or node.mainType == GATE_END:
			subType1, subType2 = 0, 1 if node.mainType == GATE_END else 0
		elif node.mainType == NODE_CURVE:
			subType1 = 1
		elif node.mainType == NODE_3WAY or node.mainType == NODE_4WAY:
			subType1 = 2 if node.mainType == NODE_3WAY else 3
			subType2 = 1 if node.passType == PASS_SWITCH else (2 if node.passType == PASS_BLOCK else 0)
			subType3 = node.orient

			if node.mainType == NODE_4WAY:
				subType4 = 2 if node.trackType == TRACKS_2 else (1 if node.trackType == TRACKS_1ICY else 0)

		return (subType1, subType2, subType3, subType4)

	def getObjects(self):
		fmt = ">B4i3BIi"

		for obj in self.objects:	
			auxPos1, auxPos2 = 0, 0
			if obj.type == OBJECT_PORTAL:
				auxPos1, auxPos2 = self.obj.destination.x, self.obj.destination.y
			elif obj.type == OBJECT_TRACKSWITCH:
				auxPos1 = node.index(obj.pairedNode) if obj.pairedNode in node else MAX_INT
			
			pathIndex = self.addLoops(obj) if obj.path else MAX_INT

			self.objData.append(struct.pack(fmt, \
				obj.type, obj.pos.x, obj.pos.y, auxPos1, auxPos2, obj.size, \
				obj.nextDir, obj.pathState, pathIndex, obj.speed))
		self.objHeader = struct.pack(">3cI", b'O', 'B', b'J', len(objects) * struct.calcsize(fmt))

	# returns index of newly inserted path loop
	def addLoops(self, obj):
		loop = []
		for node in obj.path:
			loop.append(MAX_INT if node not in self.nodes else self.nodes.index(node))
		self.loops.append(loop)
		return len(self.loops) - 1

	def getLoops(self):
		# Get size of data in loops
		loopSize = sum([1 + len(x) for x in self.loops])
		self.loopHeader = struct.pack(">4cI", b'L', b'O', b'O', b'P', loopSize)

		for loop in self.loops:
			self.loopData.append(struct.pack(">I", 1 + len(loop)))
			for index in loop:
				self.loopData.append(struct.pack(">I", index))

	def writeHeader(self, file):
		file.write(self.header)
		
	def writeNodes(self, file):
		file.write(self.nodeHeader)
		for data in self.nodeData
			file.write(data)

	def writeObjects(self, file):
		file.write(self.objHeader)
		for data in self.objData:
			file.write(data)

	def writeLoops(self, file):
		file.write(self.loopHeader)
		for data in self.loopData:
			file.write(data)
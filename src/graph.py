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
from src.pieces import *

from collections import defaultdict

class Graph():
	def __init__(self):
		self.pointmap = defaultdict(list)
		self.nodes = set()
		
	def init(self, pos=Point(0, 0)):
		assert(not self.pointmap and not self.nodes)
		self.addNewNode(pos, NODE_POINT)
		return list(self.nodes)[0]
	
	def addTrackAtPoint(self, point, orient):
		pass

	def addTrackFromNode(self, node, direction):
		assert(node.validDirection(direction))
		assert(node.getNext(direction) is NullNode)
		self._addTrack(node.getNext(direction), direction, Edge(node, node.getNext(direction)))

	def addTrackFromEdge(self, edge):
		assert(edge.has1NullNode())
		self._addTrack(edge.getNullNode(), edge.getDirToNullNode(), edge)

	def addNodeFromEdge(self, edge, nodeType):
		assert(edge.has1NullNode())

		point = edge.getNullNode().nthPoint(edge.getDirToNullNode(), type(newNode).RADIUS + 1)
		newNode = self._makeNode(nodeType, point)

		if self._addNode(newNode, point):
			# Link new node with the non-null
			self._replaceNullNode(edge, newNode)

	def rotateNode(self, node, clockwise=True):
		self._removeNodeToPointmap(node)
		node.rotate(clockwise)
		self._addNodeToPointmap(node)

	def addNewNode(self, point, nodeType):
		newNode = self._makeNode(nodeType, point)

		self._addNode(newNode, point)

	def deleteTrack(self, edge, point):
		assert(edge in self.pointmap[point])

		if edge.has1NullNode():
			self._deleteUnboundedTrack(edge)
		else:
			self._deleteBoundedTrack(edge, point)

	def _deleteUnboundedTrack(self, edge):
			nullnode = edge.getNullNode()
			newpoint = nullnode.pos.nthPoint(nextDir(edge.getDirToNullNode(), 2), 1)

			self.pointmap[nullnode.pos].remove(nullnode)
			self.pointmap[newpoint].remove(edge)
			self.pointmap[newpoint].append(nullnode)

			nullnode.pos = newpoint

	def _deleteBoundedTrack(self, edge, point):
			nullnode1, nullnode2 = NullNode(point), NullNode(point)
			edge1, edge2 = Edge(node.node1, nullnode1), Edge(node.node2, nullnode2)
			self._linkNodes(edge.node1, nullnode1)
			self._linkNodes(edge.node2, nullnode2)

			self.pointmap[point].remove(edge)
			self.pointmap[point].append(nullnode1)
			self.pointmap[point].append(nullnode2)

			self._relinkEdge(edge, edge1)
			self._relinkEdge(edge, edge2)

	def deleteNode(self, node):
		for point in node.allPoints():
			self.pointmap[point].remove(node)

		for direction in range(4):
			link = node.nexts[direction]
			if link:
				if type(link) is NullNode:
					self.pointmap[link.pos].remove(link)
				else:
					nullnode = NullNode(node.pos.nthPoint(direction, node.RADIUS))
					self._linkNodes(node.nexts[direction], nullnode)
					oldEege = self._findEdge(node.pos.nthPoint(direction, node.RADIUS + 1), node, link)
					if oldedge is not None:
						self._relinkEdge(oldedge)

	def moveRegion(self):
		pass

	def _findEdge(self, point, node1, node2):
		for Object in self.pointmap[point]:
			if type(Object) is Edge:
				if (Object.node1 == node1 and Object.node2 == node2) or (Object.node1 == node2 and Object.node2 == node1):
					return Object
		return None

	def _relinkEdge(self, oldedge, newedge):
		for point in self._allIntermediatePoints(newedge.node1, newedge.node2):
			self.pointmap[point].remove(oldedge)
			self.pointmap[point].append(newedge)

	def _linkNodes(self, node1, node2):
		direction = vectorDir(node1.pos, node2.pos)
		assert(direction is not None)

		node1.nexts[direction] = node2
		node2.nexts[nextDir(direction, 2)] = node1

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

	def _addNode(self, node, point):
		if self._checkCollision(node, None, point):
			return False

		self.nodes.add(node)
		self._addNodeToPointmap(node)
		self._addNullNodes(node)
		
		for nullnode in node.allNullnodes():
			if self._mergeable(nullnode):
				self._mergeNodes(nullnode)
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
		direction = vectorDir(parent.pos, node.pos)

		if not any([type(Object) is NullNode for Object in self.pointmap[node.pos.nthPoint(direction, -1)]]):
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
		return any([type(Object) is not NullNode for point in node.allPoints() for Object in self.pointmap[point]])

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

	def _makeNode(self, nodeType, pos):
		if nodeType == NODE_POINT:
			return PointNode(pos)
		elif nodeType == NODE_CURVE:
			return CurveNode(pos)
		elif nodeType == NODE_3WAY:
			return ThreewayNode(pos)
		elif nodeType == NODE_4WAY_REG:
			return FourwayRegularNode(pos)
		elif nodeType == NODE_4WAY_ICE:
			return FourwayIcyNode(pos)
		elif nodeType == NODE_4WAY_2PATH:
			return Fourway2PathsNode(pos)
		return None

	def _replaceNullNode(self, edge, newNode):
		assert(edge.has1NullNode())

		# Link nodes together
		self._linkNodes(edge.getNonNullNode(), newNode)

		# Update graph
		self.nodes.remove(edge.getNullNode())
		self.nodes.add(newNode)

		# Update edge
		edge.bridge.remove(edge.getNullNode())
		edge.bridge.add(newNode)

	def _addNodeToPointmap(self, node):
		for point in node.allPoints():
			self.pointmap[point].append(node)

	def _removeNodeToPointmap(self, node):
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
		for Object in self.pointmap[point]:
			if type(Object) is Edge and node in Object.bridge:
				self.pointmap[point].remove(Object)

	# Create a halo of NullNodes around a node
	def _addNullNodes(self, node):
		for point in node.allPorts():
			direction = vectorDir(node.pos, point)
			nullnode = NullNode(point)

			# Link nodes
			self._linkNodes(node, nullnode)

			# Add to pointmap
			self.pointmap[point].append(nullnode)
			self.nodes.add(nullnode)

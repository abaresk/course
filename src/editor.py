from src.graph import Graph
from utils.colors import Color
from constants.keycodes import *

import tkinter as tk

class GraphSpace():
	def __init__(self, point, piece):
		self.point = point
		self.piece = piece

	def __eq__(self, other):
		return self.point == other.point and self.piece is other.piece

class Editor():
	MODE_SELECT = 1
	MODE_DRAW = 2

	def __init__(self):
		self.graph = Graph()
		piece = self.graph.init()

		# State variables, along with mode and Object
		self.keys = []
		self.selected = GraphSpace(piece.pos, piece)

		# Init UI
		self.ui = tk.Tk()
		self.uiInit()

	def main(self):
		self.ui.mainloop()

	def uiInit(self):
		self.ui.title("Level Editor")
		self.ui.geometry("800x600")
		self.initMenu()
		self.initToolbar()
		self.initEditMenu()
		# self.initTilemap()
		# self.initScrollers()

		# Key bindings
		self.ui.bind("<KeyPress>", self.keydown)
		self.ui.bind("<KeyRelease>", self.keyup)

	def initMenu(self):
		menu = tk.Menu(self.ui)
		self.ui.config(menu=menu)
		
		filemenu = tk.Menu(menu)
		menu.add_cascade(label='File', menu=filemenu)
		filemenu.add_command(label='New File')
		filemenu.add_command(label='Open...')
		filemenu.add_command(label='Save')
		filemenu.add_command(label='Close File')
		filemenu.add_separator()
		filemenu.add_command(label='Exit', command=self.ui.quit)
		
		helpmenu = tk.Menu(menu)
		menu.add_cascade(label='Help', menu=helpmenu)
		helpmenu.add_command(label='Help')
		helpmenu.add_command(label='About')

	def initToolbar(self):
		toolbar = tk.Frame(self.ui)
		toolbar.pack(anchor='w')

		self.mode = tk.IntVar()
		selectButton = tk.Radiobutton(toolbar, variable=self.mode, text='Select', value=Editor.MODE_SELECT)
		selectButton.pack(side=tk.LEFT)
		
		drawButton = tk.Radiobutton(toolbar, variable=self.mode, text='Draw', value=Editor.MODE_DRAW)
		drawButton.pack(side=tk.LEFT)

		saveButton = tk.Button(toolbar, text='Save')
		saveButton.pack(side=tk.LEFT)
	
	def initEditMenu(self):
		label = tk.Label(self.ui, text='Edit Mode')
		label.pack(anchor='w')

		options = { 'Tracks': ['Rail', '3-way', '4-way'],
					'Enemies': ['Biter', 'Bomber', 'Bouncer', 'Christine', 'Zapper'],
					'Objects': ['Portal', 'Sludge', 'Slowdown', 'Phaser', 'Repel']}
		
		lists = DepLists(self.ui, options)
		self.editMode = lists.var_a
		self.Object = lists.var_b
		lists.pack(anchor='w')

	def initTilemap(self):
		self.tm = CanvasWindow(self.ui, cellSize=16, editor=self)
		self.tm.pack(fill=tk.BOTH, expand=tk.YES)

	def initScrollers(self):
		frame = tk.Frame(self.ui)
		frame.pack()
		
		leftButton = tk.Button(frame, text='←', command=self.tm.scrollLeft)
		leftButton.pack(side=tk.LEFT)
		
		rightButton = tk.Button(frame, text='→', command=self.tm.scrollRight)
		rightButton.pack(side=tk.LEFT)
		
		upButton = tk.Button(frame, text='↑', command=self.tm.scrollUp)
		upButton.pack(side=tk.LEFT)
		
		downButton = tk.Button(frame, text='↓', command=self.tm.scrollDown)
		downButton.pack(side=tk.LEFT)

	def keydown(self, event):
		self.keys.append(event.keysym)

	def keyup(self, event):
		self.keys.remove(event.keysym)

	# use current state to decide behavior
	def doAction(self):
		if self.mode == Editor.MODE_SELECT:
			if KEY_SPACE in self.keys:
				self.placeObject()
			elif KEY_BACKSPACE in self.keys:
				self.deleteObject()
			elif KEY_TAB in self.keys:
				self.scroll()
			elif KEY_SHIFT in self.keys:
				self.modifySelected()
			
	# handle arrows pressed (creates curves)
	def placeObject(self):
		if self.editMode == 'Tracks':
			self._placeTrack()

	def _placeTrack(self):
		if self.Object == 'Rail':
			if self._canAddRail(self.selected):
				nullnodeDirs = self.selected.piece.validNewTrackDirections()
				direction, curve = self._getRailConfig(nullnodeDirs)
				
				if direction is not None:
					if curve is not None:
						self.graph.addNodeFromPiece(self.piece, direction, NODE_CURVE, curve)
					else:
						self.graph.addTrackFromPiece(self.piece, direction)

		elif self.Object == '3-way':
			pass

	def deleteObject(self):
		pass

	def scroll(self):
		if KEY_UP in self.keys:
			self.tm.scrollUp()
		if KEY_LEFT in self.keys:
			self.tm.scrollLeft()
		if KEY_DOWN in self.keys:
			self.tm.scrollDown()
		if KEY_RIGHT in self.keys:
			self.tm.scrollRight()

	# TODO: flesh this out when rest works
	# push and pop from a stack of points
	def modifySelected(self):
		pass

	def _canAddRail(self, selected):
		if type(selected.piece) is Node and selected.piece.allNullnodes():
			return True
		elif type(selected.piece) is Edge:
			return any([(node.pos - selected.point).distSquared() == 1 \
						for node in selected.piece.allNullnodes()])	

	def _getRailConfig(self, nullnodeDirs):
		direction, curve = None, None
		if len(nullnodeDirs) == 1:
			direction = nullnodeDirs[0]
			curve = validKeyPressed(self.keys, validCurves(direction))
		else:
			direction = validKeyPressed(self.keys, nullnodeDirs)
			if direction is not None:
				curve = validKeyPressed(self.keys[self.keys.index(dir2key(direction))+1:], \
										validCurves(direction))
		return direction, curve

def validKeyPressed(keys, valid):
	for key in keys:
		if key2Dir(key) in valid:
			return key2Dir(key)
	return None

def validCurves(direction):
	return [nextDir(direction, -1), nextDir(direction, 1)]

def key2Dir(key):
	if key == KEY_UP:
		return DIR_UP
	elif key == KEY_LEFT:
		return DIR_LEFT
	elif key == KEY_DOWN:
		return DIR_DOWN
	elif key == KEY_RIGHT:
		return DIR_RIGHT
	return None

def dir2key(direction):
	if direction == DIR_UP:
		return KEY_UP
	elif direction == DIR_LEFT:
		return KEY_LEFT
	elif direction == DIR_DOWN:
		return KEY_DOWN
	elif direction == DIR_RIGHT:
		return KEY_RIGHT
	return None

class DepLists(tk.Frame):
	def __init__(self, master, options):
		tk.Frame.__init__(self, master)
		self.options = options
		self.var_a = tk.StringVar(self)
		self.var_b = tk.StringVar(self)

		self.var_a.trace('w', self.update_options)

		self.options_a = tk.OptionMenu(self, self.var_a, *self.options.keys())
		self.options_b = tk.OptionMenu(self, self.var_b, '')

		self.options_a.config(width=8)
		self.options_b.config(width=8)

		self.options_a.pack(side=tk.LEFT)
		self.options_b.pack(side=tk.LEFT)

	def update_options(self, *args):
		d = self.options[self.var_a.get()]
		self.var_b.set(d[0])

		menu = self.options_b['menu']
		menu.delete(0, 'end')

		for key in d:
			menu.add_command(label=key, command=lambda val=key: self.var_b.set(val))

class Cell():
	TRACK_COLOR_BG = Color().fromHex('#9ffca2')
	NODE_COLOR_BG = Color().fromHex('#6aabe8')
	EMPTY_COLOR_BG = Color().fromHex('#ffffff')
	SELECTED_BG = Color().fromHex("#f71b34")

	CELL_BORDER = Color().fromHex("#000000")

	def __init__(self, master, col, row, size, Objects):
		self.master = master
		self.Objects = Objects
		self.col = col
		self.row = row
		self.size = size
		self.selected = False

	def draw(self):
		""" order to the cell to draw its representation on the canvas """
		# # TODO: send this as a parameter
		# point = self.anchor + Point(self.x, -self.y)
		if self.master != None:
			if not self.Objects:
				fill = Cell.EMPTY_COLOR_BG
			elif len(self.Objects) == 1 and type(self.Objects[0]) is Node:
				fill = Cell.NODE_COLOR_BG
			else:
				fill = Cell.TRACK_COLOR_BG

			if self.selected:
				fill = fill.blend(SELECTED_BG, 0.5)

			xmin = self.x * self.size
			xmax = xmin + self.size
			ymin = self.y * self.size
			ymax = ymin + self.size

			self.master.create_rectangle(xmin, ymin, xmax, ymax, fill=fill.asString(), outline=Cell.CELL_BORDER)

class CanvasWindow(tk.Canvas):
	def __init__(self, master, cellSize, editor, *args, **kwargs):
		tk.Canvas.__init__(self, master, *args, **kwargs)
		self.cellSize = cellSize
		self.editor = editor

		self.height, self.width = super().winfo_reqheight(), super().winfo_reqwidth()

		self.bind("<Configure>", self.on_resize)
		self.bind("<Button-1>", self.handleMouseClick)

		self.updateCells()
		self.display()

	def updateCells(self):
		self.nRows, self.nCols = self.height // self.cellSize, self.width // self.cellSize
		self.anchor = Point(-self.nCols // 2, self.nRows // 2)
		self.cells = [[Cell(self, col, row, self.cellSize, self.graph.pointmap[self.anchor + Point(self.col, -self.row)]) for col in range(self.nCols)] for row in range(self.nRows)]

	# Redraw grid on resizing
	def on_resize(self, event):
		self.width = event.width
		self.height = event.height
		self.config(scrollregion=self.bbox(tk.ALL))
		self.scale("all", 0, 0, float(event.width) / self.width, float(event.height) / self.height)
		self.updateCells()

	def display(self):
		self.delete("all")
		for row in self.cells:
			for cell in row:
				cell.draw()

	def scrollUnit(self, delta):
		self.anchor += delta
		self.updateCells()

	def scrollUp(self):
		self.scrollUnit(Point(0, 1))

	def scrollLeft(self):
		self.scrollUnit(Point(-1, 0))

	def scrollDown(self):
		self.scrollUnit(Point(0, -1))

	def scrollRight(self):
		self.scrollUnit(Point(1, 0))

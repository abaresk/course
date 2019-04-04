
from src.maps import Graph
from src.maps import Point

import tkinter as tk
from collections import defaultdict

class Editor():	
	CURSOR_SELECT = 1
	CURSOR_DRAW = 2
	CURSOR_ERASE = 3

	def __init__(self):
		self.r = self.init()

	def init(self):
		r = tk.Tk()
		r.title('Level Editor')
		r.geometry("800x600")	
		self.initMenu(r)
		self.initToolbar(r)
		self.initEditMenu(r)
		self.initTilemap(r)
		self.initScrollers(r)
		return r

	def initMenu(self, r):
		menu = tk.Menu(r)
		r.config(menu=menu)
		filemenu = tk.Menu(menu)
		menu.add_cascade(label='File', menu=filemenu)
		filemenu.add_command(label='New File')
		filemenu.add_command(label='Open...')
		filemenu.add_command(label='Save')
		filemenu.add_command(label='Close File')
		filemenu.add_separator()
		filemenu.add_command(label='Exit', command=r.quit)
		helpmenu = tk.Menu(menu)
		menu.add_cascade(label='Help', menu=helpmenu)
		helpmenu.add_command(label='Help')
		helpmenu.add_command(label='About')

	def initToolbar(self, r):
		toolbar = tk.Frame(r)
		toolbar.pack(anchor='w')

		cursorState = tk.IntVar()
		selectButton = tk.Radiobutton(toolbar, variable=cursorState, text='Select', value=self.CURSOR_SELECT)
		selectButton.pack(side=tk.LEFT)
		drawButton = tk.Radiobutton(toolbar, variable=cursorState, text='Draw', value=self.CURSOR_DRAW)
		drawButton.pack(side=tk.LEFT)
		eraseButton = tk.Radiobutton(toolbar, variable=cursorState, text='Erase', value=self.CURSOR_ERASE)
		eraseButton.pack(side=tk.LEFT)
		saveButton = tk.Button(toolbar, text='Save')
		saveButton.pack(side=tk.LEFT)
		runButton = tk.Button(toolbar, text='Run!')
		runButton.pack(side=tk.LEFT)

	def initEditMenu(self, r):
		label = tk.Label(r, text='Edit Mode')
		label.pack(anchor='w')

		options = { 'Tracks': ['Single', '3-way', '4-way'],
					'Enemies': ['Biter', 'Bomber', 'Bouncer', 'Christine', 'Zapper'],
					'Objects': ['Portal', 'Sludge', 'Slowdown', 'Phaser', 'Repel']}
		
		lists = DepLists(r, options)
		self.editMode = lists.var_a
		self.editObject = lists.var_b
		lists.pack(anchor='w')

	def initTilemap(self, r):
		# tmFrame = tk.Frame(r)
		# tmFrame.grid_rowconfigure(0, weight=1)
		# tmFrame.grid_columnconfigure(0, weight=1)

		# xScroll = tk.Scrollbar(tmFrame, orient=tk.HORIZONTAL)
		# xScroll.grid(row=1, column=0, sticky="ew")

		# yScroll = tk.Scrollbar(tmFrame, orient=tk.VERTICAL)
		# yScroll.grid(row=0, column=1, sticky="ns")

		# tm = CanvasWindow(r, cellSize=16, bd=0, xscrollcommand=xScroll.set, yscrollcommand=yScroll.set)
		# tm.grid(row=0, column=0, sticky="nsew")

		# tmFrame.pack()

		self.tm = CanvasWindow(r, cellSize=16)
		self.tm.pack(fill=tk.BOTH, expand=tk.YES)

	def initScrollers(self, r):
		frame = tk.Frame(r)
		frame.pack()
		leftButton = tk.Button(frame, text='←', command=self.tm.moveLeft)
		leftButton.pack(side=tk.LEFT)
		rightButton = tk.Button(frame, text='→', command=self.tm.moveRight)
		rightButton.pack(side=tk.LEFT)
		upButton = tk.Button(frame, text='↑', command=self.tm.moveUp)
		upButton.pack(side=tk.LEFT)
		downButton = tk.Button(frame, text='↓', command=self.tm.moveDown)
		downButton.pack(side=tk.LEFT)

	def main(self):
		self.r.mainloop()

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

class CanvasWindow(tk.Canvas):
	def __init__(self, master, cellSize, *args, **kwargs):
		tk.Canvas.__init__(self, master, *args, **kwargs)
		self.cellSize = cellSize
		self.height, self.width = super().winfo_reqheight(), super().winfo_reqwidth()
		self.nRows, self.nCols = self.height // self.cellSize, self.width // self.cellSize


		self.anchor = Point(-self.nCols // 2, -self.nRows // 2)
		self.graph = Graph()

		self.cells = [[Cell(self, i, j, cellSize, self.anchor, self.graph.grid) for i in range(self.nCols)] for j in range(self.nRows)]

		self.selected = None		

		self.bind("<Configure>", self.on_resize)

		self.display()

	def updateCells(self):
		self.nRows, self.nCols = self.height // self.cellSize, self.width // self.cellSize
		self.cells = [[Cell(self, i, j, self.cellSize, self.anchor, self.graph.grid) for i in range(self.nCols)] for j in range(self.nRows)]

	# TODO: redraw grid on scrolling

	# Redraw grid on resizing
	def on_resize(self, event):
		wscale = float(event.width) / self.width
		hscale = float(event.height) / self.height
		self.width = event.width
		self.height = event.height
		self.config(scrollregion=self.bbox(tk.ALL))
		self.scale("all", 0, 0, wscale, hscale)
		self.updateCells()
		self.delete("all")
		self.display()

	def display(self):
		for row in self.cells:
			for cell in row:
				cell.draw()

	def moveUnit(self, delta):
		self.anchor += delta
		self.updateCells()
		self.delete("all")
		self.display()

	def moveLeft(self):
		self.moveUnit(Point(-1, 0))

	def moveRight(self):
		self.moveUnit(Point(1, 0))

	def moveDown(self):
		self.moveUnit(Point(0, 1))

	def moveUp(self):
		self.moveUnit(Point(0, -1))

class Cell():
	TRACK_COLOR_BG = '#9ffca2'
	EMPTY_COLOR_BG = 'white'
	TRACK_COLOR_BORDER = '#9ffca2'
	EMPTY_COLOR_BORDER = 'black'

	def __init__(self, master, x, y, size, anchor, grid):
		""" 
		Constructor of the object called by Cell(...) 
		"""
		self.master = master
		self.grid = grid
		self.x = x
		self.y = y
		self.size = size
		self.anchor = anchor
	# TODO: draw correct tile
	def draw(self):
		""" order to the cell to draw its representation on the canvas """
		gridPoint = self.anchor + Point(self.x, self.y)
		if self.master != None :
			fill = Cell.TRACK_COLOR_BG
			outline = Cell.EMPTY_COLOR_BORDER

			if not self.grid.get(gridPoint):
				fill = Cell.EMPTY_COLOR_BG
				outline = Cell.EMPTY_COLOR_BORDER

			xmin = self.x * self.size
			xmax = xmin + self.size
			ymin = self.y * self.size
			ymax = ymin + self.size

			self.master.create_rectangle(xmin, ymin, xmax, ymax, fill = fill, outline = outline)

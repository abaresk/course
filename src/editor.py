
import tkinter as tk
from collections import defaultdict

class Editor():	
	def __init__(self):
		self.r = self.init()

	def init(self):
		r = tk.Tk()
		r.title('Level Editor')
		# r.geometry("600x400")	
		self.initMenu(r)
		self.initToolbar(r)
		self.initEditMenu(r)
		self.initTilemap(r)

		# Vars:
		# self.editMode
		# self.editObject
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
		selectButton = tk.Button(toolbar, text='Select')
		selectButton.pack(side=tk.LEFT)
		drawButton = tk.Button(toolbar, text='Draw')
		drawButton.pack(side=tk.LEFT)
		eraseButton = tk.Button(toolbar, text='Erase')
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
		tm = CanvasWindow(r, width=850, height=400)
		tm.pack(fill=tk.BOTH, expand=tk.YES)


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
	def __init__(self, parent, **kwargs):
		tk.Canvas.__init__(self, parent, **kwargs)
		self.bind("<Configure>", self.on_resize)
		self.height = self.winfo_reqheight()
		self.width = self.winfo_reqheight()

	def on_resize(self, event):
		# wscale = float(event.width) / self.width
		# hscale = float(event.height) / self.height
		# self.width = event.width
		# self.height = event.height
		# self.config(width=self.width, height=self.height)
		# self.scale("all", 0, 0, wscale, hscale)
		pass

# class Cell():
# 	TRACK_COLOR_BG = '#9ffca2'
# 	EMPTY_COLOR_BG = 'white'
# 	CONNECTED_COLOR_BORDER = '#9ffca2'
# 	EMPTY_COLOR_BORDER = 'black'

# 	def __init__(self):
# 		pass




class Tilemap():
	def __init__(self, width, height):
		self.cells = [[None for i in range(width)] for j in range(height)]
		self.graph = Graph()



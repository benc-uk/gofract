
type fractWidget struct {
	canvas *canvas.Raster
	fractType string

	size     fyne.Size
	position fyne.Position
	hidden   bool
}

func (f fractWidget) MinSize() fyne.Size {
	return fyne.Size{300, 300}
}

func (f fractWidget) Visible() bool {
	return !f.hidden
}

func (f fractWidget) Show() {
	f.hidden = false
}

func (f fractWidget) Hide() {
	f.hidden = true
}

func (f fractWidget) Move(position fyne.Position) {
	f.position = position
	widget.Renderer(f).Layout(f.size)
}

func (f fractWidget) Position() fyne.Position {
	return f.position
}

func (f fractWidget) Size() fyne.Size {
	return f.size
}

func (f fractWidget) Resize(size fyne.Size) {
	fmt.Println("---- 1", size)
	f.size = size
	//widget.Renderer(g).Layout(size)
	widget.Renderer(f).Layout(f.size)
	f.canvas.Resize(size)
}

func (f *fractWidget) CreateRenderer() fyne.WidgetRenderer {

	return renderer
}

func (f fractWidget) Layout(size fyne.Size) {
	fmt.Println("----")
	//widget.Renderer(g).Layout(size)
	f.canvas.Resize(size)
}
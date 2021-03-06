package mod_author

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kiamev/moogle-mod-manager/mods"
	cw "github.com/kiamev/moogle-mod-manager/ui/custom-widgets"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"strings"
)

type downloadsDef struct {
	*entryManager
	list *cw.DynamicList
}

func newDownloadsDef() *downloadsDef {
	d := &downloadsDef{
		entryManager: newEntryManager(),
	}
	d.list = cw.NewDynamicList(cw.Callbacks{
		GetItemKey:    d.getItemKey,
		GetItemFields: d.getItemFields,
		OnEditItem:    d.onEditItem,
	})
	return d
}

func (d *downloadsDef) compile() []*mods.Download {
	downloads := make([]*mods.Download, len(d.list.Items))
	for i, item := range d.list.Items {
		downloads[i] = item.(*mods.Download)
	}
	return downloads
}

func (d *downloadsDef) getItemKey(item interface{}) string {
	return item.(*mods.Download).Name
}

func (d *downloadsDef) getItemFields(item interface{}) []string {
	return []string{
		item.(*mods.Download).Name,
		strings.Join(item.(*mods.Download).Sources, ", "),
		string(item.(*mods.Download).InstallType),
	}
}

func (d *downloadsDef) onEditItem(item interface{}) {
	d.createItem(item)
}

func (d *downloadsDef) createItem(item interface{}, done ...func(interface{})) {
	m := item.(*mods.Download)
	d.createFormItem("Name", m.Name)
	d.createFormMultiLine("Sources", strings.Join(m.Sources, "\n"))
	d.createFormSelect("Install Type", mods.InstallTypes, string(m.InstallType))

	fd := dialog.NewForm("Edit Downloadable", "Save", "Cancel", []*widget.FormItem{
		d.getFormItem("Name"),
		d.getFormItem("Sources"),
		d.getFormItem("Install Type"),
	}, func(ok bool) {
		if ok {
			m.Name = d.getString("Name")
			m.Sources = d.getStrings("Sources", "\n")
			m.InstallType = mods.InstallType(d.getString("Install Type"))
			if len(done) > 0 {
				done[0](m)
			}
		}
	}, state.Window)
	fd.Resize(fyne.NewSize(400, 400))
	fd.Show()
}

func (d *downloadsDef) draw() fyne.CanvasObject {
	return container.NewVBox(container.NewHBox(
		widget.NewLabelWithStyle("Downloadables", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewButton("Add", func() {
			d.createItem(&mods.Download{}, func(result interface{}) {
				d.list.AddItem(result)
			})
		})),
		d.list.Draw())
}

func (d *downloadsDef) set(downloadables []*mods.Download) {
	d.list.Clear()
	for _, i := range downloadables {
		d.list.AddItem(i)
	}
}

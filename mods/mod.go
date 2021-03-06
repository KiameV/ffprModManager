package mods

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/kiamev/moogle-mod-manager/config"
	"github.com/kiamev/moogle-mod-manager/ui/state"
	"os"
	"path/filepath"
	"strings"
)

type SelectType string

const (
	Auto   SelectType = "Auto"
	Select SelectType = "Select"
	Radio  SelectType = "Radio"
)

var SelectTypes = []string{string(Auto), string(Select), string(Radio)}

type Mod struct {
	ID                  string            `json:"ID" xml:"ID"`
	Name                string            `json:"Name" xml:"Name"`
	Author              string            `json:"Author" xml:"Author"`
	Version             string            `json:"Version" xml:"Version"`
	ReleaseDate         string            `json:"ReleaseDate" xml:"ReleaseDate"`
	Category            string            `json:"Category" xml:"Category"`
	Description         string            `json:"Description" xml:"Description"`
	ReleaseNotes        string            `json:"ReleaseNotes" xml:"ReleaseNotes"`
	Link                string            `json:"Link" xml:"Link"`
	ModFileLinks        []string          `json:"ModFileLink" xml:"ModFileLink"`
	Preview             *Preview          `json:"Preview,omitempty" xml:"Preview,omitempty"`
	ModCompatibility    *ModCompatibility `json:"Compatibility,omitempty" xml:"ModCompatibility,omitempty"`
	Downloadables       []*Download       `json:"Downloadable" xml:"Downloadables"`
	DonationLinks       []*DonationLink   `json:"DonationLink" xml:"DonationLinks"`
	Games               []*Game           `json:"Games" xml:"Games"`
	DownloadFiles       *DownloadFiles    `json:"DownloadFile,omitempty" xml:"DownloadFiles,omitempty"`
	Configurations      []*Configuration  `json:"Configuration,omitempty" xml:"Configurations,omitempty"`
	ConfigSelectionType SelectType        `json:"ConfigSelectionType" xml:"ConfigSelectionType"`
}

type Preview struct {
	Url   *string       `json:"Url,omitempty" xml:"Url,omitempty"`
	Local *string       `json:"Local,omitempty" xml:"Local,omitempty"`
	Size  Size          `json:"Size,omitempty" xml:"Size,omitempty"`
	img   *canvas.Image `json:"-" xml:"-"`
}

type Size struct {
	X int `json:"X" xml:"X"`
	Y int `json:"Y" xml:"Y"`
}

func (p *Preview) Get() *canvas.Image {
	if p == nil {
		return nil
	}
	if p.img == nil {
		var (
			r   fyne.Resource
			err error
		)
		if p.Local != nil {
			f := filepath.Join(state.GetBaseDir(), *p.Local)
			if _, err = os.Stat(f); err == nil {
				r, err = fyne.LoadResourceFromPath(f)
			}
		}
		if r == nil && p.Url != nil {
			r, err = fyne.LoadResourceFromURLString(*p.Url)
		}
		if r == nil || err != nil {
			return nil
		}
		p.img = canvas.NewImageFromResource(r)
		size := fyne.Size{Width: float32(p.Size.X), Height: float32(p.Size.Y)}
		p.img.SetMinSize(size)
		p.img.Resize(size)
		p.img.FillMode = canvas.ImageFillContain
	}
	return p.img
}

type ModCompatibility struct {
	Requires []*ModCompat `json:"Require" xml:"Requires"`
	Forbids  []*ModCompat `json:"Forbid" xml:"Forbids"`
	//OrderConstraints []ModCompat `json:"OrderConstraint"`
}

func (c *ModCompatibility) HasItems() bool {
	return len(c.Requires) > 0 || len(c.Forbids) > 0
}

type ModCompat struct {
	ModID    string          `json:"ModID" xml:"ModID"`
	Name     string          `json:"Name" xml:"Name"`
	Versions []string        `json:"Version,omitempty" xml:"Versions,omitempty"`
	Source   string          `json:"Source" xml:"Source"`
	Order    *ModCompatOrder `json:"Order,omitempty" xml:"Order,omitempty"`
}

type ModCompatOrder string

const (
	None   ModCompatOrder = ""
	Before ModCompatOrder = "Before"
	After  ModCompatOrder = "After"
)

var ModCompatOrders = []string{string(None), string(Before), string(After)}

type InstallType string

const (
	Bundles  InstallType = "Bundles"
	Memoria  InstallType = "Memoria"
	Magicite InstallType = "Magicite"
	BepInEx  InstallType = "BepInEx"
	// DLL Patcher https://discord.com/channels/371784427162042368/518331294858608650/863930606446182420
	//DllPatch   InstallType = "DllPatch"
	Compressed InstallType = "Compressed"
)

var InstallTypes = []string{string(Bundles), string(Memoria), string(Magicite), string(BepInEx) /*string(DllPatch),*/, string(Compressed)}

type Game struct {
	Name     config.GameName `json:"Name" xml:"Name"`
	Versions []string        `json:"Version,omitempty" xml:"GameVersions,omitempty"`
}

type Download struct {
	Name        string      `json:"Name" xml:"Name"`
	Sources     []string    `json:"Source" xml:"Sources"`
	InstallType InstallType `json:"InstallType" xml:"InstallType"`
}

type DownloadFiles struct {
	DownloadName string     `json:"DownloadName" xml:"DownloadName"`
	Files        []*ModFile `json:"File,omitempty" xml:"Files,omitempty"`
	Dirs         []*ModDir  `json:"Dir,omitempty" xml:"Dirs,omitempty"`
}

func (f *DownloadFiles) IsEmpty() bool {
	return len(f.Files) == 0 && len(f.Dirs) == 0
}

type ModFile struct {
	From string `json:"From" xml:"From"`
	To   string `json:"To" xml:"To"`
}

type ModDir struct {
	From      string `json:"From" xml:"From"`
	To        string `json:"To" xml:"To"`
	Recursive bool   `json:"Recursive" xml:"Recursive"`
}

type Configuration struct {
	Name        string    `json:"Name" xml:"Name"`
	Description string    `json:"Description" xml:"Description"`
	Preview     *Preview  `json:"Preview,omitempty" xml:"Preview, omitempty"`
	Root        bool      `json:"Root" xml:"Root"`
	Choices     []*Choice `json:"Choice" xml:"Choices"`
}

type Choice struct {
	Name                  string         `json:"Name" xml:"Name"`
	Description           string         `json:"Description" xml:"Description"`
	Preview               *Preview       `json:"Preview,omitempty" xml:"Preview,omitempty"`
	DownloadFiles         *DownloadFiles `json:"DownloadFiles,omitempty" xml:"DownloadFiles,omitempty"`
	NextConfigurationName *string        `json:"NextConfigurationName,omitempty" xml:"NextConfigurationName"`
}

type DonationLink struct {
	Name string `json:"Name" xml:"Name"`
	Link string `json:"Link" xml:"Link"`
}

func (m Mod) Validate() string {
	sb := strings.Builder{}
	if m.ID == "" {
		sb.WriteString("ID is required\n")
	}
	if m.Name == "" {
		sb.WriteString("Name is required\n")
	}
	if m.Author == "" {
		sb.WriteString("Author is required\n")
	}
	if m.ReleaseDate == "" {
		sb.WriteString("Release Date is required\n")
	}
	if m.Category == "" {
		sb.WriteString("Category is required\n")
	}
	if m.Description == "" {
		sb.WriteString("Description is required\n")
	}
	if m.Link == "" {
		sb.WriteString("Link is required\n")
	}
	if len(m.ModFileLinks) == 0 {
		sb.WriteString("ModFileLinks is required\n")
	}

	if m.Preview != nil {
		if m.Preview.Size.X <= 50 || m.Preview.Size.Y <= 50 {
			sb.WriteString("Preview size must be greater than 50\n")
		}
	}

	for _, mfl := range m.ModFileLinks {
		if strings.HasSuffix(mfl, ".json") == false && strings.HasSuffix(mfl, ".xml") == false {
			sb.WriteString(fmt.Sprintf("Mod File Link [%s] must be json or xml\n", mfl))
		}
	}

	if len(m.Downloadables) == 0 {
		sb.WriteString("Must have at least one Downloadables\n")
	}
	for _, d := range m.Downloadables {
		if d.Name == "" {
			sb.WriteString("Downloadables' name is required\n")
		}
		if len(d.Sources) == 0 {
			sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Source is required\n", d.Name))
		}
		if d.InstallType == "" {
			sb.WriteString(fmt.Sprintf("Downloadables [%s]'s Install Type is required\n", d.Name))
		}
	}

	if (m.DownloadFiles == nil || m.DownloadFiles.IsEmpty()) && len(m.Configurations) == 0 {
		sb.WriteString("One \"Always Download\", at least one \"Configuration\" or both are required\n")
	}

	if m.DownloadFiles != nil {
		if m.DownloadFiles.IsEmpty() {
			sb.WriteString(fmt.Sprintf("DownloadFiles [%s]' Must have at least one File or Dir specified\n", m.DownloadFiles.DownloadName))
		}
	}

	roots := 0
	for _, c := range m.Configurations {
		if c.Name == "" {
			sb.WriteString("Configuration's Name is required\n")
		}
		if c.Description == "" {
			sb.WriteString(fmt.Sprintf("Configuration's [%s] Description is required\n", c.Name))
		}
		if len(c.Choices) == 0 {
			sb.WriteString(fmt.Sprintf("Configuration's [%s] must have Choices\n", c.Name))
		}
		for _, ch := range c.Choices {
			if ch.Name == "" {
				sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice's Name is required\n", c.Name))
			}
			if ch.NextConfigurationName != nil && *ch.NextConfigurationName == c.Name {
				sb.WriteString(fmt.Sprintf("Configuration's [%s] Choice's Next Configuration Name must not be the same as the Configuration's Name\n", c.Name))
			}
		}
		if c.Root {
			roots++
		}
	}
	if len(m.Configurations) > 1 && roots == 0 {
		sb.WriteString("Must have at least one 'Root' Configuration\n")
	} else if roots > 1 {
		sb.WriteString("Only one 'Root' Configuration is allowed\n")
	}
	return sb.String()
}

func (m *Mod) Supports(game config.Game) error {
	gs := " " + config.String(game)
	for _, g := range m.Games {
		if strings.HasSuffix(string(g.Name), gs) {
			return nil
		}
	}
	return fmt.Errorf("%s does not support %s", m.Name, config.GameNameString(game))
}

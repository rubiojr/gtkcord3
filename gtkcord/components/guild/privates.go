package guild

import (
	"github.com/diamondburned/gtkcord3/gtkcord/gtkutils"
	"github.com/diamondburned/gtkcord3/gtkcord/ningen"
	"github.com/diamondburned/gtkcord3/gtkcord/semaphore"
	"github.com/diamondburned/gtkcord3/internal/log"
	"github.com/gotk3/gotk3/gtk"
)

type DMButton struct {
	*gtk.ListBoxRow
	Style *gtk.StyleContext

	OnClick func()

	class string
}

// thread-safe
func NewPMButton(s *ningen.State) (dm *DMButton) {
	semaphore.IdleMust(func() {
		r, _ := gtk.ListBoxRowNew()
		r.Show()
		r.SetSizeRequest(IconSize+IconPadding*2, IconSize+IconPadding*2)
		r.SetHAlign(gtk.ALIGN_FILL)
		r.SetVAlign(gtk.ALIGN_CENTER)
		r.SetTooltipMarkup("<b>Private Messages</b>")
		r.SetActivatable(true)

		s, _ := r.GetStyleContext()
		s.AddClass("dmbutton")

		i, _ := gtk.ImageNew()
		i.Show()
		i.SetSizeRequest(IconSize, IconSize)
		i.SetHAlign(gtk.ALIGN_CENTER)
		i.SetVAlign(gtk.ALIGN_CENTER)
		gtkutils.ImageSetIcon(i, "system-users-symbolic", IconSize/3*2)
		r.Add(i)

		dm = &DMButton{
			ListBoxRow: r,
			Style:      s,
		}
	})

	go func() {
		// Find and detect any unread channels:
		chs, err := s.PrivateChannels()
		if err != nil {
			log.Errorln("Failed to get private channels for DMButton:", err)
			return
		}

		for _, ch := range chs {
			rs := s.FindLastRead(ch.ID)
			if rs == nil {
				continue
			}

			// Snowflakes have timestamps, which allow us to do this:
			if ch.LastMessageID.Time().After(rs.LastMessageID.Time()) {
				dm.setUnread(true)
				break
			}
		}
	}()

	return
}

func (dm *DMButton) setUnread(unread bool) {
	var class string
	if unread {
		class = "pinged"
	}
	gtkutils.DiffClass(&dm.class, class, dm.Style)
}

package main

import (
	"fmt"

	"github.com/webview/webview"

	"github.com/ghostiam/systray"
	"github.com/ghostiam/systray/example/icon"
)

func main() {
	w := webview.New(true)
	defer w.Destroy()

	systray.Register(onReady(w))
	w.Run()
}

func onReady(w webview.WebView) func() {
	return func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Awesome App")
		systray.SetTooltip("Lantern")
		mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
		go func() {
			<-mQuitOrig.ClickedCh
			fmt.Println("Requesting quit")
			w.Terminate()
			fmt.Println("Finished quitting")
		}()

		// We can manipulate the systray in other goroutines
		go func() {
			systray.SetTemplateIcon(icon.Data, icon.Data)
			systray.SetTitle("Awesome App")
			systray.SetTooltip("Pretty awesome棒棒嗒")
			mChange := systray.AddMenuItem("Change Me", "Change Me")
			mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
			mEnabled := systray.AddMenuItem("Enabled", "Enabled")
			// Sets the icon of a menu item. Only available on Mac.
			mEnabled.SetTemplateIcon(icon.Data, icon.Data)

			systray.AddMenuItem("Ignored", "Ignored")

			subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
			subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
			subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
			subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

			mQuit := systray.AddMenuItem("退出", "Quit the whole app")

			// Sets the icon of a menu item. Only available on Mac.
			mQuit.SetIcon(icon.Data)

			systray.AddSeparator()
			mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
			shown := true
			toggle := func() {
				if shown {
					subMenuBottom.Check()
					subMenuBottom2.Hide()
					mQuitOrig.Hide()
					mEnabled.Hide()
					shown = false
				} else {
					subMenuBottom.Uncheck()
					subMenuBottom2.Show()
					mQuitOrig.Show()
					mEnabled.Show()
					shown = true
				}
			}

			for {
				select {
				case <-mChange.ClickedCh:
					mChange.SetTitle("I've Changed")
				case <-mChecked.ClickedCh:
					if mChecked.Checked() {
						mChecked.Uncheck()
						mChecked.SetTitle("Unchecked")
					} else {
						mChecked.Check()
						mChecked.SetTitle("Checked")
					}
				case <-mEnabled.ClickedCh:
					mEnabled.SetTitle("Disabled")
					mEnabled.Disable()
				case <-subMenuBottom2.ClickedCh:
					panic("panic button pressed")
				case <-subMenuBottom.ClickedCh:
					toggle()
				case <-mToggle.ClickedCh:
					toggle()
				case <-mQuit.ClickedCh:
					fmt.Println("Quit2 now...")
					w.Terminate()
				}
			}
		}()
	}
}

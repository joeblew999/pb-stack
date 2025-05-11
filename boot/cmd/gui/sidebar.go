// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 The Guigui Authors (modified for pb-stack)

package gui // Changed from 'main'

import (
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

// SidebarWidget represents the sidebar UI component.
type SidebarWidget struct {
	guigui.DefaultWidget

	sidebar        basicwidget.Sidebar
	sidebarContent sidebarContent
}

// SetModel passes the application model to the sidebar content.
func (s *SidebarWidget) SetModel(model *Model) { // Uses gui.Model
	s.sidebarContent.SetModel(model)
}

// Build constructs the sidebar widget.
func (s *SidebarWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	context.SetSize(&s.sidebarContent, context.Size(s))
	s.sidebar.SetContent(&s.sidebarContent)

	appender.AppendChildWidgetWithBounds(&s.sidebar, context.Bounds(s))

	return nil
}

type sidebarContent struct {
	guigui.DefaultWidget

	list basicwidget.TextList[string]

	model *Model // Uses gui.Model
}

func (s *sidebarContent) SetModel(model *Model) { // Uses gui.Model
	s.model = model
}

func (s *sidebarContent) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	s.list.SetStyle(basicwidget.ListStyleSidebar)

	// Define items for our main application's sidebar
	items := []basicwidget.TextListItem[string]{
		{
			Text: "Home",
			ID:   "home",
		},
		{
			Text: "Packages",
			ID:   "packages",
		},
		{
			Text: "Settings",
			ID:   "settings",
		},
	}

	s.list.SetItems(items)
	s.list.SetItemHeight(basicwidget.UnitSize(context))
	if s.model != nil {
		s.list.SelectItemByID(s.model.Mode())
		s.list.SetOnItemSelected(func(index int) {
			item, ok := s.list.ItemByIndex(index)
			if !ok {
				s.model.SetMode("home") // Default to home if item not found
				return
			}
			s.model.SetMode(item.ID)
		})
	}

	appender.AppendChildWidgetWithBounds(&s.list, context.Bounds(s))

	return nil
}

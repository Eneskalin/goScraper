package ui

import (
	"fmt"
	"webScraper/handlers"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	inputState sessionState = iota
	checkingStatus
	fetchingHTML
	takingScreenshot
	listingURLs
	saving
	done
)

type Model struct {
	textInput textinput.Model

	htmlSpinner spinner.Model
	shotSpinner spinner.Model
	urlSpinner  spinner.Model

	loading bool
	err     error
	url     string
	state   sessionState
	results []string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "https://example.com"
	ti.CharLimit = 156
	 ti.Width = 20
	ti.Focus()

	s1 := spinner.New()
	s1.Spinner = spinner.Globe
	s1.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	s2 := spinner.New()
	s2.Spinner = spinner.Dot
	s2.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	s3 := spinner.New()
	s3.Spinner = spinner.Jump
	s3.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	return Model{
		textInput:   ti,
		htmlSpinner: s1,
		shotSpinner: s2,
		urlSpinner:  s3,
		state:       inputState,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q", "esc":
            return m, tea.Quit
        case "enter":
            if m.state == inputState && m.textInput.Value() != "" {
                m.url = m.textInput.Value()
                m.state = checkingStatus
                m.results = []string{} 
                return m, handlers.CheckStatusCmd(m.url)
            }
        }

    case handlers.StatusOkMsg:

        m.state = fetchingHTML
        return m, tea.Batch(m.htmlSpinner.Tick, handlers.FetchHTMLCmd(m.url))

    case handlers.FetchSuccessMsg:
        m.results = append(m.results, "HTML içeriği başarıyla çekildi") 
        m.state = takingScreenshot
        return m, tea.Batch(m.shotSpinner.Tick, handlers.TakeScreenshotCmd(m.url))

    case handlers.ScreenshotSuccessMsg:
        m.results = append(m.results, "Ekran görüntüsü kaydedildi") 
        m.state = listingURLs
        return m, tea.Batch(m.urlSpinner.Tick, handlers.ListURLsCmd())

    case handlers.FetchUrListSuccessMsg:
        m.results = append(m.results, "URL adresleri listelendi")
        m.state = saving
        return m, handlers.SavingHandlers(m.url)

    case handlers.SaveSuccessMsg:
        m.state = done
        return m, tea.Quit

    case error:
        m.err = msg
        return m, nil
    }

    var cmd tea.Cmd
    m.htmlSpinner, cmd = m.htmlSpinner.Update(msg)
    cmds = append(cmds, cmd)
    m.shotSpinner, cmd = m.shotSpinner.Update(msg)
    cmds = append(cmds, cmd)
    m.urlSpinner, cmd = m.urlSpinner.Update(msg)
    cmds = append(cmds, cmd)
    m.textInput, cmd = m.textInput.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("\n Hata: %v\n", m.err))
	}

	if m.state == inputState {
		return fmt.Sprintf("\nURL Girin:\n%s\n\n%s",
			m.textInput.View(),
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("(esc: çıkış)"),
		)
	}

	doneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))    
	pendingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) 
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true) 

	output := fmt.Sprintf("\nTarget: %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render(m.url))

	// --- HTML Çekme ---
	if len(m.results) > 0 {
		output += fmt.Sprintf("HTML Verileri Alındı\n")
	} else if m.state == fetchingHTML {
		output += fmt.Sprintf(" %s %s...\n", m.htmlSpinner.View(), activeStyle.Render("HTML Verileri Alınıyor"))
	} else {
		output += fmt.Sprintf("   %s\n", pendingStyle.Render("->HTML Verileri "))
	}

	// ---  Ekran Görüntüsü ---
	if len(m.results) > 1 {
		output += fmt.Sprintf("Ekran Görüntüsü Oluşturuldu\n")
	} else if m.state == takingScreenshot {
		output += fmt.Sprintf(" %s %s...\n", m.shotSpinner.View(), activeStyle.Render("Ekran Görüntüsü Alınıyor"))
	} else {
		output += fmt.Sprintf("   %s\n", pendingStyle.Render("-> Ekran Görüntüsü "))
	}

	// ---  Link Listeleme ---
	if len(m.results) > 2 {
		output += fmt.Sprintf(" Linkler Listelendi\n")
	} else if m.state == listingURLs {
		output += fmt.Sprintf(" %s %s...\n", m.urlSpinner.View(), activeStyle.Render("Linkler Ayıklanıyor"))
	} else {
		output += fmt.Sprintf("   %s\n", pendingStyle.Render("-> Linkler "))
	}


	if m.state == done {
		output += doneStyle.Bold(true).Render("\n İşlem başarıyla tamamlandı!\n")
	}

	return output
}


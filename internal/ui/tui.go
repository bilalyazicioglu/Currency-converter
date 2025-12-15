package ui

import (
	"fmt"
	"strconv"
	"strings"

	"Currency-Converter/internal/api"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#01eeffff")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#19efe8ff"))

	resultStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			Padding(1, 0)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000")).
			Padding(1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
)

type State int

const (
	StateSelectFrom State = iota
	StateSelectTo
	StateEnterAmount
	StateShowResult
	StateError
)

type Model struct {
	state        State
	currencies   []string
	fromCursor   int
	toCursor     int
	fromCurrency string
	toCurrency   string
	amountInput  textinput.Model
	result       *api.ConversionResult
	err          error
	client       *api.Client
	width        int
	height       int
	filter       string
	filteredList []string
}

func NewModel(client *api.Client) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter amount"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	currencies := api.CommonCurrencies()

	return Model{
		state:        StateSelectFrom,
		currencies:   currencies,
		fromCursor:   0,
		toCursor:     0,
		amountInput:  ti,
		client:       client,
		filteredList: currencies,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type ConversionResultMsg struct {
	Result *api.ConversionResult
	Err    error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case ConversionResultMsg:
		if msg.Err != nil {
			m.state = StateError
			m.err = msg.Err
		} else {
			m.state = StateShowResult
			m.result = msg.Result
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.state == StateShowResult || m.state == StateError {
				m.state = StateSelectFrom
				m.fromCursor = 0
				m.toCursor = 0
				m.filter = ""
				m.filteredList = m.currencies
				return m, nil
			}
			if m.state == StateEnterAmount {
				m.state = StateSelectTo
				return m, nil
			}
			if m.state == StateSelectTo {
				m.state = StateSelectFrom
				return m, nil
			}

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.state == StateSelectFrom {
				if m.fromCursor > 0 {
					m.fromCursor--
				}
			} else if m.state == StateSelectTo {
				if m.toCursor > 0 {
					m.toCursor--
				}
			}

		case "down", "j":
			if m.state == StateSelectFrom {
				if m.fromCursor < len(m.filteredList)-1 {
					m.fromCursor++
				}
			} else if m.state == StateSelectTo {
				if m.toCursor < len(m.filteredList)-1 {
					m.toCursor++
				}
			}

		case "backspace":
			if m.state == StateSelectFrom || m.state == StateSelectTo {
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					m.updateFilteredList()
				}
			}

		default:
			if m.state == StateSelectFrom || m.state == StateSelectTo {
				if len(msg.String()) == 1 {
					char := strings.ToUpper(msg.String())
					if char >= "A" && char <= "Z" {
						m.filter += char
						m.updateFilteredList()
					}
				}
			}
		}
	}

	if m.state == StateEnterAmount {
		var cmd tea.Cmd
		m.amountInput, cmd = m.amountInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) updateFilteredList() {
	if m.filter == "" {
		m.filteredList = m.currencies
	} else {
		filtered := make([]string, 0)
		for _, c := range m.currencies {
			if strings.HasPrefix(c, m.filter) {
				filtered = append(filtered, c)
			}
		}
		if len(filtered) > 0 {
			m.filteredList = filtered
		}
	}
	if m.state == StateSelectFrom {
		m.fromCursor = 0
	} else {
		m.toCursor = 0
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateSelectFrom:
		if len(m.filteredList) > 0 {
			m.fromCurrency = m.filteredList[m.fromCursor]
			m.state = StateSelectTo
			m.filter = ""
			m.filteredList = m.currencies
		}

	case StateSelectTo:
		if len(m.filteredList) > 0 {
			m.toCurrency = m.filteredList[m.toCursor]
			m.state = StateEnterAmount
			m.filter = ""
			m.filteredList = m.currencies
			m.amountInput.SetValue("")
			m.amountInput.Focus()
		}

	case StateEnterAmount:
		amountStr := m.amountInput.Value()
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			m.err = fmt.Errorf("invalid amount: %s", amountStr)
			m.state = StateError
			return m, nil
		}

		return m, m.doConversion(amount)

	case StateShowResult, StateError:
		m.state = StateSelectFrom
		m.fromCursor = 0
		m.toCursor = 0
	}

	return m, nil
}

func (m Model) doConversion(amount float64) tea.Cmd {
	return func() tea.Msg {
		result, err := m.client.Convert(m.fromCurrency, m.toCurrency, amount)
		return ConversionResultMsg{Result: result, Err: err}
	}
}

func (m Model) View() string {
	var s strings.Builder

	title := titleStyle.Render("Currency Converter")
	s.WriteString(title + "\n\n")

	switch m.state {
	case StateSelectFrom:
		s.WriteString(m.renderCurrencySelector("Select source currency:", m.filteredList, m.fromCursor))

	case StateSelectTo:
		s.WriteString(fmt.Sprintf("From: %s %s\n\n", selectedStyle.Render(m.fromCurrency), getCurrencyName(m.fromCurrency)))
		s.WriteString(m.renderCurrencySelector("Select target currency:", m.filteredList, m.toCursor))

	case StateEnterAmount:
		s.WriteString(fmt.Sprintf("From: %s %s\n", selectedStyle.Render(m.fromCurrency), getCurrencyName(m.fromCurrency)))
		s.WriteString(fmt.Sprintf("To:   %s %s\n\n", selectedStyle.Render(m.toCurrency), getCurrencyName(m.toCurrency)))
		s.WriteString("Enter amount:\n")
		s.WriteString(m.amountInput.View())

	case StateShowResult:
		s.WriteString(m.renderResult())

	case StateError:
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		s.WriteString("\n\nPress Enter or Esc to try again")
	}

	s.WriteString("\n")
	s.WriteString(m.renderHelp())

	return boxStyle.Render(s.String())
}

func (m Model) renderCurrencySelector(title string, currencies []string, cursor int) string {
	var s strings.Builder
	s.WriteString(title + "\n")

	if m.filter != "" {
		s.WriteString(fmt.Sprintf("Filter: %s\n", selectedStyle.Render(m.filter)))
	}
	s.WriteString("\n")

	start := 0
	end := len(currencies)
	maxVisible := 8

	if len(currencies) > maxVisible {
		start = cursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(currencies) {
			end = len(currencies)
			start = end - maxVisible
		}
	}

	for i := start; i < end; i++ {
		currency := currencies[i]
		name := getCurrencyName(currency)

		if i == cursor {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s - %s", currency, name)))
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("  %s - %s", currency, name)))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func (m Model) renderResult() string {
	var s strings.Builder

	r := m.result
	s.WriteString(resultStyle.Render("Conversion Result"))
	s.WriteString("\n\n")

	s.WriteString(fmt.Sprintf("  %s %s\n", selectedStyle.Render(fmt.Sprintf("%.2f %s", r.Amount, r.FromCurrency)), getCurrencyName(r.FromCurrency)))
	s.WriteString("         ↓\n")
	s.WriteString(fmt.Sprintf("  %s %s\n", selectedStyle.Render(fmt.Sprintf("%.2f %s", r.Result, r.ToCurrency)), getCurrencyName(r.ToCurrency)))
	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("  Exchange Rate: 1 %s = %.4f %s\n", r.FromCurrency, r.Rate, r.ToCurrency))
	s.WriteString(fmt.Sprintf("  Last Updated: %s\n", r.LastUpdated.Format("2006-01-02 15:04:05")))
	s.WriteString("\n")
	s.WriteString("Press Enter or Esc for new conversion")

	return s.String()
}

func (m Model) renderHelp() string {
	var help string
	switch m.state {
	case StateSelectFrom, StateSelectTo:
		help = "↑/↓: Navigate • Enter: Select • Type: Filter • Esc: Back • q: Quit"
	case StateEnterAmount:
		help = "Enter: Convert • Esc: Back • q: Quit"
	case StateShowResult, StateError:
		help = "Enter/Esc: New conversion • q: Quit"
	}
	return helpStyle.Render(help)
}

func getCurrencyName(code string) string {
	if name, ok := api.CurrencyNames[code]; ok {
		return name
	}
	return ""
}

func Run(client *api.Client) error {
	model := NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

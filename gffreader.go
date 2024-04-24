package main

import (
    "bufio"
    "os"
    "strings"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

func main() {
    app := tview.NewApplication()
    table := tview.NewTable().SetFixed(1, 1)

	filePath := os.Args[1]
    file, _ := os.Open(filePath)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    row := 0


    headers := []string{"seqname", "source", "feature", "start", "end", "score", "strand", "frame", "attribute"}
    for index, header := range headers {
        table.SetCell(0, index, tview.NewTableCell(header).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed))
    }
    row++

    lines := []string{}

for scanner.Scan() {
    line := scanner.Text()

    if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "##") {
        continue
    }

    lines = append(lines, line)
    fields := strings.Split(line, "\t")

    // Determine color
    var color tcell.Color
    if len(fields) > 2 && fields[2] == "CDS" {
      color = tcell.ColorBlue
    } else {
      color = tcell.ColorWhite
    }

    // Apply color
    for col, field := range fields {
        table.SetCell(row, col, tview.NewTableCell(field).SetTextColor(color))
    }

    row++
}
// search
    searchModal := tview.NewInputField().SetLabel("Search: ")
    searchModal.SetDoneFunc(func(key tcell.Key) {
       if key == tcell.KeyEnter {
           searchTerm := searchModal.GetText()

           // clear the table
           for i := 1; i < row; i++ {
               for j := 0; j < len(headers); j++ {
                   table.GetCell(i, j).SetText("")
               }
           }

           // refill the table with the search result
           row = 1
           for _, line := range lines {
               if strings.Contains(line, searchTerm) {
                   fields := strings.Split(line, "\t")
                   for col, field := range fields {
                       table.SetCell(row, col, tview.NewTableCell(field))
                    }
                   row++
               }
           }
           app.SetRoot(table, true)
        }
    })

    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlF {
            // show search modal
            app.SetRoot(searchModal, true)
        }
        return event
    })

    app.SetRoot(table, true).Run()
}

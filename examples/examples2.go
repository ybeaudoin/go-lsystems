/*==============================================================================================================================
 * Purpose: Demonstrates the use of the "lsystems" package for 2 examples of stochastic context-free rules using different rule
 *          weights. Available plot formats are GIF, JPEG, PNG, SVG and HP-GL/2.
 * Remarks: See http://algorithmicbotany.org/papers/abop/abop.pdf, Section 1.7
 *============================================================================================================================*/
package main

import (
    "bufio"
    "fmt"
    "github.com/ybeaudoin/go-lsystems"
    "os"
    "strconv"
    "strings"
)
func main() {
    const( plotWidth    = 1800
           plotHeight   = 900
           numRuns      = 6
           lineColor    = "forest-green"
    )
    var(   bgColor      = lsystems.EncodeBgColorName("gray90")
           rules        = []string{ "F[+F]F[-F]F", "F[+F]F", "F[-F]F" }
           weights      = [][]int{ []int{1, 1, 1}, []int{4, 1, 2} }

           labels       = make([]string, numRuns)
           turtleCmds   = make([]string, numRuns)
           turtleAngles = make([]float64, numRuns)
    )
    //Select an output format
    format    := 0
    stdin     := bufio.NewReader(os.Stdin)
    terminals := []string{"gif", "jpeg", "png", "svg","hp-gl/2"}
    fmt.Println("\n>>> Enter a number to select the corresponding output format <<<")
    for k,v := range terminals {
        fmt.Printf("(%d) %s", k+1, strings.ToUpper(v))
        if k < len(terminals) - 1 { fmt.Print(" | ") } else { fmt.Print(" : ") }
    }
    input, err1  := stdin.ReadString('\n')
    format, err2 := strconv.Atoi(strings.Trim(input, " \r\n"))
    for err1 != nil || err2 != nil || ! (format > 0 && format <= len(terminals)) {
        fmt.Println("\aTry again.")
        input, err1  = stdin.ReadString('\n')
        format, err2 = strconv.Atoi(strings.Trim(input, " \r\n"))
    }
    terminal := terminals[format-1]
    //Render examples
    outFile := ""
    for k, v := range weights {
        fmt.Printf("\nEXAMPLE %d / %d:\n", k+1, len(weights))
        plotTitle := fmt.Sprintf("Stochastic(weights %d,%d,%d)", v[0], v[1], v[2])
        //Generate turtle commands
        for run := 0; run < numRuns; run++ {
            lsystems.Stochastic(5, "$F", rules, v)
            turtleCmds[run]   = lsystems.TurtleCmds
            turtleAngles[run] = 25.7
        }
        //Render
        switch terminal {
            case "hp-gl/2":
                outFile = plotTitle + ".plt"
                lsystems.HpglMultiPlot(turtleCmds, turtleAngles, plotTitle, labels, 0.35, outFile)
            case "svg":
                outFile      = plotTitle + ".svg"
                terminalCmd := fmt.Sprintf("set terminal svg lw %d size %d,%d", 1, plotWidth, plotHeight)
                outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
                lsystems.MultiPlot(turtleCmds, turtleAngles, terminalCmd, outputCmd, plotTitle, labels, lineColor)
            default:
                outFile      = plotTitle + "." + terminal
                terminalCmd := fmt.Sprintf("set terminal %s lw %d size %d,%d %s",
                                           terminal, 1, plotWidth, plotHeight, bgColor)
                outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
                lsystems.MultiPlot(turtleCmds, turtleAngles, terminalCmd, outputCmd, plotTitle, labels, lineColor)
        }
        fmt.Println("output written to", outFile)
    }
}

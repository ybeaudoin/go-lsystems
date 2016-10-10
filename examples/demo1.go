/*Dragon curve of order 10*/
package main
import (
    "fmt"
    "github.com/ybeaudoin/go-lsystems"
    "strings"
)
func main() {
    order       := 10
    angle       := 90.
    axiom       := "$FX"
    rules       := strings.NewReplacer("X", "X-YF-", "Y", "+FX+Y")
    /*Plot parameters*/
    penWidth    := 0.35 //mm
    plotTitle   := "demo1-Dragon(10)"
    plotWidth   := 300 //pixels
    plotHeight  := 300 //pixels
    lineWidth   := 1
    lineColor   := "red"
    bgColor     := lsystems.EncodeBgColorName("gray70") //i.e., "xb3b3b3"
    pngFile     := plotTitle + ".png"
    pltFile     := plotTitle + ".plt"
    terminalCmd := fmt.Sprintf("set terminal png lw %d size %d,%d %s",
                               lineWidth, plotWidth, plotHeight, bgColor)
    outputCmd   := fmt.Sprintf(`set output "%s"`, pngFile)
    /*Generate turtle commands*/
    lsystems.Deterministic(order, axiom, rules)
    //fmt.Println(lsystems.TurtleCmds) //to see the turtle commands
    /*Output PNG file*/
    lsystems.Plot(angle, terminalCmd, outputCmd, plotTitle, lineColor)
    //lsystems.Plot(angle, terminalCmd, outputCmd, plotTitle, lineColor, "gnuplot.cmds") //to capture the gnuplot commands
    fmt.Println("output written to " + pngFile)
    /*Output HP-GL/2 file*/
    lsystems.HpglPlot(angle, plotTitle, penWidth, pltFile)
    fmt.Println("output written to " + pltFile)
}

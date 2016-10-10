/*Stochastic context-free rules: see http://algorithmicbotany.org/papers/abop/abop.pdf, Section 1.7*/
package main
import (
    "fmt"
    "github.com/ybeaudoin/go-lsystems"
)
func main() {
    order       := 5
    angle       := 25.7
    axiom       := "$F"
    rules       := []string{ "F[+F]F[-F]F", "F[+F]F", "F[-F]F" }
    weights     := []int{ 3, 4, 4 } //probability of a rule being chosen: 100%*3/11, 100%*4/11 & 100%*4/11 respectively
    numPlants   := 6
    /*Plot parameters*/
    penWidth    := 0.35 //mm
    plotTitle   := "demo2-Stochastic"
    plotWidth   := 600 //pixels
    plotHeight  := 300 //pixels
    lineWidth   := 1
    lineColor   := "forest-green"
    bgColor     := lsystems.EncodeBgColorName("gray90")
    pngFile     := plotTitle + ".png"
    pltFile     := plotTitle + ".plt"
    terminalCmd := fmt.Sprintf("set terminal png lw %d size %d,%d %s",
                               lineWidth, plotWidth, plotHeight, bgColor)
    outputCmd   := fmt.Sprintf(`set output "%s"`, pngFile)
    /*Generate & save the turtle commands*/
    labels       := make([]string, numPlants)
    turtleCmds   := make([]string, numPlants)
    turtleAngles := make([]float64, numPlants)
    for run := 0; run < numPlants; run++ {
        lsystems.Stochastic(order, axiom, rules, weights)
        turtleCmds[run]   = lsystems.TurtleCmds
        turtleAngles[run] = angle
    }
    /*Output PNG file*/
    lsystems.MultiPlot(turtleCmds, turtleAngles, terminalCmd, outputCmd, plotTitle, labels, lineColor)
    fmt.Println("output written to " + pngFile)
    /*Output HP-GL/2 file*/
    lsystems.HpglMultiPlot(turtleCmds, turtleAngles, plotTitle, labels, penWidth, pltFile)
    fmt.Println("output written to " + pltFile)
}

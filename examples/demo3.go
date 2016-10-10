/*Hogeweg and Hesper context rules: see http://algorithmicbotany.org/papers/abop/abop.pdf*/
package main
import (
    "fmt"
    "github.com/ybeaudoin/go-lsystems"
)
func main() {
    order       := 30
    angle       := 15.
    axiom       := "$F1F1F1"
    rules       := map[string]string{
                    "0 < 0 > 0" : "1",
                    "0 < 0 > 1" : "1[-F1F1]",
                    "0 < 1 > 0" : "1",
                    "0 < 1 > 1" : "1",
                    "1 < 0 > 0" : "0",
                    "1 < 0 > 1" : "1F1",
                    "1 < 1 > 0" : "1",
                    "1 < 1 > 1" : "0" }
    /*Plot parameters*/
    plotTitle   := "demo3-Hogeweg&Hesper"
    plotWidth   := 300 //pixels
    plotHeight  := 300 //pixels
    lineWidth   := 1
    lineColor   := "forest-green"
    bgColor     := lsystems.EncodeBgColorName("gray90")
    pngFile     := plotTitle + ".png"
    terminalCmd := fmt.Sprintf("set terminal png lw %d size %d,%d %s",
                               lineWidth, plotWidth, plotHeight, bgColor)
    outputCmd   := fmt.Sprintf(`set output "%s"`, pngFile)
    /*Generate turtle commands*/
    lsystems.HogewegHesper(order, axiom, rules)
    /*Output PNG file*/
    lsystems.Plot(angle, terminalCmd, outputCmd, plotTitle, lineColor)
    fmt.Println("output written to " + pngFile)
}

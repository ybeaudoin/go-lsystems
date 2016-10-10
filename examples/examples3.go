/*==============================================================================================================================
 * Purpose: Demonstrates the use of the "lsystems" package for 5 Hogeweg and Hesper context rules. Available plot formats are
 *          GIF, JPEG, PNG, SVG and HP-GL/2.
 * Remarks: See http://algorithmicbotany.org/papers/abop/abop.pdf, Fig 1.31. Note here that each image was produced/published
 *          using a diferent vertical scale, making them appear to have the same height.
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
           lineColor    = "forest-green"
           plotTitle    = "Hogeweg and Hesper(ABOP Fig 1.31)"
           basename     = "Hogeweg&Hesper"
    )
    var(   bgColor      = lsystems.EncodeBgColorName("gray90")
           labels       []string
           turtleCmds   []string
           turtleAngles []float64

           examples     = []map[string]interface{} {
                           map[string]interface{} {
                               "LABEL" : "a",
                               "ORDER" : 30,
                               "ANGLE" : 22.5,
                               "AXIOM" : "$F1F1F1",
                               "RULES" : map[string]string{
                                           "0 < 0 > 0" : "0",
                                           "0 < 0 > 1" : "1[+F1F1]",
                                           "0 < 1 > 0" : "1",
                                           "0 < 1 > 1" : "1",
                                           "1 < 0 > 0" : "0",
                                           "1 < 0 > 1" : "1F1",
                                           "1 < 1 > 0" : "0",
                                           "1 < 1 > 1" : "0"} },
                           map[string]interface{} {
                               "LABEL" : "b",
                               "ORDER" : 30,
                               "ANGLE" : 22.5,
                               "AXIOM" : "$F1F1F1",
                               "RULES" : map[string]string{
                                           "0 < 0 > 0" : "1",
                                           "0 < 0 > 1" : "1[-F1F1]",
                                           "0 < 1 > 0" : "1",
                                           "0 < 1 > 1" : "1",
                                           "1 < 0 > 0" : "0",
                                           "1 < 0 > 1" : "1F1",
                                           "1 < 1 > 0" : "1",
                                           "1 < 1 > 1" : "0"} },
                           map[string]interface{} { //some of the top side branches differ slightly with the published
                                                    //rendering; two extra twigs on the right versus two on the left
                               "LABEL" : "c",
                               "ORDER" : 26,
                               "ANGLE" : 25.75,
                               "AXIOM" : "$F1F1F1",
                               "RULES" : map[string]string{
                                           "0 < 0 > 0" : "0",
                                           "0 < 0 > 1" : "1",
                                           "0 < 1 > 0" : "0",
                                           "0 < 1 > 1" : "1[+F1F1]",
                                           "1 < 0 > 0" : "0",
                                           "1 < 0 > 1" : "1F1",
                                           "1 < 1 > 0" : "0",
                                           "1 < 1 > 1" : "0"} },
                           map[string]interface{} {
                               "LABEL" : "d",
                               "ORDER" : 24,
                               "ANGLE" : 25.75,
                               "AXIOM" : "$F0F1F1",
                               "RULES" : map[string]string{
                                           "0 < 0 > 0" : "1",
                                           "0 < 0 > 1" : "0",
                                           "0 < 1 > 0" : "0",
                                           "0 < 1 > 1" : "1F1",
                                           "1 < 0 > 0" : "1",
                                           "1 < 0 > 1" : "1[+F1F1]",
                                           "1 < 1 > 0" : "1",
                                           "1 < 1 > 1" : "0"} },
                           map[string]interface{} { //order needs to be 30 and not 26 as stated to match published image
                               "LABEL" : "e",
                               "ORDER" : 30,
                               "ANGLE" : 22.5,
                               "AXIOM" : "$F1F1F1",
                               "RULES" : map[string]string{
                                           "0 < 0 > 0" : "0",
                                           "0 < 0 > 1" : "1[-F1F1]",
                                           "0 < 1 > 0" : "1",
                                           "0 < 1 > 1" : "1",
                                           "1 < 0 > 0" : "0",
                                           "1 < 0 > 1" : "1F1",
                                           "1 < 1 > 0" : "1",
                                           "1 < 1 > 1" : "0"} },
                          }
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
    //Generate turtle commands
    for k, v := range examples {
        fmt.Printf("EXAMPLE %d / %d:\n", k+1, len(examples))
        lsystems.HogewegHesper(v["ORDER"].(int), v["AXIOM"].(string), v["RULES"].(map[string]string))
        labels       = append(labels, v["LABEL"].(string))
        turtleCmds   = append(turtleCmds, lsystems.TurtleCmds)
        turtleAngles = append(turtleAngles, v["ANGLE"].(float64))
    }
    //Render
    outFile := ""
    switch terminal {
        case "hp-gl/2":
            outFile = basename + ".plt"
            lsystems.HpglMultiPlot(turtleCmds, turtleAngles, plotTitle, labels, 0.35, outFile)
        case "svg":
            outFile      = basename + ".svg"
            terminalCmd := fmt.Sprintf("set terminal svg lw %d size %d,%d", 1, plotWidth, plotHeight)
            outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
            lsystems.MultiPlot(turtleCmds, turtleAngles, terminalCmd, outputCmd, plotTitle, labels, lineColor)
        default:
            outFile      = basename + "." + terminal
            terminalCmd := fmt.Sprintf("set terminal %s lw %d size %d,%d %s",
                                       terminal, 1, plotWidth, plotHeight, bgColor)
            outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
            lsystems.MultiPlot(turtleCmds, turtleAngles, terminalCmd, outputCmd, plotTitle, labels, lineColor)
    }
    fmt.Println("output written to", outFile)
}

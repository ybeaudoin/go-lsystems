/*==============================================================================================================================
 * Purpose: Demonstrates the use of the "lsystems" package with 30 deterministic context-free examples. Available plot formats
 *          are GIF, JPEG, PNG, SVG and HP-GL/2.
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
    const( plotWidth  = 900
           plotHeight = 900
    )
    //Deterministic context-free examples
    examples :=map[string]map[string]interface{} {

      //Some classics
      "Dragon" : map[string]interface{} { //see https://en.wikipedia.org/wiki/L-system
        "ORDER"      : 10,
        "ANGLE"      : 90.,
        "AXIOM"      : "$FX",
        "RULES"      : strings.NewReplacer("X", "X-YF-",
                                           "Y", "+FX+Y"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "red"},
      "Hilbert" : map[string]interface{} { //see https://en.wikipedia.org/wiki/Hilbert_curve
        "ORDER"      : 6,
        "ANGLE"      : 90.,
        "AXIOM"      : "A",
        "RULES"      : strings.NewReplacer("A", "-BF+AFA+FB-",
                                           "B", "+AF-BFB-FA+"),
        "LINE_WIDTH" : 4,
        "LINE_COLOR" : "dark-khaki"},
      "Koch Snowflake" : map[string]interface{} { //see http://mathforum.org/advanced/robertd/lsys2d.html
        "ORDER"      : 4,
        "ANGLE"      : 60.,
        "AXIOM"      : "F++F++F",
        "RULES"      : strings.NewReplacer("F", "F-F++F-F"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "white"},
      "Koch Snowflake II" : map[string]interface{} { //order 1 -> outlines the logo of Mitsubishi Motors
        "ORDER"      : 4,
        "ANGLE"      : 60.,
        "AXIOM"      : "F++F++F",
        "RULES"      : strings.NewReplacer("F", "F+F--F+F"), //just opposite turn directions to the Koch Snowflake
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "white"},
      "Moore" : map[string]interface{} { //variant of the Hilbert curve: see https://en.wikipedia.org/wiki/Moore_curve
        "ORDER"      : 5,
        "ANGLE"      : 90.,
        "AXIOM"      : "$LFL+F+LFL",
        "RULES"      : strings.NewReplacer("L", "-RF+LFL+FR-",
                                           "R", "+LF-RFR-FL+"),
        "LINE_WIDTH" : 4,
        "LINE_COLOR" : "dark-khaki"},
      "Peano" : map[string]interface{} { //original published curve; see https://en.wikipedia.org/wiki/Space-filling_curve
                                         //http://mathforum.org/advanced/robertd/lsys2d.html
        "ORDER"      : 4,
        "ANGLE"      : 90.,
        "AXIOM"      : "$X",
        "RULES"      : strings.NewReplacer("X", "XFYFX-F-YFXFY+F+XFYFX",
                                           "Y", "YFXFY+F+XFYFX-F-YFXFY"),
        "LINE_WIDTH" : 3,
        "LINE_COLOR" : "steelblue"},
      "Penrose P3" : map[string]interface{} { //see http://www.cs.brandeis.edu/~storer/JimPuzzles/PACK/CzechFarms/PenroseTilingWikipedia.pdf
        "ORDER"      : 5,
        "ANGLE"      : 36.,
        "AXIOM"      : "(18)[X]++[X]++[X]++[X]++[X]",
        "RULES"      : strings.NewReplacer("W", "YF++ZF----XF[-YF----WF]++",
                                           "X", "+YF--ZF[---WF--XF]+",
                                           "Y", "-WF++XF[+++YF++ZF]-",
                                           "Z", "--YF++++WF[+ZF++++XF]--XF",
                                           "F", ""),
        "LINE_WIDTH" : 3,
        "LINE_COLOR" : "black",
        "BGCOLOR"    : lsystems.EncodeBgColorName("dark-goldenrod")},
      "Quadratic Koch Island" : map[string]interface{} { //see http://mathforum.org/advanced/robertd/lsys2d.html
        "ORDER"      : 3,
        "ANGLE"      : 90.,
        "AXIOM"      : "F+F+F+F",
        "RULES"      : strings.NewReplacer("F", "F-F+F+FFF-F-F+F"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "royalblue"},
      "Sierpinski" : map[string]interface{} { //see http://cas.bethel.edu/faculty/projects/gossett/sierpinski/
        "ORDER"      : 5,
        "ANGLE"      : 45.,
        "AXIOM"      : "$A--FB--FC--FD--F",
        "RULES"      : strings.NewReplacer("A", "A(-45)FB+FFD(45)FA",
                                           "B", "B(-135)FC+FFA(-45)FB",
                                           "C", "C(135)FD+FFB(-135)FC",
                                           "D", "D(45)FA+FFC(135)FD"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "goldenrod"},
      "Sierpinski II" : map[string]interface{} {
        "ORDER"      : 11,
        "ANGLE"      : 45.,
        "AXIOM"      : "A--F--A--F",
        "RULES"      : strings.NewReplacer("A", "+B-F-B+",
                                           "B", "-A+F+A-"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "goldenrod"},
      "Sierpinski Arrowhead" : map[string]interface{} { //see https://en.wikipedia.org/wiki/Sierpi%C5%84ski_arrowhead_curve
        "ORDER"      : 6,
        "ANGLE"      : 60.,
        "AXIOM"      : "XF",
        "RULES"      : strings.NewReplacer("X", "YF-XF-Y",
                                           "Y", "XF+YF+X"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "goldenrod"},
      "Sierpinski Carpet" : map[string]interface{} { //see http://ecademy.agnesscott.edu/~lriddle/ifs/carpet/carpet.htm
        "ORDER"      : 5,
        "ANGLE"      : 90.,
        "AXIOM"      : "(45)F",
        "RULES"      : strings.NewReplacer("F", "F+F-F-F-f+F+F+F-F",
                                           "f", "fff"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "goldenrod"},
      "Sierpinski Triangle" : map[string]interface{} { //also called Sierpinski Gasket or Sierpinski Sieve;
                                                       //see http://mathforum.org/advanced/robertd/lsys2d.html
        "ORDER"      : 6,
        "ANGLE"      : 60.,
        "AXIOM"      : "FXF++FF++FF",
        "RULES"      : strings.NewReplacer("X", "++FXF--FXF--FXF++",
                                           "F", "FF"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "goldenrod"},
      "Square" : map[string]interface{} { //see http://mathforum.org/advanced/robertd/lsys2d.html
        "ORDER"      : 5,
        "ANGLE"      : 90.,
        "AXIOM"      : "F+XF+F+XF",
        "RULES"      : strings.NewReplacer("X", "XF-F+F-XF+F+XF-F+F-X"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "violet"},
      "LTPL Fig 1.2" : map[string]interface{} { // http://www.vcn.bc.ca/~griffink/lisp_lsystems.pdf by KA Erstad - 2002
        "ORDER"      : 3,
        "ANGLE"      : 72.,
        "AXIOM"      : "F+F+F+F+F",
        "RULES"      : strings.NewReplacer("F", "FF+F+F+F+F+FF"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "black",
        "BGCOLOR"    : lsystems.EncodeBgColorName("steelblue")},

      //Prusinkiewicz, Przemyslaw; Lindenmayer, Aristid (1990). The Algorithmic Beauty of Plants. Springer-Verlag. pp. 101-107.
      //ISBN 978-0-387-97297-8. (http://algorithmicbotany.org/papers/abop/abop.pdf)
      "ABOP Fig 1.8" : map[string]interface{} {
        "ORDER"      : 2,
        "ANGLE"      : 90.,
        "AXIOM"      : "F+F+F+F",
        "RULES"      : strings.NewReplacer( "F", "F+f-FF+F+FF+Ff+FF-f+FF-F-FF-Ff-FFF",
                                            "f", "ffffff"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.9a" : map[string]interface{} {
        "ORDER"      : 4,
        "ANGLE"      : 90.,
        "AXIOM"      : "F-F-F-F",
        "RULES"      : strings.NewReplacer("F", "FF-F-F-F-F-F+F"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.9b" : map[string]interface{} {
        "ORDER"      : 4,
        "ANGLE"      : 90.,
        "AXIOM"      : "F-F-F-F",
        "RULES"      : strings.NewReplacer("F", "FF-F-F-F-FF"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.9c" : map[string]interface{} {
        "ORDER"      : 3,
        "ANGLE"      : 90.,
        "AXIOM"      : "F-F-F-F",
        "RULES"      : strings.NewReplacer("F", "FF-F+F-F-FF"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.9e" : map[string]interface{} {
        "ORDER"      : 5,
        "ANGLE"      : 90.,
        "AXIOM"      : "F-F-F-F",
        "RULES"      : strings.NewReplacer("F", "F-FF--F-F"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.11a(pseudo)" : map[string]interface{} { //hexagonal Gosper curve as a pseudo-L-system
        "ORDER"      : 4,
        "ANGLE"      : 60.,
        "AXIOM"      : "$Fl",
        "RULES"      : strings.NewReplacer("Fl", "Fl+Fr++Fr-Fl--FlFl-Fr+",
                                           "Fr", "-Fl+FrFr++Fr+Fl--Fl-Fr"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "slategray"},
      "ABOP Fig 1.24c" : map[string]interface{} { //branching -> stack ops
        "ORDER"      : 4,
        "ANGLE"      : 22.5,
        "AXIOM"      : "$F",
        "RULES"      : strings.NewReplacer("F", "FF-[-F+F+F]+[+F-F-F]"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "forest-green"},
      "ABOP Fig 1.24e" : map[string]interface{} {
        "ORDER"      : 7,
        "ANGLE"      : 25.7,
        "AXIOM"      : "$X",
        "RULES"      : strings.NewReplacer("X", "F[+X][-X]FX",
                                           "F", "FF"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "forest-green"},

      //Prusinkiewicz, Przemyslaw (1986): Graphical applications of L--systems. In:
      //Graphics Interface 86 May 26-30, 1986, Vancouver, British Columbia, Canada.
      //pp. 247-253. (http://graphicsinterface.org/wp-content/uploads/gi1986-44.pdf)
      "GALS Fig 1d" : map[string]interface{} { //use of the turn-away symbol "|"
        "ORDER"      : 6,
        "ANGLE"      : 90.,
        "AXIOM"      : "$F",
        "RULES"      : strings.NewReplacer("F", "F-FF|F-F"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "black",
        "BGCOLOR"    : "xffffff"},
      "GALS Fig 3c" : map[string]interface{} {
        "ORDER"      : 6,
        "ANGLE"      : 25.7,
        "AXIOM"      : "$G",
        "RULES"      : strings.NewReplacer("G", "GFX[-G][+G]",
                                           "X", "X[+FFF][-FFF]FX"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "forest-green",
        "BGCOLOR"    : lsystems.EncodeBgColorName("gray90")},
      "GALS Fig 3f" : map[string]interface{} { //conifer like
        "ORDER"      : 9,
        "ANGLE"      : 18.,
        "AXIOM"      : "$SLFFF",
        "RULES"      : strings.NewReplacer("S", "[---G][+++G]TS",
                                           "G", "-H[+G]L",
                                           "H", "+G[-H]L",
                                           "T", "TL",
                                           "L", "[+FFF][-FFF]F"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "forest-green",
        "BGCOLOR"    : lsystems.EncodeBgColorName("gray90")},
      "GALS Fig 4" : map[string]interface{} { //use of polygons to create leaves
        "ORDER"      : 4,
        "ANGLE"      : 30.,
        "AXIOM"      : "$T",
        "RULES"      : strings.NewReplacer("T",  "R-[T]++[++L]R[--L]+[T]--T",
                                           "R",  "F[++L][--L]F",
                                           "L",  "[{-FX+FX+FX-|-FX+FX+FX}]",
                                           "FX", "FX",
                                           "F",  "FF"),
        "LINE_WIDTH" : 3,
        "LINE_COLOR" : "forest-green",
        "BGCOLOR"    : lsystems.EncodeBgColorName("gray90")},

      //Prusinkiewicz, P. and Hanan, J. (2013) "Lindenmayer Systems, Fractals, and Plants", Volume 79 of Lecture Notes in
      //Biomathematics, Springer Science & Business Media, Springer Science & Business Media, ISBN 1475714289, 9781475714289
      "LSFP Fig 2.9" : map[string]interface{} { //fractal consisting of a single filled polygon [Szillard and Quinton 1979]
        "ORDER"      : 3,
        "ANGLE"      : 60.,
        "AXIOM"      : "(-60){XF-F-XF-F-XF-F}",
        "RULES"      : strings.NewReplacer("X", "XF-F-XF+F+F+XF+F-F-F+F-F-F+X"),
        "LINE_WIDTH" : 1,
        "LINE_COLOR" : "black",
        "BGCOLOR"    : "xffffff"},
      "LSFP Fig 6.1c" : map[string]interface{} { //spiral tiling
        "ORDER"      : 5,
        "ANGLE"      : 15.,
        "AXIOM"      : "AAAA",
        "RULES"      : strings.NewReplacer("A", "X-X-X-X-X-X-",
                                           "X", "[F-F-F-F[+++X+Y]-----F--------F+F+F+F]",
                                           "Y", "[F-F-F-F[+++Y]-----F--------F+F+F+F]"),
        "LINE_WIDTH" : 3,
        "LINE_COLOR" : "orange",
        "BGCOLOR"    : "xffffff"},
      "LSFP Fig 6.6a" : map[string]interface{} { //mango leaves kolam
        "ORDER"      : 7,
        "ANGLE"      : 60.,
        "AXIOM"      : "$A+++A",
        "RULES"      : strings.NewReplacer("A", "f+F-Z-F+fA",
                                           "Z", "F+FF+F++[++Z]F+FF+F++F+FF+F++"),
        "LINE_WIDTH" : 2,
        "LINE_COLOR" : "green"},
    }
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
    counter := 0
    outFile := ""
    for k, v := range examples {
        counter++
        fmt.Printf("\nEXAMPLE %d / %d (%s):\n", counter, len(examples), k)
        plotTitle := fmt.Sprintf("%s(%d)", k, v["ORDER"].(int))
        bgColor   := "x000000"
        if color, ok := v["BGCOLOR"].(string); ok { bgColor = color }
        //Generate turtle commands
        lsystems.Deterministic(v["ORDER"].(int), v["AXIOM"].(string), v["RULES"].(*strings.Replacer))
        //Render
        switch terminal {
            case "hp-gl/2":
                outFile = plotTitle + ".plt"
                lsystems.HpglPlot(v["ANGLE"].(float64), plotTitle, 0.35, outFile)
            case "svg":
                outFile      = plotTitle + ".svg"
                terminalCmd := fmt.Sprintf("set terminal svg lw %d size %d,%d", v["LINE_WIDTH"].(int),
                                           plotWidth, plotHeight)
                outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
                lsystems.Plot(v["ANGLE"].(float64), terminalCmd, outputCmd, plotTitle, "black")
            default:
                outFile      = plotTitle + "." + terminal
                terminalCmd := fmt.Sprintf("set terminal %s lw %d size %d,%d %s",
                                           terminal, v["LINE_WIDTH"].(int), plotWidth, plotHeight, bgColor)
                outputCmd   := fmt.Sprintf(`set output "%s"`, outFile)
                lsystems.Plot(v["ANGLE"].(float64), terminalCmd, outputCmd, plotTitle, v["LINE_COLOR"].(string))
        }
        fmt.Println("output written to", outFile)
    }
}

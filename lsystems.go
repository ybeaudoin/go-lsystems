/*===== Copyright (c) 2016 Yves Beaudoin - All rights reserved - MIT LICENSE (MIT) - Email: webpraxis@gmail.com ================
 *  Package:
 *      lsystems
 *  Overview:
 *      package for processing and drawing Lindenmayer systems (L-systems) using gnuplot and HP-GL/2 graphics.
 *  Variables:
 *      TurtleCmds string
 *          Generated turtle-graphics commands
 *  Functions:
 *      Deterministic(order int, axiom string, rules *strings.Replacer)
 *          Generates the required turtle commands for the specified deterministic and context-free production parameters.
 *      Stochastic(order int, axiom string, rules []string, weights []int)
 *          Generates the required turtle commands for the specified stochastic and context-free production parameters.
 *      HogewegHesper(order int, axiom string, rules map[string]string)
 *          Generates the required turtle commands for the specified Hogeweg and Hesper production parameters as
 *          described in http://algorithmicbotany.org/papers/abop/abop.pdf
 *      EncodeBgColorName(bgColorName string) string
 *          Encodes a color name into an hex string, prefixed with the character "x", for use as the specification
 *          of a gnuplot terminal's background color.
 *      Plot(angle float64, terminalCmd, outputCmd, plotTitle, lineColor string, cmdsFile ...string)
 *          Plots the latest generated turtle commands with the given parameters using gnuplot. The result will be
 *          isometrically scaled and centered. The underlying gnuplot commands can be saved optionally to a text file.
 *      MultiPlot(turtleCmds []string, turtleAngles []float64, terminalCmd, outputCmd, plotTitle string, labels []string,
 *                lineColor string, cmdsFile ...string)
 *          Plots a set of turtle commands with the given parameters using gnuplot. The result will be anisometrically scaled
 *          with the subplots generated left to right. The underlying gnuplot commands can be saved optionally to a text file.
 *      HpglPlot(angle float64, plotTitle string, penWidth float64, hpglPath string)
 *          Converts the latest generated turtle commands with the given parameters to an HP-GL/2 command set.
 *          The resulting plot will be isometrically scaled and centered.
 *      HpglMultiPlot(turtleCmds []string, turtleAngles []float64, plotTitle string, labels []string,
 *                    penWidth float64, hpglPath string)
 *          Converts a set of turtle commands with the given parameters to an HP-GL/2 command set.
 *          The resulting plot will be anisometrically scaled with the subplots generated left to right in landscape mode.
 *  Remarks: L-system symbols:
 *            Variables : any symbol that does not conflict with the constants below,
 *            Constants : F f + - | $ ( ) [ ] { }
 *              with the following turtle-graphics interpretations:
 *                "F" means "move forward drawing a line",
 *                "f" means "move forward without drawing a line",
 *                "+" means "turn left",
 *                "-" means "turn right",
 *                "|" means "turn away" (adds 180 degrees to the heading),
 *                "$" sets the heading to 90 degrees (default initial heading is 0),
 *                "(" marks the start of a new heading declaration in degrees,
 *                ")" marks the end of a heading declaration,
 *                "[" starts a branch; saves the current turtle's status,
 *                "]" ends a branch; restores the turtle's status with the last saved value (as in a LIFO stack),
 *                "{" starts filled polygon mode (line segments define the edges),
 *                "}" ends polygon mode.
 *              All other symbols will be ignored during drawing.
 *  History: v1.0.0 - September 28, 2016 - Original release.
 *============================================================================================================================*/
package lsystems

import(
    "bitbucket.org/binet/go-gnuplot/pkg/gnuplot"
    "fmt"
    "log"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "regexp"
    "runtime"
    "sort"
    "strconv"
    "strings"
    "time"
)
/*Exported -------------------------------------------------------------------------------------------------------------------*/
var TurtleCmds string //generated turtle commands

func Deterministic(order int, axiom string, rules *strings.Replacer) {
/*         Purpose : Generates the required turtle commands for the specified deterministic and context-free production
 *                   parameters.
 *       Arguments : order = order of the curve, that is, the derivation length of the production rules.
 *                           (The zeroth order corresponds to the axiom.)
 *                   axiom = production axiom.
 *                   rules = production rules.
 *         Returns : None.
 * Externals -  In : None.
 * Externals - Out : TurtleCmds
 *       Functions : halt
 *         Remarks : - Supported L-system constants are F f + - | $ ( ) [ ] { }
 *                   - Pseudo-L-systems are supported.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if order < 0   { halt("curve order must be non-negative") }
    if axiom == "" { halt("axiom was not specified") }

    //Apply the production rules
    TurtleCmds = axiom
    for n := 1; n <= order; n++ {
        TurtleCmds = rules.Replace(TurtleCmds)
    }
    return
} //end func Deterministic
func Stochastic(order int, axiom string, rules []string, weights []int) {
/*         Purpose : Generates the required turtle commands for the specified stochastic and context-free production
 *                   parameters.
 *       Arguments : order   = order of the curve, that is, the derivation length of the production rules.
 *                             (The zeroth order corresponds to the axiom.)
 *                   axiom   = production axiom.
 *                   rules   = slice of production rules for the constant "F".
 *                   weights = slice of weights for the production rules governing their chances of being chosen.
 *         Returns : None.
 * Externals -  In : None.
 * Externals - Out : TurtleCmds
 *       Functions : halt
 *         Remarks : - Supported L-system constants are F f + - | $ ( ) [ ] { }
 *                   - Only supports rules that rewrite the constant "F".
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if order        < 0   { halt("curve order must be non-negative") }
    if axiom        == "" { halt("axiom was not specified") }
    if len(rules)   == 0  { halt("rules were not specified") }
    if len(weights) == 0  { halt("weights were not specified") }
    if len(weights) < len(rules) {
        halt("fewer weights specified than the number of rules")
    }

    var selectors []int
    //Set up the chances of picking a rule (like assigning a certain number of balls in a lottery machine for each rule)
    for k, v := range weights {
        if !(v > 0) { halt("weights must be postive") }
        for n := 1; n <= v; n++ {
            selectors = append(selectors, k)
        }
    }
    //Apply the production rules
    numSelectors := len(selectors)
    TurtleCmds    = axiom
    for n := 1; n <= order; n++ {
        newCmds := ""
        for _, symbol := range strings.Split(TurtleCmds, "") {
            if symbol == "F" {
                newCmds += rules[selectors[rand.Intn(numSelectors)]]
            } else {
                newCmds += symbol
            }
        }
        TurtleCmds = newCmds
    }
    return
} //end func Stochastic
func HogewegHesper(order int, axiom string, rules map[string]string) {
/*         Purpose : Generates the required turtle commands for the specified Hogeweg and Hesper production parameters as
 *                   described in http://algorithmicbotany.org/papers/abop/abop.pdf
 *       Arguments : order = order of the curve, that is, the derivation length of the production rules.
 *                           (The zeroth order corresponds to the axiom.)
 *                   axiom = production axiom.
 *                   rules = map of production rules.
 *         Returns : None.
 * Externals -  In : None.
 * Externals - Out : TurtleCmds
 *       Functions : getContextLHS, getContextRHS, halt
 *         Remarks : - Supported L-system constants are F + - $ [ ]
 *                   - Variables are "0" and "1"
 *                   - The format for each rule is "L < a > R" : "replacemnt" where "L" denotes the left context, "a" the strict
 *                     predecessor and "R" the right context. For example, "0 < 0 > 1" : "1[+F1F1]" corresponds to the rule
 *                     0 < 0 > 1 -> 1[+F1F1].
 *                   - The rules for the turn constants "+" and "-" need not by supplied as they have been coded to alternate.
 *                   - This is "...a restricted case where daughter branches do not belong to the context of the mother branch."
 *                     (Prusinkiewicz, P. and Hanan, J. (2013) "Lindenmayer Systems, Fractals, and Plants", Volume 79 of Lecture
 *                     Notes in Biomathematics, Springer Science & Business Media, Springer Science & Business Media,
 *                     ISBN 1475714289, 9781475714289, p.42)
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if order < 0   { halt("curve order must be non-negative") }
    if axiom == "" { halt("axiom was not specified") }
    for k, v := range rules {
        if matched, _ := regexp.MatchString("^(0|1) < (0|1) > (0|1)$", k); !matched || v == "" {
            halt("the rules were not specified correctly")
         }
    }

    var( chan4LHS = make(chan string) //io channel for Coroutine getContextLHS
         chan4RHS = make(chan string) //io channel for Coroutine getContextRHS
    )
    //Initialize
    go getContextLHS(chan4LHS) //launch coroutines to find the l-h-s and r-h-s contexts respectively
    go getContextRHS(chan4RHS)
    //Apply the production rules
    TurtleCmds = axiom
    for n := 1; n <= order; n++ {
        newCmds := ""
        for pos := 0; pos < len(TurtleCmds); pos++ {
            symbol := string(TurtleCmds[pos])
            switch symbol  {
                case "F", "[", "]", "$":
                    newCmds += symbol
                case "+":
                    newCmds += "-"
                case "-":
                    newCmds += "+"
                case "0", "1":
                    chan4LHS<- TurtleCmds[:pos]
                    chan4RHS<- TurtleCmds[pos+1:]
                    context := fmt.Sprintf("%s < %s > %s", <-chan4LHS, symbol, <-chan4RHS)
                    if replacement, ok := rules[context]; ok {
                        newCmds += replacement
                    } else {
                        newCmds += symbol
                    }
                default:
                    halt(fmt.Sprintf("the symbol '%s' is not supported", symbol))
            }
        }
        TurtleCmds = newCmds
    }
    //Terminate coroutines
    close(chan4LHS)
    close(chan4RHS)
} //end func HogewegHesper
func EncodeBgColorName(bgColorName string) string {
/*         Purpose : Encodes a color name into an hex string, prefixed with the character "x", for use as the specification
 *                   of a gnuplot terminal's background color.
 *       Arguments : bgColorName = color name recognized by gnuplot.
 *         Returns : hex encoding
 * Externals -  In : _colorNames, _validColors
 * Externals - Out : None.
 *       Functions : halt
 *         Remarks : None.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    hexRGB, ok := _colorNames[strings.ToLower(bgColorName)]
    if ! ok { halt( fmt.Sprintf("the color name '%s' is not valid. Recognized names are:\n\n%s", bgColorName, _validColors)) }
    return "x" + hexRGB
}
func Plot(angle float64, terminalCmd, outputCmd, plotTitle, lineColor string, cmdsFile ...string) {
/*         Purpose : Plots the latest generated turtle commands with the given parameters using gnuplot. The result will be
 *                   isometrically scaled and centered. The underlying gnuplot commands can be saved optionally to a text file.
 *       Arguments : angle       = production angle in degrees.
 *                   terminalCmd = gnuplot terminal command.
 *                   outputCmd   = gnuplot output command.
 *                   plotTitle   = title to be centered at the top of the plot.
 *                   lineColor   = color of the line segments, specified as either a name (as recognized by gnuplot)
 *                                 or a 6-digit X11 hex rgb code prefixed with the "#" character.
 *                   cmdsFile    = optional file path for the gnuplot commands.
 *         Returns : None.
 * Externals -  In : TurtleCmds, _degs2rads, _turtleHistory, _turtleStatus, _validColors
 * Externals - Out : TurtleCmds
 *       Functions : execPlot, fileWrite, halt, makeLogo2Gnuplot, validFgColor
 *         Remarks : - Supported L-system constants are F f + - | $ ( ) [ ] { }
 *                     All other symbols will be ignored.
 *                   - The default turtle heading is 0 degrees.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if TurtleCmds == ""          { halt("the turtle commands were not generated") }
    if angle      == 0.          { halt("the production angle is zero") }
    if ! validFgColor(lineColor) { halt( fmt.Sprintf("the color name '%s' is not valid. Recognized names are:\n\n%s",
                                                     lineColor, _validColors)) }

    const( minMargin    = "1"
           maxMargin    = "2"
    )
    var(   rmargin      = minMargin
           lmargin      = minMargin
           bmargin      = minMargin
           tmargin      = map[bool]string{true: maxMargin, false: minMargin} [plotTitle != ""]

           logo2Gnuplot = makeLogo2Gnuplot(lineColor)
           plotCmds     []string
    )
    //Initialize
    plotCmds   = append(plotCmds,
                  terminalCmd,
                  outputCmd,
                  "unset border",
                  "unset border",
                  "unset tics",
                  "set bmargin " + bmargin,
                  "set tmargin " + tmargin,
                  "set rmargin " + rmargin,
                  "set lmargin " + lmargin,
                  "set size square",
                  "set autoscale fix",
                  fmt.Sprintf(`set style fill solid 1.0 border rgb "%s"`, lineColor),
                  fmt.Sprintf(`set style arrow 1 nohead lc rgb "%s"`, lineColor))
    if plotTitle != "" { plotCmds = append(plotCmds, fmt.Sprintf(`set title "%s" tc rgb "%s"`, plotTitle, lineColor)) }
   //Convert the turtle commands to headless arrows using unit turtle strides
    drawCmds, xMin, xMax, yMin, yMax := logo2Gnuplot(&TurtleCmds, 0., angle)
    plotCmds = append(plotCmds, drawCmds...)
    //Compute offsets so as to center the plot in a square bounding box
    xSpan, ySpan := xMax - xMin, yMax - yMin
    maxSpan      := math.Max(xSpan, ySpan)
    xOffset      := 0.5 * (maxSpan - xSpan)
    yOffset      := 0.5 * (maxSpan - ySpan)
    //Compose the remaining commands
    plotCmds = append(plotCmds,
                fmt.Sprintf("set xrange [%f:%f]", xMin, xMax),
                fmt.Sprintf("set yrange [%f:%f]", yMin, yMax),
                fmt.Sprintf("set offset %f,%f,%f,%f", xOffset, xOffset, yOffset, yOffset),
                "set parametric",
                fmt.Sprintf(`plot 0,0 notitle lc rgb "%s" lw 0`, lineColor),
                "quit")
    //Send the commands to the gnuplot executable
    execPlot(terminalCmd, &plotCmds)
    //Save the commands if requested
    if len(cmdsFile) != 0 { fileWrite(cmdsFile[0], strings.Join(plotCmds, "\n")) }
    return
} //end func Plot
func MultiPlot(turtleCmds []string, turtleAngles []float64, terminalCmd, outputCmd, plotTitle string, labels []string,
               lineColor string, cmdsFile ...string) {
/*         Purpose : Plots a set of turtle commands with the given parameters using gnuplot. The result will be
 *                   anisometrically scaled with the subplots generated left to right. The underlying
 *                   gnuplot commands can be saved optionally to a text file.
 *       Arguments : turtleCmds   = slice of turtle commands.
 *                   turtleAngles = slice of production angles in degrees.
 *                   terminalCmd  = gnuplot terminal command.
 *                   outputCmd    = gnuplot output command.
 *                   plotTitle    = title to be centered at the top of the plot.
 *                   labels       = slice of labels to be centered below each subplot.
 *                   lineColor    = color of the line segments, specified as either a name (as recognized by gnuplot)
 *                                  or a 6-digit X11 hex rgb code prefixed with the "#" character.
 *                   cmdsFile     = optional file path for the gnuplot commands.
 *         Returns : None.
 * Externals -  In : _degs2rads, _turtleHistory, _turtleStatus, _validColors
 * Externals - Out : None.
 *       Functions : calcXoffset, execPlot, fileWrite, halt, makeLogo2Gnuplot, validFgColor
 *         Remarks : - Supported L-system constants are F f + - | $ ( ) [ ] { }
 *                     All other symbols will be ignored.
 *                   - The default turtle heading is 0 degrees.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if len(turtleCmds)   == 0 { halt("the turtle commands were not specified") }
    if len(turtleAngles) == 0 { halt("the turtle angles were not stated") }
    if len(labels)       == 0 { halt("the labels were not specified") }
    if len(turtleAngles) < len(turtleCmds) {
        halt("fewer turtle angles specified than the number of commands")
    }
     if len(labels) < len(turtleCmds) {
        halt("fewer labels specified than the number of commands")
    }
    if ! validFgColor(lineColor) { halt( fmt.Sprintf("the color name '%s' is not valid. Recognized names are:\n\n%s",
                                                     lineColor, _validColors)) }

    const( minMargin    = "1"
           maxMargin    = "2"
    )
    var(   rmargin      = minMargin
           lmargin      = minMargin
           bmargin      = map[bool]string{true: maxMargin, false: minMargin} [strings.Join(labels, "") != ""]
           tmargin      = map[bool]string{true: maxMargin, false: minMargin} [plotTitle                != ""]

           logo2Gnuplot = makeLogo2Gnuplot(lineColor)
           drawCmds     []string
           plotCmds     []string
           xMin         float64
           xMax         float64
           yMin         float64
           yMax         float64

           commands     = make(chan string)  //io channels for Coroutine calcXoffset
           angle        = make(chan float64)
           xOffset      = make(chan float64)
    )
    //Initialize
    go calcXoffset(commands, angle, xOffset) //launch helper coroutine to calc the x-offsets of a subplot
    plotCmds = append(plotCmds,
                terminalCmd,
                outputCmd,
                "unset border",
                "unset tics",
                "set bmargin " + bmargin,
                "set tmargin " + tmargin,
                "set rmargin " + rmargin,
                "set lmargin " + lmargin,
                "set autoscale fix",
                fmt.Sprintf(`set style fill solid 1.0 border rgb "%s"`, lineColor),
                fmt.Sprintf(`set style arrow 1 nohead lc rgb "%s"`, lineColor))
    if plotTitle != "" { plotCmds = append(plotCmds, fmt.Sprintf(`set title "%s" tc rgb "%s"`, plotTitle, lineColor)) }
   //Convert the turtle commands to headless arrows using unit turtle strides
    xOrigin := 0.
    for k, v := range turtleCmds {
        if k + 1 < len(turtleCmds) {
            commands<- turtleCmds[k+1]
            angle<-    turtleAngles[k+1]
        }
        if labels[k] != "" {
            plotCmds  = append(plotCmds, fmt.Sprintf(`set label "%s" at %f,character 1 center front tc rgb "%s"`,
                                                     labels[k], xOrigin, lineColor))
        }
        drawCmds, xMin, xMax, yMin, yMax = logo2Gnuplot(&v, xOrigin, turtleAngles[k])
        plotCmds = append(plotCmds, drawCmds...)
        if k + 1 < len(turtleCmds) { xOrigin = xMax + <-xOffset }
    }
    //Compose the remaining gnuplot commands
    plotCmds = append(plotCmds,
                fmt.Sprintf("set xrange [%f:%f]", xMin, xMax),
                fmt.Sprintf("set yrange [%f:%f]", yMin, yMax),
                "set parametric",
                fmt.Sprintf(`plot 0,0 notitle lc rgb "%s" lw 0`, lineColor),
                "quit")
    //Send the commands to the gnuplot executable
    execPlot(terminalCmd, &plotCmds)
    //Save the commands if requested
    if len(cmdsFile) != 0 { fileWrite(cmdsFile[0], strings.Join(plotCmds, "\n")) }
    //Terminate coroutine
    close(commands)
    return
} //end func MultiPlot
func HpglPlot(angle float64, plotTitle string, penWidth float64, hpglPath string) {
/*         Purpose : Converts the latest generated turtle commands with the given parameters to an HP-GL/2 command set.
 *                   The resulting plot will be isometrically scaled and centered.
 *       Arguments : angle     = production angle in degrees.
 *                   plotTitle = title to be centered at the top of the plot.
 *                   penWidth  = line-width in millimeters.
 *                   hpglPath  = file path or device port for the HP-GL/2 commands.
 *         Returns : None.
 * Externals -  In : TurtleCmds, _degs2rads, _turtleHistory, _turtleStatus
 * Externals - Out : TurtleCmds
 *       Functions : fileWrite, getHeading, halt, makeLogo2Hpgl
 *         Remarks : Supported L-system constants are F f + - | $ ( ) [ ]
 *                   All other symbols will be ignored.
 *                   The default turtle heading is 0 degrees.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if TurtleCmds == "" { halt("the turtle commands were not generated") }
    if angle      == 0. { halt("the production angle is zero") }
    if hpglPath == ""   { halt("the path for the plot was not specified") }

    const( esc       = 27 //Escape code
           ext       = 3  //End of Text code
           minMargin = 0.1 //% - prevent clipping of wide pen strokes
           maxMargin = 3.0 //% - prevent clipping of title
    )
    var(
           rmargin   = minMargin
           lmargin   = 100. - minMargin
           bmargin   = minMargin
           tmargin   = map[bool]float64{true: 100. - maxMargin, false: 100. - minMargin} [plotTitle != ""]

           logo2Hpgl = makeLogo2Hpgl()
    )
    //Convert the turtle commands to HP-GL/2 commands using unit turtle strides
    plotCmds, xMin, xMax, yMin, yMax := logo2Hpgl(&TurtleCmds, 0., angle)
    //Compute offsets so as to center the plot in a square bounding box
    xSpan, ySpan := xMax - xMin, yMax - yMin
    maxSpan      := math.Max(xSpan, ySpan)
    xOffset      := 0.5 * (maxSpan - xSpan)
    yOffset      := 0.5 * (maxSpan - ySpan)
    //Compose the remaining HP-GL/2 commands
    xMin -= xOffset; xMax += xOffset
    yMin -= yOffset; yMax += yOffset
    plotCmds = //HP RTL: enter HP-GL/2 mode, begin a plot and initialize HP-GL/2
               fmt.Sprintf("%c%%-1BBPIN;\n", esc) +
               //set the margins for the drawing
               fmt.Sprintf("IR%f,%f,%f,%f;\n", rmargin, bmargin, lmargin, tmargin) +
               //set the scaling as isotropic
               fmt.Sprintf("SC%f,%f,%f,%f,1;\n", xMin, xMax, yMin, yMax) +
               //select Pen 1 (black) and set its width in millimeters
               fmt.Sprintf("SP1;WU0;PW%f;\n", penWidth) +
               //add the previous turtle pen commands
               plotCmds + "\n"
    if plotTitle != "" {
        plotCmds += //reset the margin settings for the title
                    fmt.Sprintf("IR;IR%f,%f,%f,%f;\n", rmargin, bmargin, lmargin, 100. - minMargin) +
                    //set the scaling as anisotropic
                    fmt.Sprintf("SC%f,%f,%f,%f,0;\n", xMin, xMax, yMin, yMax) +
                    //draw the plot title centered at the top
                    fmt.Sprintf("PU%f,%f;LO6;LB%s%c;\n", 0.5*(xMin + xMax), yMax, plotTitle, ext)
    }
    //end (page advance)
    plotCmds += "PG;\n"
    //Output the commands to the specified destination
    fileWrite(hpglPath, plotCmds)
    return
} //end func HpglPlot
func HpglMultiPlot(turtleCmds []string, turtleAngles []float64, plotTitle string, labels []string,
                   penWidth float64, hpglPath string) {
/*         Purpose : Converts a set of turtle commands with the given parameters to an HP-GL/2 command set.
 *                   The resulting plot will be anisometrically scaled with the subplots generated left to right in
 *                   landscape mode.
 *       Arguments : turtleCmds   = slice of turtle commands.
 *                   turtleAngles = slice of production angles in degrees.
 *                   plotTitle    = title to be centered at the top of the plot.
 *                   labels       = slice of labels to be centered below each subplot.
 *                   penWidth     = line-width in millimeters.
 *                   hpglPath     = file path or device port for the HP-GL/2 commands.
 *         Returns : None.
 * Externals -  In : _degs2rads, _turtleHistory, _turtleStatus
 * Externals - Out : None.
 *       Functions : calcXoffset, fileWrite, getHeading, halt, makeLogo2Hpgl
 *         Remarks : Supported L-system constants are F f + - | $ ( ) [ ]
 *                   All other symbols will be ignored.
 *                   The default turtle heading is 0 degrees.
 *         History : v1.0.0 - September 28, 2016 - Original release.
 */
    if len(turtleCmds)   == 0 { halt("the turtle commands were not specified") }
    if len(turtleAngles) == 0 { halt("the turtle angles were not stated") }
    if len(labels)       == 0 { halt("the labels were not specified") }
    if len(turtleAngles) < len(turtleCmds) {
        halt("fewer turtle angles specified than the number of commands")
    }
    if len(labels) < len(turtleCmds) {
        halt("fewer labels specified than the number of commands")
    }
    if hpglPath == "" { halt("the path for the plot was not specified") }

    const( esc       = 27  //Escape code
           ext       = 3   //End of Text code
           minMargin = 0.1 //% - prevent clipping of wide pen strokes
           maxMargin = 3.0 //% - prevent clipping of title and subplot labels
           yNudge    = 1.0 //user unit - increase gap between drawing and title and subplot labels
    )
    var(   rmargin   = minMargin
           lmargin   = 100. - minMargin
           bmargin   = map[bool]float64{true: maxMargin,        false: minMargin       } [strings.Join(labels, "") != ""]
           tmargin   = map[bool]float64{true: 100. - maxMargin, false: 100. - minMargin} [plotTitle                != ""]

           drawCmds  string
           logo2Hpgl = makeLogo2Hpgl()
           plotCmds  string
           xMin      float64
           xMax      float64
           yMin      float64
           yMax      float64

           commands  = make(chan string)  //io channels for Coroutine calcXoffset
           angle     = make(chan float64)
           xOffset   = make(chan float64)
    )
    //Initialize
    go calcXoffset(commands, angle, xOffset) //launch coroutine to calc x-offsets of subplots
    //Convert the turtle commands to HP-GL/2 commands using unit turtle strides
    xOrigin := 0.
    for k, v := range turtleCmds {
        if k + 1 < len(turtleCmds) {
            commands<- turtleCmds[k+1]
            angle<-    turtleAngles[k+1]
        }
        if labels[k] != "" { plotCmds += fmt.Sprintf("PU%f,%f;LO16;LB%s%c;\n", xOrigin, -yNudge, labels[k], ext) }
        drawCmds, xMin, xMax, yMin, yMax = logo2Hpgl(&v, xOrigin, turtleAngles[k])
        plotCmds += drawCmds
        if k + 1 < len(turtleCmds) { xOrigin = xMax + <-xOffset }
    }
    //Compose the remaining HP-GL/2 commands
    plotCmds = //HP RTL: enter HP-GL/2 mode, begin a plot and initialize HP-GL/2
               fmt.Sprintf("%c%%-1BBPIN;\n", esc) +
               //set the orientation to landscape
               "RO90;\n" +
               //set the margins for the drawing
               fmt.Sprintf("IR%f,%f,%f,%f;\n", rmargin, bmargin, lmargin, tmargin) +
               //set the scaling as anisotropic
               fmt.Sprintf("SC%f,%f,%f,%f,0;\n", xMin, xMax, yMin, yMax) +
               //select Pen 1 (black) and set its width in millimeters
               fmt.Sprintf("SP1;WU0;PW%f;\n", penWidth) +
               //add the previous turtle pen commands
               plotCmds + "\n"
    if plotTitle != "" { //draw the plot title centered at the top
        plotCmds += fmt.Sprintf("PU%f,%f;LO14;LB%s%c;\n", 0.5*(xMin + xMax), yMax + yNudge, plotTitle, ext)
    }
    //end (page advance)
    plotCmds += "PG;"
    //Output the commands to the specified destination
    fileWrite(hpglPath, plotCmds)
    //Terminate coroutine
    close(commands)
    return
} //end func HpglMultiPlot
/*Private  -------------------------------------------------------------------------------------------------------------------*/
type _turtleStatus struct {
    HEADING float64 //turtle's heading in degrees
    X       float64 //turtle's x ordinate
    Y       float64 //turtle's y ordinate
}
type _turtleHistory []_turtleStatus
const _progressBarLen = 50
var( _colorNames      = map[string]string{}
     _degs2rads       = math.Pi / 180.
     _reHeading       = regexp.MustCompile(`^\(([\+\-0-9.]+?)\)`)
     _reTerminal      = regexp.MustCompile(`^\s*set\s+terminal\s+(.+?)\s+`)
     _validColors     string
)
////Package initialization
func init() {
    var names []string
    //Seeding for stochastic processes
    rand.Seed(time.Now().UnixNano())
    //gnuplot predefined color names
    out, _ := exec.Command("gnuplot", "-e", "show colornames").CombinedOutput()
    for _, v := range regexp.MustCompile(`(?m)^\s+(.+?)\s+#(.+?) =`).FindAllStringSubmatch(string(out), -1) {
        _colorNames[v[1]] = v[2] //name => hex code
        names             = append(names, v[1])
    }
    sort.Strings(names)
    _validColors = strings.Join(names, ", ")
} //end func init
////Reporting
func halt(msg string) {
    pc, _, _, ok := runtime.Caller(1)
    details      := runtime.FuncForPC(pc)
    if ok && details != nil {
        log.Fatalln(fmt.Sprintf("\a%s: %s", details.Name(), msg))
    }
    log.Fatalln("\alsystems: FATAL ERROR!")
} //end func halt
func updateProgressBar(title string, current, total int) {
    //code derived from Graham King's post "Pretty command line / console output on Unix in Python and Go Lang"
    //(http://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/)
    prefix := fmt.Sprintf("%s: %d / %d ", title, current, total)
    amount := int(0.1 + float32(_progressBarLen) * float32(current) / float32(total))
    remain := _progressBarLen - amount
    bar    := strings.Repeat("\u2588", amount) + strings.Repeat("\u2591", remain)
    os.Stdout.WriteString(prefix + bar + "\r")
    if current == total { os.Stdout.WriteString(strings.Repeat(" ", len(prefix) + _progressBarLen) + "\r") }
    os.Stdout.Sync()
    return
} //end func updateProgressBar
////Stack operations
func(lifoStack *_turtleHistory) push(turtle _turtleStatus) {
    (*lifoStack) = append((*lifoStack), turtle)
    return
} //end func push
func(lifoStack *_turtleHistory) pop() (turtle _turtleStatus){
    lastIdx      := len(*lifoStack) - 1
    turtle        = (*lifoStack)[lastIdx]
    (*lifoStack)  = (*lifoStack)[:lastIdx]
    return turtle
} //end func pop
////Plot operations
func calcXoffset(commands <-chan string, productionAngle <-chan float64, xOffset chan<- float64) {
    var( stack  _turtleHistory
         turtle _turtleStatus
         xNudge = 2.
    )
    for turtleCmds := range commands {
        turtleCmds  = strings.NewReplacer("+-", "", "-+", "").Replace(turtleCmds) //remove pointless turns
        angle      := <-productionAngle
        turtle.HEADING, turtle.X = 0., 0.
        xMin := 0.
        pos  := 0
        for pos < len(turtleCmds) {
            symbol := string(turtleCmds[pos])
            switch symbol {
                case "F", "f": //draw or move forward
                    switch math.Mod(turtle.HEADING, 360.) {
                        case 0.:
                            turtle.X++
                        case 90., -270.:
                        case 180., -180.:
                            turtle.X--
                            xMin = math.Min(xMin, turtle.X)
                        case 270., -90.:
                        default:
                            turtle.X += math.Cos(turtle.HEADING * _degs2rads)
                            xMin      = math.Min(xMin, turtle.X)
                    }
                case "+": //turn left
                    turtle.HEADING += angle
                case "-": //turn right
                    turtle.HEADING -= angle
                case "|": //turn away
                    turtle.HEADING += 180.
                case "$": //head due north
                    turtle.HEADING = 90.
                case "(": //set arbitrary heading
                    turtle.HEADING, pos = getHeading(&turtleCmds, pos)
                case "[": //store status
                    stack.push(turtle)
                case "]": //restore status
                    turtle = stack.pop()
                case "{", "}": //ignore polygon mode
                default:  //remove all instances of a production variable
                    turtleCmds = strings.NewReplacer(symbol, "").Replace(turtleCmds)
                    pos--
            }
            pos++
        }
        xOffset<- math.Abs(xMin) + xNudge
    }
    return
} //end coroutine calcXoffset
func execPlot(terminalCmd string, plotCmds *[]string) {
    title := "gnuplot"
    matches := _reTerminal.FindStringSubmatch(terminalCmd)
    if matches != nil {
        title += " -> " + strings.ToUpper(matches[1])
    }
    plotter, err := gnuplot.NewPlotter("", false, false)
    if err != nil { halt("execPlot - " +err.Error()) }
    kMax := len(*plotCmds) - 1
    for k, v := range *plotCmds {
        updateProgressBar(title, k, kMax)
        plotter.CheckedCmd("%s", v)
    }
    plotter.Close()
    return
} //end func execPlot
func fileWrite(filepath string, content string) {
    writer, err := os.Create(filepath)
    if err != nil { halt("os.Create - " + err.Error()) }
    if _, err = writer.WriteString(content); err != nil { halt("writer.WriteString - " + err.Error()) }
    if err = writer.Sync(); err != nil { halt("writer.Sync - " + err.Error()) }
    if err = writer.Close(); err != nil { halt("writer.Close - " + err.Error()) }
    return
} //end func fileWrite
func getHeading(cmds *string, posLP int) (heading float64, pos int) {
    matches := _reHeading.FindStringSubmatch((*cmds)[posLP:])
    if matches == nil { halt("the specified angle is not syntactically well-formed") }
    heading, err := strconv.ParseFloat(matches[1], 64)
    if err != nil { halt("the specified angle is not syntactically well-formed") }
    pos = posLP + len(matches[0]) - 1
    return
} //end func getHeading
func makeLogo2Gnuplot(lineColor string) func(turtleCmds *string, xOrigin, angle float64) (plotCmds []string,
                                                                                          xMini, xMaxi, yMini, yMaxi float64) {
    xMin, xMax := 0., 0.
    yMin, yMax := 0., 0.
    return func(turtleCmds *string, xOrigin, angle float64) (plotCmds []string, xMini, xMaxi, yMini, yMaxi float64) {
            var( convert2Gnuplot = makeConvert2Gnuplot(lineColor)
                 stack           _turtleHistory
                 turtle          = _turtleStatus{0., xOrigin, 0.}
            )
            //Initialize
            *turtleCmds = strings.NewReplacer("+-", "", "-+", "").Replace(*turtleCmds) //remove pointless turns
            //Convert the turtle commands to gnuplot line segments using unit turtle strides
            pos := 0
            for pos < len(*turtleCmds) {
                updateProgressBar("logo -> gnuplot", pos, len(*turtleCmds)-1)
                symbol := string((*turtleCmds)[pos])
                switch symbol {
                    case "F", "f": //draw or move forward
                        xFrom, yFrom := turtle.X, turtle.Y
                        switch math.Mod(turtle.HEADING, 360.) {
                            case 0.:
                                turtle.X++
                                xMax  = math.Max(xMax, turtle.X)
                            case 90., -270.:
                                turtle.Y++
                                yMax  = math.Max(yMax, turtle.Y)
                            case 180., -180.:
                                turtle.X--
                                xMin  = math.Min(xMin, turtle.X)
                            case 270., -90.:
                                turtle.Y--
                                yMin  = math.Min(yMin, turtle.Y)
                            default:
                                radians    := turtle.HEADING * _degs2rads
                                turtle.X   += math.Cos(radians)
                                turtle.Y   += math.Sin(radians)
                                xMin, xMax  = math.Min(xMin, turtle.X), math.Max(xMax, turtle.X)
                                yMin, yMax  = math.Min(yMin, turtle.Y), math.Max(yMax, turtle.Y)
                        }
                        if cmd := convert2Gnuplot(symbol, xFrom, yFrom, turtle.X, turtle.Y); cmd != "" {
                            plotCmds = append(plotCmds, cmd)
                        }
                    case "+": //turn left
                        turtle.HEADING += angle
                    case "-": //turn right
                        turtle.HEADING -= angle
                    case "|": //turn away
                        turtle.HEADING += 180.
                    case "$": //head due north
                        turtle.HEADING = 90.
                    case "(": //set arbitrary heading
                        turtle.HEADING, pos = getHeading(turtleCmds, pos)
                    case "[": //store status
                        stack.push(turtle)
                    case "]": //restore status
                        turtle = stack.pop()
                    case "{": //start polygon mode
                        convert2Gnuplot(symbol, turtle.X, turtle.Y)
                    case "}": //end  polygon mode
                        plotCmds = append(plotCmds, convert2Gnuplot(symbol))
                    default:  //remove all instances of a production variable
                        *turtleCmds = strings.NewReplacer(symbol, "").Replace(*turtleCmds)
                        pos--
                }
                pos++
            }
            xMini, xMaxi, yMini, yMaxi = xMin, xMax, yMin, yMax
            return
           }
} //end func makeLogo2Gnuplot
func makeConvert2Gnuplot(lineColor string) func(symbol string, coords ...float64) (pltCmd string) {
    fillColor := lineColor
    polygon   := ""
    return func(symbol string, coords ...float64) (pltCmd string) {
            switch {
                case symbol == "{": //initiate polygon mode
                    polygon = fmt.Sprintf(`set object polygon fc rgb "%s" from %f,%f`, fillColor, coords[0], coords[1])
                case symbol == "}": //terminate polygon mode
                    pltCmd, polygon = polygon, ""
                case polygon != "": //continue with polygon
                    polygon += fmt.Sprintf(" to %f,%f", coords[2], coords[3])
                case symbol == "F": //draw line segment
                    pltCmd = fmt.Sprintf("set arrow as 1 from %f,%f to %f,%f", coords[0], coords[1], coords[2], coords[3])
            }
            return
           }
} //end func makeConvert2Gnuplot
func makeLogo2Hpgl() func(turtleCmds *string, xOrigin, angle float64) (plotCmds string, xMini, xMaxi, yMini, yMaxi float64) {
    xMin, xMax := 0., 0.
    yMin, yMax := 0., 0.
    return func(turtleCmds *string, xOrigin, angle float64) (plotCmds string, xMini, xMaxi, yMini, yMaxi float64) {
            var( convert2Hpgl = makeConvert2Hpgl()
                 stack        _turtleHistory
                 turtle       = _turtleStatus{0., xOrigin, 0.}
            )
            //Initialize
            *turtleCmds = strings.NewReplacer("+-", "", "-+", "").Replace(*turtleCmds) //remove pointless turns
            plotCmds    = convert2Hpgl("f", xOrigin, 0.)
            //Convert the turtle commands to HP-GL/2 commands using unit turtle strides
            pos := 0
            for pos < len(*turtleCmds) {
                updateProgressBar("logo -> HP-GL/2", pos, len(*turtleCmds)-1)
                symbol := string((*turtleCmds)[pos])
                switch symbol {
                    case "F", "f": //draw or move forward
                        switch math.Mod(turtle.HEADING, 360.) {
                            case 0.:
                                turtle.X++
                                xMax = math.Max(xMax, turtle.X)
                            case 90., -270.:
                                turtle.Y++
                                yMax = math.Max(yMax, turtle.Y)
                            case 180., -180.:
                                turtle.X--
                                xMin = math.Min(xMin, turtle.X)
                            case 270., -90.:
                                turtle.Y--
                                yMin = math.Min(yMin, turtle.Y)
                            default:
                                radians    := turtle.HEADING * _degs2rads
                                turtle.X   += math.Cos(radians)
                                turtle.Y   += math.Sin(radians)
                                xMin, xMax  = math.Min(xMin, turtle.X), math.Max(xMax, turtle.X)
                                yMin, yMax  = math.Min(yMin, turtle.Y), math.Max(yMax, turtle.Y)
                        }
                        plotCmds += convert2Hpgl(symbol, turtle.X, turtle.Y)
                    case "+": //turn left
                        turtle.HEADING += angle
                    case "-": //turn right
                        turtle.HEADING -= angle
                    case "|": //turn away
                        turtle.HEADING += 180.
                    case "$": //head due north
                        turtle.HEADING = 90.
                    case "(": //set arbitrary heading
                        turtle.HEADING, pos = getHeading(turtleCmds, pos)
                    case "[": //store status
                        stack.push(turtle)
                    case "]": //restore status
                        turtle    = stack.pop()
                        plotCmds += convert2Hpgl("f", turtle.X, turtle.Y)
                    case "{", "}": //start or end polygon mode
                        plotCmds += convert2Hpgl(symbol, turtle.X, turtle.Y)
                    default:  //remove all instances of a production variable
                        *turtleCmds = strings.NewReplacer(symbol, "").Replace(*turtleCmds)
                        pos--
                }
                pos++
            }
            plotCmds += ";"
            xMini, xMaxi, yMini, yMaxi = xMin, xMax, yMin, yMax
            return
           }
} //end func makeLogo2Hpgl
func makeConvert2Hpgl() func(symbol string, x, y float64) (pltCmd string) {
    penCmdPrev := ""
    return func(symbol string, x, y float64) (pltCmd string) {
            var penCmd = map[bool]string{true: "PD", false: "PU"} [symbol == "F"]
            switch {
                case symbol == "{": //terminate pen sequence & initiate polygon mode
                    penCmdPrev, pltCmd = "", ";\nPM0;\n"
                case symbol == "}": //terminate pen sequence, end polygon mode, fill & edge buffered polygon
                    penCmdPrev, pltCmd = "", ";\nPM2;EP;FP;\n"
                case penCmdPrev == penCmd: //continue
                    pltCmd = fmt.Sprintf(",%f,%f", x, y)
                case penCmdPrev != "": //terminate pen sequence & initiate new one
                    penCmdPrev, pltCmd = penCmd, fmt.Sprintf(";\n%s%f,%f", penCmd, x, y)
                default: //initiate pen sequence
                    penCmdPrev, pltCmd = penCmd, fmt.Sprintf("%s%f,%f", penCmd, x, y)
            }
            return
           }
} //end func makeConvert2Hpgl
func validFgColor(fgColor string) bool {
    if fgColor == "" {
        return false
    } else if strings.Contains(fgColor, "#") {
        matched, _ := regexp.MatchString("^#[a-fA-F0-9]{6}$", fgColor)
        return matched
    }
    _, ok := _colorNames[fgColor]
    return ok
} //end func validFgColor
////Context searches
func getContextLHS(io chan string) {
    for lhs := range io {
        context := ""
        pos     := len(lhs) - 1
        searchLoop: for pos >= 0 {
            symbol := string(lhs[pos])
            switch symbol {
                case "F", "+", "-", "$": //ignore constant
                case "[":                //keep looking for parent's symbol
                case "]":                //skip over branch symbol(s)
                    unmatchedBrakets := 1
                    for unmatchedBrakets != 0 {
                        pos--
                        char := string(lhs[pos])
                        if char == "[" { unmatchedBrakets-- }
                        if char == "]" { unmatchedBrakets++ }
                    }
                case "0", "1":
                    context = symbol
                    break searchLoop
                default:
                    halt(fmt.Sprintf("the symbol '%s' is not supported", symbol))
            }
            pos--
        }
        io<- context
    }
    return
} //end coroutine getContextLHS
func getContextRHS(io chan string) {
    for rhs := range io {
        context := ""
        pos     := 0
        searchLoop: for pos < len(rhs) {
            symbol := string(rhs[pos])
            switch symbol {
                case "F", "+", "-", "$": //ignore constant
                case "]":                //reached top of branch
                    break searchLoop
                case "[":                //skip over side branch symbol(s)
                    unmatchedBrakets := 1
                    for unmatchedBrakets != 0 {
                        pos++
                        char := string(rhs[pos])
                        if char == "[" { unmatchedBrakets++ }
                        if char == "]" { unmatchedBrakets-- }
                    }
                case "0", "1":
                    context = symbol
                    break searchLoop
                default:
                    halt(fmt.Sprintf("the symbol '%s' is not supported", symbol))
            }
            pos++
        }
        io<- context
    }
    return
} //end coroutine getContextRHS
//===== Copyright (c) 2016 Yves Beaudoin - All rights reserved - MIT LICENSE (MIT) - Email: webpraxis@gmail.com ================
//end of Package lsystems

# lsystems

Go package for processing and drawing Lindenmayer systems (L-systems) using gnuplot and HP-GL/2 graphics.

## Install

To install the package and its examples:
```sh
go get -u github.com/ybeaudoin/go-lsystems
```

The package imports binet's gnuplot package which can be installed as follows:
```sh
go get -u bitbucket.org/binet/go-gnuplot/pkg/gnuplot
```
It in turn requires that a gnuplot executable be installed and be findable via the environment path statement.
See http://www.gnuplot.info/download.html for available versions.

## At a glance

The package exports the following:

 * Variable
   * `TurtleCmds string`  
     Generated turtle-graphics commands
 * Functions:
   * `Deterministic(order int, axiom string, rules *strings.Replacer)`  
     Generates the required turtle commands for the specified deterministic and context-free production parameters.
   * `Stochastic(order int, axiom string, rules []string, weights []int)`  
     Generates the required turtle commands for the specified stochastic and context-free production parameters.
   * `HogewegHesper(order int, axiom string, rules map[string]string)`  
     Generates the required turtle commands for the specified Hogeweg and Hesper production parameters as
     described in http://algorithmicbotany.org/papers/abop/abop.pdf
   * `EncodeBgColorName(bgColorName string) string`  
     Encodes a color name into an hex string, prefixed with the character "x", for use as the specification
     of a gnuplot terminal's background color.
   * `Plot(angle float64, terminalCmd, outputCmd, plotTitle, lineColor string, cmdsFile ...string)`  
     Plots the latest generated turtle commands with the given parameters using gnuplot. The result will be
     isometrically scaled and centered. The underlying gnuplot commands can be saved optionally to a text file.
   * `MultiPlot(turtleCmds []string, turtleAngles []float64, terminalCmd, outputCmd, plotTitle string, labels []string, lineColor string, cmdsFile ...string)`  
     Plots a set of turtle commands with the given parameters using gnuplot. The result will be anisometrically scaled
     with the subplots generated left to right. The underlying gnuplot commands can be saved optionally to a text file.
   * `HpglPlot(angle float64, plotTitle string, penWidth float64, hpglPath string)`  
     Converts the latest generated turtle commands with the given parameters to an HP-GL/2 command set.
     The resulting plot will be isometrically scaled and centered.
   * `HpglMultiPlot(turtleCmds []string, turtleAngles []float64, plotTitle string, labels []string, penWidth float64, hpglPath string)`  
     Converts a set of turtle commands with the given parameters to an HP-GL/2 command set.
     The resulting plot will be anisometrically scaled with the subplots generated left to right in landscape mode.

All plot routines auto scale to achieve the best fit possible given the canvas or media size. The HP-GL/2 functions are provided
for users not having any joy with older versions of gnuplot's hpgl-supported terminals and newer compliant output devices.

Both `Plot` and `MultiPlot` offer the option of saving the gnuplot commands to a file. This can facilitate debugging the terminal
and output declarations by allowing the user to feed the commands directly to the gnuplot executable, viz. `gnuplot debug.cmds`,
and then view the error messages.

## L-system symbols

 * Variables  
   Any symbol that does not conflict with the constants below,
 * Constants  
   **F f + - | $ ( ) \[ \] { }** with the following turtle-graphics interpretations:  
   **F** means "move forward drawing a line",  
   **f** means "move forward without drawing a line",  
   **+** means "turn left",  
   **-** means "turn right",  
   **|** means "turn away" (adds 180 degrees to the heading),  
   **$** sets the heading to 90 degrees (default initial heading is 0),  
   **(** marks the start of a new heading declaration in degrees,  
   **)** marks the end of a heading declaration,  
   **\[** starts a branch; saves the current turtle's status,  
   **\]** ends a branch; restores the turtle's status with the last saved value (as in a LIFO stack),  
   **{** starts filled polygon mode (line segments define the edges),  
   **}** ends polygon mode.  
   All other symbols will be ignored during drawing.

## MIT License

Copyright (c) 2016 Yves Beaudoin webpraxis@gmail.com

See the file LICENSE for copying permission.


















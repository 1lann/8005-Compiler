# eightc
Compiles a C like language into instructions for the 8005 microprocessor. For the UNSW HS1917 Computing course.

## Install
Install Go, setup your $GOPATH and run
`go get github.com/1lann/eightc/eightc`

## Usage
`eightc [file name]`

## How to use
See [test-programs/](test-programs/) for examples.

A list of internal functions:
- `:printChar` prints r0 as an ASCII character
- `:printInt` prints r0 as an integer
- `:swap` swaps r0 and r1
- `:ring` rings the bell

## Notices
- Note that a lot of functions can only be performed on r0, not r1.
- r0 is used a lot internally and will most likely not persist through loops, function calls,
and exiting if statements. Use r1 instead which is not used internally by the compiler when compiling.

## License
eightc is licensed under the MIT license, see [LICENSE.txt](LICENSE.txt)

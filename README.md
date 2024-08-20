# Cubez

Cubez is a 3d physics library written in the [Go][golang] programming language. It is
 mostly a port of [cyclone-physics][cyclone] by Ian Millington and using his book
"Game Physics Engine Design" as a reference.

## Current Features

* Full 3d rigid body real-time physics simulation suitable for games; meaning
  both linear velocity as well as angular velocity are calculated.
* Collision detection between collider primitives.
* Primitives supported: planes, spheres, cubes.
* Math library defaults to 64-bit floats but can easily be tuned down to 32-bit.

## Examples

Ballistic: shoot spheres at a cube by pressing the space bar.

![ballistic][ballistic_ss]

Cubedrop: hit the space bar to drop some cubes onto the ground

![cubedrop][cubedrop_ss]

## OS Support

Cubez is known to work on the following:

* Windows 7 x64 with mingw-w64 (see this [tutorial][am_mingw64] if necessary)
* Linux (Ubuntu 14.04)
* MacOS (14.X)

At present, I suspect it should work on any Windows or Linux 64-bit system for which
there is an acceptable Go x64 and gcc x64 compiler set available.

Support for 32-bit systems is untested.

## Dependencies

The only dependency on the core `cubez` package is the math package included in `cubez`.

For the examples, you will need GLFW 3.3.x installed on your system, and you will need
to install the go-gl project's [gl][gogl_gl], [glfw][gogl_glfw] and [mathgl][gogl_mgl]
libraries. Your system will also need to be OpenGL 3.3 capable.

The examples use a basic OpenGL *framework-in-a-file* inspired by [@tbogdala](https://github.com/tbogdala)'s graphics engine
called [fizzle][fizzle]. This way the full [fizzle][fizzle] library is not a dependency.

## Installation

Clone or download this project, then go to the examples directory and download the dependencies.

```bash
cd cubez/examples
go mod download
```

Each example can then be run from its folder within `examples`:

```bash
cd cubez/examples/ballistic
go run ballistic.go
```

```bash
cd cubez/examples/cubedrop
go run cubedrop.go
```

## Documentation

Currently, you'll have to use godoc to read the API documentation and check
out the examples to figure out how to use the library.


## License

Cubez is released under the BSD license. See the [LICENSE][license-link] file for more details.


[golang]: https://golang.org/
[license-link]: https://raw.githubusercontent.com/tbogdala/cubez/master/LICENSE
[cyclone]: https://github.com/idmillington/cyclone-physics
[fizzle]: https://github.com/tbogdala/fizzle
[gogl_gl]: https://github.com/go-gl/gl
[gogl_glfw]: https://github.com/go-gl/glfw
[gogl_mgl]: https://github.com/go-gl/mathgl
[am_mingw64]: http://animal-machine.com/blog/150723_mingw-w64_and_Go.md

[ballistic_ss]: https://raw.githubusercontent.com/tbogdala/cubez/master/examples/screenshots/ballistic-150912.jpg
[cubedrop_ss]: https://raw.githubusercontent.com/tbogdala/cubez/master/examples/screenshots/cubedrop-150912.jpg

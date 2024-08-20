// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package main

import (
	gl "github.com/go-gl/gl/v3.3-core/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/rverpillot/cubez"
	ex "github.com/rverpillot/cubez/examples"
	m "github.com/rverpillot/cubez/math"
)

var (
	app *ex.ExampleApp

	cube      *ex.Entity
	backboard *ex.Entity
	bullets   []*ex.Entity

	colorShader uint32
	groundPlane *cubez.CollisionPlane
	ground      *ex.Renderable
)

// update object locations
func updateObjects(delta float64) {
	// for now there's only one box to update
	body := cube.Collider.GetBody()
	body.Integrate(m.Real(delta))
	cube.Collider.CalculateDerivedData()

	// for now we hack in the position and rotation of the collider into the renderable
	ex.SetGlVector3(&cube.Node.Location, &body.Position)
	ex.SetGlQuat(&cube.Node.LocalRotation, &body.Orientation)

	for _, bullet := range bullets {
		bulletBody := bullet.Collider.GetBody()
		bulletBody.Integrate(m.Real(delta))
		bullet.Collider.CalculateDerivedData()
		ex.SetGlVector3(&bullet.Node.Location, &bulletBody.Position)
		ex.SetGlQuat(&bullet.Node.LocalRotation, &bulletBody.Orientation)
	}
}

// see if any of the rigid bodys contact
func generateContacts(delta float64) (bool, []*cubez.Contact) {
	var returnFound bool

	// create the ground plane
	groundPlane := cubez.NewCollisionPlane(m.Vector3{0.0, 1.0, 0.0}, 0.0)

	// see if we have a collision with the ground
	found, contacts := cubez.CheckForCollisions(cube.Collider, groundPlane, 0.9, 0.1, nil)
	if found {
		returnFound = true
	}
	// see if there's a collision against the backboard
	found, contacts = cubez.CheckForCollisions(cube.Collider, backboard.Collider, 0.9, 0.1, contacts)
	if found {
		returnFound = true
	}

	// run collision checks on bullets
	for _, bullet := range bullets {
		// check against the ground
		found, contacts = cubez.CheckForCollisions(bullet.Collider, groundPlane, 0.9, 0.1, contacts)
		if found {
			returnFound = true
		}

		// check against the cube
		found, contacts = cubez.CheckForCollisions(cube.Collider, bullet.Collider, 0.9, 0.1, contacts)
		if found {
			returnFound = true
		}

		// check against the backboard
		found, contacts = cubez.CheckForCollisions(backboard.Collider, bullet.Collider, 0.9, 0.1, contacts)
		if found {
			returnFound = true
		}

		// check against other bullets
		for _, bullet2 := range bullets {
			if bullet2 == bullet {
				continue
			}
			found, contacts = cubez.CheckForCollisions(bullet2.Collider, bullet.Collider, 0.9, 0.1, contacts)
			if found {
				returnFound = true
			}
		}
	}

	return returnFound, contacts
}

func updateCallback(delta float64) {
	updateObjects(delta)
	foundContacts, contacts := generateContacts(delta)
	if foundContacts {
		cubez.ResolveContacts(len(contacts)*8, contacts, m.Real(delta))
	}
}

func renderCallback(delta float64) {
	gl.Viewport(0, 0, int32(app.Width), int32(app.Height))
	gl.ClearColor(0.196078, 0.6, 0.8, 1.0) // some pov-ray sky blue
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// make the projection and view matrixes
	projection := mgl.Perspective(mgl.DegToRad(60.0), float32(app.Width)/float32(app.Height), 1.0, 200.0)
	view := app.CameraRotation.Mat4()
	view = view.Mul4(mgl.Translate3D(-app.CameraPos[0], -app.CameraPos[1], -app.CameraPos[2]))

	// draw the cube
	cube.Node.Draw(projection, view)

	// draw all of the bullets
	for _, bullet := range bullets {
		bullet.Node.Draw(projection, view)
	}

	// draw the backboard
	backboard.Node.Draw(projection, view)

	// draw the ground
	ground.Draw(projection, view)

	//time.Sleep(10 * time.Millisecond)
}

func main() {
	app = ex.NewApp()
	app.InitGraphics("Ballistic", 800, 600)
	app.SetKeyCallback(keyCallback)
	app.OnRender = renderCallback
	app.OnUpdate = updateCallback
	defer app.Terminate()

	// compile the shaders
	var err error
	colorShader, err = ex.LoadShaderProgram(ex.DiffuseColorVertShader, ex.DiffuseColorFragShader)
	if err != nil {
		panic("Failed to compile the shader! " + err.Error())
	}

	// create the ground plane
	groundPlane = cubez.NewCollisionPlane(m.Vector3{0.0, 1.0, 0.0}, 0.0)

	// make a ground plane to draw
	ground = ex.CreatePlaneXZ(-500.0, 500.0, 500.0, -500.0, 1.0)
	ground.Shader = colorShader
	ground.Color = mgl.Vec4{0.6, 0.6, 0.6, 1.0}

	// create a test cube to render
	cubeNode := ex.CreateCube(-1.0, -1.0, -1.0, 1.0, 1.0, 1.0)
	cubeNode.Shader = colorShader
	cubeNode.Color = mgl.Vec4{1.0, 0.0, 0.0, 1.0}

	// create the collision box for the the cube
	var cubeMass m.Real = 8.0
	var cubeInertia m.Matrix3
	cubeCollider := cubez.NewCollisionBox(nil, m.Vector3{1.0, 1.0, 1.0})
	cubeCollider.Body.Position = m.Vector3{0.0, 5.0, 0.0}
	cubeCollider.Body.SetMass(cubeMass)
	cubeInertia.SetBlockInertiaTensor(&cubeCollider.HalfSize, cubeMass)
	cubeCollider.Body.SetInertiaTensor(&cubeInertia)
	cubeCollider.Body.CalculateDerivedData()
	cubeCollider.CalculateDerivedData()

	// make the entity out of the renerable and collider
	cube = ex.NewEntity(cubeNode, cubeCollider)

	// make a slice of entities for bullets
	bullets = make([]*ex.Entity, 0, 16)

	// make the backboard to bound the bullets off of
	backboardNode := ex.CreateCube(-0.5, -2.0, -0.25, 0.5, 2.0, 0.25)
	backboardNode.Shader = colorShader
	backboardNode.Color = mgl.Vec4{0.25, 0.2, 0.2, 1.0}
	backboardCollider := cubez.NewCollisionBox(nil, m.Vector3{0.5, 2.0, 0.25})
	backboardCollider.Body.Position = m.Vector3{0.0, 2.0, -10.0}
	backboardCollider.Body.SetInfiniteMass()
	backboardCollider.Body.CalculateDerivedData()
	backboardCollider.CalculateDerivedData()
	ex.SetGlVector3(&backboardNode.Location, &backboardCollider.Body.Position)

	// make the backboard entity
	backboard = ex.NewEntity(backboardNode, backboardCollider)

	// setup the camera
	app.CameraPos = mgl.Vec3{-3.0, 3.0, 15.0}
	app.CameraRotation = mgl.QuatLookAtV(
		mgl.Vec3{-3.0, 3.0, 15.0},
		mgl.Vec3{0.0, 1.0, 0.0},
		mgl.Vec3{0.0, 1.0, 0.0})

	gl.Enable(gl.DEPTH_TEST)
	app.RenderLoop()
}

func fire() {
	var mass m.Real = 1.5
	var radius m.Real = 0.2

	// create a test sphere to render
	bullet := ex.CreateSphere(float32(radius), 16, 16)
	bullet.Shader = colorShader
	bullet.Color = mgl.Vec4{0.2, 0.2, 1.0, 1.0}

	// create the collision box for the the bullet
	bulletCollider := cubez.NewCollisionSphere(nil, radius)
	bulletCollider.Body.Position = m.Vector3{0.0, 1.5, 20.0}

	var cubeInertia m.Matrix3
	var coeff m.Real = 0.4 * mass * radius * radius
	cubeInertia.SetInertiaTensorCoeffs(coeff, coeff, coeff, 0.0, 0.0, 0.0)
	bulletCollider.GetBody().SetInertiaTensor(&cubeInertia)

	bulletCollider.Body.SetMass(mass)
	bulletCollider.Body.Velocity = m.Vector3{0.0, 0.0, -40.0}
	bulletCollider.Body.Acceleration = m.Vector3{0.0, -2.5, 0.0}

	bulletCollider.Body.CalculateDerivedData()
	bulletCollider.CalculateDerivedData()

	e := ex.NewEntity(bullet, bulletCollider)
	bullets = append(bullets, e)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// Key W == close app
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	if key == glfw.KeySpace && action == glfw.Press {
		fire()
	}
}

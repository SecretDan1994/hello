package main

import (
	"fmt"
	// "github.com/gen2brain/raylib-go/raygui"
	"math"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/gen2brain/raylib-go/raymath"
)

const (
	screenWidth  = 800
	screenHeight = 450
	PLAYER_MAX_LIFE = 5
	LINES_OF_BRICKS = 5
	BRICKS_PER_LINE = 20
	SENSITIVITY = 0.1
)

var (
	firstMouse = true
	lastX float32 = screenWidth/2
	lastY float32  = screenHeight/2
	yaw float32  = -90.0
	pitch float32  = 0.0
)

type Player struct {
	position rl.Vector2
	size     rl.Vector2
	life     int
}

type Ball struct {
	position rl.Vector2
	speed    rl.Vector2
	radius   float32
	active   bool
}

type Brick struct {
	position rl.Vector2
	active   bool
}


type Game struct {
	gameOver  bool
	pause     bool
	player    Player
	ball      Ball
	brick     [LINES_OF_BRICKS][BRICKS_PER_LINE]Brick
	brickSize rl.Vector2
}

func radians(degrees float32) float64{
	return float64(degrees * (math.Pi / 180))
}

func mouse_callback(camera rl.Camera3D) rl.Vector3{
	var mousepos = rl.GetMousePosition()
	var xpos = mousepos.X
	var ypos = mousepos.Y
	if firstMouse{
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}
	var xoffset = xpos - lastX
	var yoffset = lastY - ypos

	lastX = xpos
	lastY = ypos

	xoffset *= SENSITIVITY
	yoffset *= SENSITIVITY

	yaw   += xoffset
	pitch += yoffset

	if pitch > 89.0 {
		pitch = 89.0
	}
	if pitch < -89.0 {
		pitch = -89.0
	}

	var front rl.Vector3
	front.X = float32(math.Cos(radians(yaw)) * math.Cos(radians(pitch)))
	front.Y = float32(math.Sin(radians(pitch)))
	front.Z = float32(math.Sin(radians(yaw)) * math.Cos(radians(pitch)))
	raymath.Vector3Normalize(&front)

	if rl.IsMouseButtonPressed(rl.MouseRightButton){
		fmt.Printf("\nxoffset: %f | yoffset: %f", xoffset, yoffset)
		fmt.Printf("\nyaw: %f | pitch: %f", yaw, pitch)
		fmt.Printf("\nfrontX: %f | frontY: %f | frontZ: %f", front.X, front.Y, front.Z)
	}

	return front
}  

func main() {
	rl.InitWindow(screenWidth, screenHeight, "test game")
	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(-10.0, 20.0, -10.0)
	camera.Target = rl.NewVector3(0,0, 0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 45.0
	camera.Type = rl.CameraPerspective

	cubePosition := rl.NewVector3(0.0, 1.0, 0.0)
	cubeSize := rl.NewVector3(2.0, 2.0, 2.0)

	var ray rl.Ray

	collision := false

	rl.SetCameraMode(camera, rl.CameraFree) // Set a free camera mode

	rl.SetTargetFPS(60)

	// Initialize 
	// var oldRay = rl.GetMouseRay(rl.GetMousePosition(), camera)

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera) // Update camera

		ray = rl.GetMouseRay(rl.GetMousePosition(), camera)

		// var vec = mouse_callback(camera)
		// camera.Target = vec
		// rl.SetMousePosition(int(vec.X), int(vec.Y))
		// Moving mouse along x-axis.
		// if ray.Direction.X != oldRay.Direction.X && ray.Direction.Y != oldRay.Direction.Y && ray.Direction.Z != oldRay.Direction.Z{
		// 	if ray.Direction.X < oldRay.Direction.X {				//Turning left
		// 		camera.Target.X += SENSITIVITY
		// 		camera.Target.Y *= SENSITIVITY
		// 		camera.Target.Z *= SENSITIVITY				
		// 	} else {							//Turning right
		// 		camera.Target.X += SENSITIVITY
		// 		camera.Target.Y *= SENSITIVITY
		// 		camera.Target.Z *= SENSITIVITY
		// 	}
		// 	// // Moving mouse along y-axis.
		// 	// if ray.Direction.Y > oldRay.Direction.Y {
		// 	// 	camera.Target.Y += SENSITIVITY	// Turning up
		// 	// } else {
		// 	// 	camera.Target.Y -= SENSITIVITY	// Turning down
		// 	// }
		// 	// // Moving mouse along y-axis.
		// 	// if ray.Direction.Z > oldRay.Direction.Z {
		// 	// 	camera.Target.Z += SENSITIVITY	// Turning up
		// 	// } else {
		// 	// 	camera.Target.Z -= SENSITIVITY	// Turning down
		// 	// }
		// 	oldRay = ray
		// }
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			// Check collision between ray and box
			min := rl.NewVector3(cubePosition.X-cubeSize.X/2, cubePosition.Y-cubeSize.Y/2, cubePosition.Z-cubeSize.Z/2)
			max := rl.NewVector3(cubePosition.X+cubeSize.X/2, cubePosition.Y+cubeSize.Y/2, cubePosition.Z+cubeSize.Z/2)
			collision = rl.CheckCollisionRayBox(ray, rl.NewBoundingBox(min, max))
		}

		//Combo keys.

		// if{

		// }
		//Individually Pressed.
		// else{	
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Position.X-=10
			camera.Position.Z+=10
			var oldz = camera.Target.Z
			camera.Target.Z = camera.Target.X
			camera.Target.X = oldz
			// camera.Position.Y+=1
			// camera.Position.Z+=1
		}
		if rl.IsKeyDown(rl.KeyUp) {
			camera.Position.Y-=10
			camera.Position.X-=10
		}
		if rl.IsKeyDown(rl.KeyDown) {
			camera.Position.Y+=10
			camera.Position.Z+=10
		}
		if rl.IsKeyDown(rl.KeyRight) {
			camera.Position.X+=10			
			camera.Position.Z-=10
			// camera.Position.Y+=1
			// camera.Position.Z-=1
		}

		if rl.IsKeyDown(rl.KeyD) {
			camera.Target.X-=1
			camera.Target.Z+=1
			// camera.Target.Y+=1
			// camera.Target.Z+=1
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.Target.X-=1
			camera.Target.Y+=1
			camera.Target.Z-=1
		}
		if rl.IsKeyDown(rl.KeyW) {
			camera.Target.X+=1
			camera.Target.Y-=1
			camera.Target.Z+=1
		}
		if rl.IsKeyDown(rl.KeyA) {
			camera.Target.X+=1			
			camera.Target.Z-=1
			// camera.Target.Y+=1
			// camera.Target.Z-=1
		}			
		// }




		/* Formula I got online for camera direction and position.
		cameraDirection = cameraRot * glm::vec3(-1.0f, 0.0f, 0.0f);
		cameraUp = cameraRot * glm::vec3(0.0f, 1.0f, 0.0f);
		cameraPos = cameraTarget - cameraDistance * cameraDirection;*/
		// if camera.Target.X == 0 && camera.Target.Y == 0 && camera.Target.Y == 0{
		// 	camera.Target.X = 1000 * ray.Direction.X
		// 	camera.Target.Y = 1000 * ray.Direction.Y
		// 	camera.Target.Z = 1000 * ray.Direction.Z
		// } else {
		// 	camera.Target.X = ray.Direction.X
		// 	camera.Target.Y = ray.Direction.Y
		// 	camera.Target.Z = ray.Direction.Z
		// }	

		matrix := rl.GetCameraMatrix(camera)
		
		if rl.IsMouseButtonPressed(rl.MouseRightButton){
			fmt.Println()
			fmt.Printf("M0: %f | M4: %f | M8: %f | M12: %f\n",matrix.M0,matrix.M4,matrix.M8,matrix.M12)
			fmt.Printf("M1: %f | M5: %f | M9: %f | M13: %f\n",matrix.M1,matrix.M5,matrix.M9,matrix.M13)
			fmt.Printf("M2: %f | M6: %f | M10: %f | M14: %f\n",matrix.M2,matrix.M6,matrix.M10,matrix.M14)
			fmt.Printf("M3: %f | M7: %f | M11: %f | M15: %f\n",matrix.M3,matrix.M7,matrix.M11,matrix.M15)
			fmt.Printf("Camera Position X: %f | Camera Position Y: %f| Camera Position Z: %f", camera.Position.X, camera.Position.Y, camera.Position.Z)
			fmt.Printf("Camera Target X: %f | Camera Target Y: %f| Camera Target Z: %f", camera.Target.X, camera.Target.Y, camera.Target.Z)
			fmt.Printf("Mouse Direction X: %f | Mouse Direction Y: %f| Mouse Direction Z: %f", ray.Direction.X, ray.Direction.Y, ray.Direction.Z)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(camera)

		rl.DrawCube(rl.NewVector3(0,0,2.5),15,5,0, rl.Gold)			//Right Wall
		rl.DrawCube(rl.NewVector3(0,0,-2.5),15,5,0, rl.Red)			//Left Wall
		rl.DrawCube(rl.NewVector3(7.5,0,0),0,5,5, rl.Blue)			//Front Wall
		rl.DrawCube(rl.NewVector3(-7.5,0,0),0,5,5, rl.Green)		//Back Wall
		rl.DrawCube(rl.NewVector3(0,0,0),15,0,5, rl.Beige)			//Floor
		rl.DrawCube(rl.NewVector3(0,2.5,0),15,0,5, rl.Purple)		//Ceiling


		if collision {
			rl.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Red)
			rl.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Maroon)

			rl.DrawCubeWires(cubePosition, cubeSize.X+0.2, cubeSize.Y+0.2, cubeSize.Z+0.2, rl.Green)
		} else {
			rl.DrawCube(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.Gray)
			rl.DrawCubeWires(cubePosition, cubeSize.X, cubeSize.Y, cubeSize.Z, rl.DarkGray)
		}

		// rl.DrawRay(ray, rl.Maroon)

		// rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.DrawText("Try selecting the box with mouse!", 240, 10, 20, rl.DarkGray)

		if collision {
			rl.DrawText("BOX SELECTED", (screenWidth-rl.MeasureText("BOX SELECTED", 30))/2, int32(float32(screenHeight)*0.1), 30, rl.Green)
		}

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
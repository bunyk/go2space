
package main

import (
	"fmt"
    "time"
	"os"
	"image"
    "math"
	"math/rand"
	_ "image/png"
	_ "image/jpeg"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const ROCKET_THRUST = 15.0
const ROTATION_SPEED = 5.0
const WIDHT = 800
const HEIGHT = 600
const STAR_COUNT = 100

func run() {
    cfg := pixelgl.WindowConfig{
		Title:  "Go 2 space!",
		Bounds: pixel.R(0, 0, WIDHT, HEIGHT),
        VSync: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Blue
	imd.EndShape = imdraw.SharpEndShape
	imd.Push(pixel.V(-100.0, 0.0), pixel.V(100.0, 0.0))
	imd.Push(pixel.V(0.0, -100.0), pixel.V(0.0, 100.0))
	imd.Line(3)

	rocket := Rocket{
		angle: math.Pi /2,
		pos: pixel.V(0, 100),
		vel: pixel.ZV,
		sprite: pic2sprite("images/rocket.png"),
	}

	var (
		stars    [STAR_COUNT]int
		matrices [STAR_COUNT]pixel.Matrix
	)
	stars_sprites := loadStars()
	for i := 0; i < STAR_COUNT; i++ {
		stars[i] = rand.Intn(len(stars_sprites))
		x := rand.NormFloat64() * 10000
		y := rand.NormFloat64() * 10000
		fmt.Println(x, y)
		matrices[i] = pixel.IM.Scaled(pixel.ZV, 4).Moved(pixel.V(x, y))
	}


	last := time.Now()
	zoom := 0.0
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// cam_pos := win.Bounds().Center().Sub(rocket.pos)
		cam := pixel.IM.
			Rotated(rocket.pos, -rocket.angle + math.Pi / 2).
			Moved(pixel.ZV.Sub(rocket.pos)).
			Moved(win.Bounds().Center())

        // matrix = matrix.Rotated(pixel.ZV, r.angle - math.Pi / 2)
        // matrix = matrix.Moved(r.pos)

		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyEscape) {
			break
		}
		if win.Pressed(pixelgl.KeyUp) {
			zoom += 1.0
		}
		if win.Pressed(pixelgl.KeyDown) {
			zoom -= 1.0
		}
		if win.Pressed(pixelgl.KeyLeft) {
			rocket.Turn(dt)
		}
		if win.Pressed(pixelgl.KeyRight) {
			rocket.Turn(-dt)
		}
		rocket.Move(dt)


        win.Clear(colornames.Black)

		for i, star := range stars {
			stars_sprites[star].Draw(win, matrices[i])
		}

        rocket.Draw(win)
		imd.Draw(win)

		win.Update()
	}
}

type Rocket struct {
	angle float64
	pos pixel.Vec
	vel pixel.Vec
	sprite *pixel.Sprite
}
func (r *Rocket) Turn(angle float64) {
	r.angle += angle * ROTATION_SPEED
}

func (r *Rocket) Speed() float64 {
	return math.Sqrt(sqr(r.vel.X) + sqr(r.vel.Y))
}

func (r *Rocket) Move(dt float64) {
	r.vel.X += math.Cos(r.angle) * dt * ROCKET_THRUST
	r.vel.Y += math.Sin(r.angle) * dt * ROCKET_THRUST

	r.pos.X += r.vel.X * dt
	r.pos.Y += r.vel.Y * dt
}

func (r *Rocket) Draw(win *pixelgl.Window) {
        matrix := pixel.IM.Scaled(pixel.ZV, 0.05)
        matrix = matrix.Rotated(pixel.ZV, r.angle - math.Pi / 2)
        matrix = matrix.Moved(r.pos)

        r.sprite.Draw(win, matrix)
}

func sqr(x float64) float64 {
	return x*x
}

func loadStars() (stars [5]*pixel.Sprite) {
	for i, fn := range [...]string{
		"crab_nebula.jpg",
		"fomalhaut.jpg",
		"galaxy.jpg",
		"mira.jpg",
		"trappist1.jpg",
	} {
		fmt.Println("Loading", fn)
		stars[i] = pic2sprite("images/" + fn)
	}
	return
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func pic2sprite(fn string) *pixel.Sprite {
	pic, err := loadPicture(fn)
	if err != nil {
		panic(err)
	}
	return pixel.NewSprite(pic, pic.Bounds())
}

func main() {
	pixelgl.Run(run)
}

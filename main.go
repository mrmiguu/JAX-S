package main

import (
	"math"

	"github.com/gopherjs/gopherjs/js"
)

var (
	phaser *js.Object

	game,
	gameDebug,
	gameLoad,
	gameAdd,
	gameWorld,
	gamePhysics,
	gameCamera,
	gameScale,
	gameInput *js.Object

	bg    *js.Object
	bgs   []*js.Object
	floor *js.Object

	player,
	body,
	bodyBody,
	bodyBodyVelocity *js.Object

	cursors   *js.Object
	spotlight *js.Object
	bgms      []*js.Object

	walk            []*js.Object
	walking         bool
	up, left, right bool

	started     bool
	level       int
	chasing     float64
	pendingPlay bool
)

func init() {
	document := js.Global.Get("document")

	style := document.Get("body").Get("style")
	style.Set("background", "#000000")
	style.Set("margin", 0)

	phaserjs := document.Call("createElement", "script")
	phaserjs.Call("setAttribute", "src", "phaser.min.js")
	phaserjs.Set("onload", func() {
		phaser = js.Global.Get("Phaser")
		game = phaser.Get("Game").New(
			1280, 720,
			phaser.Get("AUTO"),
			"JAX'S",
			js.M{
				"preload": preload,
				"create":  create,
				"update":  update,
				"render":  render,
			},
		)
	})
	document.Get("head").Call("appendChild", phaserjs)
}

func preload() {
	gameLoad = game.Get("load")
	gameAdd = game.Get("add")
	gameWorld = game.Get("world")
	gamePhysics = game.Get("physics")
	gameCamera = game.Get("camera")
	gameDebug = game.Get("debug")

	game.Get("canvas").Set("oncontextmenu", func(e *js.Object) { e.Call("preventDefault") })

	gameScale = game.Get("scale")
	gameScale.Set("scaleMode", phaser.Get("ScaleManager").Get("SHOW_ALL"))
	gameScale.Set("fullScreenScaleMode", phaser.Get("ScaleManager").Get("SHOW_ALL"))
	gameScale.Set("pageAlignHorizontally", true)
	gameScale.Set("pageAlignVertically", true)

	gameLoad.Call("image", "bg", "assets/bg.jpg")

	gameLoad.Call("image", "bg1", "assets/bg1.jpg")
	gameLoad.Call("image", "bg2", "assets/bg2.jpg")
	gameLoad.Call("image", "bg3", "assets/bg3.jpg")
	gameLoad.Call("image", "bg4", "assets/bg4.jpg")
	gameLoad.Call("image", "bg5", "assets/bg5.jpg")
	gameLoad.Call("image", "bg6", "assets/bg6.jpg")
	gameLoad.Call("image", "bg7", "assets/bg7.jpg")
	gameLoad.Call("image", "bg8", "assets/bg8.jpg")
	gameLoad.Call("image", "bg9", "assets/bg9.jpg")
	gameLoad.Call("image", "bg10", "assets/bg10.jpg")

	gameLoad.Call("image", "floor", "assets/floor.png")
	gameLoad.Call("image", "spotlight", "assets/spotlight.png")
	gameLoad.Call("spritesheet", "udlr", "assets/udlr.png", 171, 252, 4)

	gameLoad.Call("audio", "10", "assets/10.mp3")
	gameLoad.Call("audio", "20", "assets/20.mp3")
	gameLoad.Call("audio", "30", "assets/30.mp3")
	gameLoad.Call("audio", "40", "assets/40.mp3")
	gameLoad.Call("audio", "50", "assets/50.mp3")
	gameLoad.Call("audio", "60", "assets/60.mp3")
	gameLoad.Call("audio", "70", "assets/70.mp3")
	gameLoad.Call("audio", "80", "assets/80.mp3")
	gameLoad.Call("audio", "90", "assets/90.mp3")
	gameLoad.Call("audio", "96", "assets/96.mp3")
	gameLoad.Call("audio", "97", "assets/97.mp3")
	gameLoad.Call("audio", "98", "assets/98.mp3")
	gameLoad.Call("audio", "99", "assets/99.mp3")

	gameLoad.Call("audio", "walk0", "assets/walk0.mp3")
	gameLoad.Call("audio", "walk1", "assets/walk1.mp3")
	gameLoad.Call("audio", "walk2", "assets/walk2.mp3")
	gameLoad.Call("audio", "walk3", "assets/walk3.mp3")
	gameLoad.Call("audio", "walk4", "assets/walk4.mp3")
	gameLoad.Call("audio", "walk5", "assets/walk5.mp3")
	gameLoad.Call("audio", "walk6", "assets/walk6.mp3")
}

func create() {
	bg = gameAdd.Call("sprite", 0, 0, "bg")

	walk = []*js.Object{
		gameAdd.Call("audio", "walk0"),
		gameAdd.Call("audio", "walk1"),
		gameAdd.Call("audio", "walk2"),
		gameAdd.Call("audio", "walk3"),
		gameAdd.Call("audio", "walk4"),
		gameAdd.Call("audio", "walk5"),
		gameAdd.Call("audio", "walk6"),
	}
	for n, w := range walk {
		n := n
		w.Get("onStop").Call("add", func() {
			walking = up || left || right
			if walking {
				walk[(n+1)%len(walk)].Call("play")
			}
		})
	}

	bgms = []*js.Object{
		gameAdd.Call("audio", "10"),
		gameAdd.Call("audio", "20"),
		gameAdd.Call("audio", "30"),
		gameAdd.Call("audio", "40"),
		gameAdd.Call("audio", "50"),
		gameAdd.Call("audio", "60"),
		gameAdd.Call("audio", "70"),
		gameAdd.Call("audio", "80"),
		gameAdd.Call("audio", "90"),
		gameAdd.Call("audio", "96"),
		gameAdd.Call("audio", "97"),
		gameAdd.Call("audio", "98"),
		gameAdd.Call("audio", "99"),
	}
	for _, m := range bgms {
		m.Set("volume", 0)
		m.Get("onStop").Call("add", func() {
		})
		m.Get("onPlay").Call("add", func() {
		})
	}

	bgs = []*js.Object{
		gameAdd.Call("sprite", 0, 0, "bg1"),
		gameAdd.Call("sprite", 0, 0, "bg2"),
		gameAdd.Call("sprite", 0, 0, "bg3"),
		gameAdd.Call("sprite", 0, 0, "bg4"),
		gameAdd.Call("sprite", 0, 0, "bg5"),
		gameAdd.Call("sprite", 0, 0, "bg6"),
		gameAdd.Call("sprite", 0, 0, "bg7"),
		gameAdd.Call("sprite", 0, 0, "bg8"),
		gameAdd.Call("sprite", 0, 0, "bg9"),
		gameAdd.Call("sprite", 0, 0, "bg10"),
	}
	for _, b := range bgs {
		b.Set("visible", false)
	}

	floor = gameAdd.Call("sprite", 0, 0, "floor")
	floor.Set("visible", false)

	gameWorld.Call("setBounds", 0, 0, 1920, 1080)
	arcade := phaser.Get("Physics").Get("ARCADE")
	gamePhysics.Call("startSystem", arcade)

	player = gameAdd.Call("sprite", 0, 0, "udlr")
	player.Set("visible", false)
	player.Get("anchor").Set("x", 0.5)
	player.Get("anchor").Set("y", 1)

	gamePhysics.Call("enable", floor, arcade)
	floor.Get("body").Call("setSize", 1920, 780, 0, 0)

	body = gameAdd.Call("sprite", gameWorld.Get("centerX"), gameWorld.Get("height").Int()-250)
	gamePhysics.Call("enable", body, arcade)

	spotlight = gameAdd.Call("sprite", gameWorld.Get("centerX"), gameWorld.Get("centerY"), "spotlight")
	spotlight.Set("visible", false)
	spotlight.Get("anchor").Call("set", 0.5)

	gameInput = game.Get("input")
	cursors = gameInput.Get("keyboard").Call("createCursorKeys")
	gameCamera.Call("follow", player, phaser.Get("Camera").Get("FOLLOW_LOCKON"), 0.1, 0.1)

	bodyBody = body.Get("body")
	bodyBodyVelocity = bodyBody.Get("velocity")

	gameInput.Get("onDown").Call("add", func() {
		if started {
			return
		}

		gameScale.Call("startFullScreen")

		bgs[0].Set("visible", true)
		floor.Set("visible", true)
		player.Set("visible", true)
		spotlight.Set("visible", true)

		started = true
	})
}

func update() {
	if !started {
		return
	}

	gameCamera.Call("shake", float64(level)/1000, 250*float64(level)/float64(len(bgs)))

	bodyBodyVelocity.Set("x", 0)
	bodyBodyVelocity.Set("y", 0)

	spotlight.Set("x", body.Get("x"))
	spotlight.Set("y", body.Get("y"))
	player.Set("x", bodyBody.Get("x"))
	player.Set("y", bodyBody.Get("y"))

	up = cursors.Get("up").Get("isDown").Bool()
	left = cursors.Get("left").Get("isDown").Bool()
	right = cursors.Get("right").Get("isDown").Bool()

	startWalk := (up || left || right) && !walking

	bodyBody.Set("angularVelocity", 0)

	if left {
		bodyBody.Set("angularVelocity", -125)
	} else if right {
		bodyBody.Set("angularVelocity", 125)
	}

	O := body.Get("angle").Float()

	if -45 <= O && O <= 45 {
		player.Set("frame", 3)
	} else if -135 < O && O < -45 {
		player.Set("frame", 0)
	} else if -180 <= O && O <= -135 || 135 <= O && O <= 180 {
		player.Set("frame", 2)
	} else if 45 < O && O < 135 {
		player.Set("frame", 1)
	}

	dist := 250 / math.Max(float64(level)/3, 1)
	pt := phaser.Get("Point").New()

	if up {
		gamePhysics.Get("arcade").Call("velocityFromAngle", O, dist, pt)
	}

	if bodyBodyY() < 800 {
		pt.Set("y", math.Max(pt.Get("y").Float(), 0))
	}

	if bodyBodyY() > 1240 {
		pt.Set("y", math.Min(pt.Get("y").Float(), 0))

		// bgs[level].Set("visible", false)
		// if level < len(bgs)-1 {
		// 	startHorror()
		// }
		// level = int(math.Min(float64(level+1), float64(len(bgs)-1)))
		// bodyBodyY(800)
		// bgs[level].Set("visible", true)
	}

	if bodyBodyX() < 0 {
		pt.Set("x", math.Max(pt.Get("x").Float(), 0))

		bgs[level].Set("visible", false)
		if level < len(bgs)-1 {
			startHorror()
		}
		level = int(math.Min(float64(level+1), float64(len(bgs)-1)))
		bodyBodyX(1900)
		bgs[level].Set("visible", true)
	}

	if bodyBodyX() > 1900 {
		pt.Set("x", math.Min(pt.Get("x").Float(), 0))

		// bgs[level].Set("visible", false)
		// if level < len(bgs)-1 {
		// 	startHorror()
		// }
		// level = int(math.Min(float64(level+1), float64(len(bgs)-1)))
		// bodyBodyX(0)
		// bgs[level].Set("visible", true)
	}

	pt.Set("y", pt.Get("y").Float()/2)

	bodyBody.Set("velocity", pt)

	if startWalk {
		walking = true
		walk[0].Call("play")
	}
}

func render() {
	// gameDebug.Call("text", "x: "+strconv.Itoa(bodyBodyX()), 50, 50)
	// gameDebug.Call("text", "y: "+strconv.Itoa(bodyBodyY()), 50, 70)
	// gameDebug.Call("text", "vy: "+bodyBody.Get("velocity").Get("y").String(), 50, 90)
}

func startHorror() {
	bgms[int(math.Max(float64(level-1), 0))].Call("stop")
	bgms[level].Call("play")
	bgms[level].Call("fadeTo",
		fadeInMS(),
		volume(),
	)
}

func fadeInMS() int {
	return 100000 - 100000*level/len(bgs)
}

func volume() float64 {
	return float64(level+1) / float64(len(bgs))
}

func bodyBodyX(x ...int) int {
	if len(x) > 0 {
		bodyBody.Set("x", x[0])
	}
	return bodyBody.Get("x").Int()
}
func bodyBodyY(y ...int) int {
	if len(y) > 0 {
		bodyBody.Set("y", y[0])
	}
	return bodyBody.Get("y").Int()
}

func main() {}

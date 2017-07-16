package main

import "github.com/gopherjs/gopherjs/js"

var (
	phaser                                                      *js.Object
	game, gameLoad, gameAdd, gameWorld, gamePhysics, gameCamera *js.Object
	bg1                                                         *js.Object
	floor                                                       *js.Object
	player, playerBody                                          *js.Object
	cursors                                                     *js.Object
	spotlight                                                   *js.Object
	bgm10                                                       *js.Object
)

func init() {
	style := js.Global.Get("document").Get("body").Get("style")
	style.Set("margin", 0)

	phaser = js.Global.Get("Phaser")
	game = phaser.Get("Game").New(
		1280, 720,
		phaser.Get("AUTO"),
		"JAX'S",
		js.M{
			"preload": preload,
			"create":  create,
			"update":  update,
		},
	)
}

func preload() {
	gameLoad = game.Get("load")
	gameAdd = game.Get("add")
	gameWorld = game.Get("world")
	gamePhysics = game.Get("physics")
	gameCamera = game.Get("camera")

	scale := game.Get("scale")
	scale.Set("scaleMode", phaser.Get("ScaleManager").Get("SHOW_ALL"))
	scale.Set("fullScreenScaleMode", phaser.Get("ScaleManager").Get("SHOW_ALL"))
	scale.Set("pageAlignHorizontally", true)
	scale.Set("pageAlignVertically", true)

	gameLoad.Call("image", "bg1", "assets/bg1.jpg")
	gameLoad.Call("image", "floor", "assets/floor.png")
	gameLoad.Call("image", "spotlight", "assets/spotlight.png")
	gameLoad.Call("image", "me", "assets/me.png")
	gameLoad.Call("audio", "10", "assets/10.mp3")
}

func create() {
	bgm10 = gameAdd.Call("audio", "10")
	bgm10.Call("loopFull")
	bg1 = gameAdd.Call("sprite", 0, 0, "bg1")
	floor = gameAdd.Call("sprite", 0, 0, "floor")

	gameWorld.Call("setBounds", 0, 0, 1920, 1080)
	gamePhysics.Call("startSystem", phaser.Get("Physics").Get("P2JS"))

	player = gameAdd.Call("sprite", gameWorld.Get("centerX"), gameWorld.Get("height").Int()-250, "me")
	player.Get("anchor").Call("set", 0.5)
	gamePhysics.Get("p2").Call("enable", player)
	spotlight = gameAdd.Call("sprite", gameWorld.Get("centerX"), gameWorld.Get("centerY"), "spotlight")
	spotlight.Get("anchor").Call("set", 0.5)

	cursors = game.Get("input").Get("keyboard").Call("createCursorKeys")
	gameCamera.Call("follow", player, phaser.Get("Camera").Get("FOLLOW_LOCKON"), 0.1, 0.1)

	playerBody = player.Get("body")
	playerBody.Set("fixedRotation", true)
}

func update() {
	gameCamera.Call("shake", 0.05, 5)
	playerBody.Call("setZeroVelocity")
	spotlight.Set("x", player.Get("x"))
	spotlight.Set("y", player.Get("y"))
	if cursors.Get("up").Get("isDown").Bool() {
		playerBody.Call("moveUp", 125)
	} else if cursors.Get("down").Get("isDown").Bool() {
		playerBody.Call("moveDown", 125)
	}

	if cursors.Get("left").Get("isDown").Bool() {
		playerBody.Call("moveLeft", 250)
	} else if cursors.Get("right").Get("isDown").Bool() {
		playerBody.Call("moveRight", 250)
	}
}

func main() {

}

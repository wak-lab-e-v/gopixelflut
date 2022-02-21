package main

const (
	display_y = 33
	display_x = 60
	bar_x = 3
	bar_y = 13
)
import (
	"fmt"
)

var output_array [[2][3]]
int left_height = 10
int right_height = 10
var black[3]int64
var white[3]int64
white[0],white[1],white[2] = int64(255)

//position the ball in the middle of the playground
int ball_x = int(display_x/2)
int ball_y = int(display_y/2)

bool new_game = true

//draws the left player box
func draw_left_box(left int) {
	left_height = left_height + left
	if left_height+bar_y > display_y {
		left_height = display_y-bar_y
	}
	if left_height <0 {
		left_height = 02
	}
	var arr [[2][3]]
	for x:= 0; x<=bar_x; x++ {
		for y:=left_height; y<left_height+bar_y {
			//add values to the output array
			arr[1] = white
			arr[0] = [x,y]
			append(output_array,arr)
			//reset arr
			arr = nil
		}
	}
}

//draws the right player box
func draw_right_box(right int) {
	right_height = right_height + right
	if right_height+bar_y > display_y {
		right_height = display_y-bar_y
	}
	if right_height <0 {
		right_height = 0
	}
	var arr [[2][3]]
	for x:= display_x-bar_x; x<display_x; x++ {
		for y:=right_height; y<right_height+bar_y {
			//add values to the output array
			arr[0] = [x,y]
			arr[1] = white
			append(output_array,arr)
			//reset arr
			arr = nil
		}
	}
}

//puts the ball at the given coordinates
func draw_ball(x int, y int) {
	var arr = arr [[2][3]]
	for 
}

//restarts the game
func restart_game(left int, right int) {
	left_height = 10
	right_height = 10

	//drawing the boxes
	draw_left_box(left)
	draw_right_box(right)

	new_game = false
}

//main function
//returns the array
func pong_pixel_array(left int, right int) {
	output_array = nil
	if new_game {
		restart_game(left,right)
	} else {

	}

	return output_array
}

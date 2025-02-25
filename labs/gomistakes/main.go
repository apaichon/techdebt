package main

import (
	"gomistakes/lesson"
)

func main() {
	// lesson.GoroutineLeak()
	// lesson.GoroutineLeakWithContext()
	// lesson.SliceLeak()
	// lesson.ClosureLeak()
	// lesson.MapLeak()
	// lesson.ChannelLeak()
	// lesson.TimerLeak()
	// lesson.GlobalVariableLeak()
	// lesson.HTTPBodyLeak()
	lesson.DeferInLoopLeak()
}

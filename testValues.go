package main

import "time"

var targetIsDown = []time.Duration{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
var targetIsUp = []time.Duration{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

var targetIsAlreadyDown = []time.Duration{11, 12, 13, 14, 15, 16, -1, -1, -1, -1}
var targetIsAlreadyUp = []time.Duration{-1, -1, -1, -1, -1, -1, -1, -1, 19, 20}

var targetToDown = []time.Duration{11, 12, 13, 14, 15, 16, 17, -1, -1, -1}
var targetToUp = []time.Duration{-1, -1, -1, -1, -1, -1, -1, -1, -1, 20}

var targetFlappyToDown = []time.Duration{11, -1, 13, -1, -1, -1, 17, -1, -1, -1}
var targetFlappyToUp = []time.Duration{-1, 12, -1, 14, 15, 16, -1, -1, -1, 20}

var targetEmpty = []time.Duration{}
var targetTooShort = []time.Duration{11, 12, 13}

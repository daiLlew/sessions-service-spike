package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/daiLlew/sessions-service-spike/sessions"
)

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	idGenerator := sessions.NewIDGenerator(r)

	sess := sessions.New("test@test.com", idGenerator)

	b, _ := json.MarshalIndent(sess, "", "  ")

	fmt.Printf("%s", string(b))
}
